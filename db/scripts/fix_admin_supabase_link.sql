-- Script to link users_identity records with Supabase Auth users
-- This is needed for the custom claims hook to work

-- IMPORTANT: This script assumes you're using Supabase and have access to auth.users table
-- If you get permission errors, you may need to run this as a superuser or through Supabase Dashboard

-- Step 1: Check current state
SELECT 
    ui.user_id,
    ui.full_name,
    ui.auth_email,
    ui.supabase_user_id,
    ui.role,
    au.id as actual_supabase_id,
    au.email as supabase_email,
    CASE 
        WHEN ui.supabase_user_id IS NULL THEN '❌ NOT LINKED'
        WHEN ui.supabase_user_id = au.id THEN '✅ CORRECTLY LINKED'
        ELSE '⚠️ MISMATCH'
    END as link_status
FROM users_identity ui
LEFT JOIN auth.users au ON au.email = ui.auth_email
WHERE ui.auth_email = 'baylesinfo@gmail.com';

-- Step 2: Update supabase_user_id if it's NULL or incorrect
-- This links the users_identity record to the Supabase Auth user
UPDATE users_identity ui
SET supabase_user_id = au.id
FROM auth.users au
WHERE ui.auth_email = au.email
    AND ui.auth_email = 'baylesinfo@gmail.com'
    AND (ui.supabase_user_id IS NULL OR ui.supabase_user_id != au.id);

-- Step 3: Verify the update
SELECT 
    ui.user_id,
    ui.full_name,
    ui.auth_email,
    ui.supabase_user_id,
    ui.role,
    au.id as actual_supabase_id,
    CASE 
        WHEN ui.supabase_user_id = au.id THEN '✅ LINKED SUCCESSFULLY'
        ELSE '❌ STILL NOT LINKED'
    END as link_status
FROM users_identity ui
LEFT JOIN auth.users au ON au.email = ui.auth_email
WHERE ui.auth_email = 'baylesinfo@gmail.com';

-- Step 4: If the above doesn't work, manually set it
-- Uncomment and replace YOUR_SUPABASE_USER_ID with the actual ID from auth.users
/*
UPDATE users_identity
SET supabase_user_id = 'YOUR_SUPABASE_USER_ID'
WHERE auth_email = 'baylesinfo@gmail.com';
*/
