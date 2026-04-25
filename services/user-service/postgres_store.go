package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
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
	
	// Disable prepared statements for compatibility with PgBouncer/Supabase pooler
	// This prevents "prepared statement already exists" errors
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	
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
SELECT user_id, full_name, phone_e164, auth_email, contact_email, supabase_user_id, COALESCE(supabase_login_email,''),
password_hash, phone_verified_at, contact_email_verified_at, created_at, date_of_birth, nationality,
COALESCE(role, 'user'), COALESCE(status, 'active'), role_assigned_at, role_assigned_by, last_login
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
SELECT user_id, full_name, phone_e164, auth_email, contact_email, supabase_user_id, COALESCE(supabase_login_email,''),
password_hash, phone_verified_at, contact_email_verified_at, created_at, date_of_birth, nationality,
COALESCE(role, 'user'), COALESCE(status, 'active'), role_assigned_at, role_assigned_by, last_login
FROM users_identity
WHERE LOWER(TRIM(contact_email)) = LOWER(TRIM($1))
   OR LOWER(TRIM(auth_email)) = LOWER(TRIM($1))
   OR LOWER(TRIM(COALESCE(supabase_login_email,''))) = LOWER(TRIM($1))
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
SELECT user_id, full_name, phone_e164, auth_email, contact_email, supabase_user_id, COALESCE(supabase_login_email,''),
password_hash, phone_verified_at, contact_email_verified_at, created_at, date_of_birth, nationality,
COALESCE(role, 'user'), COALESCE(status, 'active'), role_assigned_at, role_assigned_by, last_login
FROM users_identity WHERE user_id = $1`, userID)
	u, err := scanUser(row)
	if err != nil {
		return User{}, false
	}
	return u, true
}

func scanUser(row pgx.Row) (User, error) {
	var u User
	var contact, supabaseID *string
	var supabaseLogin string
	var pwdHash *string
	var phoneV, emailV *time.Time
	var created time.Time
	var dob *time.Time
	var nat *string
	var role, status string
	var roleAssignedAt, lastLogin *time.Time
	var roleAssignedBy *string
	
	err := row.Scan(&u.ID, &u.FullName, &u.PhoneE164, &u.AuthEmail, &contact, &supabaseID, &supabaseLogin,
		&pwdHash, &phoneV, &emailV, &created, &dob, &nat, &role, &status, &roleAssignedAt, &roleAssignedBy, &lastLogin)
	if err != nil {
		return User{}, err
	}
	if contact != nil {
		u.ContactEmail = *contact
	}
	if supabaseID != nil {
		u.SupabaseUserID = *supabaseID
	}
	u.SupabaseLoginEmail = supabaseLogin
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
	u.Role = role
	u.Status = status
	u.RoleAssignedAt = roleAssignedAt
	if roleAssignedBy != nil {
		u.RoleAssignedBy = *roleAssignedBy
	}
	u.LastLogin = lastLogin
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
	var supabaseID any
	if u.SupabaseUserID != "" {
		supabaseID = u.SupabaseUserID
	}
	var supabaseLogin any
	if u.SupabaseLoginEmail != "" {
		supabaseLogin = u.SupabaseLoginEmail
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
INSERT INTO users_identity (user_id, full_name, phone_e164, auth_email, contact_email, supabase_user_id, supabase_login_email, password_hash, date_of_birth, nationality, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, now())`,
		u.ID, u.FullName, u.PhoneE164, u.AuthEmail, contact, supabaseID, supabaseLogin, ph, dob, nat)
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
	if strings.TrimSpace(u.SupabaseUserID) != "" {
		iid = u.SupabaseUserID
	}
	var ifLogin any
	if strings.TrimSpace(u.SupabaseLoginEmail) != "" {
		ifLogin = u.SupabaseLoginEmail
	}
	var ph any
	if strings.TrimSpace(u.PasswordHash) != "" {
		ph = u.PasswordHash
	}
	_, err := s.pool.Exec(ctx, `
UPDATE users_identity SET
  full_name = $2,
  contact_email = COALESCE($3, contact_email),
  supabase_user_id = COALESCE($4, supabase_user_id),
  supabase_login_email = COALESCE($5, supabase_login_email),
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

func (s *PostgresStore) UpsertUserConsents(c UserConsents) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `
INSERT INTO user_consents (user_id, terms_accepted_at, privacy_accepted_at, terms_version, privacy_version, updated_at)
VALUES ($1,$2,$3,$4,$5, now())
ON CONFLICT (user_id) DO UPDATE SET
  terms_accepted_at = EXCLUDED.terms_accepted_at,
  privacy_accepted_at = EXCLUDED.privacy_accepted_at,
  terms_version = EXCLUDED.terms_version,
  privacy_version = EXCLUDED.privacy_version,
  updated_at = now()`,
		c.UserID, c.TermsAcceptedAt, c.PrivacyAcceptedAt, nullStr(c.TermsVersion), nullStr(c.PrivacyVersion))
	return err
}

func (s *PostgresStore) GetUserConsents(userID string) (UserConsents, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
SELECT user_id, terms_accepted_at, privacy_accepted_at, COALESCE(terms_version,''), COALESCE(privacy_version,''), updated_at
FROM user_consents WHERE user_id = $1`, userID)
	var c UserConsents
	var tv, pv string
	if err := row.Scan(&c.UserID, &c.TermsAcceptedAt, &c.PrivacyAcceptedAt, &tv, &pv, &c.LastUpdatedAt); err != nil {
		return UserConsents{}, false
	}
	c.TermsVersion = tv
	c.PrivacyVersion = pv
	return c, true
}

