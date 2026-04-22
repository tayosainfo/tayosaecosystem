package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type ctxKey int

const ctxAuthedUserID ctxKey = 1

type authResult struct {
	UserID string
}

// authFromRequest validates the Bearer token by calling InsForge.
// dev-token-* and local fallbacks have been removed — InsForge is required.
func authFromRequest(r *http.Request) (authResult, error) {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer ") {
		return authResult{}, errors.New("missing bearer token")
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	if token == "" {
		return authResult{}, errors.New("missing bearer token")
	}

	if !insforgeConfigured() {
		return authResult{}, errors.New("auth backend not configured")
	}

	sess, _, err := insforgeUserGet("/api/auth/sessions/current", token)
	if err != nil {
		return authResult{}, errors.New("invalid or expired session")
	}
	userObj, _ := sess["user"].(map[string]any)
	id := mapGetString(userObj, "id")
	if id == "" {
		return authResult{}, errors.New("invalid session")
	}
	return authResult{UserID: id}, nil
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac, err := authFromRequest(r)
		if err != nil {
			respond(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
			return
		}
		ctx := context.WithValue(r.Context(), ctxAuthedUserID, ac.UserID)
		next(w, r.WithContext(ctx))
	}
}

func authedUserID(r *http.Request) string {
	v, _ := r.Context().Value(ctxAuthedUserID).(string)
	return v
}

func bearerToken(r *http.Request) string {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
}
