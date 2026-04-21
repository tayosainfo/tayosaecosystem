package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

func envURL(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		v = fallback
	}
	return strings.TrimRight(v, "/")
}

func postJSON(url string, body map[string]any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{Timeout: 8 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func emitAudit(actor, action, resource, detail string) {
	_ = postJSON(envURL("AUDIT_SERVICE_URL", "http://localhost:8014")+"/api/v1/audit/events", map[string]any{
		"actor":    actor,
		"action":   action,
		"service":  "user-service",
		"resource": resource,
		"detail":   detail,
	})
}

func emitNotification(channel, recipient, template string) {
	_ = postJSON(envURL("NOTIFICATION_SERVICE_URL", "http://localhost:8010")+"/api/v1/notifications/send", map[string]any{
		"channel":   channel,
		"recipient": recipient,
		"template":  template,
	})
}

func createAffiliateReferral(referralCode, referrerID, refereeID string) {
	if strings.TrimSpace(referralCode) == "" || strings.TrimSpace(referrerID) == "" || strings.TrimSpace(refereeID) == "" {
		return
	}
	_ = postJSON(envURL("AFFILIATE_SERVICE_URL", "http://localhost:8016")+"/api/v1/affiliate/referrals", map[string]any{
		"referralCode": referralCode,
		"referrerId":   referrerID,
		"refereeId":    refereeID,
	})
}