func (s *PostgresStore) UpsertKYCProfile(k KYCProfile) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `
INSERT INTO kyc_profiles (
  user_id, status, date_of_birth, gender, nationality, occupation_status, id_type, id_number,
  nok_full_name, nok_relationship, nok_phone, nok_email, source_of_funds, pep_status,
  sacco_membership_disclosures, submitted_at, reviewed_at, review_note, reviewed_by, updated_at
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19, now())
ON CONFLICT (user_id) DO UPDATE SET
  status = EXCLUDED.status,
  date_of_birth = EXCLUDED.date_of_birth,
  gender = EXCLUDED.gender,
  nationality = EXCLUDED.nationality,
  occupation_status = EXCLUDED.occupation_status,
  id_type = EXCLUDED.id_type,
  id_number = EXCLUDED.id_number,
  nok_full_name = EXCLUDED.nok_full_name,
  nok_relationship = EXCLUDED.nok_relationship,
  nok_phone = EXCLUDED.nok_phone,
  nok_email = EXCLUDED.nok_email,
  source_of_funds = EXCLUDED.source_of_funds,
  pep_status = EXCLUDED.pep_status,
  sacco_membership_disclosures = EXCLUDED.sacco_membership_disclosures,
  submitted_at = EXCLUDED.submitted_at,
  reviewed_at = EXCLUDED.reviewed_at,
  review_note = EXCLUDED.review_note,
  reviewed_by = EXCLUDED.reviewed_by,
  updated_at = now()`,
		k.UserID, k.Status, k.DateOfBirth, nullStr(k.Gender), nullStr(k.Nationality), nullStr(k.OccupationStatus),
		nullStr(k.IDType), nullStr(k.IDNumber), nullStr(k.NOKFullName), nullStr(k.NOKRelationship), nullStr(k.NOKPhone),
		nullStr(k.NOKEmail), nullStr(k.SourceOfFunds), k.PEPStatus, nullStr(k.SACCOMembershipDisclosures),
		k.SubmittedAt, k.ReviewedAt, nullStr(k.ReviewNote), nullStr(k.ReviewedBy))
	return err
}

