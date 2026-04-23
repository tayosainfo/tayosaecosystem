package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type notification struct {
	Channel   string    `json:"channel"`
	Recipient string    `json:"recipient"`
	Template  string    `json:"template"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

var outbox = struct {
	mu    sync.RWMutex
	items []notification
}{}

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
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

		// Add user ID to request headers for downstream use
		r.Header.Set("X-User-Id", userID)
		next(w, r)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "notification-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/notifications/send", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		var req notification
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		req.Status = "queued"
		req.CreatedAt = time.Now()
		outbox.mu.Lock()
		outbox.items = append(outbox.items, req)
		outbox.mu.Unlock()
		writeJSON(w, http.StatusAccepted, map[string]any{
			"message":      "Notification queued",
			"notification": req,
		})
	}))

	mux.HandleFunc("/api/v1/notifications/outbox", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		outbox.mu.RLock()
		defer outbox.mu.RUnlock()
		writeJSON(w, http.StatusOK, map[string]any{"items": outbox.items})
	}))

	log.Printf("Notification Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
