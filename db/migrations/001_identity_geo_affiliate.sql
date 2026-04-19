-- Identity model for dual identifier login
CREATE TABLE IF NOT EXISTS users_identity (
  user_id TEXT PRIMARY KEY,
  full_name TEXT NOT NULL,
  phone_e164 TEXT NOT NULL UNIQUE,
  auth_email TEXT NOT NULL UNIQUE,
  contact_email TEXT UNIQUE,
  insforge_user_id TEXT UNIQUE,
  phone_verified_at TIMESTAMP NULL,
  contact_email_verified_at TIMESTAMP NULL,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Onboarding profile state (phase 1-4)
CREATE TABLE IF NOT EXISTS onboarding_profiles (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id),
  phase INTEGER NOT NULL DEFAULT 1,
  referral_code TEXT NULL,
  district TEXT NULL,
  county TEXT NULL,
  sub_county TEXT NULL,
  parish TEXT NULL,
  village TEXT NULL,
  trust_score_seed INTEGER NOT NULL DEFAULT 10,
  phase_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Canonical Uganda geo hierarchy data
CREATE TABLE IF NOT EXISTS uganda_geo_units (
  id BIGSERIAL PRIMARY KEY,
  district TEXT NOT NULL,
  county TEXT NOT NULL,
  sub_county TEXT NOT NULL,
  parish TEXT NOT NULL,
  village TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_geo_district ON uganda_geo_units(district);
CREATE INDEX IF NOT EXISTS idx_geo_county ON uganda_geo_units(county);
CREATE INDEX IF NOT EXISTS idx_geo_sub_county ON uganda_geo_units(sub_county);
CREATE INDEX IF NOT EXISTS idx_geo_parish ON uganda_geo_units(parish);

-- Parish SACCO to village Kibiina model
CREATE TABLE IF NOT EXISTS parish_saccos (
  id BIGSERIAL PRIMARY KEY,
  parish TEXT NOT NULL,
  district TEXT NOT NULL,
  sub_county TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS village_kibiina_groups (
  id BIGSERIAL PRIMARY KEY,
  parish_sacco_id BIGINT NOT NULL REFERENCES parish_saccos(id),
  village TEXT NOT NULL,
  group_name TEXT NOT NULL,
  cycle_frequency TEXT NOT NULL,
  contribution_amount NUMERIC(14,2) NOT NULL,
  payout_method TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_kibiina_parish_sacco ON village_kibiina_groups(parish_sacco_id);

-- Affiliate referral rewards lifecycle
CREATE TABLE IF NOT EXISTS affiliate_referrals (
  id BIGSERIAL PRIMARY KEY,
  referral_code TEXT NOT NULL,
  referrer_user_id TEXT NOT NULL REFERENCES users_identity(user_id),
  referee_user_id TEXT NOT NULL REFERENCES users_identity(user_id),
  reward_points INTEGER NOT NULL DEFAULT 100,
  status TEXT NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_affiliate_referrer ON affiliate_referrals(referrer_user_id);