func (s *PostgresStore) GetKYCProfile(userID string) (KYCProfile, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
SELECT user_id, status, date_of_birth, COALESCE(gender,''), COALESCE(nationality,''), COALESCE(occupation_status,''),
COALESCE(id_type,''), COALESCE(id_number,''), COALESCE(nok_full_name,''), COALESCE(nok_relationship,''),
COALESCE(nok_phone,''), COALESCE(nok_email,''), COALESCE(source_of_funds,''), pep_status,
COALESCE(sacco_membership_disclosures,''), submitted_at, reviewed_at, COALESCE(review_note,''), COALESCE(reviewed_by,''), updated_at
FROM kyc_profiles WHERE user_id = $1`, userID)
	var k KYCProfile
	if err := row.Scan(&k.UserID, &k.Status, &k.DateOfBirth, &k.Gender, &k.Nationality, &k.OccupationStatus,
		&k.IDType, &k.IDNumber, &k.NOKFullName, &k.NOKRelationship, &k.NOKPhone, &k.NOKEmail, &k.SourceOfFunds, &k.PEPStatus,
		&k.SACCOMembershipDisclosures, &k.SubmittedAt, &k.ReviewedAt, &k.ReviewNote, &k.ReviewedBy, &k.LastUpdatedAt); err != nil {
		return KYCProfile{}, false
	}
	return k, true
}

func (s *PostgresStore) ReplaceKYCDocuments(userID string, docs []KYCDocument) error {
	ctx := context.Background()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, `DELETE FROM kyc_documents WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, d := range docs {
		if _, err := tx.Exec(ctx, `INSERT INTO kyc_documents (user_id, doc_type, doc_side, storage_key) VALUES ($1,$2,$3,$4)`,
			userID, d.DocType, nullStr(d.DocSide), d.StorageKey); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *PostgresStore) GetKYCDocuments(userID string) ([]KYCDocument, error) {
	ctx := context.Background()
	rows, err := s.pool.Query(ctx, `SELECT doc_type, COALESCE(doc_side,''), storage_key, created_at FROM kyc_documents WHERE user_id=$1 ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []KYCDocument
	for rows.Next() {
		var d KYCDocument
		if err := rows.Scan(&d.DocType, &d.DocSide, &d.StorageKey, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *PostgresStore) SetKYCDecision(userID, status, reviewedBy, reviewNote string) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `
UPDATE kyc_profiles
SET status=$2, reviewed_by=$3, review_note=$4, reviewed_at=now(), updated_at=now()
WHERE user_id=$1`, userID, status, reviewedBy, reviewNote)
	return err
}

func (s *PostgresStore) UpsertSaccoMembership(m SaccoMembership) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `
INSERT INTO sacco_memberships (
  user_id, status, district, county, sub_county, parish, village, street_plot, mobile_money_provider,
  mobile_money_number, secondary_momo_number, contribution_frequency, savings_goal_amount, savings_goal_purpose,
  shares_to_purchase, entrance_fee_payment_method, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16, now())
ON CONFLICT (user_id) DO UPDATE SET
  status=EXCLUDED.status,
  district=EXCLUDED.district,
  county=EXCLUDED.county,
  sub_county=EXCLUDED.sub_county,
  parish=EXCLUDED.parish,
  village=EXCLUDED.village,
  street_plot=EXCLUDED.street_plot,
  mobile_money_provider=EXCLUDED.mobile_money_provider,
  mobile_money_number=EXCLUDED.mobile_money_number,
  secondary_momo_number=EXCLUDED.secondary_momo_number,
  contribution_frequency=EXCLUDED.contribution_frequency,
  savings_goal_amount=EXCLUDED.savings_goal_amount,
  savings_goal_purpose=EXCLUDED.savings_goal_purpose,
  shares_to_purchase=EXCLUDED.shares_to_purchase,
  entrance_fee_payment_method=EXCLUDED.entrance_fee_payment_method,
  updated_at=now()`,
		m.UserID, m.Status, nullStr(m.District), nullStr(m.County), nullStr(m.SubCounty), nullStr(m.Parish), nullStr(m.Village),
		nullStr(m.StreetPlot), nullStr(m.MobileMoneyProvider), nullStr(m.MobileMoneyNumber), nullStr(m.SecondaryMoMoNumber),
		nullStr(m.ContributionFrequency), m.SavingsGoalAmount, nullStr(m.SavingsGoalPurpose), m.SharesToPurchase, nullStr(m.EntranceFeePaymentMethod))
	return err
}

func (s *PostgresStore) GetSaccoMembership(userID string) (SaccoMembership, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
SELECT user_id, status, COALESCE(district,''), COALESCE(county,''), COALESCE(sub_county,''), COALESCE(parish,''),
COALESCE(village,''), COALESCE(street_plot,''), COALESCE(mobile_money_provider,''), COALESCE(mobile_money_number,''),
COALESCE(secondary_momo_number,''), COALESCE(contribution_frequency,''), COALESCE(savings_goal_amount,0),
COALESCE(savings_goal_purpose,''), COALESCE(shares_to_purchase,0), COALESCE(entrance_fee_payment_method,''), updated_at
FROM sacco_memberships WHERE user_id=$1`, userID)
	var m SaccoMembership
	if err := row.Scan(&m.UserID, &m.Status, &m.District, &m.County, &m.SubCounty, &m.Parish, &m.Village, &m.StreetPlot,
		&m.MobileMoneyProvider, &m.MobileMoneyNumber, &m.SecondaryMoMoNumber, &m.ContributionFrequency,
		&m.SavingsGoalAmount, &m.SavingsGoalPurpose, &m.SharesToPurchase, &m.EntranceFeePaymentMethod, &m.LastUpdatedAt); err != nil {
		return SaccoMembership{}, false
	}
	return m, true
}

func (s *PostgresStore) EnsureSharesLedger(userID string, sharesUnits int) error {
	ctx := context.Background()
	if sharesUnits <= 0 {
		sharesUnits = 1
	}
	_, err := s.pool.Exec(ctx, `
INSERT INTO shares_ledger (user_id, balance_units, shares_balance)
VALUES ($1,$2,$2)
ON CONFLICT (user_id) DO UPDATE SET
  balance_units = GREATEST(shares_ledger.balance_units, EXCLUDED.balance_units),
  shares_balance = GREATEST(shares_ledger.shares_balance, EXCLUDED.shares_balance),
  updated_at = now()`, userID, sharesUnits)
	return err
}

func (s *PostgresStore) GetSharesUnits(userID string) (int, bool, error) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `SELECT COALESCE(balance_units, 0) FROM shares_ledger WHERE user_id=$1`, userID)
	var v int
	if err := row.Scan(&v); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return v, true, nil
}

