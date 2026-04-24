# Admin Role Assignment Guide

## Overview

This guide explains how to assign and manage admin roles in the system. Admin roles grant users full access to administrative features including user management, KYC review, and system configuration.

## Prerequisites

- Database access (PostgreSQL)
- Supabase Dashboard access (for hook configuration)
- Admin credentials (for UI-based role management)

## Initial Admin Role Assignment

### Method 1: Direct Database Update (Recommended for Initial Setup)

Use this method to assign the first admin users during system setup.

#### Step 1: Identify User IDs

```sql
-- Find users by email
SELECT user_id, full_name, auth_email, role 
FROM users_identity 
WHERE auth_email IN ('admin1@example.com', 'admin2@example.com');

-- Find users by phone
SELECT user_id, full_name, phone_e164, role 
FROM users_identity 
WHERE phone_e164 IN ('+256700000001', '+256700000002');
```

#### Step 2: Assign Admin Role

```sql
-- Assign admin role to specific users
UPDATE users_identity 
SET 
  role = 'admin',
  role_assigned_at = NOW(),
  role_assigned_by = 'system',  -- or use admin user_id
  status = 'active',
  updated_at = NOW()
WHERE user_id IN (
  'user-id-1',
  'user-id-2',
  'user-id-3'
);
```

#### Step 3: Verify Assignment

```sql
-- Verify admin roles were assigned
SELECT 
  user_id,
  full_name,
  auth_email,
  role,
  status,
  role_assigned_at,
  role_assigned_by
FROM users_identity 
WHERE role = 'admin'
ORDER BY role_assigned_at DESC;
```

#### Step 4: Check Audit Log

```sql
-- Verify audit trail was created
SELECT 
  user_id,
  action,
  previous_role,
  new_role,
  assigned_by,
  created_at
FROM admin_role_audit 
ORDER BY created_at DESC 
LIMIT 10;
```

### Method 2: Using Admin Dashboard (Recommended for Ongoing Management)

Use this method after initial setup to manage roles through the UI.

#### Step 1: Login as Admin

1. Navigate to the application
2. Login with admin credentials
3. Access the admin dashboard at `/admin/users`

#### Step 2: Find User

1. Use the search bar to find user by name, email, or phone
2. Or use filters to narrow down the user list
3. Click on the user to view their details

#### Step 3: Assign Admin Role

1. Click the "Change Role" button
2. Select "Admin" from the dropdown
3. Review the warning message about admin privileges
4. Click "Continue"
5. Confirm the admin role assignment in the confirmation dialog
6. Click "Yes, Grant Admin Access"

#### Step 4: Verify Assignment

1. User's role badge should update to "admin"
2. Check the activity log for the role change event
3. User should now have access to admin features

## Revoking Admin Role

### Method 1: Using Admin Dashboard

1. Navigate to `/admin/users/{userId}`
2. Click "Change Role" button
3. Select "User" from the dropdown
4. Click "Update Role"
5. Verify the role change in the activity log

### Method 2: Direct Database Update

```sql
-- Revoke admin role
UPDATE users_identity 
SET 
  role = 'user',
  role_assigned_at = NOW(),
  role_assigned_by = 'admin-user-id',  -- ID of admin performing action
  updated_at = NOW()
WHERE user_id = 'user-id-to-revoke';
```

## Verification Queries

### List All Admins

```sql
SELECT 
  user_id,
  full_name,
  auth_email,
  phone_e164,
  role,
  status,
  role_assigned_at,
  role_assigned_by,
  last_login
FROM users_identity 
WHERE role = 'admin'
ORDER BY role_assigned_at DESC;
```

### Count Admins

```sql
SELECT 
  COUNT(*) as total_admins,
  COUNT(CASE WHEN status = 'active' THEN 1 END) as active_admins,
  COUNT(CASE WHEN status = 'suspended' THEN 1 END) as suspended_admins
FROM users_identity 
WHERE role = 'admin';
```

### Recent Role Changes

```sql
SELECT 
  a.user_id,
  u.full_name,
  u.auth_email,
  a.action,
  a.previous_role,
  a.new_role,
  a.assigned_by,
  a.created_at
FROM admin_role_audit a
JOIN users_identity u ON a.user_id = u.user_id
ORDER BY a.created_at DESC
LIMIT 20;
```

### Admin Activity Summary

```sql
SELECT 
  assigned_by,
  COUNT(*) as role_changes,
  COUNT(CASE WHEN action = 'granted' THEN 1 END) as grants,
  COUNT(CASE WHEN action = 'revoked' THEN 1 END) as revocations,
  MIN(created_at) as first_action,
  MAX(created_at) as last_action
FROM admin_role_audit
WHERE assigned_by IS NOT NULL
GROUP BY assigned_by
ORDER BY role_changes DESC;
```

## Audit Trail Review

### View All Role Changes for a User

```sql
SELECT 
  action,
  previous_role,
  new_role,
  assigned_by,
  reason,
  created_at
FROM admin_role_audit
WHERE user_id = 'specific-user-id'
ORDER BY created_at DESC;
```

### View All Actions by an Admin

```sql
SELECT 
  a.user_id,
  u.full_name as affected_user,
  a.action,
  a.previous_role,
  a.new_role,
  a.created_at
FROM admin_role_audit a
JOIN users_identity u ON a.user_id = u.user_id
WHERE a.assigned_by = 'admin-user-id'
ORDER BY a.created_at DESC;
```

### Export Audit Log

