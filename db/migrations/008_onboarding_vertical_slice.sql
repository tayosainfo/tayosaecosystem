-- Vertical-slice onboarding/compliance model

CREATE TABLE IF NOT EXISTS user_consents (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id) ON DELETE CASCADE,
  terms_accepted_at TIMESTAMP NOT NULL,
  privacy_accepted_at TIMESTAMP NOT NULL,
  terms_version TEXT NULL,
  privacy_version TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_referral_codes (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id) ON DELETE CASCADE,
  referral_code TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS kyc_profiles (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'not_started',
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
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_kyc_profiles_status ON kyc_profiles(status);

CREATE TABLE IF NOT EXISTS kyc_documents (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users_identity(user_id) ON DELETE CASCADE,
  doc_type TEXT NOT NULL,
  doc_side TEXT NULL,
  storage_key TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_kyc_documents_user ON kyc_documents(user_id);

CREATE TABLE IF NOT EXISTS sacco_memberships (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'not_started',
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
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS kibiina_preferences (
  user_id TEXT PRIMARY KEY REFERENCES users_identity(user_id) ON DELETE CASCADE,
  action TEXT NOT NULL DEFAULT 'none',
  invite_code TEXT NULL,
  group_name TEXT NULL,
  contribution_amount NUMERIC(14,2) NULL,
  cycle_frequency TEXT NULL,
  max_group_size INTEGER NULL,
  payout_order_preference TEXT NULL,
  notification_preference TEXT NULL,
  language_preference TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Keep shares ledger initialization explicit and idempotent.
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'uq_shares_ledger_user'
  ) THEN
    ALTER TABLE shares_ledger
      ADD CONSTRAINT uq_shares_ledger_user UNIQUE (user_id);
  END IF;
END $$;
