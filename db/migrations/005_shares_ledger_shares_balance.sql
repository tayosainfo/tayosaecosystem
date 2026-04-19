-- Triggers in some environments reference shares_balance (not balance_units).
ALTER TABLE public.shares_ledger
  ADD COLUMN IF NOT EXISTS shares_balance NUMERIC(20, 4) NOT NULL DEFAULT 0;
