package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// AdminClaims represents the structure of admin-related JWT claims
type AdminClaims struct {
	UserRole string `json:"user_role"`
}

// extractRoleFromJWT extracts the user role from JWT token by calling Supabase
func extractRoleFromJWT(token string) (string, error) {
	if !supabaseConfigured() {
		return "", errors.New("Supabase not configured")
	}

	url := supabaseBaseURL() + "/auth/v1/user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", supabaseAnonKey())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to validate token: status %d", resp.StatusCode)
	}

	var userData struct {
		AppMetadata struct {
			UserRole string `json:"user_role"`
		} `json:"app_metadata"`
	}

	if err := json.NewDecoder(strings.NewReader(string(body))).Decode(&userData); err != nil {
		return "", err
	}

	// Default to 'user' if role not found
	if userData.AppMetadata.UserRole == "" {
		return "user", nil
	}

	return userData.AppMetadata.UserRole, nil
}

// requireAdmin middleware checks if user has admin role
func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(auth, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "authentication required",
			})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if token == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "authentication required",
			})
			return
		}

		// Extract role from JWT
		role, err := extractRoleFromJWT(token)
		if err != nil {
			log.Printf("Failed to extract role from JWT: %v", err)
			writeJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "authentication failed",
			})
			return
		}

		// Check if user has admin role
		if role != "admin" {
			log.Printf("Access denied: user role '%s' is not admin", role)
			writeJSON(w, http.StatusForbidden, map[string]any{
				"error": "insufficient permissions",
			})
			return
		}

		// User is admin, proceed with request
		next(w, r)
	}
}

// isMigrationMode checks if migration mode is enabled
func isMigrationMode() bool {
	return strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_MIGRATION_MODE"))) == "true"
}

// requireAdminWithFallback supports both JWT and API key authentication during migration
func requireAdminWithFallback(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try JWT authentication first
		auth := strings.TrimSpace(r.Header.Get("Authorization"))
		if strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			token = strings.TrimSpace(token)
			
			if token != "" {
				role, err := extractRoleFromJWT(token)
				
				if err == nil && role == "admin" {
					log.Printf("AUTH_METHOD: jwt, user_role: admin")
					next(w, r)
					return
				}
				
				// JWT auth failed, try fallback if in migration mode
				if !isMigrationMode() {
					writeJSON(w, http.StatusForbidden, map[string]any{
						"error": "insufficient permissions",
					})
					return
				}
			}
		}
		
		// Fallback to API key (only in migration mode)
		if isMigrationMode() {
			adminSecret := strings.TrimSpace(os.Getenv("ADMIN_API_KEY"))
			if adminSecret != "" && strings.TrimSpace(r.Header.Get("X-Admin-Secret")) == adminSecret {
				log.Printf("AUTH_METHOD: api_key (fallback)")
				next(w, r)
				return
			}
		}
		
		// Both methods failed
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "authentication required",
		})
	}
}
