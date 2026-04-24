-- Script: Assign admin roles to initial users
-- Description: Assigns admin role to specified users by email
-- Usage: Update the email list below, then execute this script
-- Date: 2026-04-24

-- ============================================================================
-- CONFIGURATION: Update this list with actual admin user emails
-- ============================================================================

-- Assign admin role by email
-- Replace the example emails below with actual admin user emails
UPDATE users_identity
SET 
  role = 'admin',
  role_assigned_at = CURRENT_TIMESTAMP,
  role_assigned_by = 'system_migration'
WHERE auth_email IN (
  -- Add admin emails here (one per line)
  -- Example:
  -- 'admin@tayosaecosystem.com',
  -- 'superadmin@tayosaecosystem.com'
  
  -- TODO: Replace with actual admin emails before running
  'REPLACE_WITH_ADMIN_EMAIL_1',
  'REPLACE_WITH_ADMIN_EMAIL_2'
);

-- ============================================================================
-- VERIFICATION: Check role assignments
-- ============================================================================

-- Display all users with admin role
SELECT 
  user_id,
  full_name,
  auth_email,
  role,
  role_assigned_at,
  role_assigned_by,
  created_at
FROM users_identity
WHERE role = 'admin'
ORDER BY role_assigned_at DESC;

-- Display count of admin users
SELECT 
  COUNT(*) as admin_count
FROM users_identity
WHERE role = 'admin';

-- Display audit trail of role assignments
SELECT 
  a.id,
  a.user_id,
  u.full_name,
  u.auth_email,
  a.action,
  a.previous_role,
  a.new_role,
  a.assigned_by,
  a.reason,
  a.created_at
FROM admin_role_audit a
JOIN users_identity u ON a.user_id = u.user_id
ORDER BY a.created_at DESC
LIMIT 20;

-- ============================================================================
-- NOTES
-- ============================================================================
-- 
-- To assign admin role to a specific user by email:
-- UPDATE users_identity
-- SET role = 'admin', role_assigned_at = CURRENT_TIMESTAMP, role_assigned_by = 'your_user_id'
-- WHERE auth_email = 'user@example.com';
--
-- To revoke admin role:
-- UPDATE users_identity
-- SET role = 'user', role_assigned_at = CURRENT_TIMESTAMP, role_assigned_by = 'your_user_id'
-- WHERE auth_email = 'user@example.com';
--
-- All role changes are automatically logged in admin_role_audit table
-- ============================================================================
