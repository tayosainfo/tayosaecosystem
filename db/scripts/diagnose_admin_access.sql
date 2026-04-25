-- Diagnostic script to troubleshoot admin access issues
-- Run this to verify the admin user configuration

-- 1. Check if user exists in users_identity with admin role
SELECT 
    user_id,
    full_name,
    auth_email,
    role,
    supabase_user_id,
    role_assigned_at,
    status
FROM users_identity
WHERE auth_email = 'baylesinfo@gmail.com';

-- 2. Check if custom claims hook function exists
SELECT 
    routine_name,
    routine_type,
    routine_definition
FROM information_schema.routines
WHERE routine_name = 'custom_access_token_hook'
    AND routine_schema = 'public';

-- 3. Check permissions on the custom claims hook
SELECT 
    grantee,
    privilege_type
FROM information_schema.routine_privileges
WHERE routine_name = 'custom_access_token_hook'
    AND routine_schema = 'public';

-- 4. Check if supabase_auth_admin has SELECT permission on users_identity
SELECT 
    grantee,
    privilege_type,
    table_name
FROM information_schema.table_privileges
WHERE table_name = 'users_identity'
    AND grantee = 'supabase_auth_admin';

-- 5. Test the custom claims hook function manually
-- Replace 'YOUR_SUPABASE_USER_ID' with the actual supabase_user_id from query #1
DO $$
DECLARE
    test_event jsonb;
    result jsonb;
    user_supabase_id text;
BEGIN
    -- Get the supabase_user_id for baylesinfo@gmail.com
    SELECT supabase_user_id INTO user_supabase_id
    FROM users_identity
    WHERE auth_email = 'baylesinfo@gmail.com';
    
    IF user_supabase_id IS NULL THEN
        RAISE NOTICE 'ERROR: No supabase_user_id found for baylesinfo@gmail.com';
        RAISE NOTICE 'This means the user record does not have a linked Supabase Auth account';
        RETURN;
    END IF;
    
    -- Create a test event
    test_event := jsonb_build_object(
        'user_id', user_supabase_id,
        'claims', jsonb_build_object(
            'sub', user_supabase_id,
            'email', 'baylesinfo@gmail.com'
        )
    );
    
    -- Call the custom claims hook
    result := public.custom_access_token_hook(test_event);
    
    -- Display results
    RAISE NOTICE '=== Custom Claims Hook Test Results ===';
    RAISE NOTICE 'Input user_id: %', user_supabase_id;
    RAISE NOTICE 'Output claims: %', result->'claims';
    RAISE NOTICE 'User role in claims: %', result->'claims'->'user_role';
    
    IF result->'claims'->'user_role' = '"admin"'::jsonb THEN
        RAISE NOTICE '✅ SUCCESS: Custom claims hook is working correctly';
    ELSE
        RAISE NOTICE '❌ FAILURE: Custom claims hook did not set admin role';
        RAISE NOTICE 'Expected: "admin", Got: %', result->'claims'->'user_role';
    END IF;
END $$;

-- 6. Check if there are multiple users with the same email
SELECT 
    COUNT(*) as user_count,
    auth_email
FROM users_identity
WHERE auth_email = 'baylesinfo@gmail.com'
GROUP BY auth_email;

-- 7. Verify the role enum type exists
SELECT 
    enumlabel as role_value
FROM pg_enum
WHERE enumtypid = 'user_role'::regtype
ORDER BY enumsortorder;
