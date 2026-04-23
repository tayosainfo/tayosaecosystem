package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
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

func trimBase(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), "/")
}

func supabaseBaseURL() string {
	b := strings.TrimSpace(os.Getenv("SUPABASE_URL"))
	if b == "" {
		b = "https://ablvrbnbsdqshrorhmjf.supabase.co"
	}
	return trimBase(b)
}

func storageBucket() string {
	b := strings.TrimSpace(os.Getenv("SUPABASE_STORAGE_BUCKET"))
	if b == "" {
		return "collateral_docs"
	}
	return b
}

// supabaseServiceBearer returns the Supabase anon key as the service credential
// for Supabase storage calls.
//
// IMPORTANT: We NEVER forward the caller's TAYOSA token (e.g. "dev-token-<id>")
// to Supabase. Those are internal TAYOSA tokens — Supabase will reject them.
// The caller has already been authenticated by the API gateway.
func supabaseServiceBearer() (string, bool) {
	key := strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY"))
	if key != "" {
		return "Bearer " + key, true
	}
	return "", false
}

// supabaseAdminBearer returns the Supabase service role key for admin operations
// like bucket creation. This should only be used for administrative tasks.
func supabaseAdminBearer() (string, bool) {
	key := strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	if key != "" {
		return "Bearer " + key, true
	}
	return "", false
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
	url := supabaseBaseURL() + "/auth/v1/user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	if key := strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY")); key != "" {
		req.Header.Set("apikey", key)
	}

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

// hasCallerBearer confirms the incoming request carries any Bearer token,
// so we can gate the endpoint without forwarding it to Supabase.
func hasCallerBearer(r *http.Request) bool {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}
	return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")) != ""
}

var httpClient = &http.Client{Timeout: 60 * time.Second}

// ensureBucket creates the Supabase storage bucket if it does not yet exist.
func ensureBucket(bucket, serviceBearer string) {
	url := supabaseBaseURL() + "/storage/v1/bucket"

	// Check if it exists first.
	checkBody := strings.NewReader("")
	req, err := http.NewRequest(http.MethodGet, url+"/"+bucket, checkBody)
	if err != nil {
		log.Printf("ensureBucket check error: %v", err)
		return
	}
	req.Header.Set("Authorization", serviceBearer)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("ensureBucket check request error: %v", err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		log.Printf("ensureBucket: bucket %q already exists", bucket)
		return
	}

	// Try to create it.
	createPayload := fmt.Sprintf(`{"id":%q,"name":%q,"public":false}`, bucket, bucket)
	req2, err := http.NewRequest(http.MethodPost, url, strings.NewReader(createPayload))
	if err != nil {
		log.Printf("ensureBucket create error: %v", err)
		return
	}
	req2.Header.Set("Authorization", serviceBearer)
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := httpClient.Do(req2)
	if err != nil {
		log.Printf("ensureBucket create request error: %v", err)
		return
	}
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	log.Printf("ensureBucket create %q: status=%d body=%s", bucket, resp2.StatusCode, string(body2))
}

// handleUpload accepts multipart/form-data with a "file" field and an optional
// "category" field (defaults to "kyc"). It uploads the raw bytes to Supabase
// storage using the service anon key and returns {"key": "<path>"} on success.
func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}

	// 32 MB in-memory, remainder spills to temp files.
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid multipart form: " + err.Error()})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing 'file' field: " + err.Error()})
		return
	}
	defer file.Close()

	category := strings.TrimSpace(r.FormValue("category"))
	if category == "" {
		category = "kyc"
	}

	// Determine MIME type.
	contentType := header.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if mt := mime.TypeByExtension(ext); mt != "" {
			contentType = mt
		} else {
			contentType = "application/octet-stream"
		}
	}

	// Sanitise filename and build storage path.
	safeFilename := strings.ReplaceAll(strings.TrimSpace(header.Filename), " ", "_")
	if safeFilename == "" {
		safeFilename = "upload.bin"
	}

	// Get user ID from validated token (set by requireAuth middleware)
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "user ID not found in request"})
		return
	}

	// Get the user's JWT token from the Authorization header
	// We need to forward this to Supabase so RLS policies can check auth.uid()
	userToken := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(userToken, "Bearer ") {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing bearer token"})
		return
	}

	// Build folder-based path: {user_id}/{category}/{timestamp}-{filename}
	// This enables folder-based RLS policies for better security
	objectPath := fmt.Sprintf("%s/%s/%s-%s",
		userID,    // User's folder (Supabase user ID)
		category,  // kyc, documents, etc.
		time.Now().UTC().Format("20060102150405"),
		safeFilename,
	)

	bucket := storageBucket()
	uploadURL := supabaseBaseURL() + "/storage/v1/object/" + bucket + "/" + objectPath

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, uploadURL, file)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	// CRITICAL: Use the user's JWT token, not the anon key
	// This allows Supabase RLS policies to check auth.uid()
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", contentType)
	// Supabase also requires the apikey header
	if anonKey := strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY")); anonKey != "" {
		req.Header.Set("apikey", anonKey)
	}
	req.Header.Set("x-upsert", "true")
	req.ContentLength = header.Size

	resp, err := httpClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "Supabase storage unreachable: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		log.Printf("Supabase storage upload failed: status=%d body=%s url=%s", resp.StatusCode, string(respBody), uploadURL)
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"error":   "Supabase storage returned an error",
			"status":  resp.StatusCode,
			"details": string(respBody),
		})
		return
	}

	log.Printf("Supabase storage upload OK: %s", objectPath)
	writeJSON(w, http.StatusOK, map[string]any{"key": objectPath})
}

func main() {
	_ = godotenv.Load(filepath.Join("..", "..", ".env"))
	_ = godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8015"
	}

	// Eagerly verify bucket exists at startup.
	if bearer, ok := supabaseAdminBearer(); ok {
		go ensureBucket(storageBucket(), bearer)
	} else {
		log.Print("WARNING: SUPABASE_SERVICE_ROLE_KEY not set — bucket creation will fail")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "object-storage-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/storage/upload", requireAuth(handleUpload))

	log.Printf("Object Storage Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
