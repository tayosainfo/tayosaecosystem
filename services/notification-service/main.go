package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

	mux.HandleFunc("/api/v1/notifications/send", func(w http.ResponseWriter, r *http.Request) {
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
	})

	mux.HandleFunc("/api/v1/notifications/outbox", func(w http.ResponseWriter, r *http.Request) {
		outbox.mu.RLock()
		defer outbox.mu.RUnlock()
		writeJSON(w, http.StatusOK, map[string]any{"items": outbox.items})
	})

	log.Printf("Notification Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
