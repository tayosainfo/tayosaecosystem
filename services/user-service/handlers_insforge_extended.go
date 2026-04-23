package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var errOAuthEmailAccountMismatch = errors.New("an account with this email already exists with a different identity; sign in with password or use the same OAuth provider")

func syntheticOAuthPhone(supabaseUserID string) string {
	h := sha256.Sum256([]byte(supabaseUserID + ":tayosa-oauth-phone"))
	n := binary.BigEndian.Uint64(h[:8]) % 100000000
	return fmt.Sprintf("+25678%07d", n)
}

func respondRawJSON(w http.ResponseWriter, status int, raw []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(raw) > 0 {
		_, _ = w.Write(raw)
	}
}

func oauthStartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	if !supabaseConfigured() {
		respond(w, http.StatusServiceUnavailable, map[string]any{"error": "Supabase is not configured"})
		return
	}
	customKey := strings.TrimSpace(r.URL.Query().Get("customKey"))
	provider := strings.TrimSpace(r.URL.Query().Get("provider"))
	redirectURI := strings.TrimSpace(r.URL.Query().Get("redirect_uri"))
	challenge := strings.TrimSpace(r.URL.Query().Get("code_challenge"))
	if redirectURI == "" || challenge == "" {
		respond(w, http.StatusBadRequest, map[string]any{
			"error": "redirect_uri and code_challenge are required (PKCE)",
		})
		return
	}
	q := url.Values{}
	q.Set("redirect_uri", redirectURI)
	q.Set("code_challenge", challenge)
	var path string
	if customKey != "" {
		path = "/auth/v1/authorize?provider=custom:" + url.PathEscape(customKey)
	} else {
		if provider == "" {
			respond(w, http.StatusBadRequest, map[string]any{"error": "provider is required unless customKey is set"})
			return
		}
		path = "/auth/v1/authorize?provider=" + url.PathEscape(provider)
	}
	raw, status, err := supabaseGetAnon(path, q)
	if err != nil {
		code, msg := supabaseUpstreamHTTP(err)
		respond(w, code, map[string]any{"error": msg})
		return
	}
	respondRawJSON(w, status, raw)
}

func oauthExchangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	if !supabaseConfigured() {
		respond(w, http.StatusServiceUnavailable, map[string]any{"error": "Supabase is not configured"})
		return
	}
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	q := url.Values{}
	if ct := strings.TrimSpace(r.URL.Query().Get("client_type")); ct != "" {
		q.Set("client_type", ct)
	}
	// Supabase token exchange: POST /auth/v1/token?grant_type=pkce
	q.Set("grant_type", "pkce")
	resp, status, err := supabasePostWithQuery("/auth/v1/token", q, body)
	if err != nil {
		code, msg := supabaseUpstreamHTTP(err)
		respond(w, code, map[string]any{"error": msg})
		return
	}
	if err := syncOAuthUserFromExchange(resp); err != nil {
		if errors.Is(err, errOAuthEmailAccountMismatch) {
			respond(w, http.StatusConflict, map[string]any{"error": err.Error()})
			return
		}
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	respond(w, status, resp)
}

func syncOAuthUserFromExchange(resp map[string]any) error {
	userObj, ok := resp["user"].(map[string]any)
	if !ok {
		return nil
	}
	id := mapGetString(userObj, "id")
	email := normalizeEmail(mapGetString(userObj, "email"))
	// Supabase stores name in user_metadata
	meta, _ := userObj["user_metadata"].(map[string]any)
	name := strings.TrimSpace(mapGetString(meta, "name"))
	if name == "" {
		name = strings.TrimSpace(mapGetString(meta, "full_name"))
	}
	if name == "" && email != "" {
		if i := strings.IndexByte(email, '@'); i > 0 {
			name = email[:i]
		} else {
			name = email
		}
	}
	if id == "" {
		return nil
	}
	if u, ok := activeStore.FindByUserID(id); ok {
		if name != "" {
			u.FullName = name
		}
		if email != "" {
			u.ContactEmail = email
			u.SupabaseLoginEmail = email
		}
		u.SupabaseUserID = id
		return activeStore.UpdateIdentity(u)
	}
	if email != "" {
		if existing, ok := activeStore.FindByEmailKey(email); ok && existing.ID != id {
			return errOAuthEmailAccountMismatch
		}
	}
	phone := syntheticOAuthPhone(id)
	u := User{
		ID:                  id,
		FullName:            name,
		PhoneE164:           phone,
		AuthEmail:           email,
		ContactEmail:        email,
		SupabaseLoginEmail:  email,
		SupabaseUserID:      id,
		ContactEmailChecked: true,
		CreatedAt:           time.Now(),
	}
	ob := OnboardingProfile{
		UserID:         id,
		Phase:          1,
		TrustScoreSeed: 10,
		LastUpdatedAt:  time.Now(),
	}
	return activeStore.CreateIdentityWithOnboarding(u, ob, nil)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	if !supabaseConfigured() {
		respond(w, http.StatusServiceUnavailable, map[string]any{"error": "Supabase is not configured"})
		return
	}
	q := url.Values{}
	if ct := strings.TrimSpace(r.URL.Query().Get("client_type")); ct != "" {
		q.Set("client_type", ct)
	}
	raw, status, err := supabaseRefreshForward(r, q)
	if err != nil {
		code, msg := supabaseUpstreamHTTP(err)
		respond(w, code, map[string]any{"error": msg})
		return
	}
	respondRawJSON(w, status, raw)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	if !supabaseConfigured() {
		respond(w, http.StatusServiceUnavailable, map[string]any{"error": "Supabase is not configured"})
		return
	}
	raw, status, err := supabaseLogoutForward(r)
	if err != nil {
		code, msg := supabaseUpstreamHTTP(err)
		respond(w, code, map[string]any{"error": msg})
		return
	}
	respondRawJSON(w, status, raw)
}

func publicConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	// Return a static config since Supabase doesn't have a direct equivalent
	respond(w, http.StatusOK, map[string]any{
		"oAuthProviders":            []any{},
		"customOAuthProviders":      []any{},
		"requireEmailVerification":  true,
		"passwordMinLength":         6,
		"verifyEmailMethod":         "code",
		"resetPasswordMethod":       "code",
	})
}

func profilePatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		respond(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	if !supabaseConfigured() {
		respond(w, http.StatusServiceUnavailable, map[string]any{"error": "Supabase is not configured"})
		return
	}
	token := bearerToken(r)
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	// Supabase: PUT /auth/v1/user to update user metadata
	sbBody := map[string]any{"data": body}
	out, status, err := supabaseUserPut("/auth/v1/user", token, sbBody)
	if err != nil {
		code, msg := supabaseUpstreamHTTP(err)
		respond(w, code, map[string]any{"error": msg})
		return
	}
	uid := authedUserID(r)
	if u, ok := activeStore.FindByUserID(uid); ok {
		if meta, ok := out["user_metadata"].(map[string]any); ok {
			if n := strings.TrimSpace(mapGetString(meta, "name")); n != "" {
				u.FullName = n
			}
			_ = activeStore.UpdateIdentity(u)
		}
	}
	respond(w, status, out)
}
