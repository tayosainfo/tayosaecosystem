package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

var proxyHTTP = &http.Client{Timeout: 25 * time.Second}

func trimBase(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), "/")
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return trimBase(v)
	}
	return trimBase(fallback)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id, X-Request-Id, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func hasBearer(r *http.Request) bool {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.HasPrefix(auth, "Bearer ") && strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")) != ""
}

func passthrough(w http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()
	ct := resp.Header.Get("Content-Type")
	if ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func forward(w http.ResponseWriter, r *http.Request, targetBase string) {
	target := targetBase + r.URL.Path
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}
	req, err := http.NewRequestWithContext(r.Context(), r.Method, target, r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	if auth := r.Header.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	if rid := r.Header.Get("X-Request-Id"); rid != "" {
		req.Header.Set("X-Request-Id", rid)
	}
	resp, err := proxyHTTP.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	passthrough(w, resp)
}

func healthReadyHandler(w http.ResponseWriter, r *http.Request) {
	userBase := envOr("USER_SERVICE_URL", "http://localhost:8081")
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, userBase+"/health/ready", nil)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "not_ready", "error": err.Error()})
		return
	}
	resp, err := proxyHTTP.Do(req)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "not_ready", "user_service": err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "not_ready", "user_service": resp.Status})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ready",
		"service": "api-gateway-service",
		"time":    time.Now().Format(time.RFC3339),
	})
}

type route struct {
	path   string
	base   string
	public bool
}

func main() {
	_ = godotenv.Load(filepath.Join("..", "..", ".env"))
	_ = godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	userBase := envOr("USER_SERVICE_URL", "http://localhost:8081")
	storageBase := envOr("OBJECT_STORAGE_SERVICE_URL", "http://localhost:8015")
	affiliateBase := envOr("AFFILIATE_SERVICE_URL", "http://localhost:8016")
	notifyBase := envOr("NOTIFICATION_SERVICE_URL", "http://localhost:8010")
	auditBase := envOr("AUDIT_SERVICE_URL", "http://localhost:8014")
	loanBase := envOr("LOAN_SERVICE_URL", "http://localhost:8013")
	feeBase := envOr("FEE_SERVICE_URL", "http://localhost:8004")
	kibiinaBase := envOr("KIBIINA_SERVICE_URL", "http://localhost:8086")

	routes := []route{
		// user-service
		{path: "/api/v1/auth/register", base: userBase, public: true},
		{path: "/api/v1/auth/login", base: userBase, public: true},
		{path: "/api/v1/auth/resend-verification-email", base: userBase, public: true},
		{path: "/api/v1/auth/verify-email", base: userBase, public: true},
		{path: "/api/v1/auth/send-reset-password-email", base: userBase, public: true},
		{path: "/api/v1/auth/exchange-reset-password-token", base: userBase, public: true},
		{path: "/api/v1/auth/reset-password", base: userBase, public: true},
		{path: "/api/v1/auth/oauth/start", base: userBase, public: true},
		{path: "/api/v1/auth/oauth/exchange", base: userBase, public: true},
		{path: "/api/v1/auth/refresh", base: userBase, public: true},
		{path: "/api/v1/auth/logout", base: userBase, public: true},
		{path: "/api/v1/auth/public-config", base: userBase, public: true},
		{path: "/api/v1/auth/profile", base: userBase, public: false},
		{path: "/api/v1/users/me", base: userBase, public: false},
		{path: "/api/v1/onboarding/phase", base: userBase, public: false},
		{path: "/api/v1/geo", base: userBase, public: true},
		{path: "/api/v1/groups/policy", base: userBase, public: true},
		// object storage
		{path: "/api/v1/storage/upload-url", base: storageBase, public: false},
		{path: "/api/v1/storage/object", base: storageBase, public: false},
		// affiliate
		{path: "/api/v1/affiliate/referrals", base: affiliateBase, public: false},
		{path: "/api/v1/affiliate/rewards", base: affiliateBase, public: false},
		// notifications
		{path: "/api/v1/notifications/send", base: notifyBase, public: false},
		{path: "/api/v1/notifications/outbox", base: notifyBase, public: false},
		// audit
		{path: "/api/v1/audit/events", base: auditBase, public: false},
		// loans
		{path: "/api/v1/loans/score", base: loanBase, public: false},
		// fees & kibiina
		{path: "/api/v1/fees/transactions", base: feeBase, public: false},
		{path: "/api/v1/kibiina/groups", base: kibiinaBase, public: false},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "api-gateway-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("/health/ready", healthReadyHandler)

	for _, rt := range routes {
		rt := rt
		mux.HandleFunc(rt.path, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != rt.path {
				http.NotFound(w, r)
				return
			}
			if !rt.public && !hasBearer(r) {
				writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing bearer token"})
				return
			}
			forward(w, r, rt.base)
		})
	}

	log.Printf("API gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
