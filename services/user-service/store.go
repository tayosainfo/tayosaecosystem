package main

import (
	"context"
	"log"
	"os"
	"strings"
)

// Store persists identity, onboarding, and answers geo queries.
type Store interface {
	FindByPhone(phoneE164 string) (User, bool)
	FindByEmailKey(normalizedEmail string) (User, bool)
	FindByUserID(userID string) (User, bool)
	CreateIdentityWithOnboarding(u User, ob OnboardingProfile, passwordHash *string) error
	UpdateIdentity(u User) error
	SetContactEmailVerified(phoneE164 string, verified bool) error
	UpsertOnboarding(p OnboardingProfile) error
	GetOnboarding(userID string) (OnboardingProfile, bool)
	GeoDistinct(level, parent string) ([]string, error)
	GeoRecordExists(district, county, subCounty, parish, village string) (bool, error)
	GroupPolicyStats() (GroupPolicyStats, error)
	UpsertUserConsents(c UserConsents) error
	GetUserConsents(userID string) (UserConsents, bool)
	SetUserReferralCode(userID, referralCode string) error
	GetUserReferralCode(userID string) (string, bool)
	FindUserIDByReferralCode(referralCode string) (string, bool)
	UpsertKYCProfile(k KYCProfile) error
	GetKYCProfile(userID string) (KYCProfile, bool)
	ReplaceKYCDocuments(userID string, docs []KYCDocument) error
	GetKYCDocuments(userID string) ([]KYCDocument, error)
	SetKYCDecision(userID, status, reviewedBy, reviewNote string) error
	UpsertSaccoMembership(m SaccoMembership) error
	GetSaccoMembership(userID string) (SaccoMembership, bool)
	EnsureSharesLedger(userID string, sharesUnits int) error
	GetSharesUnits(userID string) (int, bool, error)
	UpsertKibiinaPreference(k KibiinaPreference) error
	GetKibiinaPreference(userID string) (KibiinaPreference, bool)
	ListAdminKYCQueue(status string, limit int) ([]AdminKYCQueueItem, error)
	GetAdminSetting(key string) (map[string]any, bool, error)
	SetAdminSetting(key string, value map[string]any) error
	EnsureGeoSeeded() error
	Ping() error
	Close()
}

var activeStore Store

// databaseDSN returns Postgres URL for user-service. DATABASE_URL is preferred;
// CONNECTION_STRING is accepted for compatibility with various .env exports.
func databaseDSN() string {
	if v := strings.TrimSpace(os.Getenv("DATABASE_URL")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("CONNECTION_STRING"))
}

func initStore() {
	dsn := databaseDSN()
	if dsn != "" {
		ctx := context.Background()
		pg, err := NewPostgresStore(ctx, dsn)
		if err == nil {
			activeStore = pg
			if err := activeStore.EnsureGeoSeeded(); err != nil {
				log.Printf("geo seed warning: %v", err)
			}
			log.Print("user-service: using PostgreSQL store")
			return
		}
		if strings.TrimSpace(os.Getenv("USER_SERVICE_ALLOW_MEMORY_FALLBACK")) != "1" {
			log.Fatalf(`postgres store: %v

user-service cannot open DATABASE_URL. Common causes on Windows:
  • Managed Postgres hostnames often do NOT accept direct TCP from your PC — use Supabase HTTP APIs for auth, not raw Postgres from your laptop.
  • For a real local DB: docker compose up -d postgres (repo root), apply db/migrations/*.sql, then set DATABASE_URL to localhost (see .env.example).
  • user-service reads DATABASE_URL or CONNECTION_STRING (either name works).
  • For quick UI tests without Postgres: remove or comment DATABASE_URL in .env (in-memory store; data resets when the process exits).
  • To keep DATABASE_URL in .env but still run when the cloud DB is down: set USER_SERVICE_ALLOW_MEMORY_FALLBACK=1 (uses in-memory store; not for production).`, err)
		}
		log.Printf("user-service: postgres unavailable (%v); USER_SERVICE_ALLOW_MEMORY_FALLBACK=1 — using in-memory store", err)
	}
	ms := NewMemoryStore()
	activeStore = ms
	if err := activeStore.EnsureGeoSeeded(); err != nil {
		log.Printf("geo seed warning: %v", err)
	}
	log.Print("user-service: using in-memory store (set DATABASE_URL for Postgres)")
}