func (s *PostgresStore) UpsertKibiinaPreference(k KibiinaPreference) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `
INSERT INTO kibiina_preferences (
  user_id, action, invite_code, group_name, contribution_amount, cycle_frequency, max_group_size,
  payout_order_preference, notification_preference, language_preference, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, now())
ON CONFLICT (user_id) DO UPDATE SET
  action=EXCLUDED.action,
  invite_code=EXCLUDED.invite_code,
  group_name=EXCLUDED.group_name,
  contribution_amount=EXCLUDED.contribution_amount,
  cycle_frequency=EXCLUDED.cycle_frequency,
  max_group_size=EXCLUDED.max_group_size,
  payout_order_preference=EXCLUDED.payout_order_preference,
  notification_preference=EXCLUDED.notification_preference,
  language_preference=EXCLUDED.language_preference,
  updated_at=now()`,
		k.UserID, k.Action, nullStr(k.InviteCode), nullStr(k.GroupName), k.ContributionAmount, nullStr(k.CycleFrequency),
		k.MaxGroupSize, nullStr(k.PayoutOrderPreference), nullStr(k.NotificationPreference), nullStr(k.LanguagePreference))
	return err
}

func (s *PostgresStore) GetKibiinaPreference(userID string) (KibiinaPreference, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
SELECT user_id, action, COALESCE(invite_code,''), COALESCE(group_name,''), COALESCE(contribution_amount,0),
COALESCE(cycle_frequency,''), COALESCE(max_group_size,0), COALESCE(payout_order_preference,''),
COALESCE(notification_preference,''), COALESCE(language_preference,''), updated_at
FROM kibiina_preferences WHERE user_id=$1`, userID)
	var k KibiinaPreference
	if err := row.Scan(&k.UserID, &k.Action, &k.InviteCode, &k.GroupName, &k.ContributionAmount, &k.CycleFrequency,
		&k.MaxGroupSize, &k.PayoutOrderPreference, &k.NotificationPreference, &k.LanguagePreference, &k.LastUpdatedAt); err != nil {
		return KibiinaPreference{}, false
	}
	return k, true
}

func (s *PostgresStore) ListAdminKYCQueue(status string, limit int) ([]AdminKYCQueueItem, error) {
	ctx := context.Background()
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	q := `
SELECT
  u.user_id,
  u.full_name,
  u.phone_e164,
  COALESCE(u.contact_email, ''),
  COALESCE(k.status, 'not_started'),
  COALESCE(k.id_type, ''),
  COALESCE(k.id_number, ''),
  k.submitted_at,
  k.reviewed_at
FROM users_identity u
LEFT JOIN kyc_profiles k ON k.user_id = u.user_id
WHERE ($1 = '' OR $1 = 'all' OR COALESCE(k.status, 'not_started') = $1)
ORDER BY COALESCE(k.submitted_at, u.created_at) DESC
LIMIT $2`
	rows, err := s.pool.Query(ctx, q, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AdminKYCQueueItem
	for rows.Next() {
		var k AdminKYCQueueItem
		if err := rows.Scan(&k.UserID, &k.FullName, &k.PhoneE164, &k.ContactEmail, &k.Status, &k.IDType, &k.IDNumber, &k.SubmittedAt, &k.ReviewedAt); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

func (s *PostgresStore) GetAdminSetting(key string) (map[string]any, bool, error) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `SELECT value FROM admin_settings WHERE key = $1`, key)
	var raw []byte
	if err := row.Scan(&raw); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, false, err
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, true, nil
}

func (s *PostgresStore) SetAdminSetting(key string, value map[string]any) error {
	ctx := context.Background()
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
INSERT INTO admin_settings (key, value, updated_at)
VALUES ($1, $2::jsonb, now())
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()`, key, string(raw))
	return err
}

