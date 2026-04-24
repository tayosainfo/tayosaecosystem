-- KYC Data Verification Queries
-- Run these in Supabase SQL Editor to verify KYC data is stored correctly

-- ============================================================================
-- 1. CHECK ALL KYC PROFILES (Most Recent First)
-- ============================================================================
SELECT 
    user_id,
    status,
    date_of_birth,
    gender,
    nationality,
    occupation_status,
    id_type,
    id_number,
    nok_full_name,
    nok_relationship,
    nok_phone,
    source_of_funds,
    pep_status,
    submitted_at,
    reviewed_at
FROM public.kyc_profiles
ORDER BY submitted_at DESC
LIMIT 10;

-- ============================================================================
-- 2. CHECK KYC DOCUMENTS (File References)
-- ============================================================================
SELECT 
    kd.user_id,
    kd.doc_type,
    kd.doc_side,
    kd.storage_key,
    kd.uploaded_at,
    ui.full_name,
    ui.contact_email
FROM public.kyc_documents kd
LEFT JOIN public.users_identity ui ON kd.user_id = ui.id
ORDER BY kd.uploaded_at DESC
LIMIT 20;

-- ============================================================================
-- 3. CHECK UPLOADED FILES IN STORAGE
-- ============================================================================
SELECT 
    name AS file_path,
    bucket_id,
    created_at,
    updated_at,
    (metadata->>'size')::bigint AS file_size_bytes,
    metadata->>'mimetype' AS mime_type
FROM storage.objects
WHERE bucket_id = 'collateral_docs'
ORDER BY created_at DESC
LIMIT 20;

-- ============================================================================
-- 4. COMPLETE KYC VIEW (Profile + Documents + User Info)
-- ============================================================================
SELECT 
    ui.id AS user_id,
    ui.full_name,
    ui.contact_email,
    ui.phone_e164,
    kp.status AS kyc_status,
    kp.id_type,
    kp.id_number,
    kp.submitted_at,
    kp.reviewed_at,
    COUNT(kd.id) AS document_count
FROM public.users_identity ui
LEFT JOIN public.kyc_profiles kp ON ui.id = kp.user_id
LEFT JOIN public.kyc_documents kd ON ui.id = kd.user_id
WHERE kp.user_id IS NOT NULL
GROUP BY ui.id, ui.full_name, ui.contact_email, ui.phone_e164, 
         kp.status, kp.id_type, kp.id_number, kp.submitted_at, kp.reviewed_at
ORDER BY kp.submitted_at DESC;

-- ============================================================================
-- 5. CHECK SPECIFIC USER'S KYC DATA (Replace YOUR_USER_ID)
-- ============================================================================
-- First, find your user ID:
SELECT id, full_name, contact_email, phone_e164 
FROM public.users_identity 
WHERE contact_email = 'YOUR_EMAIL@example.com';

-- Then check KYC profile:
SELECT * FROM public.kyc_profiles 
WHERE user_id = 'YOUR_USER_ID';

-- And documents:
SELECT * FROM public.kyc_documents 
WHERE user_id = 'YOUR_USER_ID';

-- And storage files:
SELECT name, created_at, metadata 
FROM storage.objects 
WHERE bucket_id = 'collateral_docs' 
AND name LIKE 'YOUR_USER_ID/%';

-- ============================================================================
-- 6. KYC STATISTICS
-- ============================================================================
SELECT 
    status,
    COUNT(*) AS count,
    MIN(submitted_at) AS first_submission,
    MAX(submitted_at) AS last_submission
FROM public.kyc_profiles
GROUP BY status
ORDER BY count DESC;

-- ============================================================================
-- 7. RECENT KYC SUBMISSIONS (Last 24 Hours)
-- ============================================================================
SELECT 
    ui.full_name,
    ui.contact_email,
    kp.status,
    kp.id_type,
    kp.submitted_at,
    COUNT(kd.id) AS documents_uploaded
FROM public.kyc_profiles kp
JOIN public.users_identity ui ON kp.user_id = ui.id
LEFT JOIN public.kyc_documents kd ON kp.user_id = kd.user_id
WHERE kp.submitted_at > NOW() - INTERVAL '24 hours'
GROUP BY ui.full_name, ui.contact_email, kp.status, kp.id_type, kp.submitted_at
ORDER BY kp.submitted_at DESC;

-- ============================================================================
-- 8. VERIFY FILE STORAGE PATHS ARE CORRECT
-- ============================================================================
-- Check if storage keys match the expected format: {user_id}/kyc/{timestamp}-{filename}
SELECT 
    kd.user_id,
    kd.doc_type,
    kd.storage_key,
    CASE 
        WHEN kd.storage_key LIKE kd.user_id || '/kyc/%' THEN 'CORRECT'
        ELSE 'INCORRECT'
    END AS path_format_check
FROM public.kyc_documents kd
ORDER BY kd.uploaded_at DESC
LIMIT 20;

-- ============================================================================
-- 9. CHECK FOR ORPHANED DOCUMENTS (Documents without KYC Profile)
-- ============================================================================
SELECT 
    kd.user_id,
    kd.doc_type,
    kd.storage_key,
    kd.uploaded_at
FROM public.kyc_documents kd
LEFT JOIN public.kyc_profiles kp ON kd.user_id = kp.user_id
WHERE kp.user_id IS NULL;

-- ============================================================================
-- 10. CHECK FOR INCOMPLETE KYC SUBMISSIONS (Profile but Missing Documents)
-- ============================================================================
SELECT 
    kp.user_id,
    ui.full_name,
    ui.contact_email,
    kp.status,
    kp.submitted_at,
    COUNT(kd.id) AS document_count
FROM public.kyc_profiles kp
JOIN public.users_identity ui ON kp.user_id = ui.id
LEFT JOIN public.kyc_documents kd ON kp.user_id = kd.user_id
GROUP BY kp.user_id, ui.full_name, ui.contact_email, kp.status, kp.submitted_at
HAVING COUNT(kd.id) < 3  -- Should have 3 documents (ID front, ID back, selfie)
ORDER BY kp.submitted_at DESC;
