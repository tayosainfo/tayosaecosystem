-- Migration: Restore custom_access_token_hook as a no-op
-- Description: Supabase is configured to call this hook, so we need it to exist
--              but since we removed the admin system, it just passes through without changes
-- Date: 2026-04-25

-- Create a simple pass-through custom access token hook
-- This doesn't add any custom claims, it just returns the event unchanged
CREATE OR REPLACE FUNCTION public.custom_access_token_hook(event jsonb)
RETURNS jsonb
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  -- Just return the event unchanged
  -- No custom claims needed since we removed the admin system
  RETURN event;
END;
$$;

-- Grant necessary permissions to supabase_auth_admin role
GRANT USAGE ON SCHEMA public TO supabase_auth_admin;
GRANT EXECUTE ON FUNCTION public.custom_access_token_hook TO supabase_auth_admin;

-- Add comment
COMMENT ON FUNCTION public.custom_access_token_hook IS 'Pass-through hook - no custom claims added (admin system removed)';
