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

type referralReward struct {
	ReferralCode string    `json:"referralCode"`
	ReferrerID   string    `json:"referrerId"`
	RefereeID    string    `json:"refereeId"`
	Status       string    `json:"status"`
	RewardPoints int       `json:"rewardPoints"`
	CreatedAt    time.Time `json:"createdAt"`
}

var rewards = struct {
	mu    sync.RWMutex
	items []referralReward
}{}

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8016"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "affiliate-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/affiliate/referrals", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		var req struct {
			ReferralCode string `json:"referralCode"`
			ReferrerID   string `json:"referrerId"`
			RefereeID    string `json:"refereeId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		if strings.TrimSpace(req.ReferralCode) == "" || req.ReferrerID == "" || req.RefereeID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "referralCode, referrerId, refereeId required"})
			return
		}

		record := referralReward{
			ReferralCode: req.ReferralCode,
			ReferrerID:   req.ReferrerID,
			RefereeID:    req.RefereeID,
			Status:       "pending",
			RewardPoints: 100,
			CreatedAt:    time.Now(),
		}
		rewards.mu.Lock()
		rewards.items = append(rewards.items, record)
		rewards.mu.Unlock()

		writeJSON(w, http.StatusAccepted, map[string]any{
			"reward":    record,
			"notify":    true,
			"auditEvent": "affiliate_referral_pending",
		})
	})

	mux.HandleFunc("/api/v1/affiliate/rewards", func(w http.ResponseWriter, r *http.Request) {
		rewards.mu.RLock()
		defer rewards.mu.RUnlock()
		writeJSON(w, http.StatusOK, map[string]any{
			"items": rewards.items,
		})
	})

	log.Printf("Affiliate Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
