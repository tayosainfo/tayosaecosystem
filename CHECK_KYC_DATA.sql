-- Quick check for KYC data in your database
-- Run this in Supabase SQL Editor

-- 1. Check if there are any users
SELECT COUNT(*) as total_users FROM public.users_identity;

-- 2. Check if there are any KYC profiles
SELECT COUNT(*) as total_kyc_profiles FROM public.kyc_profiles;

-- 3. Check KYC profiles by status
SELECT 
    status,
    COUNT(*) as count
FROM public.kyc_profiles
GROUP BY status;

-- 4. List all users with their KYC status
SELECT
  u.user_id,
  u.full_name,
  u.phone_e164,
  u.contact_email,
  COALESCE(k.status, 'not_started') as kyc_status,
  k.submitted_at
FROM public.users_identity u
LEFT JOIN public.kyc_profiles k ON k.user_id = u.user_id
ORDER BY COALESCE(k.submitted_at, u.created_at) DESC
LIMIT 20;

-- 5. Check if tables exist
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN ('users_identity', 'kyc_profiles', 'kyc_documents')
ORDER BY table_name;
