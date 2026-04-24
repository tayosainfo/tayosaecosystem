# Requirements Document

## Introduction

This document specifies requirements for implementing role-based admin authentication using Supabase Auth to replace the current insecure shared API key approach. The system will add admin role management to user records, validate admin privileges through JWT tokens on the backend, and remove the exposed frontend API key that creates a security vulnerability.

## Glossary

- **Admin_User**: A user account with administrative privileges granted through a role field in the database
- **JWT_Token**: JSON Web Token issued by Supabase Auth containing user identity and claims
- **Backend_Service**: Go services (user-service, api-gateway) that process API requests
- **Admin_Endpoint**: API endpoints requiring administrative privileges (e.g., `/api/v1/admin/kyc`, `/api/v1/admin/settings`)
- **Supabase_Auth**: Authentication service managing user sessions and JWT token generation
- **Role_Field**: Database column storing user's administrative status (boolean or enum)
- **Frontend_Client**: React + Vite application making API requests

## Requirements

### Requirement 1: Admin Role Storage

**User Story:** As a system administrator, I want user admin status stored in the database, so that admin privileges can be managed and verified.

#### Acceptance Criteria

1. THE Backend_Service SHALL add a role field to the user table schema
2. THE role field SHALL support at least two values: admin and regular user
3. THE role field SHALL default to regular user for new accounts
4. WHEN a user record is created, THE Backend_Service SHALL initialize the role field
5. THE Backend_Service SHALL provide a mechanism to assign admin role to existing users

### Requirement 2: JWT Token Admin Claim

**User Story:** As a backend developer, I want admin status included in JWT tokens, so that authentication and authorization happen in a single verification step.

#### Acceptance Criteria

1. WHEN Supabase_Auth generates a JWT_Token, THE token SHALL include the user's role information
2. THE Backend_Service SHALL extract role claims from the JWT_Token without additional database queries
3. THE JWT_Token SHALL remain valid according to Supabase_Auth expiration policies
4. IF the JWT_Token lacks role information, THEN THE Backend_Service SHALL treat the user as a regular user

### Requirement 3: Backend Admin Authorization

**User Story:** As a security engineer, I want admin endpoints protected by role verification, so that only authorized users can access administrative functions.

#### Acceptance Criteria

1. WHEN a request arrives at an Admin_Endpoint, THE Backend_Service SHALL verify the JWT_Token
2. WHEN the JWT_Token is verified, THE Backend_Service SHALL extract the role claim
3. IF the role claim indicates admin privileges, THEN THE Backend_Service SHALL process the request
4. IF the role claim does not indicate admin privileges, THEN THE Backend_Service SHALL return a 403 Forbidden error
5. IF the JWT_Token is invalid or missing, THEN THE Backend_Service SHALL return a 401 Unauthorized error
6. THE Backend_Service SHALL apply role verification to all endpoints under `/api/v1/admin/*`

### Requirement 4: Remove Shared API Key

**User Story:** As a security engineer, I want the shared API key removed from the system, so that the security vulnerability is eliminated.

#### Acceptance Criteria

1. THE Backend_Service SHALL remove all code that validates the `ADMIN_API_KEY` environment variable
2. THE Frontend_Client SHALL remove all code that sends the `VITE_ADMIN_API_KEY` in request headers
3. THE Backend_Service SHALL remove the `ADMIN_API_KEY` environment variable from configuration
4. THE Frontend_Client SHALL remove the `VITE_ADMIN_API_KEY` environment variable from configuration
5. WHEN admin requests are made, THE Frontend_Client SHALL send only the JWT_Token for authentication

### Requirement 5: Frontend Authentication Integration

**User Story:** As a frontend developer, I want to use existing Supabase auth tokens for admin requests, so that no additional authentication mechanism is needed.

#### Acceptance Criteria

1. WHEN making admin requests, THE Frontend_Client SHALL include the JWT_Token from Supabase_Auth
2. THE Frontend_Client SHALL use the same authentication mechanism for admin and regular endpoints
3. IF the JWT_Token is expired, THEN THE Frontend_Client SHALL refresh the token before retrying
4. THE Frontend_Client SHALL handle 403 Forbidden responses by displaying an access denied message
5. THE Frontend_Client SHALL handle 401 Unauthorized responses by redirecting to login

### Requirement 6: Admin Role Assignment

**User Story:** As a system administrator, I want a secure method to assign admin roles, so that I can grant administrative access to trusted users.

#### Acceptance Criteria

1. THE Backend_Service SHALL provide a SQL script to assign admin role to users by email or user ID
2. WHERE a secure admin assignment endpoint is implemented, THE endpoint SHALL require existing admin authentication
3. THE admin assignment mechanism SHALL validate that the target user exists before assignment
4. WHEN admin role is assigned, THE Backend_Service SHALL update the user's role field
5. THE Backend_Service SHALL log all admin role assignments for audit purposes

### Requirement 7: Backward Compatibility During Migration

**User Story:** As a DevOps engineer, I want a safe migration path, so that the system remains operational during the transition.

#### Acceptance Criteria

1. THE Backend_Service SHALL support a migration phase where both authentication methods work
2. WHILE in migration mode, IF JWT_Token authentication fails, THEN THE Backend_Service SHALL fall back to API key validation
3. THE Backend_Service SHALL log which authentication method was used for each request
4. THE migration mode SHALL be controlled by a feature flag or environment variable
5. WHEN migration is complete, THE Backend_Service SHALL remove the fallback authentication code

