-- Set admin role for baylesinfo@gmail.com
-- This script updates the user's role to 'admin' in the database

-- First, find the user ID for baylesinfo@gmail.com
SELECT user_id, full_name, contact_email, auth_email, role, status
FROM public.users_identity
WHERE LOWER(contact_email) = LOWER('baylesinfo@gmail.com')
   OR LOWER(auth_email) = LOWER('baylesinfo@gmail.com');

-- Update the user's role to admin
UPDATE public.users_identity
SET 
    role = 'admin',
    role_assigned_at = NOW(),
    role_assigned_by = 'system',
    updated_at = NOW()
WHERE LOWER(contact_email) = LOWER('baylesinfo@gmail.com')
   OR LOWER(auth_email) = LOWER('baylesinfo@gmail.com');

-- Verify the update
SELECT user_id, full_name, contact_email, auth_email, role, status, role_assigned_at, role_assigned_by
FROM public.users_identity
WHERE LOWER(contact_email) = LOWER('baylesinfo@gmail.com')
   OR LOWER(auth_email) = LOWER('baylesinfo@gmail.com');
