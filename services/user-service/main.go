package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id, X-Request-Id, X-CSRF-Token, X-Admin-Secret")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loadDotEnv() {
	wd, err := os.Getwd()
	if err != nil {
		_ = godotenv.Overload(".env")
		return
	}
	var paths []string
	for d := wd; ; d = filepath.Dir(d) {
		p := filepath.Join(d, ".env")
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			paths = append(paths, p)
		}
		if filepath.Dir(d) == d {
			break
		}
	}
	if len(paths) == 0 {
		_ = godotenv.Overload(".env")
		return
	}
	// Furthest ancestor first, then toward cwd — Overload so more specific .env wins.
	for i := len(paths) - 1; i >= 0; i-- {
		if err := godotenv.Overload(paths[i]); err != nil {
			log.Printf("user-service: loading %s: %v", paths[i], err)
		} else {
			log.Printf("user-service: loaded env file %s", paths[i])
		}
	}
}

func main() {
	loadDotEnv()

	initStore()
	defer activeStore.Close()

	if supabaseConfigured() {
		log.Printf("user-service: Supabase auth enabled (%s)", supabaseBaseURL())
		if supabaseServiceRoleConfigured() {
			log.Print("user-service: SUPABASE_SERVICE_ROLE_KEY is set (admin operations available)")
		} else {
			log.Print("user-service: SUPABASE_SERVICE_ROLE_KEY not set — admin user lookup disabled")
		}
	} else {
		log.Print("user-service: Supabase auth disabled (set SUPABASE_URL and SUPABASE_ANON_KEY)")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]any{
			"status":  "active",
			"service": "user-service",
			"time":    time.Now().Format(time.RFC3339),
		}
		if supabaseConfigured() {
			payload["supabase"] = map[string]any{
				"baseUrl":            supabaseBaseURL(),
				"serviceRoleLoaded":  supabaseServiceRoleConfigured(),
			}
		}
		respond(w, http.StatusOK, payload)
	})
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := activeStore.Ping(); err != nil {
			respond(w, http.StatusServiceUnavailable, map[string]any{
				"status": "not_ready",
				"error":  err.Error(),
			})
			return
		}
		respond(w, http.StatusOK, map[string]any{"status": "ready", "service": "user-service"})
	})

	mux.HandleFunc("POST /api/v1/auth/register", registerHandler)
	mux.HandleFunc("POST /api/v1/auth/login", loginHandler)
	mux.HandleFunc("POST /api/v1/auth/resend-verification-email", resendVerificationHandler)
	mux.HandleFunc("POST /api/v1/auth/verify-email", verifyEmailHandler)
	mux.HandleFunc("POST /api/v1/auth/send-reset-password-email", sendResetPasswordEmailHandler)
	mux.HandleFunc("POST /api/v1/auth/exchange-reset-password-token", exchangeResetPasswordTokenHandler)
	mux.HandleFunc("POST /api/v1/auth/reset-password", resetPasswordHandler)
	mux.HandleFunc("GET /api/v1/auth/oauth/start", oauthStartHandler)
	mux.HandleFunc("POST /api/v1/auth/oauth/exchange", oauthExchangeHandler)
	mux.HandleFunc("POST /api/v1/auth/refresh", refreshHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", logoutHandler)
	mux.HandleFunc("GET /api/v1/auth/public-config", publicConfigHandler)
	mux.HandleFunc("PATCH /api/v1/auth/profile", requireAuth(profilePatchHandler))
	mux.HandleFunc("GET /api/v1/users/me", requireAuth(meHandler))
	mux.HandleFunc("GET /api/v1/onboarding/phase", requireAuth(onboardingHandler))
	mux.HandleFunc("POST /api/v1/onboarding/kyc", requireAuth(onboardingKYCHandler))
	mux.HandleFunc("POST /api/v1/onboarding/sacco", requireAuth(onboardingSaccoHandler))
	mux.HandleFunc("POST /api/v1/onboarding/kibiina", requireAuth(onboardingKibiinaHandler))
	mux.HandleFunc("GET /api/v1/admin/kyc", adminKYCDecisionHandler)
	mux.HandleFunc("POST /api/v1/admin/kyc", adminKYCDecisionHandler)
	mux.HandleFunc("GET /api/v1/admin/settings", adminSettingsHandler)
	mux.HandleFunc("PATCH /api/v1/admin/settings", adminSettingsHandler)
	mux.HandleFunc("POST /api/v1/admin/check-status", adminCheckStatusHandler)
	// User management admin endpoints - wrapped with requireAuth
	mux.HandleFunc("GET /api/v1/admin/users", requireAuth(adminUsersListHandler))
	mux.HandleFunc("GET /api/v1/admin/users/{userId}", requireAuth(adminUserDetailHandler))
	mux.HandleFunc("PATCH /api/v1/admin/users/{userId}/status", requireAuth(adminUserStatusHandler))
	mux.HandleFunc("PATCH /api/v1/admin/users/{userId}/role", requireAuth(adminUserRoleHandler))
	mux.HandleFunc("POST /api/v1/admin/users/{userId}/reset-password", requireAuth(adminUserResetPasswordHandler))
	mux.HandleFunc("GET /api/v1/admin/users/{userId}/activity", requireAuth(adminUserActivityHandler))
	mux.HandleFunc("GET /api/v1/geo", geoLookupHandler)
	mux.HandleFunc("GET /api/v1/groups/policy", parishGroupPolicyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("user-service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, allowCORS(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
