package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Configuration helpers
// ---------------------------------------------------------------------------

func supabaseConfigured() bool {
	return supabaseBaseURL() != "" && supabaseAnonKey() != ""
}

func supabaseBaseURL() string {
	return strings.TrimRight(strings.TrimSpace(os.Getenv("SUPABASE_URL")), "/")
}

func supabaseAnonKey() string {
	return strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY"))
}

func supabaseServiceRoleKey() string {
	s := strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	s = strings.TrimPrefix(s, "\ufeff") // UTF-8 BOM
	return s
}

func supabaseServiceRoleConfigured() bool {
	return supabaseServiceRoleKey() != "" && supabaseBaseURL() != ""
}

func supabaseHTTPClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

// ---------------------------------------------------------------------------
// Error type for Supabase API errors
// ---------------------------------------------------------------------------

type SupabaseRequestError struct {
	Status  int
	Message string
}

func (e *SupabaseRequestError) Error() string {
	return e.Message
}

func supabaseErrorMessage(body []byte, status int) string {
	var top map[string]any
	if json.Unmarshal(body, &top) == nil {
		if m, ok := top["msg"].(string); ok && strings.TrimSpace(m) != "" {
			return m
		}
		if m, ok := top["message"].(string); ok && strings.TrimSpace(m) != "" {
			return m
		}
		if m, ok := top["error_description"].(string); ok && strings.TrimSpace(m) != "" {
			return m
		}
		if s, ok := top["error"].(string); ok && strings.TrimSpace(s) != "" {
			return s
		}
	}
	if len(body) > 0 && len(body) < 512 {
		return strings.TrimSpace(string(body))
	}
	return fmt.Sprintf("Supabase request failed (%d)", status)
}

// ---------------------------------------------------------------------------
// URL builder
// ---------------------------------------------------------------------------

func supabaseURL(path string, q url.Values) string {
	u := supabaseBaseURL() + path
	if q != nil && len(q) > 0 {
		u += "?" + q.Encode()
	}
	return u
}

// ---------------------------------------------------------------------------
// Core HTTP dispatcher
// ---------------------------------------------------------------------------

// supabaseDoJSON performs a Supabase request.
// auth: "anon" (apikey + Bearer anon), "user" (apikey + Bearer user token),
//       "service" (apikey + Bearer service role), "none" (apikey only, no Bearer).
func supabaseDoJSON(method, path string, q url.Values, body any, auth string, userToken string, extra http.Header) ([]byte, int, error) {
	base := supabaseBaseURL()
	key := supabaseAnonKey()
	if base == "" || key == "" {
		return nil, 0, errors.New("Supabase is not configured")
	}

	var bodyReader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, supabaseURL(path, q), bodyReader)
	if err != nil {
		return nil, 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	// Supabase always requires the apikey header.
	req.Header.Set("apikey", key)

	switch auth {
	case "anon":
		req.Header.Set("Authorization", "Bearer "+key)
	case "user":
		tok := strings.TrimSpace(userToken)
		if tok == "" {
			return nil, 0, errors.New("missing access token")
		}
		req.Header.Set("Authorization", "Bearer "+tok)
	case "service":
		svc := supabaseServiceRoleKey()
		if svc == "" {
			return nil, 0, errors.New("Supabase service role key is not configured")
		}
		req.Header.Set("apikey", svc)
		req.Header.Set("Authorization", "Bearer "+svc)
	}
	if extra != nil {
		for k, vv := range extra {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}

	resp, err := supabaseHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, &SupabaseRequestError{
			Status:  resp.StatusCode,
			Message: supabaseErrorMessage(respBody, resp.StatusCode),
		}
	}
	return respBody, resp.StatusCode, nil
}

// ---------------------------------------------------------------------------
// Convenience wrappers
// ---------------------------------------------------------------------------

func supabasePost(path string, body any) (map[string]any, int, error) {
	return supabasePostWithQuery(path, nil, body)
}

func supabasePostWithQuery(path string, q url.Values, body any) (map[string]any, int, error) {
	raw, status, err := supabaseDoJSON(http.MethodPost, path, q, body, "anon", "", nil)
	if err != nil {
		return nil, status, err
	}
	var out map[string]any
	if len(raw) == 0 {
		return map[string]any{}, status, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, status, err
	}
	return out, status, nil
}

