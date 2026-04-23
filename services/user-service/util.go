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

// supabaseLoginEmail returns the email used to authenticate with Supabase.
// Email is now always the real contact email — phone aliases have been removed.
func supabaseLoginEmail(u User) string {
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

// topLevelJSONKeys lists response keys for debugging shape mismatches.
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
