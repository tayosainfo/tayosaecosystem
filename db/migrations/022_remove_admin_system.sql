-- Migration: Remove admin system
-- Description: Removes admin roles, audit tables, and custom claims since KYC is now auto-approved
-- Date: 2026-04-25

-- First, drop RLS policies that depend on the role column
DROP POLICY IF EXISTS "Users can view own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can insert own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can update own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can delete own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Service role full access to KYC documents" ON public.kyc_documents;

-- Recreate simple RLS policies without admin role checks
CREATE POLICY "Users can view own KYC documents"
ON public.kyc_documents
FOR SELECT
TO authenticated
USING (user_id = auth.uid()::text);

CREATE POLICY "Users can insert own KYC documents"
ON public.kyc_documents
FOR INSERT
TO authenticated
WITH CHECK (user_id = auth.uid()::text);

CREATE POLICY "Users can update own KYC documents"
ON public.kyc_documents
FOR UPDATE
TO authenticated
USING (user_id = auth.uid()::text)
WITH CHECK (user_id = auth.uid()::text);

CREATE POLICY "Users can delete own KYC documents"
ON public.kyc_documents
FOR DELETE
TO authenticated
USING (user_id = auth.uid()::text);

CREATE POLICY "Service role full access to KYC documents"
ON public.kyc_documents
FOR ALL
TO service_role
USING (true)
WITH CHECK (true);

-- Drop admin-related triggers
DROP TRIGGER IF EXISTS trigger_log_role_change ON users_identity;

-- Drop admin-related functions
DROP FUNCTION IF EXISTS log_role_change();
DROP FUNCTION IF EXISTS custom_access_token_hook(jsonb);

-- Drop admin audit table
DROP TABLE IF EXISTS admin_role_audit CASCADE;

-- Drop the admin-related index
DROP INDEX IF EXISTS idx_users_identity_role;

-- Remove role columns from users_identity
ALTER TABLE users_identity 
DROP COLUMN IF EXISTS role CASCADE,
DROP COLUMN IF EXISTS role_assigned_at,
DROP COLUMN IF EXISTS role_assigned_by;

-- Drop the user_role enum type
DROP TYPE IF EXISTS user_role CASCADE;

-- Revoke permissions from supabase_auth_admin (if they were granted)
REVOKE ALL ON public.users_identity FROM supabase_auth_admin;

-- Add comment
COMMENT ON TABLE users_identity IS 'User identity table - admin system removed, KYC is auto-approved';
