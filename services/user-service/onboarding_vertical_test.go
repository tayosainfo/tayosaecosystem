package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func seedTestUser(t *testing.T, userID string) {
	t.Helper()
	activeStore = NewMemoryStore()
	ms := activeStore.(*MemoryStore)
	ms.geoRows = []geoRow{{District: "Kampala", County: "Nakawa", SubCounty: "Nakawa", Parish: "Naguru", Village: "Kiwatule"}}
	err := activeStore.CreateIdentityWithOnboarding(User{
		ID:        userID,
		FullName:  "Test User",
		PhoneE164: "+256700000001",
		AuthEmail: "700000001@tayosa.local",
		CreatedAt: time.Now(),
	}, OnboardingProfile{UserID: userID, Phase: 1, TrustScoreSeed: 10}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func authReq(method, path string, body any, token string) *http.Request {
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func TestVerticalFlowKycApproveSaccoAndReadModel(t *testing.T) {
	seedTestUser(t, "u_123")

	kycPayload := map[string]any{
		"dateOfBirth":                "1990-06-14",
		"gender":                     "female",
		"nationality":                "UG",
		"occupationStatus":           "self_employed",
		"idType":                     "national_id",
		"idNumber":                   "CF92001003",
		"idDocumentFrontKey":         "collateral_docs/front.jpg",
		"idDocumentBackKey":          "collateral_docs/back.jpg",
		"selfieKey":                  "collateral_docs/selfie.jpg",
		"nokFullName":                "Jane Doe",
		"nokRelationship":            "sister",
		"nokPhone":                   "+256700000002",
		"sourceOfFunds":              "business",
		"pepStatus":                  false,
		"saccoMembershipDisclosures": "none",
	}
	rr := httptest.NewRecorder()
	requireAuth(onboardingKYCHandler).ServeHTTP(rr, authReq(http.MethodPost, "/api/v1/onboarding/kyc", kycPayload, "dev-token-u_123"))
	if rr.Code >= 300 {
		t.Fatalf("kyc submit failed: %d %s", rr.Code, rr.Body.String())
	}

	_ = os.Setenv("ADMIN_API_KEY", "secret")
	adminReq := authReq(http.MethodPatch, "/api/v1/admin/kyc?userId=u_123", map[string]any{
		"status":     "approved",
		"reviewedBy": "admin_1",
	}, "dev-token-u_123")
	adminReq.Header.Set("X-Admin-Secret", "secret")
	rr = httptest.NewRecorder()
	adminKYCDecisionHandler(rr, adminReq)
	if rr.Code >= 300 {
		t.Fatalf("admin kyc decision failed: %d %s", rr.Code, rr.Body.String())
	}

	saccoPayload := map[string]any{
		"district":            "Kampala",
		"county":              "Nakawa",
		"subCounty":           "Nakawa",
		"parish":              "Naguru",
		"village":             "Kiwatule",
		"mobileMoneyProvider": "MTN",
		"mobileMoneyNumber":   "+256700000001",
		"sharesToPurchase":    3,
	}
	rr = httptest.NewRecorder()
	requireAuth(onboardingSaccoHandler).ServeHTTP(rr, authReq(http.MethodPost, "/api/v1/onboarding/sacco", saccoPayload, "dev-token-u_123"))
	if rr.Code >= 300 {
		t.Fatalf("sacco submit failed: %d %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	requireAuth(meHandler).ServeHTTP(rr, authReq(http.MethodGet, "/api/v1/users/me", nil, "dev-token-u_123"))
	if rr.Code >= 300 {
		t.Fatalf("me failed: %d %s", rr.Code, rr.Body.String())
	}
	if !bytes.Contains(rr.Body.Bytes(), []byte(`"canTransact":true`)) {
		t.Fatalf("expected canTransact true in read model: %s", rr.Body.String())
	}
}

