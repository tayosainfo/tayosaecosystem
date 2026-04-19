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

func insforgeConfigured() bool {
	base := strings.TrimSpace(os.Getenv("INSFORGE_BASE_URL"))
	key := strings.TrimSpace(os.Getenv("INSFORGE_ANON_KEY"))
	return base != "" && key != ""
}

func insforgeBaseURL() string {
	return strings.TrimRight(strings.TrimSpace(os.Getenv("INSFORGE_BASE_URL")), "/")
}

func insforgeAnonKey() string {
	return strings.TrimSpace(os.Getenv("INSFORGE_ANON_KEY"))
}

func insforgeHTTPClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

// InsforgeRequestError is returned when InsForge responds with a non-2xx status.
type InsforgeRequestError struct {
	Status  int
	Message string
}

func (e *InsforgeRequestError) Error() string {
	return e.Message
}

func insforgeErrorMessage(body []byte, status int) string {
	var top map[string]any
	if json.Unmarshal(body, &top) == nil {
		if m, ok := top["message"].(string); ok && strings.TrimSpace(m) != "" {
			return m
		}
		if s, ok := top["error"].(string); ok && strings.TrimSpace(s) != "" {
			return s
		}
		if errObj, ok := top["error"].(map[string]any); ok {
			if m, ok := errObj["message"].(string); ok && strings.TrimSpace(m) != "" {
				return m
			}
		}
	}
	if len(body) > 0 && len(body) < 512 {
		return strings.TrimSpace(string(body))
	}
	return fmt.Sprintf("InsForge request failed (%d)", status)
}

func insforgeURL(path string, q url.Values) string {
	u := insforgeBaseURL() + path
	if q != nil && len(q) > 0 {
		u += "?" + q.Encode()
	}
	return u
}

// insforgePost calls an InsForge auth API route with the anon key (same as web SDK).
func insforgePost(path string, body any) (map[string]any, int, error) {
	return insforgePostWithQuery(path, nil, body)
}

