-- FK checks on INSERT can fail if a BEFORE INSERT trigger on users_identity inserts into
-- shares_ledger before the parent row is visible. Defer validation to end of transaction.
ALTER TABLE public.shares_ledger DROP CONSTRAINT IF EXISTS shares_ledger_user_id_fkey;

ALTER TABLE public.shares_ledger
  ADD CONSTRAINT shares_ledger_user_id_fkey
  FOREIGN KEY (user_id) REFERENCES public.users_identity (user_id) ON DELETE CASCADE
  DEFERRABLE INITIALLY DEFERRED;
