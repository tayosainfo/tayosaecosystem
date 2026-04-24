# Implementation Plan: Role-Based Admin Authentication

## Overview

This implementation plan converts the role-based admin authentication design into actionable coding tasks. The feature replaces the insecure shared API key approach with JWT-based role verification using Supabase Auth. Implementation follows a phased migration strategy to ensure zero downtime.

## Tasks

- [x] 1. Database schema setup and configuration
  - [x] 1.1 Create database migration for user roles
    - Create `db/migrations/012_add_user_roles.sql` with user_role enum type
    - Add role, role_assigned_at, role_assigned_by columns to users_identity table
    - Create admin_role_audit table for tracking role changes
    - Add indexes for efficient role-based queries
    - Create log_role_change() trigger function for automatic audit logging
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 6.4, 6.5_

  - [x] 1.2 Create Supabase custom claims hook
    - Create `db/migrations/013_configure_custom_claims.sql`
    - Implement custom_access_token_hook() function to add user_role to JWT claims
    - Grant necessary permissions to supabase_auth_admin role
    - _Requirements: 2.1, 2.2, 2.3_

  - [x] 1.3 Create admin role assignment script
    - Create `db/scripts/assign_admin_roles.sql` for initial admin assignments
    - Include verification query to confirm role assignments
    - _Requirements: 6.1, 6.3, 6.4_

  - [ ]* 1.4 Create database schema verification script
    - Create `db/migrations/verify_schema.sql` to validate schema changes
    - Check role column exists with correct type
    - Verify audit table structure
    - _Requirements: 1.1, 1.2_

