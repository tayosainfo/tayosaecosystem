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

// authFromRequest validates the Bearer token by calling Supabase.
func authFromRequest(r *http.Request) (authResult, error) {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer ") {
		return authResult{}, errors.New("missing bearer token")
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	if token == "" {
		return authResult{}, errors.New("missing bearer token")
	}

	if !supabaseConfigured() {
		return authResult{}, errors.New("auth backend not configured")
	}

	// Supabase: GET /auth/v1/user returns user object directly (not wrapped in { user: ... })
	user, _, err := supabaseUserGet("/auth/v1/user", token)
	if err != nil {
		return authResult{}, errors.New("invalid or expired session")
	}
	id := mapGetString(user, "id")
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
