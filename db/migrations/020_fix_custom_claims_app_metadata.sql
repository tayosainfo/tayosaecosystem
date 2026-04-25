-- Migration: Fix custom claims hook to add user_role to app_metadata
-- Description: Updates the custom JWT claims hook to add user_role to both JWT claims and app_metadata
-- Date: 2026-04-25

-- Drop the existing function
DROP FUNCTION IF EXISTS public.custom_access_token_hook(jsonb);

-- Create updated custom access token hook function
-- This function is called by Supabase Auth when generating JWT tokens
CREATE OR REPLACE FUNCTION public.custom_access_token_hook(event jsonb)
RETURNS jsonb
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
DECLARE
  claims jsonb;
  user_role text;
BEGIN
  -- Fetch the user role from users_identity table
  SELECT role INTO user_role
  FROM public.users_identity
  WHERE supabase_user_id = (event->>'user_id')::text;

  -- Get existing claims from the event
  claims := event->'claims';

  IF user_role IS NOT NULL THEN
    -- Set custom claim for user role in JWT claims
    claims := jsonb_set(claims, '{user_role}', to_jsonb(user_role));
    
    -- Also add to app_metadata for easier access via /auth/v1/user endpoint
    claims := jsonb_set(claims, '{app_metadata, user_role}', to_jsonb(user_role));
  ELSE
    -- Default to 'user' if no role found
    claims := jsonb_set(claims, '{user_role}', '"user"');
    claims := jsonb_set(claims, '{app_metadata, user_role}', '"user"');
  END IF;

  -- Update the 'claims' object in the original event
  event := jsonb_set(event, '{claims}', claims);

  RETURN event;
END;
$$;

-- Grant necessary permissions to supabase_auth_admin role
GRANT USAGE ON SCHEMA public TO supabase_auth_admin;
GRANT SELECT ON public.users_identity TO supabase_auth_admin;
GRANT EXECUTE ON FUNCTION public.custom_access_token_hook TO supabase_auth_admin;

-- Add comment for documentation
COMMENT ON FUNCTION public.custom_access_token_hook IS 'Supabase Auth hook to add user_role to JWT claims and app_metadata';
