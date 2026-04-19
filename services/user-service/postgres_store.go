package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, dsn string) (*PostgresStore, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) Ping() error {
	return s.pool.Ping(context.Background())
}

func (s *PostgresStore) FindByPhone(phoneE164 string) (User, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
SELECT user_id, full_name, phone_e164, auth_email, contact_email, insforge_user_id, COALESCE(insforge_login_email,''),
password_hash, phone_verified_at, contact_email_verified_at, created_at, date_of_birth, nationality
FROM users_identity WHERE phone_e164 = $1`, phoneE164)
	u, err := scanUser(row)
	if err != nil {
		return User{}, false
	}
	return u, true
}

func (s *PostgresStore) FindByEmailKey(normalizedEmail string) (User, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
SELECT user_id, full_name, phone_e164, auth_email, contact_email, insforge_user_id, COALESCE(insforge_login_email,''),
password_hash, phone_verified_at, contact_email_verified_at, created_at, date_of_birth, nationality
FROM users_identity
WHERE LOWER(TRIM(contact_email)) = LOWER(TRIM($1))
   OR LOWER(TRIM(auth_email)) = LOWER(TRIM($1))
   OR LOWER(TRIM(COALESCE(insforge_login_email,''))) = LOWER(TRIM($1))
LIMIT 1`, normalizedEmail)
	u, err := scanUser(row)
	if err != nil {
		return User{}, false
	}
	return u, true
}

