-- If a trigger inserts into shares_ledger from another transaction (e.g. auth user creation)
-- before public.users_identity has the same user_id, a strict FK cannot succeed.
-- Drop the FK so registration is not blocked; enforce linkage in application logic or
-- recreate the FK once triggers only fire AFTER INSERT on users_identity in the same DB.
ALTER TABLE public.shares_ledger DROP CONSTRAINT IF EXISTS shares_ledger_user_id_fkey;
