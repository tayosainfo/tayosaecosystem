-- Admin settings key-value store for platform configuration (fees, charges, etc.)

CREATE TABLE IF NOT EXISTS admin_settings (
  key TEXT PRIMARY KEY,
  value JSONB NOT NULL DEFAULT '{}'::jsonb,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Seed defaults (idempotent).
INSERT INTO admin_settings (key, value)
VALUES
  ('fees', jsonb_build_object(
    'registrationFeeUGX', 0,
    'saccoEntranceFeeUGX', 0,
    'transactionFeePct', 0
  ))
ON CONFLICT (key) DO NOTHING;

