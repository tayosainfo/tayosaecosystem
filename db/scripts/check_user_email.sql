-- Check what emails exist in users_identity table
SELECT 
    user_id,
    full_name,
    auth_email,
    role,
    status,
    created_at
FROM users_identity
ORDER BY created_at DESC
LIMIT 20;

-- Check specifically for baylesinfo@gmail.com (exact match)
SELECT 
    user_id,
    full_name,
    auth_email,
    role,
    status
FROM users_identity
WHERE auth_email = 'baylesinfo@gmail.com';

-- Check for any similar emails (case-insensitive)
SELECT 
    user_id,
    full_name,
    auth_email,
    role,
    status
FROM users_identity
WHERE LOWER(auth_email) = LOWER('baylesinfo@gmail.com');

-- Count total users
SELECT COUNT(*) as total_users FROM users_identity;

-- Check if there are any admin users
SELECT 
    user_id,
    full_name,
    auth_email,
    role
FROM users_identity
WHERE role = 'admin';
