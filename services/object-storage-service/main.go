package main

import (
	"encoding/json"
	"fmt"
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

func trimBase(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), "/")
}

func insforgeBaseURL() string {
	b := strings.TrimSpace(os.Getenv("INSFORGE_BASE_URL"))
	if b == "" {
		b = "https://74qj9u5z.us-east.insforge.app"
	}
	return trimBase(b)
}

func storageBucket() string {
	b := strings.TrimSpace(os.Getenv("INSFORGE_STORAGE_BUCKET"))
	if b == "" {
		return "collateral_docs"
	}
	return b
}

// bearerForInsForge prefers the incoming Authorization header (user JWT through the gateway),
// otherwise falls back to INSFORGE_ANON_KEY for direct calls to this service.
func bearerForInsForge(r *http.Request) (string, bool) {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(auth, "Bearer ") {
		tok := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if tok != "" {
			return auth, true
		}
	}
	key := strings.TrimSpace(os.Getenv("INSFORGE_ANON_KEY"))
	if key != "" {
		return "Bearer " + key, true
	}
	return "", false
}

var httpClient = &http.Client{Timeout: 60 * time.Second}

// handleUpload accepts multipart/form-data with a "file" field and an optional
// "category" field (defaults to "kyc"), then proxies the raw bytes directly to
// InsForge's storage API at POST /storage/v1/object/{bucket}/{path}.
// On success it returns {"key": "<category>/<timestamp>-<filename>"}.
func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}

	auth, ok := bearerForInsForge(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing bearer token"})
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

	objectPath := fmt.Sprintf("%s/%s-%s",
		category,
		time.Now().UTC().Format("20060102150405"),
		header.Filename,
	)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	bucket := storageBucket()
	insforgeURL := insforgeBaseURL() + "/storage/v1/object/" + bucket + "/" + objectPath

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, insforgeURL, file)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = header.Size

	resp, err := httpClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "InsForge storage unreachable: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		writeJSON(w, resp.StatusCode, map[string]any{
			"error":   "InsForge storage returned an error",
			"status":  resp.StatusCode,
			"details": string(respBody),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"key": objectPath})
}

func main() {
	_ = godotenv.Load(filepath.Join("..", "..", ".env"))
	_ = godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8015"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "object-storage-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/storage/upload", handleUpload)

	log.Printf("Object Storage Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
