# Remove Admin System - Make KYC Auto-Approved

## Overview
This document outlines the changes needed to remove the admin system and make KYC verification automatic upon submission.

## Changes Required

### 1. Database Migrations

#### Migration 021: Auto-approve KYC (CREATED ✅)
- File: `db/migrations/021_auto_approve_kyc.sql`
- Creates trigger to auto-approve KYC submissions
- Updates existing pending KYC to approved

#### Migration 022: Remove Admin System (CREATED ✅)
- File: `db/migrations/022_remove_admin_system.sql`
- Removes admin role columns from users_identity
- Drops admin_role_audit table
- Drops custom_access_token_hook function
- Drops user_role enum type

### 2. Backend Changes

#### A. Remove Admin Endpoints (user-service/main.go)
Remove these lines:
```go
mux.HandleFunc("GET /api/v1/admin/kyc", adminKYCDecisionHandler)
mux.HandleFunc("POST /api/v1/admin/kyc", adminKYCDecisionHandler)
mux.HandleFunc("GET /api/v1/admin/settings", adminSettingsHandler)
mux.HandleFunc("PATCH /api/v1/admin/settings", adminSettingsHandler)
mux.HandleFunc("POST /api/v1/admin/check-status", adminCheckStatusHandler)
mux.HandleFunc("GET /api/v1/admin/users", requireAdminAuth(adminUsersListHandler))
mux.HandleFunc("GET /api/v1/admin/users/{userId}", requireAdminAuth(adminUserDetailHandler))
mux.HandleFunc("PATCH /api/v1/admin/users/{userId}/status", requireAdminAuth(adminUserStatusHandler))
mux.HandleFunc("PATCH /api/v1/admin/users/{userId}/role", requireAdminAuth(adminUserRoleHandler))
mux.HandleFunc("POST /api/v1/admin/users/{userId}/reset-password", requireAdminAuth(adminUserResetPasswordHandler))
mux.HandleFunc("GET /api/v1/admin/users/{userId}/activity", requireAdminAuth(adminUserActivityHandler))
```

#### B. Update KYC Handler (user-service/handlers.go)
Modify `onboardingKYCHandler` to auto-approve:
- Change status from "pending" to "approved"
- Set reviewed_at to current timestamp
- Set reviewed_by to "system_auto_approval"

#### C. Remove Admin Handler Functions
Delete these files:
- `services/user-service/admin_handlers.go`
- `services/user-service/auth.go` (admin auth middleware)

#### D. Update Store Interface (user-service/store.go)
Remove these methods:
- `SetKYCDecision`
- `ListAdminKYCQueue`
- `GetAdminSetting`
- `SetAdminSetting`
- `ListUsersWithFilters`
- `UpdateUserStatus`
- `UpdateUserRole`
- `GetUserActivity`

#### E. Update Store Implementations
- `services/user-service/postgres_store.go` - Remove admin methods
- `services/user-service/memory_store.go` - Remove admin methods

#### F. Update API Gateway (api-gateway-service/main.go)
Remove admin routes:
```go
{path: "/api/v1/admin/kyc", base: userBase, public: false},
{path: "/api/v1/admin/settings", base: userBase, public: false},
{path: "/api/v1/admin/users", base: userBase, public: false},
// ... all other admin routes
```

Remove admin middleware:
- Delete `requireAdmin` function
- Delete `requireAdminWithFallback` function
- Delete `extractRoleFromJWT` function

### 3. Frontend Changes

#### A. Remove Admin Pages
Delete these files:
- `src/pages/Admin.tsx`
- `src/pages/admin/Users.tsx`
- `src/pages/admin/UserDetail.tsx`
- Any other admin-related pages

#### B. Remove Admin Routes (src/App.tsx or router config)
Remove routes like:
- `/admin`
- `/admin/users`
- `/admin/users/:id`

#### C. Update Home Page (src/pages/Home.tsx)
- Remove KYC status checks for "pending" state
- Update UI to show KYC as automatically approved after submission
- Remove "under review" messaging

#### D. Remove Admin API Functions (src/lib/platformApi.ts)
Remove:
- `adminListKyc`
- Any other admin-related API functions

#### E. Remove Admin Utilities (src/utils/api.ts)
Remove:
- `makeAdminRequest` function

### 4. Environment Variables to Remove
- `ADMIN_API_KEY` - No longer needed
- `AUTH_MIGRATION_MODE` - No longer needed

### 5. Documentation Updates
- Update README.md to remove admin setup instructions
- Update API documentation to remove admin endpoints

## Implementation Steps

### Step 1: Run Database Migrations
```sql
-- In Supabase SQL Editor, run in order:
-- 1. db/migrations/021_auto_approve_kyc.sql
-- 2. db/migrations/022_remove_admin_system.sql
```

### Step 2: Update Backend Code
1. Remove admin endpoints from main.go
2. Update onboardingKYCHandler to auto-approve
3. Delete admin_handlers.go
4. Delete auth.go (admin middleware)
5. Remove admin methods from store interface and implementations
6. Update API gateway to remove admin routes

### Step 3: Update Frontend Code
1. Delete admin pages
2. Remove admin routes
3. Update Home.tsx to remove pending state
4. Remove admin API functions

### Step 4: Deploy
1. Commit all changes
2. Push to GitHub
3. Deploy backend services (Render)
4. Deploy frontend (Vercel)

### Step 5: Test
1. Register a new user
2. Submit KYC documents
3. Verify KYC is automatically approved
4. Verify SACCO enrollment works immediately after KYC

## Rollback Plan
If needed, you can restore admin functionality by:
1. Reverting the database migrations
2. Restoring the deleted code files from git history
3. Re-running migrations 013-020

## Benefits
- ✅ Simpler system architecture
- ✅ Faster user onboarding (no waiting for admin approval)
- ✅ Reduced maintenance burden
- ✅ No admin authentication complexity
- ✅ Automatic verification upon document submission
