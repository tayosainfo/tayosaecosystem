-- Share ledger scaffold (SACCO / member equity tracking).
-- Some environments define a TRIGGER that inserts into this table when a row is added to
-- users_identity. If the table is missing, registration fails with:
--   relation "public.shares_ledger" does not exist

CREATE TABLE IF NOT EXISTS public.shares_ledger (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES public.users_identity(user_id) ON DELETE CASCADE,
  balance_units NUMERIC(20, 4) NOT NULL DEFAULT 0,
  shares_balance NUMERIC(20, 4) NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_shares_ledger_user_id ON public.shares_ledger(user_id);