func (s *PostgresStore) FindByUserID(userID string) (User, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
SELECT user_id, full_name, phone_e164, auth_email, contact_email, insforge_user_id, COALESCE(insforge_login_email,''),
password_hash, phone_verified_at, contact_email_verified_at, created_at, date_of_birth, nationality
FROM users_identity WHERE user_id = $1`, userID)
	u, err := scanUser(row)
	if err != nil {
		return User{}, false
	}
	return u, true
}

func scanUser(row pgx.Row) (User, error) {
	var u User
	var contact, insforgeID *string
	var insforgeLogin string
	var pwdHash *string
	var phoneV, emailV *time.Time
	var created time.Time
	var dob *time.Time
	var nat *string
	err := row.Scan(&u.ID, &u.FullName, &u.PhoneE164, &u.AuthEmail, &contact, &insforgeID, &insforgeLogin,
		&pwdHash, &phoneV, &emailV, &created, &dob, &nat)
	if err != nil {
		return User{}, err
	}
	if contact != nil {
		u.ContactEmail = *contact
	}
	if insforgeID != nil {
		u.InsforgeUserID = *insforgeID
	}
	u.InsforgeEmail = insforgeLogin
	if pwdHash != nil {
		u.PasswordHash = *pwdHash
	}
	if phoneV != nil {
		u.PhoneVerifiedAt = *phoneV
	}
	u.ContactEmailChecked = emailV != nil
	u.CreatedAt = created
	if dob != nil {
		t := *dob
		u.DateOfBirth = &t
	}
	if nat != nil {
		u.Nationality = *nat
	}
	return u, nil
}

func (s *PostgresStore) CreateIdentityWithOnboarding(u User, ob OnboardingProfile, passwordHash *string) error {
	ctx := context.Background()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Defer FK checks to commit so BEFORE INSERT triggers that write shares_ledger can succeed.
	_, _ = tx.Exec(ctx, `SET CONSTRAINTS ALL DEFERRED`)

	var contact any
	if u.ContactEmail != "" {
		contact = u.ContactEmail
	}
	var insforgeID any
	if u.InsforgeUserID != "" {
		insforgeID = u.InsforgeUserID
	}
	var ifLogin any
	if u.InsforgeEmail != "" {
		ifLogin = u.InsforgeEmail
	}
	var ph any
	if passwordHash != nil {
		ph = *passwordHash
	}
	var dob any
	if u.DateOfBirth != nil {
		dob = *u.DateOfBirth
	}
	var nat any
	if strings.TrimSpace(u.Nationality) != "" {
		nat = strings.TrimSpace(u.Nationality)
	}

	_, err = tx.Exec(ctx, `
INSERT INTO users_identity (user_id, full_name, phone_e164, auth_email, contact_email, insforge_user_id, insforge_login_email, password_hash, date_of_birth, nationality, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, now())`,
		u.ID, u.FullName, u.PhoneE164, u.AuthEmail, contact, insforgeID, ifLogin, ph, dob, nat)
	if err != nil {
		return err
	}

	payload := onboardingPayloadJSON(ob)
	_, err = tx.Exec(ctx, `
INSERT INTO onboarding_profiles (user_id, phase, referral_code, district, county, sub_county, parish, village, trust_score_seed, phase_payload, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::jsonb, now())`,
		ob.UserID, ob.Phase, nullStr(ob.ReferralCode),
		nullStr(ob.Geo["district"]), nullStr(ob.Geo["county"]), nullStr(ob.Geo["sub_county"]),
		nullStr(ob.Geo["parish"]), nullStr(ob.Geo["village"]), ob.TrustScoreSeed, payload)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func nullStr(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func onboardingPayloadJSON(p OnboardingProfile) string {
	m := map[string]any{
		"kyc":          p.KYC,
		"membership":   p.Membership,
		"kibiina":      p.Kibiina,
		"referralCode": p.ReferralCode,
		"geo":          p.Geo,
	}
	if p.AdditionalContext != nil {
		m["additionalContext"] = p.AdditionalContext
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (s *PostgresStore) UpdateIdentity(u User) error {
	ctx := context.Background()
	var contact any
	if strings.TrimSpace(u.ContactEmail) != "" {
		contact = u.ContactEmail
	}
	var iid any
	if strings.TrimSpace(u.InsforgeUserID) != "" {
		iid = u.InsforgeUserID
	}
	var ifLogin any
	if strings.TrimSpace(u.InsforgeEmail) != "" {
		ifLogin = u.InsforgeEmail
	}
	var ph any
	if strings.TrimSpace(u.PasswordHash) != "" {
		ph = u.PasswordHash
	}
	_, err := s.pool.Exec(ctx, `
UPDATE users_identity SET
  full_name = $2,
  contact_email = COALESCE($3, contact_email),
  insforge_user_id = COALESCE($4, insforge_user_id),
  insforge_login_email = COALESCE($5, insforge_login_email),
  password_hash = COALESCE($6, password_hash),
  updated_at = now()
WHERE user_id = $1`,
		u.ID, u.FullName, contact, iid, ifLogin, ph)
	return err
}

func (s *PostgresStore) SetContactEmailVerified(phoneE164 string, verified bool) error {
	ctx := context.Background()
	var t any
	if verified {
		t = time.Now()
	}
	_, err := s.pool.Exec(ctx, `
UPDATE users_identity SET contact_email_verified_at = $2, updated_at = now() WHERE phone_e164 = $1`,
		phoneE164, t)
	return err
}

func (s *PostgresStore) UpsertOnboarding(p OnboardingProfile) error {
	ctx := context.Background()
	payload := onboardingPayloadJSON(p)
	_, err := s.pool.Exec(ctx, `
INSERT INTO onboarding_profiles (user_id, phase, referral_code, district, county, sub_county, parish, village, trust_score_seed, phase_payload, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::jsonb, now())
ON CONFLICT (user_id) DO UPDATE SET
  phase = EXCLUDED.phase,
  referral_code = EXCLUDED.referral_code,
  district = EXCLUDED.district,
  county = EXCLUDED.county,
  sub_county = EXCLUDED.sub_county,
  parish = EXCLUDED.parish,
  village = EXCLUDED.village,
  trust_score_seed = EXCLUDED.trust_score_seed,
  phase_payload = EXCLUDED.phase_payload,
  updated_at = now()`,
		p.UserID, p.Phase, nullStr(p.ReferralCode),
		nullStr(p.Geo["district"]), nullStr(p.Geo["county"]), nullStr(p.Geo["sub_county"]),
		nullStr(p.Geo["parish"]), nullStr(p.Geo["village"]), p.TrustScoreSeed, payload)
	return err
}

func (s *PostgresStore) GetOnboarding(userID string) (OnboardingProfile, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
SELECT user_id, phase, referral_code, district, county, sub_county, parish, village, trust_score_seed, COALESCE(phase_payload, '{}'::jsonb), updated_at
FROM onboarding_profiles WHERE user_id = $1`, userID)
	var p OnboardingProfile
	var ref *string
	var dist, county, subc, par, vill *string
	var payload []byte
	var updated time.Time
	if err := row.Scan(&p.UserID, &p.Phase, &ref, &dist, &county, &subc, &par, &vill, &p.TrustScoreSeed, &payload, &updated); err != nil {
		return OnboardingProfile{}, false
	}
	if ref != nil {
		p.ReferralCode = *ref
	}
	p.Geo = map[string]string{}
	if dist != nil {
		p.Geo["district"] = *dist
	}
	if county != nil {
		p.Geo["county"] = *county
	}
	if subc != nil {
		p.Geo["sub_county"] = *subc
	}
	if par != nil {
		p.Geo["parish"] = *par
	}
	if vill != nil {
		p.Geo["village"] = *vill
	}
	p.LastUpdatedAt = updated
	var aux map[string]json.RawMessage
	if json.Unmarshal(payload, &aux) == nil {
		unpackPhasePayload(&p, aux)
	}
	return p, true
}