func (s *PostgresStore) SetUserReferralCode(userID, referralCode string) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `
INSERT INTO user_referral_codes (user_id, referral_code)
VALUES ($1,$2)
ON CONFLICT (user_id) DO UPDATE SET referral_code = EXCLUDED.referral_code`, userID, referralCode)
	return err
}

func (s *PostgresStore) GetUserReferralCode(userID string) (string, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `SELECT referral_code FROM user_referral_codes WHERE user_id = $1`, userID)
	var code string
	if err := row.Scan(&code); err != nil {
		return "", false
	}
	return code, true
}

func (s *PostgresStore) FindUserIDByReferralCode(referralCode string) (string, bool) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `SELECT user_id FROM user_referral_codes WHERE referral_code = $1`, referralCode)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", false
	}
	return userID, true
}

// ListUsersWithFilters returns paginated list of users with search and filtering
func (s *PostgresStore) ListUsersWithFilters(search, statusFilter, kycFilter string, limit, offset int) ([]User, int, error) {
	ctx := context.Background()
	
	// Build WHERE clause dynamically
	whereClauses := []string{}
	args := []any{}
	argIdx := 1
	
	// Search filter (name, email, phone)
	if strings.TrimSpace(search) != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(`(
			LOWER(u.full_name) LIKE LOWER($%d) OR
			LOWER(u.auth_email) LIKE LOWER($%d) OR
			LOWER(COALESCE(u.contact_email, '')) LIKE LOWER($%d) OR
			u.phone_e164 LIKE $%d
		)`, argIdx, argIdx, argIdx, argIdx))
		searchPattern := "%" + strings.TrimSpace(search) + "%"
		args = append(args, searchPattern)
		argIdx++
	}
	
	// Status filter
	if strings.TrimSpace(statusFilter) != "" && statusFilter != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf(`COALESCE(u.status, 'active') = $%d`, argIdx))
		args = append(args, statusFilter)
		argIdx++
	}
	
	// KYC status filter
	if strings.TrimSpace(kycFilter) != "" && kycFilter != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf(`COALESCE(k.status, 'not_started') = $%d`, argIdx))
		args = append(args, kycFilter)
		argIdx++
	}
	
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}
	
	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM users_identity u
		LEFT JOIN kyc_profiles k ON k.user_id = u.user_id
		%s`, whereClause)
	
	var total int
	if err := s.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT u.user_id, u.full_name, u.phone_e164, u.auth_email, u.contact_email, 
		       u.supabase_user_id, COALESCE(u.supabase_login_email,''),
		       u.password_hash, u.phone_verified_at, u.contact_email_verified_at, 
		       u.created_at, u.date_of_birth, u.nationality,
		       COALESCE(u.role, 'user'), COALESCE(u.status, 'active'), 
		       u.role_assigned_at, u.role_assigned_by, u.last_login
		FROM users_identity u
		LEFT JOIN kyc_profiles k ON k.user_id = u.user_id
		%s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var u User
		var contact, supabaseID *string
		var supabaseLogin string
		var pwdHash *string
		var phoneV, emailV *time.Time
		var created time.Time
		var dob *time.Time
		var nat *string
		var role, status string
		var roleAssignedAt *time.Time
		var roleAssignedBy *string
		var lastLogin *time.Time
		
		err := rows.Scan(&u.ID, &u.FullName, &u.PhoneE164, &u.AuthEmail, &contact, &supabaseID, &supabaseLogin,
			&pwdHash, &phoneV, &emailV, &created, &dob, &nat, &role, &status, &roleAssignedAt, &roleAssignedBy, &lastLogin)
		if err != nil {
			return nil, 0, err
		}
		
		if contact != nil {
			u.ContactEmail = *contact
		}
		if supabaseID != nil {
			u.SupabaseUserID = *supabaseID
		}
		u.SupabaseLoginEmail = supabaseLogin
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
		u.Role = role
		u.Status = status
		if roleAssignedAt != nil {
			t := *roleAssignedAt
			u.RoleAssignedAt = &t
		}
		if roleAssignedBy != nil {
			u.RoleAssignedBy = *roleAssignedBy
		}
		if lastLogin != nil {
			t := *lastLogin
			u.LastLogin = &t
		}
		
		users = append(users, u)
	}
	
	return users, total, rows.Err()
}

// UpdateUserStatus updates user account status
func (s *PostgresStore) UpdateUserStatus(userID, status, adminID, reason string) error {
	ctx := context.Background()
	
	// Update user status
	_, err := s.pool.Exec(ctx, `
		UPDATE users_identity 
		SET status = $2, updated_at = now()
		WHERE user_id = $1`,
		userID, status)
	
	if err != nil {
		return err
	}
	
	// Log status change in audit table (if exists)
	_, _ = s.pool.Exec(ctx, `
		INSERT INTO admin_audit_log (admin_id, action, target_user_id, details, created_at)
		VALUES ($1, 'user_status_change', $2, $3, now())`,
		adminID, userID, fmt.Sprintf("status changed to %s: %s", status, reason))
	
	return nil
}

// UpdateUserRole updates user role
func (s *PostgresStore) UpdateUserRole(userID, role, adminID string) error {
	ctx := context.Background()
	
	_, err := s.pool.Exec(ctx, `
		UPDATE users_identity 
		SET role = $2, role_assigned_at = now(), role_assigned_by = $3, updated_at = now()
		WHERE user_id = $1`,
		userID, role, adminID)
	
	return err
}

// GetUserActivity returns user activity log
func (s *PostgresStore) GetUserActivity(userID string, since time.Time) ([]ActivityLog, error) {
	ctx := context.Background()
	
	// Query activity from multiple sources
	query := `
		SELECT 
			'login' as action,
			'User logged in' as details,
			created_at as timestamp,
			'' as ip_address,
			'' as device_info
		FROM user_sessions
		WHERE user_id = $1 AND created_at >= $2
		
		UNION ALL
		
		SELECT 
			'status_change' as action,
			details,
			created_at as timestamp,
			'' as ip_address,
			'' as device_info
		FROM admin_audit_log
		WHERE target_user_id = $1 AND created_at >= $2
		
		UNION ALL
		
		SELECT 
			action,
			CONCAT('Role changed from ', previous_role, ' to ', new_role) as details,
			created_at as timestamp,
			'' as ip_address,
			'' as device_info
		FROM admin_role_audit
		WHERE user_id = $1 AND created_at >= $2
		
		ORDER BY timestamp DESC
		LIMIT 100`
	
	rows, err := s.pool.Query(ctx, query, userID, since)
	if err != nil {
		// If tables don't exist, return empty array
		return []ActivityLog{}, nil
	}
	defer rows.Close()
	
	var activities []ActivityLog
	for rows.Next() {
		var a ActivityLog
		err := rows.Scan(&a.Action, &a.Details, &a.Timestamp, &a.IPAddress, &a.DeviceInfo)
		if err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	
	return activities, rows.Err()
}
