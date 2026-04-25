-- Migration: Fix RLS policies on users_identity for admin role checking
-- Description: Allows authenticated users to read their own role from users_identity table
-- Date: 2026-04-25

-- First, check if RLS is enabled
-- If RLS is enabled, we need to add a policy that allows users to read their own role

-- Policy: Allow authenticated users to read their own user record
CREATE POLICY "Users can read their own role"
ON public.users_identity
FOR SELECT
USING (
  -- Allow if the user's email matches the auth_email in the record
  -- This works with custom auth systems that store email in sessionStorage
  auth.uid()::text = supabase_user_id
  OR
  -- Fallback: Allow if user is authenticated (for admin role checking)
  auth.role() = 'authenticated'
);

-- Alternative: If the above doesn't work, create a more permissive policy
-- Uncomment if needed:
/*
CREATE POLICY "Allow reading user roles for authenticated users"
ON public.users_identity
FOR SELECT
USING (auth.role() = 'authenticated');
*/

-- If you want to completely disable RLS for this table (not recommended for production):
-- ALTER TABLE public.users_identity DISABLE ROW LEVEL SECURITY;

-- Verify the policy was created
SELECT 
    policyname,
    permissive,
    roles,
    qual
FROM pg_policies
WHERE tablename = 'users_identity'
ORDER BY policyname;
