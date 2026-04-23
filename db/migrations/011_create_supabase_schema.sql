-- Create fresh Supabase database schema
-- This migration creates the complete schema for a new Supabase database
-- with correct column names (supabase_* instead of insforge_*)

-- Identity model for dual identifier login
CREATE TABLE IF NOT EXISTS users_identity (
  user_id TEXT PRIMARY KEY,
  full_name TEXT NOT NULL,
  phone_e164 TEXT NOT NULL UNIQUE,
  auth_email TEXT NOT NULL UNIQUE,
  contact_email TEXT UNIQUE,
  supabase_user_id TEXT UNIQUE,
  supabase_login_email TEXT,
  password_hash TEXT,
  phone_verified_at TIMESTAMP NULL,
  contact_email_verified_at TIMESTAMP NULL,
  date_of_birth DATE NULL,
  nationality TEXT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on supabase_login_email
CREATE INDEX IF NOT EXISTS idx_users_supabase_login_email ON users_identity(supabase_login_email);

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

-- Additional tables that may be referenced by the application
-- (Add other tables from migrations 003-010 as needed)

-- User consents tracking
CREATE TABLE IF NOT EXISTS user_consents (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id),
  terms_accepted_at TIMESTAMP NULL,
  privacy_accepted_at TIMESTAMP NULL,
  terms_version TEXT NULL,
  privacy_version TEXT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- KYC documents and profiles
CREATE TABLE IF NOT EXISTS kyc_profiles (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id),
  status TEXT NOT NULL DEFAULT 'pending',
  date_of_birth DATE NULL,
  gender TEXT NULL,
  nationality TEXT NULL,
  occupation_status TEXT NULL,
  id_type TEXT NULL,
  id_number TEXT NULL,
  nok_full_name TEXT NULL,
  nok_relationship TEXT NULL,
  nok_phone TEXT NULL,
  nok_email TEXT NULL,
  source_of_funds TEXT NULL,
  pep_status BOOLEAN NULL,
  sacco_membership_disclosures TEXT NULL,
  submitted_at TIMESTAMP NULL,
  reviewed_at TIMESTAMP NULL,
  review_note TEXT NULL,
  reviewed_by TEXT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- SACCO membership details
CREATE TABLE IF NOT EXISTS sacco_memberships (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id),
  status TEXT NOT NULL DEFAULT 'pending',
  district TEXT NULL,
  county TEXT NULL,
  sub_county TEXT NULL,
  parish TEXT NULL,
  village TEXT NULL,
  street_plot TEXT NULL,
  mobile_money_provider TEXT NULL,
  mobile_money_number TEXT NULL,
  secondary_momo_number TEXT NULL,
  contribution_frequency TEXT NULL,
  savings_goal_amount NUMERIC(14,2) NULL,
  savings_goal_purpose TEXT NULL,
  shares_to_purchase INTEGER NULL,
  entrance_fee_payment_method TEXT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Kibiina preferences
CREATE TABLE IF NOT EXISTS kibiina_preferences (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id),
  action TEXT NOT NULL,
  invite_code TEXT NULL,
  group_name TEXT NULL,
  contribution_amount NUMERIC(14,2) NULL,
  cycle_frequency TEXT NULL,
  max_group_size INTEGER NULL,
  payout_order_preference TEXT NULL,
  notification_preference TEXT NULL,
  language_preference TEXT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Shares ledger for tracking user shares
CREATE TABLE IF NOT EXISTS shares_ledger (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users_identity(user_id),
  transaction_type TEXT NOT NULL,
  shares_amount INTEGER NOT NULL,
  shares_balance INTEGER NOT NULL,
  description TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_shares_ledger_user ON shares_ledger(user_id);

-- Admin settings
CREATE TABLE IF NOT EXISTS admin_settings (
  key TEXT PRIMARY KEY,
  value JSONB NOT NULL,
  description TEXT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default admin settings
INSERT INTO admin_settings (key, value, description) VALUES
  ('shares_price_ugx', '10000', 'Price per share in UGX'),
  ('entrance_fee_ugx', '5000', 'SACCO entrance fee in UGX'),
  ('min_shares_purchase', '1', 'Minimum shares a user must purchase'),
  ('max_shares_purchase', '100', 'Maximum shares a user can purchase')
ON CONFLICT (key) DO NOTHING;
