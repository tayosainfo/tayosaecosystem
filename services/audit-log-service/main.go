package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type auditEvent struct {
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Service   string    `json:"service"`
	Resource  string    `json:"resource"`
	Timestamp time.Time `json:"timestamp"`
	Detail    string    `json:"detail"`
}

var eventStore = struct {
	mu    sync.RWMutex
	items []auditEvent
}{}

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8014"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "audit-log-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/audit/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var evt auditEvent
			if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
				return
			}
			evt.Timestamp = time.Now()
			eventStore.mu.Lock()
			eventStore.items = append(eventStore.items, evt)
			eventStore.mu.Unlock()
			writeJSON(w, http.StatusAccepted, map[string]any{"event": evt})
		case http.MethodGet:
			eventStore.mu.RLock()
			defer eventStore.mu.RUnlock()
			writeJSON(w, http.StatusOK, map[string]any{"items": eventStore.items})
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		}
	})

	log.Printf("Audit Log Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