// supabasePostAnyStatus POSTs with the anon key and returns JSON + status even on 4xx/5xx.
func supabasePostAnyStatus(path string, q url.Values, body any) (map[string]any, int, error) {
	base := supabaseBaseURL()
	key := supabaseAnonKey()
	if base == "" || key == "" {
		return nil, 0, errors.New("Supabase is not configured")
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest(http.MethodPost, supabaseURL(path, q), bytes.NewReader(raw))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := supabaseHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	var out map[string]any
	if len(respBody) > 0 {
		_ = json.Unmarshal(respBody, &out)
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, resp.StatusCode, nil
}

func supabaseGetAnon(path string, q url.Values) ([]byte, int, error) {
	return supabaseDoJSON(http.MethodGet, path, q, nil, "anon", "", nil)
}

// supabaseUserGet calls a Supabase route with the end-user access token.
func supabaseUserGet(path string, userAccessToken string) (map[string]any, int, error) {
	raw, status, err := supabaseDoJSON(http.MethodGet, path, nil, nil, "user", userAccessToken, nil)
	if err != nil {
		return nil, status, err
	}
	var out map[string]any
	if len(raw) == 0 {
		return map[string]any{}, status, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, status, err
	}
	return out, status, nil
}

// supabaseUserPut calls PUT with the user's Bearer token (e.g. password update).
func supabaseUserPut(path string, userAccessToken string, body any) (map[string]any, int, error) {
	raw, status, err := supabaseDoJSON(http.MethodPut, path, nil, body, "user", userAccessToken, nil)
	if err != nil {
		return nil, status, err
	}
	var out map[string]any
	if len(raw) == 0 {
		return map[string]any{}, status, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, status, err
	}
	return out, status, nil
}

// ---------------------------------------------------------------------------
// Verification & password reset (service-role calls)
// ---------------------------------------------------------------------------

// supabaseResendVerification calls POST /auth/v1/otp for OTP-based verification.
func supabaseResendVerification(q url.Values, body map[string]any) (map[string]any, int, error) {
	// Use OTP endpoint for better rate limits and code-based verification
	payload := map[string]any{
		"type":  "signup",
		"email": body["email"],
	}
	return supabasePostWithQuery("/auth/v1/otp", q, payload)
}

// ---------------------------------------------------------------------------
// Refresh & Logout forwarding
// ---------------------------------------------------------------------------

func supabaseRefreshForward(r *http.Request, q url.Values) ([]byte, int, error) {
	base := supabaseBaseURL()
	key := supabaseAnonKey()
	if base == "" || key == "" {
		return nil, 0, errors.New("Supabase is not configured")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, 0, err
	}
	// Supabase refresh: POST /auth/v1/token?grant_type=refresh_token
	rq := url.Values{}
	rq.Set("grant_type", "refresh_token")
	for k, v := range q {
		rq[k] = v
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, supabaseURL("/auth/v1/token", rq), bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	if ct := r.Header.Get("Content-Type"); ct != "" {
		req.Header.Set("Content-Type", ct)
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := supabaseHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, &SupabaseRequestError{
			Status:  resp.StatusCode,
			Message: supabaseErrorMessage(respBody, resp.StatusCode),
		}
	}
	return respBody, resp.StatusCode, nil
}

func supabaseLogoutForward(r *http.Request) ([]byte, int, error) {
	base := supabaseBaseURL()
	key := supabaseAnonKey()
	if base == "" || key == "" {
		return nil, 0, errors.New("Supabase is not configured")
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, supabaseURL("/auth/v1/logout", nil), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("apikey", key)
	// Forward the user's token for logout
	if auth := r.Header.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	} else {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	resp, err := supabaseHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, &SupabaseRequestError{
			Status:  resp.StatusCode,
			Message: supabaseErrorMessage(respBody, resp.StatusCode),
		}
	}
	return respBody, resp.StatusCode, nil
}

// ---------------------------------------------------------------------------
// Admin / service-role operations
// ---------------------------------------------------------------------------

func supabaseAdminPost(path string, q url.Values, body any) (map[string]any, int, error) {
	raw, status, err := supabaseDoJSON(http.MethodPost, path, q, body, "service", "", nil)
	if err != nil {
		return nil, status, err
	}
	var out map[string]any
	if len(raw) == 0 {
		return map[string]any{}, status, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, status, err
	}
	return out, status, nil
}

func supabaseAdminGET(pathWithQuery string) ([]byte, int, error) {
	svc := supabaseServiceRoleKey()
	base := supabaseBaseURL()
	if svc == "" || base == "" {
		return nil, 0, errors.New("Supabase service role key is not configured")
	}
	p := pathWithQuery
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	req, err := http.NewRequest(http.MethodGet, base+p, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("apikey", svc)
	req.Header.Set("Authorization", "Bearer "+svc)
	resp, err := supabaseHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, &SupabaseRequestError{
			Status:  resp.StatusCode,
			Message: supabaseErrorMessage(respBody, resp.StatusCode),
		}
	}
	return respBody, resp.StatusCode, nil
}

// supabaseAdminFindUserIDByEmail uses the admin API to find a user by email.
func supabaseAdminFindUserIDByEmail(email string) (string, error) {
	email = normalizeEmail(email)
	if email == "" {
		return "", errors.New("empty email")
	}
	// Supabase admin: GET /auth/v1/admin/users
	raw, _, err := supabaseAdminGET("/auth/v1/admin/users")
	if err != nil {
		return "", err
	}
	// Supabase returns { "users": [...], "aud": "...", "total": N }
	var wrap struct {
		Users []map[string]any `json:"users"`
	}
	if json.Unmarshal(raw, &wrap) == nil {
		for _, u := range wrap.Users {
			if strings.EqualFold(normalizeEmail(mapGetString(u, "email")), email) {
				return mapGetString(u, "id"), nil
			}
		}
	}
	// Fallback: try as array
	var asList []map[string]any
	if json.Unmarshal(raw, &asList) == nil {
		for _, u := range asList {
			if strings.EqualFold(normalizeEmail(mapGetString(u, "email")), email) {
				return mapGetString(u, "id"), nil
			}
		}
	}
	return "", errors.New("admin user list did not contain matching email")
}

// ---------------------------------------------------------------------------
// Supabase-specific response helpers
// ---------------------------------------------------------------------------

// supabaseUpstreamHTTP maps errors from Supabase to an HTTP status and message.
func supabaseUpstreamHTTP(err error) (status int, msg string) {
	var sbErr *SupabaseRequestError
	if errors.As(err, &sbErr) {
		s := sbErr.Status
		if s >= 400 && s <= 599 {
			return s, sbErr.Message
		}
		return http.StatusBadGateway, sbErr.Message
	}
	return http.StatusBadGateway, err.Error()
}

// supabaseSignupHasAccessToken checks if the signup response includes a session.
func supabaseSignupHasAccessToken(m map[string]any) bool {
	if m == nil {
		return false
	}
	v, ok := m["access_token"]
	if !ok || v == nil {
		return false
	}
	s := strings.TrimSpace(fmt.Sprint(v))
	return s != "" && !strings.EqualFold(s, "null")
}

// supabaseIndicatesEmailNotVerified checks if the error indicates unverified email.
func supabaseIndicatesEmailNotVerified(httpStatus int, msg string) bool {
	m := strings.ToLower(strings.TrimSpace(msg))
	if m == "" {
		return false
	}
	if strings.Contains(m, "email not confirmed") {
		return true
	}
	if strings.Contains(m, "email_not_confirmed") {
		return true
	}
	if strings.Contains(m, "confirm your email") || strings.Contains(m, "confirmation") {
		if httpStatus == 400 || httpStatus == 401 || httpStatus == 403 {
			return true
		}
	}
	return false
}

// extractSupabaseSignupUserID reads user id from varied Supabase signup response shapes.
func extractSupabaseSignupUserID(resp map[string]any) string {
	if resp == nil {
		return ""
	}
	// Supabase signup with confirmation: top-level { "id": "...", "email": "..." }
	if id := mapGetString(resp, "id"); id != "" {
		return id
	}
	// Supabase signup without confirmation: { "user": { "id": "..." }, "access_token": "..." }
	if userObj, ok := resp["user"].(map[string]any); ok {
		if id := mapGetString(userObj, "id"); id != "" {
			return id
		}
	}
	return ""
}

// sessionPayloadFromSupabase builds the gateway-facing session object from a Supabase response.
func sessionPayloadFromSupabase(resp map[string]any, userID string) map[string]any {
	out := map[string]any{
		"accessToken": mapGetString(resp, "access_token"),
		"userId":      userID,
	}
	if rt := mapGetString(resp, "refresh_token"); rt != "" {
		out["refreshToken"] = rt
	}
	return out
}
