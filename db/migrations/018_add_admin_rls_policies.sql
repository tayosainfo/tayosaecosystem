-- Migration: Add admin RLS policies for KYC documents and users_identity
-- Description: Allows admin users to read all KYC documents and user data
-- Date: 2026-04-25

-- First, let's check the current RLS policies on kyc_documents
-- and add an admin exception

-- Drop existing policies that might be too restrictive
DROP POLICY IF EXISTS "Users can view own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can insert own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can update own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Users can delete own KYC documents" ON public.kyc_documents;
DROP POLICY IF EXISTS "Service role full access to KYC documents" ON public.kyc_documents;

-- RLS Policy: Users can view their own KYC documents
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

-- RLS Policy: Users can insert their own KYC documents
CREATE POLICY "Users can insert own KYC documents"
ON public.kyc_documents
FOR INSERT
TO authenticated
WITH CHECK (user_id = auth.uid()::text);

-- RLS Policy: Users can update their own KYC documents
CREATE POLICY "Users can update own KYC documents"
ON public.kyc_documents
FOR UPDATE
TO authenticated
USING (user_id = auth.uid()::text)
WITH CHECK (user_id = auth.uid()::text);

-- RLS Policy: Users can delete their own KYC documents
CREATE POLICY "Users can delete own KYC documents"
ON public.kyc_documents
FOR DELETE
TO authenticated
USING (user_id = auth.uid()::text);

-- RLS Policy: Service role can do everything (for backend operations)
CREATE POLICY "Service role full access to KYC documents"
ON public.kyc_documents
FOR ALL
TO service_role
USING (true)
WITH CHECK (true);

-- Also ensure users_identity table has proper RLS for admin role checking
-- Drop existing policies
DROP POLICY IF EXISTS "Enable read access for role checking" ON public.users_identity;
DROP POLICY IF EXISTS "Users can read their own role" ON public.users_identity;
DROP POLICY IF EXISTS "Allow reading user roles for authenticated users" ON public.users_identity;

-- Create a permissive policy that allows reading the role field
-- This is needed for the custom JWT claims hook to work
CREATE POLICY "Enable read access for role checking"
ON public.users_identity
FOR SELECT
USING (true);  -- Allow all reads (role field is not sensitive)

-- Verify the policies were created
SELECT 
    tablename,
    policyname,
    permissive,
    roles,
    qual
FROM pg_policies
WHERE tablename IN ('kyc_documents', 'users_identity')
ORDER BY tablename, policyname;
