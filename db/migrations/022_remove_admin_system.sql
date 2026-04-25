-- Migration: Remove admin system
-- Description: Removes admin roles, audit tables, and custom claims since KYC is now auto-approved
-- Date: 2026-04-25

-- Drop admin-related triggers
DROP TRIGGER IF EXISTS trigger_log_role_change ON users_identity;

-- Drop admin-related functions
DROP FUNCTION IF EXISTS log_role_change();
DROP FUNCTION IF EXISTS custom_access_token_hook(jsonb);

-- Drop admin audit table
DROP TABLE IF EXISTS admin_role_audit CASCADE;

-- Remove role columns from users_identity
ALTER TABLE users_identity 
DROP COLUMN IF EXISTS role,
DROP COLUMN IF EXISTS role_assigned_at,
DROP COLUMN IF EXISTS role_assigned_by;

-- Drop the user_role enum type
DROP TYPE IF EXISTS user_role CASCADE;

-- Drop the admin-related index
DROP INDEX IF EXISTS idx_users_identity_role;

-- Revoke permissions from supabase_auth_admin (if they were granted)
REVOKE ALL ON public.users_identity FROM supabase_auth_admin;
REVOKE ALL ON FUNCTION public.custom_access_token_hook FROM supabase_auth_admin;

-- Add comment
COMMENT ON TABLE users_identity IS 'User identity table - admin system removed, KYC is auto-approved';
