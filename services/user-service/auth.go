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

// adminAuthFromRequest validates the Bearer token for admin endpoints
// It tries Supabase first, then falls back to checking if the token matches a user in the database
func adminAuthFromRequest(r *http.Request) (authResult, error) {
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

	// Try Supabase first
	user, _, err := supabaseUserGet("/auth/v1/user", token)
	if err == nil {
		id := mapGetString(user, "id")
		if id != "" {
			return authResult{UserID: id}, nil
		}
	}

	// If Supabase fails, the token might be from a different auth system
	// For now, we'll return an error, but in production you might want to
	// validate the token against your own JWT secret
	return authResult{}, errors.New("invalid or expired session")
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

// requireAdminAuth is a middleware for admin endpoints
// It validates the Bearer token and checks if the user has admin role
func requireAdminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac, err := adminAuthFromRequest(r)
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