func insforgePostWithQuery(path string, q url.Values, body any) (map[string]any, int, error) {
	raw, status, err := insforgeDoJSON(http.MethodPost, path, q, body, "anon", "", nil)
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

// insforgePostAnonAnyHTTPStatus POSTs with the anon key and returns JSON + status even on 4xx/5xx
// (used to read user id from EMAIL_NOT_VERIFIED login responses).
func insforgePostAnonAnyHTTPStatus(path string, q url.Values, body any) (map[string]any, int, error) {
	base := insforgeBaseURL()
	key := insforgeAnonKey()
	if base == "" || key == "" {
		return nil, 0, errors.New("InsForge is not configured")
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest(http.MethodPost, insforgeURL(path, q), bytes.NewReader(raw))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := insforgeHTTPClient().Do(req)
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

// insforgeGetAnon performs a GET with the anon key (e.g. OAuth authorize URL, public-config).
func insforgeGetAnon(path string, q url.Values) ([]byte, int, error) {
	return insforgeDoJSON(http.MethodGet, path, q, nil, "anon", "", nil)
}

// insforgeDoJSON performs an InsForge request and returns the raw JSON body on success.
// auth: "anon" (Bearer anon key), "user" (Bearer access token), "none" (no Authorization).
func insforgeDoJSON(method, path string, q url.Values, body any, auth string, userToken string, extra http.Header) ([]byte, int, error) {
	base := insforgeBaseURL()
	key := insforgeAnonKey()
	if base == "" {
		return nil, 0, errors.New("InsForge is not configured")
	}
	if auth == "anon" && key == "" {
		return nil, 0, errors.New("InsForge is not configured")
	}

	var bodyReader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, insforgeURL(path, q), bodyReader)
	if err != nil {
		return nil, 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}
	switch auth {
	case "anon":
		req.Header.Set("Authorization", "Bearer "+key)
	case "user":
		tok := strings.TrimSpace(userToken)
		if tok == "" {
			return nil, 0, errors.New("missing access token")
		}
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	if extra != nil {
		for k, vv := range extra {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}

	resp, err := insforgeHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, &InsforgeRequestError{
			Status:  resp.StatusCode,
			Message: insforgeErrorMessage(respBody, resp.StatusCode),
		}
	}
	return respBody, resp.StatusCode, nil
}

// insforgeUserGet calls an InsForge route with the end-user access token (not the anon key).
func insforgeUserGet(path string, userAccessToken string) (map[string]any, int, error) {
	raw, status, err := insforgeDoJSON(http.MethodGet, path, nil, nil, "user", userAccessToken, nil)
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

// insforgeUserPatch calls InsForge with the user's Bearer token (e.g. profile update).
func insforgeUserPatch(path string, userAccessToken string, body any) (map[string]any, int, error) {
	raw, status, err := insforgeDoJSON(http.MethodPatch, path, nil, body, "user", userAccessToken, nil)
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

// insforgeRefreshForward proxies refresh: forwards body and optional CSRF / Cookie from the client request.
func insforgeRefreshForward(r *http.Request, q url.Values) ([]byte, int, error) {
	base := insforgeBaseURL()
	key := insforgeAnonKey()
	if base == "" || key == "" {
		return nil, 0, errors.New("InsForge is not configured")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, insforgeURL("/api/auth/refresh", q), bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	if ct := r.Header.Get("Content-Type"); ct != "" {
		req.Header.Set("Content-Type", ct)
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+key)
	if x := r.Header.Get("X-CSRF-Token"); x != "" {
		req.Header.Set("X-CSRF-Token", x)
	}
	if c := r.Header.Get("Cookie"); c != "" {
		req.Header.Set("Cookie", c)
	}

	resp, err := insforgeHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, &InsforgeRequestError{
			Status:  resp.StatusCode,
			Message: insforgeErrorMessage(respBody, resp.StatusCode),
		}
	}
	return respBody, resp.StatusCode, nil
}

// insforgeLogoutForward proxies logout to InsForge (cookies for web, or optional Bearer).
func insforgeLogoutForward(r *http.Request) ([]byte, int, error) {
	base := insforgeBaseURL()
	key := insforgeAnonKey()
	if base == "" || key == "" {
		return nil, 0, errors.New("InsForge is not configured")
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, insforgeURL("/api/auth/logout", nil), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	if auth := r.Header.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	if c := r.Header.Get("Cookie"); c != "" {
		req.Header.Set("Cookie", c)
	}

	resp, err := insforgeHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, &InsforgeRequestError{
			Status:  resp.StatusCode,
			Message: insforgeErrorMessage(respBody, resp.StatusCode),
		}
	}
	return respBody, resp.StatusCode, nil
}

func insforgeAdminAPIKey() string {
	return strings.TrimSpace(os.Getenv("INSFORGE_ADMIN_API_KEY"))
}

func insforgeAdminConfigured() bool {
	return insforgeAdminAPIKey() != "" && insforgeBaseURL() != ""
}

// insforgeAdminGET calls an InsForge admin route (project admin API key or JWT).
func insforgeAdminGET(pathWithQuery string) ([]byte, int, error) {
	key := insforgeAdminAPIKey()
	base := insforgeBaseURL()
	if key == "" || base == "" {
		return nil, 0, errors.New("InsForge admin API key is not configured")
	}
	p := pathWithQuery
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	req, err := http.NewRequest(http.MethodGet, base+p, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := insforgeHTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, &InsforgeRequestError{
			Status:  resp.StatusCode,
			Message: insforgeErrorMessage(respBody, resp.StatusCode),
		}
	}
	return respBody, resp.StatusCode, nil
}

// insforgeAdminFindUserIDByEmail lists auth users and returns the id for an exact email match.
func insforgeAdminFindUserIDByEmail(email string) (string, error) {
	email = normalizeEmail(email)
	if email == "" {
		return "", errors.New("empty email")
	}
	q := url.Values{}
	q.Set("limit", "50")
	q.Set("search", email)
	raw, _, err := insforgeAdminGET("/api/auth/users?" + q.Encode())
	if err != nil {
		return "", err
	}
	var asList []map[string]any
	if json.Unmarshal(raw, &asList) == nil {
		for _, u := range asList {
			if strings.EqualFold(normalizeEmail(mapGetString(u, "email")), email) {
				return mapGetString(u, "id"), nil
			}
		}
	}
	var wrap struct {
		Users []map[string]any `json:"users"`
		Data  []map[string]any `json:"data"`
		Items []map[string]any `json:"items"`
	}
	if json.Unmarshal(raw, &wrap) == nil {
		for _, list := range [][]map[string]any{wrap.Users, wrap.Data, wrap.Items} {
			for _, u := range list {
				if strings.EqualFold(normalizeEmail(mapGetString(u, "email")), email) {
					return mapGetString(u, "id"), nil
				}
			}
		}
	}
	return "", errors.New("admin user list did not contain matching email")
}
