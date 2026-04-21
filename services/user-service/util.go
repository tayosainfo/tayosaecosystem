package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

func normalizePhone(raw string) (string, error) {
	digitsOnly := regexp.MustCompile(`\D`).ReplaceAllString(strings.TrimSpace(raw), "")
	switch {
	case strings.HasPrefix(digitsOnly, "256") && len(digitsOnly) == 12:
		return "+" + digitsOnly, nil
	case strings.HasPrefix(digitsOnly, "0") && len(digitsOnly) == 10:
		return "+256" + digitsOnly[1:], nil
	default:
		return "", errors.New("phone must be Uganda format")
	}
}

func normalizeEmail(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func makeAuthEmail(phoneE164 string) string {
	digits := strings.TrimPrefix(phoneE164, "+")
	return digits + "@tayosa.local"
}

func insforgeLoginEmail(u User) string {
	if u.InsforgeEmail != "" {
		return u.InsforgeEmail
	}
	if u.ContactEmail != "" {
		return u.ContactEmail
	}
	return u.AuthEmail
}

// parseOptionalDateOfBirth accepts YYYY-MM-DD or empty (no birth date stored).
func parseOptionalDateOfBirth(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, errors.New("dateOfBirth must be YYYY-MM-DD")
	}
	return &t, nil
}

func normalizeNationality(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", nil
	}
	if utf8.RuneCountInString(s) > 64 {
		return "", errors.New("nationality must be at most 64 characters")
	}
	return s, nil
}

func mapGetString(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return fmt.Sprint(v)
}

func mapGetBool(m map[string]any, key string) bool {
	if m == nil {
		return false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func referralCodeForUserID(userID string) string {
	sum := sha1.Sum([]byte("tayosa:" + strings.TrimSpace(userID)))
	return "TAY-" + strings.ToUpper(hex.EncodeToString(sum[:]))[:8]
}

// extractInsForgeSignupUserID reads user id from varied InsForge POST /api/auth/users response shapes.
func extractInsForgeSignupUserID(resp map[string]any) string {
	if resp == nil {
		return ""
	}
	var walkUser func(any) string
	walkUser = func(u any) string {
		switch m := u.(type) {
		case map[string]any:
			return firstNonEmpty(
				mapGetString(m, "id"),
				mapGetString(m, "userId"),
				mapGetString(m, "user_id"),
				mapGetString(m, "sub"),
			)
		case []any:
			for _, el := range m {
				if id := walkUser(el); id != "" {
					return id
				}
			}
		}
		return ""
	}
	if id := walkUser(resp["user"]); id != "" {
		return id
	}
	if d, ok := resp["data"].(map[string]any); ok {
		if id := walkUser(d["user"]); id != "" {
			return id
		}
		if id := firstNonEmpty(mapGetString(d, "id"), mapGetString(d, "userId"), mapGetString(d, "user_id")); id != "" {
			return id
		}
	}
	return firstNonEmpty(
		mapGetString(resp, "userId"),
		mapGetString(resp, "user_id"),
		mapGetString(resp, "id"),
	)
}

// topLevelJSONKeys lists response keys for debugging InsForge shape mismatches.
func topLevelJSONKeys(m map[string]any) []string {
	if m == nil {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// insforgeIndicatesEmailNotVerified is true when InsForge rejected a session because the address is not verified yet.
func insforgeIndicatesEmailNotVerified(httpStatus int, msg string) bool {
	m := strings.ToLower(strings.TrimSpace(msg))
	if m == "" {
		return false
	}
	if strings.Contains(m, "email_not_verified") || strings.Contains(m, "email not verified") {
		return true
	}
	// Avoid treating a plain wrong-password 401 as "needs email verify".
	if httpStatus == 401 {
		if (strings.Contains(m, "invalid") || strings.Contains(m, "incorrect")) &&
			(strings.Contains(m, "password") || strings.Contains(m, "credentials")) &&
			!strings.Contains(m, "verif") && !strings.Contains(m, "confirm") && !strings.Contains(m, "verified") {
			return false
		}
	}
	if strings.Contains(m, "verify your email") || strings.Contains(m, "must verify your email") {
		return true
	}
	if strings.Contains(m, "email verification") && (httpStatus == 401 || httpStatus == 403 || httpStatus == 400) {
		return true
	}
	if strings.Contains(m, "email address is not verified") || strings.Contains(m, "not been verified") {
		return true
	}
	if httpStatus == 403 && strings.Contains(m, "verif") {
		return true
	}
	return false
}

// insforgeSignupHasAccessToken is false when InsForge withheld tokens until email verification.
func insforgeSignupHasAccessToken(m map[string]any) bool {
	if m == nil {
		return false
	}
	v, ok := m["accessToken"]
	if !ok || v == nil {
		return false
	}
	s := strings.TrimSpace(fmt.Sprint(v))
	return s != "" && !strings.EqualFold(s, "null")
}

// sessionPayloadFromInsForge builds the gateway-facing session object, including refreshToken
// and csrfToken when InsForge returns them (non-web client_type).
func sessionPayloadFromInsForge(resp map[string]any, userID string) map[string]any {
	out := map[string]any{
		"accessToken": mapGetString(resp, "accessToken"),
		"userId":      userID,
	}
	if rt := mapGetString(resp, "refreshToken"); rt != "" {
		out["refreshToken"] = rt
	}
	if cs := mapGetString(resp, "csrfToken"); cs != "" {
		out["csrfToken"] = cs
	}
	return out
}