- [ ] 2. Backend authentication middleware implementation
  - [x] 2.1 Implement JWT role extraction function
    - Create `services/api-gateway-service/auth.go`
    - Implement extractRoleFromJWT() function to call Supabase /auth/v1/user endpoint
    - Extract user_role from app_metadata in JWT response
    - Default to 'user' role if role claim missing
    - Handle errors and invalid tokens appropriately
    - _Requirements: 2.2, 2.4, 3.2_

  - [x] 2.2 Implement requireAdmin middleware
    - Add requireAdmin() middleware function in `services/api-gateway-service/auth.go`
    - Validate Authorization header contains Bearer token
    - Call extractRoleFromJWT() to get user role
    - Return 401 Unauthorized for missing/invalid tokens
    - Return 403 Forbidden for non-admin roles
    - Allow request to proceed for admin role
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

  - [x] 2.3 Implement migration mode dual authentication
    - Add isMigrationMode() function checking AUTH_MIGRATION_MODE env var
    - Implement requireAdminWithFallback() middleware
    - Try JWT authentication first, fallback to API key if in migration mode
    - Log which authentication method was used
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [x] 2.4 Apply middleware to admin routes
    - Update route registration in `services/api-gateway-service/main.go`
    - Apply requireAdminWithFallback to all /api/v1/admin/* routes
    - Ensure all admin endpoints are protected
    - _Requirements: 3.6_

  - [ ]* 2.5 Write unit tests for JWT extraction
    - Create `services/api-gateway-service/auth_test.go`
    - Test extractRoleFromJWT with admin token, user token, missing role, invalid token
    - Mock Supabase API responses
    - _Requirements: 16.1_

  - [ ]* 2.6 Write unit tests for requireAdmin middleware
    - Create `services/api-gateway-service/middleware_test.go`
    - Test admin token succeeds (200), user token forbidden (403)
    - Test missing token unauthorized (401), invalid token unauthorized (401)
    - _Requirements: 16.1, 16.6_

- [x] 3. Backend user management endpoints
  - [x] 3.1 Implement user list endpoint
    - Create GET /api/v1/admin/users endpoint in user-service
    - Support pagination with page and limit query parameters
    - Support search by name, email, phone number
    - Support filtering by status (active, suspended, deactivated)
    - Support filtering by KYC status (pending, approved, rejected)
    - Return user list with pagination metadata
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

  - [x] 3.2 Implement user detail endpoint
    - Create GET /api/v1/admin/users/:userId endpoint
    - Return comprehensive user information (profile, KYC, status, role)
    - Include onboarding progress and completion status
    - Include registration date and last login timestamp
    - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5_

  - [x] 3.3 Implement user status management endpoint
    - Create PATCH /api/v1/admin/users/:userId/status endpoint
    - Accept status (active, suspended, deactivated) and reason in request body
    - Update user status in database
    - Log status change with admin ID, timestamp, and reason
    - Invalidate user sessions when suspended
    - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5_

  - [x] 3.4 Implement role assignment endpoint
    - Create PATCH /api/v1/admin/users/:userId/role endpoint
    - Accept role (user, admin) in request body
    - Require admin authentication (already protected by middleware)
    - Update user role in database
    - Trigger automatic audit logging via database trigger
    - _Requirements: 13.1, 13.2, 13.3, 13.4, 13.5_

  - [x] 3.5 Implement credential management endpoints
    - Create POST /api/v1/admin/users/:userId/reset-password endpoint
    - Trigger Supabase password reset email
    - Create POST /api/v1/admin/users/:userId/unlock endpoint
    - Reset failed login counters and unlock account
    - Log all credential management actions
    - _Requirements: 14.1, 14.2, 14.3, 14.4, 14.5_

  - [x] 3.6 Implement user activity endpoint
    - Create GET /api/v1/admin/users/:userId/activity endpoint
    - Return activity timeline (logins, transactions, status changes)
    - Include login history with timestamps, IP addresses, device info
    - Include admin actions performed on the account
    - Return at least 90 days of history
    - _Requirements: 15.1, 15.2, 15.3, 15.4, 15.5_

  - [ ]* 3.7 Write integration tests for user management endpoints
    - Create `services/api-gateway-service/integration_test.go`
    - Test admin can access all user management endpoints
    - Test non-admin receives 403 for user management endpoints
    - Test unauthenticated requests receive 401
    - _Requirements: 16.2, 16.3_

- [x] 4. Checkpoint - Backend implementation complete
  - Ensure all tests pass, ask the user if questions arise.

- [x] 5. Frontend authentication utilities
  - [x] 5.1 Create admin status check utility
    - Create `src/utils/auth.ts` file
    - Implement checkAdminStatus() function to extract role from JWT app_metadata
    - Return UserRole object with isAdmin boolean and role string
    - Default to non-admin if user not found or error occurs
    - _Requirements: 8.1, 8.2, 8.3_

  - [x] 5.2 Create useAdminStatus React hook
    - Implement useAdminStatus() hook in `src/utils/auth.ts`
    - Call checkAdminStatus() on mount and cache result
    - Listen for auth state changes and update admin status
    - Return isAdmin, role, and loading state
    - _Requirements: 8.1, 8.4, 8.5_

  - [x] 5.3 Create admin API request helper
    - Create `src/utils/api.ts` file
    - Implement makeAdminRequest() function
    - Get current session from Supabase
    - Include JWT token in Authorization header
    - Handle 401 by refreshing token and retrying
    - Handle 403 by throwing "Insufficient permissions" error
    - Handle 401 after retry by throwing "Authentication failed" error
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [ ] 6. Frontend admin page updates
  - [x] 6.1 Update Admin page to use JWT authentication
    - Update `src/pages/Admin.tsx`
    - Remove ADMIN_API_KEY and X-Admin-Secret header usage
    - Use makeAdminRequest() helper for all API calls
    - Use useAdminStatus() hook to check admin status
    - Show loading state while checking admin status
    - Show "Access Denied" message for non-admin users
    - _Requirements: 4.1, 4.2, 4.5, 5.1, 8.1, 8.5_

  - [x] 6.2 Update Home page to remove API key
    - Update `src/pages/Home.tsx`
    - Remove ADMIN_API_KEY variable and usage
    - Use makeAdminRequest() for admin API calls if any
    - _Requirements: 4.1, 4.2_

  - [x] 6.3 Remove API key from environment files
    - Remove VITE_ADMIN_API_KEY from `.env`
    - Remove VITE_ADMIN_API_KEY from `.env.example`
    - _Requirements: 4.3, 4.4_

- [ ] 7. Frontend user management dashboard
  - [x] 7.1 Create user list component
    - Create `src/pages/admin/Users.tsx`
    - Display paginated user list with 20 users per page
    - Implement search form with name, email, phone search
    - Add status filter dropdown (all, active, suspended, deactivated)
    - Add KYC status filter dropdown (all, pending, approved, rejected)
    - Display user table with columns: User, Contact, Role, Status, KYC, Joined, Actions
    - Add pagination controls (Previous, Page X of Y, Next)
    - Use makeAdminRequest() to fetch user data
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

  - [x] 7.2 Create user detail component
    - Create `src/pages/admin/UserDetail.tsx`
    - Display comprehensive user information header
    - Show user profile fields (name, email, phone, role, status, KYC status)
    - Show account metadata (joined date, last login)
    - Add action buttons (Change Status, Change Role, Reset Password)
    - Display activity log timeline with action, details, timestamp, IP
    - Use makeAdminRequest() to fetch user details and activity
    - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5, 15.1, 15.2, 15.3, 15.4, 15.5_

  - [x] 7.3 Implement status change modal
    - Add status change modal to UserDetail component
    - Provide dropdown to select new status (active, suspended, deactivated)
    - Require reason/note text input
    - Call PATCH /api/v1/admin/users/:userId/status endpoint
    - Reload user data after successful status change
    - _Requirements: 12.1, 12.2, 12.3_

  - [x] 7.4 Implement role change modal
    - Add role change modal to UserDetail component
    - Provide dropdown to select new role (user, admin)
    - Show confirmation dialog when assigning admin role
    - Call PATCH /api/v1/admin/users/:userId/role endpoint
    - Reload user data after successful role change
    - _Requirements: 13.1, 13.2, 13.3_

  - [x] 7.5 Implement password reset action
    - Add password reset button handler to UserDetail component
    - Show confirmation dialog before triggering reset
    - Call POST /api/v1/admin/users/:userId/reset-password endpoint
    - Display success message after email sent
    - _Requirements: 14.1, 14.3_

  - [ ]* 7.6 Write frontend component tests
    - Create tests for Users component using Vitest + React Testing Library
    - Test admin status check and access control
    - Test user list rendering and filtering
    - Test user detail rendering and actions
    - _Requirements: 16.1_

- [x] 8. Checkpoint - Frontend implementation complete
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. Property-based testing implementation
  - [ ]* 9.1 Write property test for admin endpoint authorization
    - Create `services/api-gateway-service/properties_test.go`
    - **Property 1: Admin Endpoint Authorization**
    - **Validates: Requirements 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 6.2**
    - Use gopter library with 100 minimum successful tests
    - Generate test cases with varying endpoints, roles, and token validity
    - Verify admin tokens get 2xx, user tokens get 403, invalid tokens get 401
    - _Requirements: 16.2_

  - [ ]* 9.2 Write property test for JWT role extraction
    - Add to `services/api-gateway-service/properties_test.go`
    - **Property 2: JWT Role Extraction Consistency**
    - **Validates: Requirements 2.2, 3.2**
    - Generate tokens with different roles (admin, user, moderator)
    - Verify extracted role matches token's app_metadata role exactly
    - Verify no database queries are made during extraction
    - _Requirements: 16.2_

  - [ ]* 9.3 Write property test for error response sanitization
    - Add to `services/api-gateway-service/properties_test.go`
    - **Property 3: Error Response Sanitization**
    - **Validates: Requirements 9.1, 9.2, 9.5**
    - Generate various invalid tokens and endpoints
    - Verify error responses use consistent JSON structure
    - Verify no sensitive data (user_id, role, database info) in responses
    - _Requirements: 16.2_

  - [ ]* 9.4 Write property test for admin route pattern protection
    - Add to `services/api-gateway-service/properties_test.go`
    - **Property 4: Admin Route Pattern Protection**
    - **Validates: Requirements 3.6**
    - Generate various subpaths under /api/v1/admin/* with different HTTP methods
    - Verify all routes require admin role regardless of path or method
    - Verify non-admin tokens always receive 403
    - _Requirements: 16.2_

- [ ] 10. Migration strategy implementation
  - [x] 10.1 Configure Supabase custom claims hook
    - Log into Supabase Dashboard
    - Navigate to Authentication > Hooks
    - Enable "Custom Access Token" hook
    - Set hook function to public.custom_access_token_hook
    - Save configuration and verify hook is active
    - _Requirements: 2.1_

  - [x] 10.2 Run database migrations
    - Execute `db/migrations/012_add_user_roles.sql` on database
    - Execute `db/migrations/013_configure_custom_claims.sql` on database
    - Verify schema changes with `db/migrations/verify_schema.sql`
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

  - [x] 10.3 Assign initial admin roles
    - Update `db/scripts/assign_admin_roles.sql` with actual admin emails
    - Execute script to assign admin roles to initial users
    - Verify role assignments with SELECT query
    - _Requirements: 6.1, 6.3, 6.4_

  - [x] 10.4 Deploy backend with migration mode enabled
    - Set AUTH_MIGRATION_MODE=true in environment
    - Keep ADMIN_API_KEY in environment temporarily
    - Deploy updated api-gateway-service with dual auth support
    - Verify both JWT and API key authentication work
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

  - [x] 10.5 Test JWT authentication in staging
    - Test admin user can access admin endpoints with JWT token
    - Test non-admin user receives 403 with JWT token
    - Test invalid token receives 401
    - Verify custom claims include user_role in JWT
    - _Requirements: 2.1, 2.2, 3.1, 3.2, 3.3, 3.4, 3.5_

  - [x] 10.6 Deploy frontend with JWT-only authentication
    - Build frontend with updated code (no API key)
    - Deploy to production
    - Test admin pages work with JWT authentication
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 5.1, 5.2, 5.3, 5.4, 5.5_

  - [x] 10.7 Monitor production and remove API key fallback
    - Monitor logs for 3-7 days to verify JWT auth working
    - Verify no API key fallback usage in logs
    - Remove requireAdminWithFallback and use requireAdmin only
    - Remove isMigrationMode function
    - Set AUTH_MIGRATION_MODE=false
    - Remove ADMIN_API_KEY from environment
    - Deploy final backend version
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

- [x] 11. Checkpoint - Migration complete
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 12. Documentation and verification
  - [x] 12.1 Create migration completion documentation
    - Create `docs/ADMIN_AUTH_MIGRATION.md`
    - Document changes made (database, backend, frontend)
    - Document admin role assignment process
    - Document verification steps
    - Include rollback procedures
    - _Requirements: 16.4, 16.5_

  - [x] 12.2 Create admin role assignment guide
    - Create `docs/ADMIN_ROLE_ASSIGNMENT.md`
    - Document SQL script for assigning admin roles
    - Document verification queries
    - Document audit trail review
    - _Requirements: 6.1, 6.4, 16.4_

  - [ ]* 12.3 Create verification smoke test script
    - Create `scripts/verify_migration.sh`
    - Check for ADMIN_API_KEY references in codebase
    - Check for VITE_ADMIN_API_KEY references in codebase
    - Check for X-Admin-Secret header usage
    - Exit with error if any references found
    - _Requirements: 4.1, 4.2, 4.3, 4.4_

  - [ ]* 12.4 Run final verification
    - Execute smoke test script to verify API key removal
    - Test admin endpoints require JWT authentication
    - Test non-admin users cannot access admin endpoints
    - Verify role changes are logged in admin_role_audit table
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 6.5, 9.1, 9.2, 9.3, 9.4, 9.5_

- [x] 13. Final checkpoint - Feature complete
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional testing and verification tasks that can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation at key milestones
- Property tests validate universal correctness properties from the design
- Unit tests validate specific examples and edge cases
- Migration strategy ensures zero downtime with phased rollout
- All database changes are reversible via rollback scripts