### Requirement 8: Admin Status Verification

**User Story:** As a frontend developer, I want to check if the current user is an admin, so that I can show or hide admin UI elements.

#### Acceptance Criteria

1. THE Frontend_Client SHALL provide a function to check current user's admin status
2. THE admin status check SHALL use the JWT_Token claims without making additional API calls
3. WHEN the user logs in, THE Frontend_Client SHALL cache the admin status
4. WHEN the JWT_Token is refreshed, THE Frontend_Client SHALL update the cached admin status
5. THE Frontend_Client SHALL hide admin navigation and UI elements for non-admin users

### Requirement 9: Error Handling and Security

**User Story:** As a security engineer, I want proper error handling that doesn't leak sensitive information, so that the system remains secure.

#### Acceptance Criteria

1. WHEN authorization fails, THE Backend_Service SHALL return generic error messages
2. THE Backend_Service SHALL NOT include role information or user details in error responses
3. THE Backend_Service SHALL log detailed authorization failures for security monitoring
4. IF multiple authorization failures occur from the same user, THEN THE Backend_Service SHALL log a security event
5. THE error responses SHALL be consistent between authentication and authorization failures

### Requirement 10: Admin User List and Search

**User Story:** As an administrator, I want to view and search all users in the system, so that I can manage user accounts effectively.

#### Acceptance Criteria

1. THE Frontend_Client SHALL display a paginated list of all users
2. THE Frontend_Client SHALL provide search functionality by name, email, or phone number
3. THE Frontend_Client SHALL provide filters for user status (active, suspended, deactivated)
4. THE Frontend_Client SHALL provide filters for KYC status (pending, approved, rejected)
5. THE Backend_Service SHALL provide an endpoint `/api/v1/admin/users` that returns user data with pagination

### Requirement 11: User Detail View

**User Story:** As an administrator, I want to view detailed information about a specific user, so that I can understand their account status and activity.

#### Acceptance Criteria

1. THE Frontend_Client SHALL display comprehensive user information including profile, KYC status, and account status
2. THE Frontend_Client SHALL display user's onboarding progress and completion status
3. THE Frontend_Client SHALL display user's role and permissions
4. THE Frontend_Client SHALL display user's registration date and last login
5. THE Backend_Service SHALL provide an endpoint `/api/v1/admin/users/:userId` that returns detailed user information

### Requirement 12: User Status Management

**User Story:** As an administrator, I want to change user account status, so that I can suspend or deactivate problematic accounts.

#### Acceptance Criteria

1. THE Frontend_Client SHALL provide controls to change user status (active, suspended, deactivated)
2. WHEN changing user status, THE Frontend_Client SHALL require a reason/note
3. THE Backend_Service SHALL update the user's status in the database
4. THE Backend_Service SHALL log all status changes with admin ID, timestamp, and reason
5. WHEN a user is suspended, THE Backend_Service SHALL invalidate their active sessions

### Requirement 13: Admin Role Assignment UI

**User Story:** As an administrator, I want to assign or revoke admin roles through the UI, so that I can manage administrative access without SQL scripts.

#### Acceptance Criteria

1. THE Frontend_Client SHALL provide a control to assign or revoke admin role
2. WHEN assigning admin role, THE Frontend_Client SHALL require confirmation
3. THE Backend_Service SHALL provide an endpoint `/api/v1/admin/users/:userId/role` to update user roles
4. THE Backend_Service SHALL require existing admin authentication for role changes
5. THE Backend_Service SHALL log all role changes in the admin_role_audit table

### Requirement 14: User Credential Management

**User Story:** As an administrator, I want to reset user credentials or unlock accounts, so that I can help users who are locked out.

#### Acceptance Criteria

1. THE Frontend_Client SHALL provide a control to trigger password reset for a user
2. THE Frontend_Client SHALL provide a control to unlock a user account after failed login attempts
3. THE Backend_Service SHALL send password reset email to the user when triggered by admin
4. THE Backend_Service SHALL unlock user accounts and reset failed login counters
5. THE Backend_Service SHALL log all credential management actions

### Requirement 15: User Activity Audit Trail

**User Story:** As an administrator, I want to view a user's activity history, so that I can investigate suspicious behavior or disputes.

#### Acceptance Criteria

1. THE Frontend_Client SHALL display a timeline of user activities (logins, transactions, status changes)
2. THE Frontend_Client SHALL display login history with timestamps, IP addresses, and device information
3. THE Frontend_Client SHALL display admin actions performed on the user account
4. THE Backend_Service SHALL provide an endpoint `/api/v1/admin/users/:userId/activity` that returns activity logs
5. THE activity logs SHALL include at least 90 days of history

### Requirement 16: Documentation and Testing

**User Story:** As a developer, I want clear documentation and test coverage, so that I can maintain and verify the authentication system.

#### Acceptance Criteria

1. THE Backend_Service SHALL include unit tests for JWT_Token validation logic
2. THE Backend_Service SHALL include integration tests for admin endpoint authorization
3. THE Backend_Service SHALL include integration tests for user management endpoints
4. THE project SHALL include documentation for assigning admin roles
5. THE project SHALL include documentation for the migration process
6. THE Backend_Service SHALL include tests verifying that non-admin users cannot access Admin_Endpoints
