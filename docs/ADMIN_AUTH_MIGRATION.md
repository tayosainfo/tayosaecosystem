# Admin Authentication Migration Guide

## Overview

This document describes the migration from insecure shared API key authentication to JWT-based role authentication using Supabase Auth for admin endpoints.

## Migration Date

**Completed:** April 24, 2026

## Changes Made

### 1. Database Changes

#### Migration 013: User Roles
- **File:** `db/migrations/013_add_user_roles.sql`
- **Changes:**
  - Created `user_role` enum type ('user', 'admin')
  - Added `role` column to `users_identity` table (default: 'user')
  - Added `role_assigned_at` timestamp column
  - Added `role_assigned_by` column (tracks admin who assigned role)
  - Added `status` column ('active', 'suspended', 'deactivated')
  - Added `last_login` timestamp column
  - Created `admin_role_audit` table for tracking role changes
  - Created `log_role_change()` trigger function for automatic audit logging
  - Added indexes for efficient role-based queries

#### Migration 014: Custom Claims Hook
- **File:** `db/migrations/014_configure_custom_claims.sql`
- **Changes:**
  - Created `custom_access_token_hook()` function
  - Configured function to add `user_role` to JWT claims
  - Granted necessary permissions to `supabase_auth_admin` role
  - Defaults to 'user' role if no role found

### 2. Backend Changes

#### API Gateway Service
- **File:** `services/api-gateway-service/auth.go`
- **New Functions:**
  - `extractRoleFromJWT()` - Extracts user role from JWT by calling Supabase /auth/v1/user endpoint
  - `requireAdmin()` - Middleware that validates admin role from JWT
  - `requireAdminWithFallback()` - Migration middleware supporting both JWT and API key (temporary)
  - `isMigrationMode()` - Checks AUTH_MIGRATION_MODE environment variable

- **File:** `services/api-gateway-service/main.go`
- **Changes:**
  - Applied `requireAdminWithFallback` middleware to all `/api/v1/admin/*` routes

#### User Service
- **File:** `services/user-service/admin_handlers.go`
- **New Endpoints:**
  - `GET /api/v1/admin/users` - List users with pagination, search, and filters
  - `GET /api/v1/admin/users/{userId}` - Get user details
  - `PATCH /api/v1/admin/users/{userId}/status` - Update user status
  - `PATCH /api/v1/admin/users/{userId}/role` - Update user role
  - `POST /api/v1/admin/users/{userId}/reset-password` - Trigger password reset
  - `GET /api/v1/admin/users/{userId}/activity` - Get user activity log

- **File:** `services/user-service/postgres_store.go` & `memory_store.go`
- **New Methods:**
  - `ListUsersWithFilters()` - Query users with search and filters
  - `UpdateUserStatus()` - Update user account status
  - `UpdateUserRole()` - Update user role with audit logging
  - `GetUserActivity()` - Retrieve user activity logs

- **File:** `services/user-service/types.go`
- **Updated Types:**
  - Added `Role`, `Status`, `RoleAssignedAt`, `RoleAssignedBy`, `LastLogin` fields to `User` struct
  - Added `ActivityLog` type for activity tracking

### 3. Frontend Changes

#### Authentication Utilities
- **File:** `src/utils/auth.ts`
- **New Functions:**
  - `checkAdminStatus()` - Extracts role from JWT app_metadata
  - `useAdminStatus()` - React hook for checking admin status

- **File:** `src/utils/api.ts`
- **New Functions:**
  - `makeAdminRequest()` - Helper for making authenticated admin API requests with JWT

#### Admin Pages
- **File:** `src/pages/Admin.tsx`
- **Changes:**
  - Removed `VITE_ADMIN_API_KEY` usage
  - Removed `X-Admin-Secret` header
  - Implemented JWT-based authentication using `makeAdminRequest()`
  - Added admin status check using `useAdminStatus()`
  - Added access control (shows "Access Denied" for non-admin users)

- **File:** `src/pages/Home.tsx`
- **Changes:**
  - Removed `ADMIN_API_KEY` variable
  - Updated admin API calls to use `makeAdminRequest()`

- **File:** `src/pages/admin/Users.tsx` (NEW)
- **Features:**
  - Paginated user list (20 users per page)
  - Search by name, email, phone
  - Filter by status and KYC status
  - View user details
  - Uses JWT authentication

- **File:** `src/pages/admin/UserDetail.tsx` (NEW)
- **Features:**
  - Comprehensive user information display
  - Status change modal with reason requirement
  - Role change modal with admin confirmation
  - Password reset functionality
  - Activity log timeline
  - Uses JWT authentication

#### Environment Configuration
- **Files:** `.env`, `.env.example`
- **Removed:**
  - `VITE_ADMIN_API_KEY` (insecure shared API key)

## Admin Role Assignment Process

### Initial Setup
1. Run database migrations:
   ```bash
   psql $DATABASE_URL -f db/migrations/013_add_user_roles.sql
   psql $DATABASE_URL -f db/migrations/014_configure_custom_claims.sql
   ```

2. Configure Supabase custom claims hook:
   - Navigate to Supabase Dashboard → Authentication → Hooks
   - Enable "Custom Access Token" hook
   - Set hook function to: `public.custom_access_token_hook`
   - Save configuration