func unpackPhasePayload(p *OnboardingProfile, aux map[string]json.RawMessage) {
	if raw, ok := aux["kyc"]; ok {
		_ = json.Unmarshal(raw, &p.KYC)
	}
	if raw, ok := aux["membership"]; ok {
		_ = json.Unmarshal(raw, &p.Membership)
	}
	if raw, ok := aux["kibiina"]; ok {
		_ = json.Unmarshal(raw, &p.Kibiina)
	}
	if raw, ok := aux["referralCode"]; ok {
		_ = json.Unmarshal(raw, &p.ReferralCode)
	}
	if raw, ok := aux["geo"]; ok {
		_ = json.Unmarshal(raw, &p.Geo)
	}
	if raw, ok := aux["additionalContext"]; ok {
		_ = json.Unmarshal(raw, &p.AdditionalContext)
	}
}

func (s *PostgresStore) GeoDistinct(level, parent string) ([]string, error) {
	ctx := context.Background()
	level = strings.ToLower(strings.TrimSpace(level))
	parent = strings.ToLower(strings.TrimSpace(parent))
	var q string
	var args []any
	switch level {
	case "district":
		q = `SELECT DISTINCT district FROM uganda_geo_units ORDER BY district`
	case "county":
		q = `SELECT DISTINCT county FROM uganda_geo_units WHERE LOWER(TRIM(district)) = $1 ORDER BY county`
		args = append(args, parent)
	case "subcounty":
		q = `SELECT DISTINCT sub_county FROM uganda_geo_units WHERE LOWER(TRIM(county)) = $1 ORDER BY sub_county`
		args = append(args, parent)
	case "parish":
		q = `SELECT DISTINCT parish FROM uganda_geo_units WHERE LOWER(TRIM(sub_county)) = $1 ORDER BY parish`
		args = append(args, parent)
	case "village":
		q = `SELECT DISTINCT village FROM uganda_geo_units WHERE LOWER(TRIM(parish)) = $1 ORDER BY village`
		args = append(args, parent)
	default:
		return nil, errors.New("invalid level")
	}
	rows, err := s.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (s *PostgresStore) EnsureGeoSeeded() error {
	ctx := context.Background()
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM uganda_geo_units`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	csvPath := filepath.Join("..", "..", "uganda_geo_data_2025-11-26.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil || len(rows) <= 1 {
		return err
	}
	batch := &pgx.Batch{}
	nq := 0
	for _, row := range rows[1:] {
		if len(row) < 5 {
			continue
		}
		batch.Queue(`INSERT INTO uganda_geo_units (district, county, sub_county, parish, village) VALUES ($1,$2,$3,$4,$5)`,
			strings.TrimSpace(row[0]), strings.TrimSpace(row[1]), strings.TrimSpace(row[2]), strings.TrimSpace(row[3]), strings.TrimSpace(row[4]))
		nq++
	}
	br := s.pool.SendBatch(ctx, batch)
	defer br.Close()
	for i := 0; i < nq; i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return br.Close()
}

func (s *PostgresStore) GeoRecordExists(district, county, subCounty, parish, village string) (bool, error) {
	ctx := context.Background()
	var n int
	err := s.pool.QueryRow(ctx, `
SELECT COUNT(*)::int FROM uganda_geo_units WHERE
  LOWER(TRIM(district)) = LOWER(TRIM($1)) AND
  LOWER(TRIM(county)) = LOWER(TRIM($2)) AND
  LOWER(TRIM(sub_county)) = LOWER(TRIM($3)) AND
  LOWER(TRIM(parish)) = LOWER(TRIM($4)) AND
  LOWER(TRIM(village)) = LOWER(TRIM($5))`,
		district, county, subCounty, parish, village).Scan(&n)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (s *PostgresStore) GroupPolicyStats() (GroupPolicyStats, error) {
	ctx := context.Background()
	var st GroupPolicyStats
	err := s.pool.QueryRow(ctx, `
SELECT
  COALESCE((SELECT COUNT(*)::int FROM parish_saccos), 0),
  COALESCE((SELECT COUNT(*)::int FROM village_kibiina_groups), 0)`).Scan(&st.ParishSaccos, &st.VillageKibiinaGroups)
	return st, err
}
