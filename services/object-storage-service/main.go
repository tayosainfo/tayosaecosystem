package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

type uploadStrategyInput struct {
	FileName    string
	Category    string
	ContentType string
	Size        int64 // 0 = omit
}

func parseUploadStrategyInput(r *http.Request) (uploadStrategyInput, error) {
	var in uploadStrategyInput
	switch r.Method {
	case http.MethodGet:
		in.FileName = r.URL.Query().Get("fileName")
		in.Category = r.URL.Query().Get("category")
		in.ContentType = r.URL.Query().Get("contentType")
		if s := strings.TrimSpace(r.URL.Query().Get("size")); s != "" {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return in, err
			}
			in.Size = n
		}
	case http.MethodPost:
		bodyBytes, _ := io.ReadAll(r.Body)
		ct := r.Header.Get("Content-Type")
		if len(bodyBytes) > 0 && strings.Contains(ct, "application/json") {
			var body struct {
				FileName    string `json:"fileName"`
				Category    string `json:"category"`
				ContentType string `json:"contentType"`
				Size        *int64 `json:"size"`
			}
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				return in, err
			}
			in.FileName = body.FileName
			in.Category = body.Category
			in.ContentType = body.ContentType
			if body.Size != nil {
				in.Size = *body.Size
			}
		} else {
			in.FileName = r.URL.Query().Get("fileName")
			in.Category = r.URL.Query().Get("category")
			in.ContentType = r.URL.Query().Get("contentType")
			if s := strings.TrimSpace(r.URL.Query().Get("size")); s != "" {
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return in, err
				}
				in.Size = n
			}
		}
	default:
		return in, errors.New("unsupported method")
	}
	if strings.TrimSpace(in.FileName) == "" {
		in.FileName = "document.bin"
	}
	if strings.TrimSpace(in.Category) == "" {
		in.Category = "kyc"
	}
	return in, nil
}

func buildObjectFilename(in uploadStrategyInput) string {
	fn := strings.TrimSpace(in.FileName)
	cat := strings.TrimSpace(in.Category)
	return cat + "/" + time.Now().Format("20060102150405") + "-" + fn
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

func proxyUploadStrategy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}

	auth, ok := bearerForInsForge(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "missing Authorization bearer token and INSFORGE_ANON_KEY is not set",
		})
		return
	}

	in, err := parseUploadStrategyInput(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request: " + err.Error()})
		return
	}

	filename := buildObjectFilename(in)
	bodyObj := map[string]any{"filename": filename}
	if strings.TrimSpace(in.ContentType) != "" {
		bodyObj["contentType"] = strings.TrimSpace(in.ContentType)
	}
	if in.Size > 0 {
		bodyObj["size"] = in.Size
	}
	raw, err := json.Marshal(bodyObj)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	bucket := storageBucket()
	url := insforgeBaseURL() + "/api/storage/buckets/" + strings.TrimPrefix(bucket, "/") + "/upload-strategy"
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	resp, err := httpClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	ct := resp.Header.Get("Content-Type")
	if ct != "" {
		w.Header().Set("Content-Type", ct)
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(respBody)
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

	mux.HandleFunc("/api/v1/storage/upload-url", proxyUploadStrategy)

	mux.HandleFunc("/api/v1/storage/object", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"message": "Object metadata lookup placeholder active.",
		})
	})

	log.Printf("Object Storage Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
