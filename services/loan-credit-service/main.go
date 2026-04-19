package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8013"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "active",
			"service": "loan-credit-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/loans/score", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID          string  `json:"userId"`
			MonthlyIncome   float64 `json:"monthlyIncome"`
			ExistingDebt    float64 `json:"existingDebt"`
			TrustScore      int     `json:"trustScore"`
			OnTimePayments  int     `json:"onTimePayments"`
		}
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		if strings.TrimSpace(req.UserID) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "userId required"})
			return
		}
		score := 500 + req.TrustScore + (req.OnTimePayments * 2)
		if req.MonthlyIncome > 0 && req.ExistingDebt > 0 && req.ExistingDebt/req.MonthlyIncome > 0.5 {
			score -= 30
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"userId":          req.UserID,
			"score":           score,
			"riskBand":        "medium",
			"maxEligibleLoan": 1000000,
		})
	})

	log.Printf("Loan/Credit Service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
