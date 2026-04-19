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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id, X-Request-Id, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Repo root .env, then optional local override (services/user-service/.env).
	_ = godotenv.Load(filepath.Join("..", "..", ".env"))
	_ = godotenv.Load(".env")

	initStore()
	defer activeStore.Close()

	if insforgeConfigured() {
		log.Printf("user-service: InsForge auth enabled (%s)", insforgeBaseURL())
	} else {
		log.Print("user-service: InsForge auth disabled (set INSFORGE_BASE_URL and INSFORGE_ANON_KEY for live codes)")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "user-service",
			"time":    time.Now().Format(time.RFC3339),
		})
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

	mux.HandleFunc("/api/v1/auth/register", registerHandler)
	mux.HandleFunc("/api/v1/auth/login", loginHandler)
	mux.HandleFunc("/api/v1/auth/resend-verification-email", resendVerificationHandler)
	mux.HandleFunc("/api/v1/auth/verify-email", verifyEmailHandler)
	mux.HandleFunc("/api/v1/auth/send-reset-password-email", sendResetPasswordEmailHandler)
	mux.HandleFunc("/api/v1/auth/exchange-reset-password-token", exchangeResetPasswordTokenHandler)
	mux.HandleFunc("/api/v1/auth/reset-password", resetPasswordHandler)
	mux.HandleFunc("GET /api/v1/auth/oauth/start", oauthStartHandler)
	mux.HandleFunc("POST /api/v1/auth/oauth/exchange", oauthExchangeHandler)
	mux.HandleFunc("POST /api/v1/auth/refresh", refreshHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", logoutHandler)
	mux.HandleFunc("GET /api/v1/auth/public-config", publicConfigHandler)
	mux.HandleFunc("PATCH /api/v1/auth/profile", requireAuth(profilePatchHandler))
	mux.HandleFunc("/api/v1/users/me", requireAuth(meHandler))
	mux.HandleFunc("/api/v1/onboarding/phase", requireAuth(onboardingHandler))
	mux.HandleFunc("/api/v1/geo", geoLookupHandler)
	mux.HandleFunc("/api/v1/groups/policy", parishGroupPolicyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("user-service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, allowCORS(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
