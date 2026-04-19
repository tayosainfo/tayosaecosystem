package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type feeRecord struct {
	UserID      string    `json:"userId"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	Reference   string    `json:"reference"`
}

var feeStore = struct {
	mu    sync.RWMutex
	items []feeRecord
}{}

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "fee-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/fees/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var rec feeRecord
			if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
				return
			}
			rec.CreatedAt = time.Now()
			if rec.Status == "" {
				rec.Status = "completed"
			}
			feeStore.mu.Lock()
			feeStore.items = append(feeStore.items, rec)
			feeStore.mu.Unlock()
			writeJSON(w, http.StatusCreated, map[string]any{"transaction": rec})
		case http.MethodGet:
			feeStore.mu.RLock()
			defer feeStore.mu.RUnlock()
			writeJSON(w, http.StatusOK, map[string]any{"items": feeStore.items})
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		}
	})

	log.Printf("Fee Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