```sql
-- Export to CSV (PostgreSQL)
COPY (
  SELECT 
    a.user_id,
    u.full_name,
    u.auth_email,
    a.action,
    a.previous_role,
    a.new_role,
    a.assigned_by,
    a.created_at
  FROM admin_role_audit a
  JOIN users_identity u ON a.user_id = u.user_id
  WHERE a.created_at >= NOW() - INTERVAL '30 days'
  ORDER BY a.created_at DESC
) TO '/tmp/admin_role_audit.csv' WITH CSV HEADER;
```

## Best Practices

### Security

1. **Principle of Least Privilege**
   - Only assign admin role to users who absolutely need it
   - Regularly review admin user list
   - Revoke admin access when no longer needed

2. **Audit Trail**
   - Always provide a reason when changing roles via UI
   - Regularly review audit logs for suspicious activity
   - Keep audit logs for compliance (minimum 90 days)

3. **Access Control**
   - Use strong passwords for admin accounts
   - Enable 2FA for admin users (if available)
   - Monitor admin login activity

### Operational

1. **Documentation**
   - Document why each user was granted admin access
   - Keep a separate record of admin users
   - Update documentation when roles change

2. **Communication**
   - Notify users when they are granted admin access
   - Provide training on admin responsibilities
   - Inform users when admin access is revoked

3. **Monitoring**
   - Set up alerts for new admin role assignments
   - Monitor admin activity logs
   - Review admin list monthly

## Troubleshooting

### Issue: User has admin role but gets 403 Forbidden

**Possible Causes:**
1. Custom claims hook not configured in Supabase
2. JWT token doesn't include user_role claim
3. User needs to logout and login again to get new token

**Solution:**
```sql
-- Verify role in database
SELECT user_id, role FROM users_identity WHERE user_id = 'user-id';

-- Check if custom claims hook is configured
SELECT routine_name 
FROM information_schema.routines 
WHERE routine_name = 'custom_access_token_hook';
```

Then:
1. Verify Supabase hook is enabled (Dashboard → Authentication → Hooks)
2. Ask user to logout and login again
3. Decode JWT token to verify user_role claim is present

### Issue: Role change not reflected in JWT

**Cause:** JWT tokens are cached and don't update immediately

**Solution:**
1. User must logout and login again to get new token with updated role
2. Or wait for token to expire (typically 1 hour)
3. Or implement token refresh mechanism

### Issue: Audit log not recording changes

**Possible Causes:**
1. Trigger not created or disabled
2. Database permissions issue

**Solution:**
```sql
-- Check if trigger exists
SELECT trigger_name, event_manipulation, event_object_table
FROM information_schema.triggers
WHERE trigger_name = 'trigger_log_role_change';

-- Recreate trigger if missing
CREATE TRIGGER trigger_log_role_change
AFTER UPDATE ON users_identity
FOR EACH ROW
EXECUTE FUNCTION log_role_change();
```

## Emergency Procedures

### Revoke All Admin Access (Emergency)

⚠️ **Warning:** This will revoke admin access for ALL users. Use only in emergency.

```sql
-- Backup current admins first
CREATE TEMP TABLE admin_backup AS
SELECT * FROM users_identity WHERE role = 'admin';

-- Revoke all admin roles
UPDATE users_identity 
SET 
  role = 'user',
  role_assigned_at = NOW(),
  role_assigned_by = 'emergency-revocation',
  updated_at = NOW()
WHERE role = 'admin';

-- Verify
SELECT COUNT(*) FROM users_identity WHERE role = 'admin';
-- Should return 0

-- To restore from backup
UPDATE users_identity u
SET 
  role = b.role,
  role_assigned_at = b.role_assigned_at,
  role_assigned_by = b.role_assigned_by
FROM admin_backup b
WHERE u.user_id = b.user_id;
```

### Grant Emergency Admin Access

If all admins are locked out:

```sql
-- Grant admin to specific user immediately
UPDATE users_identity 
SET 
  role = 'admin',
  status = 'active',
  role_assigned_at = NOW(),
  role_assigned_by = 'emergency-grant',
  updated_at = NOW()
WHERE auth_email = 'emergency-admin@example.com';
```

## Compliance and Reporting

### Monthly Admin Review Report

```sql
-- Generate monthly admin review report
SELECT 
  u.user_id,
  u.full_name,
  u.auth_email,
  u.role,
  u.status,
  u.role_assigned_at,
  u.role_assigned_by,
  u.last_login,
  COUNT(a.id) as role_changes_made
FROM users_identity u
LEFT JOIN admin_role_audit a ON u.user_id = a.assigned_by
WHERE u.role = 'admin'
GROUP BY u.user_id, u.full_name, u.auth_email, u.role, u.status, 
         u.role_assigned_at, u.role_assigned_by, u.last_login
ORDER BY u.role_assigned_at DESC;
```

### Audit Compliance Export

```sql
-- Export audit trail for compliance (last 90 days)
SELECT 
  a.id,
  a.user_id,
  u.full_name as affected_user,
  u.auth_email as affected_email,
  a.action,
  a.previous_role,
  a.new_role,
  a.assigned_by,
  admin.full_name as admin_name,
  admin.auth_email as admin_email,
  a.reason,
  a.created_at
FROM admin_role_audit a
JOIN users_identity u ON a.user_id = u.user_id
LEFT JOIN users_identity admin ON a.assigned_by = admin.user_id
WHERE a.created_at >= NOW() - INTERVAL '90 days'
ORDER BY a.created_at DESC;
```

## References

- Migration Guide: `docs/ADMIN_AUTH_MIGRATION.md`
- Design Document: `.kiro/specs/role-based-admin-auth/design.md`
- Requirements Document: `.kiro/specs/role-based-admin-auth/requirements.md`
- Database Schema: `db/migrations/013_add_user_roles.sql`
