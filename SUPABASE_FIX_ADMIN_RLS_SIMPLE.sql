-- ============================================================================
-- SUPABASE ADMIN RLS FIX - SIMPLIFIED VERSION
-- ============================================================================
-- This SQL script fixes the insufficient permission issue for admin users
-- in the KYC management dashboard.
--
-- This version doesn't use functions, just direct RLS policy updates.
--
-- ============================================================================

-- Step 1: Drop existing RLS policies on kyc_documents
DROP POLICY IF EXISTS "Users can view own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can insert own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can update own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can delete own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Service role full access to KYC documents" ON public.kyc_documents;

-- Step 2: Create new RLS policies that allow admin users to view all KYC documents
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

-- Step 3: Update RLS policies on users_identity table
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

-- ============================================================================
-- DONE! The admin user should now be able to access the KYC endpoint.
-- Test: GET https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending
-- ============================================================================
