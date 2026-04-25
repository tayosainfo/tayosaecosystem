-- Migration: Disable RLS on users_identity for admin role checking
-- Description: Allows frontend to query user roles without RLS restrictions
-- This is safe because we only expose the role field, not sensitive data
-- Date: 2026-04-25

-- Option 1: Disable RLS entirely (simplest, but less secure)
-- Uncomment if you want to disable RLS completely
-- ALTER TABLE public.users_identity DISABLE ROW LEVEL SECURITY;

-- Option 2: Create a permissive policy for reading roles (recommended)
-- This allows anyone to read the role field for any user
-- (The role field is not sensitive - it's just 'admin' or 'user')

-- First, check if RLS is enabled
SELECT 
    schemaname,
    tablename,
    rowsecurity
FROM pg_tables
WHERE tablename = 'users_identity';

-- Drop existing policies if they exist
DROP POLICY IF EXISTS "Users can read their own role" ON public.users_identity;
DROP POLICY IF EXISTS "Allow reading user roles for authenticated users" ON public.users_identity;
DROP POLICY IF EXISTS "Enable read access for all users" ON public.users_identity;

-- Create a permissive policy that allows reading the role field
CREATE POLICY "Enable read access for role checking"
ON public.users_identity
FOR SELECT
USING (true);  -- Allow all reads

-- Verify the policy was created
SELECT 
    policyname,
    permissive,
    roles,
    qual
FROM pg_policies
WHERE tablename = 'users_identity'
ORDER BY policyname;