3. Assign initial admin roles:
   ```sql
   -- Update with actual admin user IDs
   UPDATE users_identity 
   SET role = 'admin', 
       role_assigned_at = NOW(), 
       role_assigned_by = 'system'
   WHERE user_id IN ('user-id-1', 'user-id-2');
   ```

### Ongoing Role Management
- Use the admin dashboard at `/admin/users/{userId}`
- Click "Change Role" button
- Select new role (user/admin)
- Confirm admin role assignment (requires confirmation)
- Changes are automatically logged in `admin_role_audit` table

## Verification Steps

### 1. Verify Database Schema
```sql
-- Check role column exists
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users_identity' 
  AND column_name IN ('role', 'status', 'role_assigned_at', 'role_assigned_by', 'last_login');

-- Check audit table exists
SELECT table_name 
FROM information_schema.tables 
WHERE table_name = 'admin_role_audit';

-- Check custom claims hook exists
SELECT routine_name 
FROM information_schema.routines 
WHERE routine_name = 'custom_access_token_hook';
```

### 2. Verify JWT Claims
- Login as admin user
- Decode JWT token (use jwt.io)
- Verify `app_metadata.user_role` is present and set to 'admin'

### 3. Verify Admin Endpoints
```bash
# Get JWT token from Supabase Auth
TOKEN="your-jwt-token"

# Test admin endpoint with JWT
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/admin/users

# Should return 200 for admin users
# Should return 403 for non-admin users
# Should return 401 for invalid/missing tokens
```

### 4. Verify Frontend Access Control
- Login as admin user → Should see admin dashboard
- Login as regular user → Should see "Access Denied" message
- No authentication → Should redirect to login

### 5. Verify Audit Logging
```sql
-- Check role changes are logged
SELECT * FROM admin_role_audit 
ORDER BY created_at DESC 
LIMIT 10;
```

## Rollback Procedures

### Emergency Rollback (Restore API Key Auth)
If critical issues arise, you can temporarily restore API key authentication:

1. Set environment variable:
   ```bash
   AUTH_MIGRATION_MODE=true
   ADMIN_API_KEY=your-secure-api-key
   ```

2. Restart services
3. Both JWT and API key authentication will work
4. Investigate and fix issues
5. Remove API key when ready

### Full Rollback (Remove Role-Based Auth)
⚠️ **Warning:** This will remove all role-based access control

```sql
-- Remove role columns
ALTER TABLE users_identity 
  DROP COLUMN IF EXISTS role,
  DROP COLUMN IF EXISTS status,
  DROP COLUMN IF EXISTS role_assigned_at,
  DROP COLUMN IF EXISTS role_assigned_by,
  DROP COLUMN IF EXISTS last_login;

-- Drop audit table
DROP TABLE IF EXISTS admin_role_audit;

-- Drop custom claims hook
DROP FUNCTION IF EXISTS public.custom_access_token_hook;
DROP FUNCTION IF EXISTS log_role_change;

-- Drop enum type
DROP TYPE IF EXISTS user_role;
```

## Security Considerations

### Before Migration
- ❌ Shared API key exposed in frontend code
- ❌ API key could be extracted from browser
- ❌ No per-user access control
- ❌ No audit trail for admin actions

### After Migration
- ✅ JWT-based authentication (secure)
- ✅ Per-user role verification
- ✅ Automatic audit logging
- ✅ Token expiration and refresh
- ✅ No secrets in frontend code
- ✅ Supabase Auth integration

## Performance Impact

- **JWT Validation:** ~10-50ms per request (calls Supabase API)
- **Database Queries:** Minimal impact (indexed role column)
- **Frontend:** No noticeable impact
- **Caching:** JWT role is cached in token (no DB query per request)

## Monitoring

### Key Metrics to Monitor
1. **Authentication failures:** Track 401/403 responses
2. **Role changes:** Monitor `admin_role_audit` table
3. **API response times:** Ensure JWT validation doesn't slow requests
4. **User complaints:** Watch for access issues

### Logs to Review
```bash
# Check authentication method usage
grep "AUTH_METHOD" logs/api-gateway.log

# Check for JWT validation errors
grep "Failed to extract role" logs/api-gateway.log

# Check for permission denials
grep "insufficient permissions" logs/api-gateway.log
```

## Support

### Common Issues

**Issue:** Admin user gets 403 Forbidden
- **Cause:** Role not set in database or custom claims hook not configured
- **Fix:** Verify role in database, check Supabase hook configuration

**Issue:** JWT validation fails
- **Cause:** Supabase credentials incorrect or network issue
- **Fix:** Verify SUPABASE_URL and SUPABASE_ANON_KEY in environment

**Issue:** Custom claims not in JWT
- **Cause:** Hook not enabled in Supabase Dashboard
- **Fix:** Enable hook in Authentication → Hooks

### Contact
For issues or questions, contact the development team.

## References

- [Supabase Custom Claims Documentation](https://supabase.com/docs/guides/auth/auth-hooks)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- Design Document: `.kiro/specs/role-based-admin-auth/design.md`
- Requirements Document: `.kiro/specs/role-based-admin-auth/requirements.md`
