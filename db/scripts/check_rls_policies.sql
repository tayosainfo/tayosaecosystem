-- Check RLS policies on users_identity table
SELECT 
    schemaname,
    tablename,
    policyname,
    permissive,
    roles,
    qual,
    with_check
FROM pg_policies
WHERE tablename = 'users_identity'
ORDER BY policyname;

-- Check if RLS is enabled on users_identity
SELECT 
    schemaname,
    tablename,
    rowsecurity
FROM pg_tables
WHERE tablename = 'users_identity';

-- Check table grants
SELECT 
    grantee,
    privilege_type
FROM information_schema.table_privileges
WHERE table_name = 'users_identity'
ORDER BY grantee;
