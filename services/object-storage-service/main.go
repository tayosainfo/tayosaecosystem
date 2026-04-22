package main

import (
	"encoding/json"
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

// insforgeServiceBearer returns the InsForge anon key as the service credential
// for InsForge storage calls.
//
// IMPORTANT: We NEVER forward the caller's TAYOSA token (e.g. "dev-token-<id>")
// to InsForge. Those are internal TAYOSA tokens — InsForge will reject them.
// The caller has already been authenticated by the API gateway.
func insforgeServiceBearer() (string, bool) {
	key := strings.TrimSpace(os.Getenv("INSFORGE_ANON_KEY"))
	if key != "" {
		return "Bearer " + key, true
	}
	return "", false
}

// hasCallerBearer confirms the incoming request carries any Bearer token,
// so we can gate the endpoint without forwarding it to InsForge.
func hasCallerBearer(r *http.Request) bool {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}
	return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")) != ""
}

var httpClient = &http.Client{Timeout: 60 * time.Second}

// ensureBucket creates the InsForge storage bucket if it does not yet exist.
func ensureBucket(bucket, serviceBearer string) {
	url := insforgeBaseURL() + "/storage/v1/bucket"

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
// "category" field (defaults to "kyc"). It uploads the raw bytes to InsForge
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

	// Gate: caller must present a bearer (validated by gateway), but we use
	// our own service credential when talking to InsForge.
	if !hasCallerBearer(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing bearer token"})
		return
	}

	serviceBearer, ok := insforgeServiceBearer()
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "INSFORGE_ANON_KEY is not configured on this server",
		})
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
	objectPath := fmt.Sprintf("%s/%s-%s",
		category,
		time.Now().UTC().Format("20060102150405"),
		safeFilename,
	)

	bucket := storageBucket()
	uploadURL := insforgeBaseURL() + "/storage/v1/object/" + bucket + "/" + objectPath

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, uploadURL, file)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	req.Header.Set("Authorization", serviceBearer)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")
	req.ContentLength = header.Size

	resp, err := httpClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "InsForge storage unreachable: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		log.Printf("InsForge storage upload failed: status=%d body=%s url=%s", resp.StatusCode, string(respBody), uploadURL)
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"error":   "InsForge storage returned an error",
			"status":  resp.StatusCode,
			"details": string(respBody),
		})
		return
	}

	log.Printf("InsForge storage upload OK: %s", objectPath)
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
	if bearer, ok := insforgeServiceBearer(); ok {
		go ensureBucket(storageBucket(), bearer)
	} else {
		log.Print("WARNING: INSFORGE_ANON_KEY not set — KYC file uploads will fail")
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
