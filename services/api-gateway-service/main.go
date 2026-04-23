package main

import (
	"encoding/json"
	"errors"
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

// ---------------------------------------------------------------------------
// Supabase Configuration
// ---------------------------------------------------------------------------

func supabaseConfigured() bool {
	return supabaseBaseURL() != "" && supabaseAnonKey() != ""
}

func supabaseBaseURL() string {
	return strings.TrimRight(strings.TrimSpace(os.Getenv("SUPABASE_URL")), "/")
}

func supabaseAnonKey() string {
	return strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY"))
}

// ---------------------------------------------------------------------------
// Supabase Token Validation
// ---------------------------------------------------------------------------

type SupabaseRequestError struct {
	Status  int
	Message string
}

func (e *SupabaseRequestError) Error() string {
	return e.Message
}

func supabaseErrorMessage(body []byte, status int) string {
	var top map[string]any
	if json.Unmarshal(body, &top) == nil {
		if m, ok := top["msg"].(string); ok && strings.TrimSpace(m) != "" {
			return m
		}
		if m, ok := top["message"].(string); ok && strings.TrimSpace(m) != "" {
			return m
		}
		if m, ok := top["error_description"].(string); ok && strings.TrimSpace(m) != "" {
			return m
		}
		if s, ok := top["error"].(string); ok && strings.TrimSpace(s) != "" {
			return s
		}
	}
	if len(body) > 0 && len(body) < 512 {
		return strings.TrimSpace(string(body))
	}
	return "Supabase request failed"
}

// validateSupabaseToken validates a JWT token by calling Supabase /auth/v1/user endpoint
func validateSupabaseToken(token string) (userID string, err error) {
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &SupabaseRequestError{
			Status:  resp.StatusCode,
			Message: supabaseErrorMessage(body, resp.StatusCode),
		}
	}

	var user map[string]any
	if err := json.Unmarshal(body, &user); err != nil {
		return "", err
	}

	id, ok := user["id"].(string)
	if !ok || id == "" {
		return "", errors.New("invalid user response from Supabase")
	}

	return id, nil
}

// requireAuth middleware validates Bearer token before allowing request to proceed
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(auth, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing bearer token"})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if token == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing bearer token"})
			return
		}

		userID, err := validateSupabaseToken(token)
		if err != nil {
			var sbErr *SupabaseRequestError
			if errors.As(err, &sbErr) {
				if sbErr.Status == 401 || sbErr.Status == 403 {
					writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid or expired session"})
					return
				}
			}
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "authentication failed"})
			return
		}

		// Add user ID to request headers for downstream services
		r.Header.Set("X-User-Id", userID)
		next(w, r)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id, X-Request-Id, X-CSRF-Token, X-Admin-Secret")
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
	if as := r.Header.Get("X-Admin-Secret"); as != "" {
		req.Header.Set("X-Admin-Secret", as)
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
		{path: "/api/v1/onboarding/kyc", base: userBase, public: false},
		{path: "/api/v1/onboarding/sacco", base: userBase, public: false},
		{path: "/api/v1/onboarding/kibiina", base: userBase, public: false},
		{path: "/api/v1/admin/kyc", base: userBase, public: false},
		{path: "/api/v1/admin/settings", base: userBase, public: false},
		{path: "/api/v1/geo", base: userBase, public: true},
		{path: "/api/v1/groups/policy", base: userBase, public: true},
		// object storage
		{path: "/api/v1/storage/upload", base: storageBase, public: false},
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
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != rt.path {
				http.NotFound(w, r)
				return
			}
			if rt.path == "/api/v1/admin/kyc" || rt.path == "/api/v1/admin/settings" {
				adminSecret := strings.TrimSpace(os.Getenv("ADMIN_API_KEY"))
				if adminSecret == "" || strings.TrimSpace(r.Header.Get("X-Admin-Secret")) != adminSecret {
					writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "admin authentication failed"})
					return
				}
			}
			forward(w, r, rt.base)
		}
		
		// Apply authentication middleware to protected routes
		if !rt.public {
			mux.HandleFunc(rt.path, requireAuth(handler))
		} else {
			mux.HandleFunc(rt.path, handler)
		}
	}

	log.Printf("API gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
