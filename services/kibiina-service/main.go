package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type kibiinaGroup struct {
	ID                 string    `json:"id"`
	Parish             string    `json:"parish"`
	Village            string    `json:"village"`
	CycleFrequency     string    `json:"cycleFrequency"`
	ContributionAmount float64   `json:"contributionAmount"`
	PayoutMethod       string    `json:"payoutMethod"`
	CreatedAt          time.Time `json:"createdAt"`
}

var kibiinaStore = struct {
	mu     sync.RWMutex
	groups []kibiinaGroup
}{}

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "kibiina-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/kibiina/groups", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req kibiinaGroup
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
				return
			}
			if strings.TrimSpace(req.Parish) == "" || strings.TrimSpace(req.Village) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": "parish and village are required"})
				return
			}
			req.ID = "kib_" + strings.ToLower(strings.ReplaceAll(req.Parish+"_"+req.Village, " ", "_"))
			req.CreatedAt = time.Now()
			kibiinaStore.mu.Lock()
			kibiinaStore.groups = append(kibiinaStore.groups, req)
			kibiinaStore.mu.Unlock()
			writeJSON(w, http.StatusCreated, map[string]any{"group": req})
		case http.MethodGet:
			kibiinaStore.mu.RLock()
			defer kibiinaStore.mu.RUnlock()
			writeJSON(w, http.StatusOK, map[string]any{"items": kibiinaStore.groups})
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		}
	})

	log.Printf("Kibiina Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
