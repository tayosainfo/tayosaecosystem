-- ============================================================================
-- SUPABASE ADMIN RLS FIX
-- ============================================================================
-- This SQL script fixes the insufficient permission issue for admin users
-- in the KYC management dashboard.
--
-- INSTRUCTIONS:
-- 1. Go to your Supabase Dashboard: https://app.supabase.com
-- 2. Select your project: ablvrbnbsdqshrorhmjf
-- 3. Go to SQL Editor
-- 4. Click "New Query"
-- 5. Paste this entire script
-- 6. Click "Run"
-- 7. Test the endpoint: GET https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending
--
-- ============================================================================

-- Step 1: Create the exec_sql function (if it doesn't exist)
-- This allows us to execute SQL via RPC
CREATE OR REPLACE FUNCTION public.exec_sql(sql text)
RETURNS void
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  EXECUTE sql;
END;
$$;

GRANT EXECUTE ON FUNCTION public.exec_sql(text) TO service_role;

-- Step 2: Drop existing RLS policies on kyc_documents
DROP POLICY IF EXISTS "Users can view own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can insert own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can update own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can delete own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Service role full access to KYC documents" ON public.kyc_documents;

-- Step 3: Create new RLS policies that allow admin users to view all KYC documents
-- Policy: Users can view their own KYC documents OR admin users can view all
CREATE POLICY "Users can view own KYC documents"
ON public.kyc_documents
FOR SELECT
TO authenticated
USING (
  user_id = auth.uid()::text
  OR
  -- Allow admin users to view all KYC documents
  EXISTS (
    SELECT 1 FROM public.users_identity
    WHERE supabase_user_id = auth.uid()::text
    AND role = 'admin'
  )
);

-- Policy: Users can insert their own KYC documents
CREATE POLICY "Users can insert own KYC documents"
ON public.kyc_documents
FOR INSERT
TO authenticated
WITH CHECK (user_id = auth.uid()::text);

-- Policy: Users can update their own KYC documents
CREATE POLICY "Users can update own KYC documents"
ON public.kyc_documents
FOR UPDATE
TO authenticated
USING (user_id = auth.uid()::text)
WITH CHECK (user_id = auth.uid()::text);

-- Policy: Users can delete their own KYC documents
CREATE POLICY "Users can delete own KYC documents"
ON public.kyc_documents
FOR DELETE
TO authenticated
USING (user_id = auth.uid()::text);

-- Policy: Service role can do everything (for backend operations)
CREATE POLICY "Service role full access to KYC documents"
ON public.kyc_documents
FOR ALL
TO service_role
USING (true)
WITH CHECK (true);

-- Step 4: Update RLS policies on users_identity table
-- Drop existing policies
DROP POLICY IF EXISTS "Enable read access for role checking" ON public.users_identity;
DROP POLICY IF EXISTS "Users can read their own role" ON public.users_identity;
DROP POLICY IF EXISTS "Allow reading user roles for authenticated users" ON public.users_identity;

-- Create a permissive policy that allows reading the role field
-- This is needed for the custom JWT claims hook to work
CREATE POLICY "Enable read access for role checking"
ON public.users_identity
FOR SELECT
USING (true);

-- Step 5: Verify the policies were created
SELECT 
    tablename,
    policyname,
    permissive,
    roles,
    qual
FROM pg_policies
WHERE tablename IN ('kyc_documents', 'users_identity')
ORDER BY tablename, policyname;

-- ============================================================================
-- DONE! The admin user should now be able to access the KYC endpoint.
-- Test: GET https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending
-- ============================================================================
