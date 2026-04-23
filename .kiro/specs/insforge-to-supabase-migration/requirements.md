# Requirements Document

## Introduction

This document specifies the requirements for migrating the Tayosa banking ecosystem from InsForge Backend-as-a-Service (BaaS) to Supabase. This is a straightforward migration that replaces InsForge with Supabase while maintaining the same architecture, functionality, and user experience. The primary driver for this migration is InsForge's inability to send verification emails, which blocks user registration.

The system consists of:
- **Frontend**: React + Vite web application (Supabase SDK already installed)
- **Backend**: 9 Go microservices (api-gateway, user, affiliate, audit-log, fee, kibiina, loan-credit, notification, object-storage)
- **Mobile**: Flutter application
- **Database**: PostgreSQL with 10 migration files

## Glossary

- **System**: The Tayosa banking ecosystem (frontend, backend services, mobile app, database)
- **Frontend**: React + Vite web application
- **Backend_Services**: Collection of 9 Go microservices
- **User_Service**: Go microservice handling authentication and user management
- **Mobile_App**: Flutter mobile application
- **InsForge**: Current Backend-as-a-Service provider being replaced
- **Supabase**: Target Backend-as-a-Service provider
- **Supabase_Client**: Supabase SDK client instance
- **Auth_Token**: JWT access token for authenticated requests
- **Service_Role_Key**: Supabase admin key for privileged operations
- **Anon_Key**: Supabase public key for client-side operations
- **RLS**: Row Level Security (Supabase database security feature)
- **Migration**: Process of replacing InsForge with Supabase

## Requirements

### Requirement 1: Frontend Authentication Migration

**User Story:** As a frontend developer, I want the React application to use Supabase for authentication, so that users can register and log in with email verification.

#### Acceptance Criteria

1. THE Frontend SHALL replace all InsForge SDK imports with Supabase SDK imports
2. THE Frontend SHALL use the existing Supabase_Client configured in src/lib/insforge.ts
3. WHEN a user registers, THE Frontend SHALL call Supabase authentication methods
4. WHEN a user logs in, THE Frontend SHALL call Supabase authentication methods
5. WHEN a user logs out, THE Frontend SHALL call Supabase authentication methods
6. THE Frontend SHALL store Auth_Token in the same format as before migration
7. THE Frontend SHALL maintain the same authentication flow user experience
8. THE Frontend SHALL handle Supabase authentication errors with user-friendly messages

### Requirement 2: Backend User Service Migration

**User Story:** As a backend developer, I want the user-service to fully integrate with Supabase, so that authentication operations work correctly.

#### Acceptance Criteria

1. THE User_Service SHALL remove all InsForge SDK dependencies from go.mod
2. THE User_Service SHALL add Supabase Go client library to go.mod
3. THE User_Service SHALL use environment variables SUPABASE_URL, SUPABASE_ANON_KEY, and SUPABASE_SERVICE_ROLE_KEY
4. THE User_Service SHALL implement all authentication endpoints using Supabase API
5. THE User_Service SHALL maintain backward compatibility with existing API contracts
6. WHEN User_Service receives authentication requests, THE User_Service SHALL forward them to Supabase
7. THE User_Service SHALL validate Auth_Token using Supabase token verification
8. THE User_Service SHALL handle Supabase API errors and return appropriate HTTP status codes

### Requirement 3: Backend Services Authentication Update

**User Story:** As a backend developer, I want all backend services to validate tokens using Supabase, so that authentication works across all microservices.

#### Acceptance Criteria

1. THE Backend_Services SHALL update token validation logic to use Supabase
2. WHEN a service receives a request with Auth_Token, THE service SHALL validate it against Supabase
3. THE Backend_Services SHALL use Service_Role_Key for admin operations
4. THE Backend_Services SHALL use Anon_Key for public operations
5. THE Backend_Services SHALL maintain the same authorization logic as before migration
6. THE Backend_Services SHALL handle Supabase authentication failures with appropriate error responses

### Requirement 4: Mobile App Authentication Migration

**User Story:** As a mobile developer, I want the Flutter app to use Supabase for authentication, so that mobile users can register and log in.

#### Acceptance Criteria

1. THE Mobile_App SHALL add Supabase Flutter SDK to pubspec.yaml
2. THE Mobile_App SHALL remove InsForge SDK dependencies from pubspec.yaml
3. THE Mobile_App SHALL configure Supabase_Client with SUPABASE_URL and SUPABASE_ANON_KEY
4. THE Mobile_App SHALL replace all InsForge authentication calls with Supabase authentication calls
5. THE Mobile_App SHALL store Auth_Token using the same mechanism as before migration
6. THE Mobile_App SHALL maintain the same authentication UI and user experience
7. WHEN a mobile user registers, THE Mobile_App SHALL call Supabase signup methods
8. WHEN a mobile user logs in, THE Mobile_App SHALL call Supabase signin methods

### Requirement 5: Environment Configuration Migration

**User Story:** As a DevOps engineer, I want environment variables updated for Supabase, so that all components can connect to the new backend.

#### Acceptance Criteria

1. THE System SHALL replace VITE_INSFORGE_URL with VITE_SUPABASE_URL in .env files
2. THE System SHALL replace VITE_INSFORGE_ANON_KEY with VITE_SUPABASE_ANON_KEY in .env files
3. THE System SHALL add SUPABASE_SERVICE_ROLE_KEY to backend service environment variables
4. THE System SHALL update .env.example with Supabase configuration template
5. THE System SHALL document all required environment variables
6. THE System SHALL use the Supabase project URL: https://ablvrbnbsdqshrorhmjf.supabase.co
7. THE System SHALL configure RLS policies if required for database access

### Requirement 6: Database Migration Compatibility

**User Story:** As a database administrator, I want existing PostgreSQL migrations to work with Supabase, so that the database schema remains consistent.

#### Acceptance Criteria

1. THE System SHALL verify all 10 migration files are compatible with Supabase PostgreSQL
2. WHEN migrations are applied, THE System SHALL execute them in the correct order
3. THE System SHALL maintain all existing tables, columns, and constraints
4. THE System SHALL preserve all existing data during migration
5. IF a migration uses InsForge-specific features, THEN THE System SHALL replace them with Supabase equivalents
6. THE System SHALL configure database connection strings for Supabase PostgreSQL

### Requirement 7: Test File Migration

**User Story:** As a developer, I want test files updated to use Supabase, so that automated tests continue to work.

#### Acceptance Criteria

1. THE System SHALL update test_unverified_login.js to use Supabase SDK
2. THE System SHALL update test_unified_auth_flow.js to use Supabase SDK
3. THE System SHALL replace InsForge client initialization with Supabase client initialization in test files
4. WHEN tests run, THE System SHALL connect to Supabase test environment
5. THE System SHALL maintain the same test assertions and expected behaviors

### Requirement 8: Email Verification Functionality

**User Story:** As a user, I want to receive email verification when I register, so that I can activate my account.

#### Acceptance Criteria

1. WHEN a user registers, THE System SHALL send a verification email via Supabase
2. THE System SHALL configure Supabase email templates for verification emails
3. WHEN a user clicks the verification link, THE System SHALL verify the email address
4. THE System SHALL prevent unverified users from accessing protected resources
5. THE System SHALL provide a resend verification email option
6. THE System SHALL display appropriate messages for email verification status

### Requirement 9: Password Reset Functionality

**User Story:** As a user, I want to reset my password using email, so that I can regain access if I forget my password.

#### Acceptance Criteria

1. WHEN a user requests password reset, THE System SHALL send a reset email via Supabase
2. THE System SHALL configure Supabase email templates for password reset emails
3. WHEN a user clicks the reset link, THE System SHALL allow password change
4. THE System SHALL validate new passwords meet security requirements
5. THE System SHALL invalidate old Auth_Token after password reset
6. THE System SHALL confirm successful password reset to the user

### Requirement 10: OAuth Integration Migration

**User Story:** As a user, I want to log in with OAuth providers, so that I can use social authentication.

#### Acceptance Criteria

1. THE System SHALL migrate OAuth provider configurations from InsForge to Supabase
2. THE System SHALL support the same OAuth providers as before migration
3. WHEN a user initiates OAuth login, THE System SHALL redirect to Supabase OAuth flow
4. WHEN OAuth authentication succeeds, THE System SHALL create or update user records
5. THE System SHALL handle OAuth callback URLs correctly
6. THE System SHALL sync OAuth user data with local user profiles

### Requirement 11: Session Management Migration

**User Story:** As a user, I want my session to persist across page refreshes, so that I don't have to log in repeatedly.

#### Acceptance Criteria

1. THE System SHALL use Supabase session management for token refresh
2. WHEN Auth_Token expires, THE System SHALL automatically refresh it using Supabase
3. THE System SHALL store refresh tokens securely
4. THE System SHALL maintain session state across browser refreshes
5. WHEN a user logs out, THE System SHALL invalidate all session tokens
6. THE System SHALL handle session expiration gracefully with re-authentication prompts

### Requirement 12: API Gateway Service Migration

**User Story:** As a backend developer, I want the API gateway to route requests correctly with Supabase authentication, so that all services work together.

#### Acceptance Criteria

1. THE api-gateway-service SHALL validate Auth_Token using Supabase before routing requests
2. THE api-gateway-service SHALL forward authenticated requests to appropriate backend services
3. THE api-gateway-service SHALL handle Supabase authentication errors at the gateway level
4. THE api-gateway-service SHALL maintain the same routing logic as before migration
5. THE api-gateway-service SHALL add Supabase authentication headers to forwarded requests

### Requirement 13: Error Handling Consistency

**User Story:** As a developer, I want consistent error handling across all components, so that debugging is easier.

#### Acceptance Criteria

1. THE System SHALL map Supabase error codes to application error messages
2. THE System SHALL maintain the same error response format as before migration
3. WHEN Supabase returns an error, THE System SHALL log it with appropriate context
4. THE System SHALL provide user-friendly error messages for common authentication failures
5. THE System SHALL handle network errors when communicating with Supabase

### Requirement 14: Documentation and Configuration Files

**User Story:** As a developer, I want updated documentation, so that I understand the new Supabase integration.

#### Acceptance Criteria

1. THE System SHALL update README.md with Supabase setup instructions
2. THE System SHALL document all Supabase environment variables
3. THE System SHALL provide migration guide from InsForge to Supabase
4. THE System SHALL update API documentation with Supabase authentication details
5. THE System SHALL document Supabase project configuration (URL, keys, RLS policies)

### Requirement 15: Backward Compatibility During Migration

**User Story:** As a project manager, I want the system to remain functional during migration, so that users are not disrupted.

#### Acceptance Criteria

1. THE System SHALL support a migration phase where both InsForge and Supabase can coexist
2. THE System SHALL provide feature flags to toggle between InsForge and Supabase
3. WHEN migration is incomplete, THE System SHALL fall back to InsForge for unmigrated components
4. THE System SHALL log which authentication backend is being used for each request
5. WHEN migration is complete, THE System SHALL remove all InsForge dependencies

### Requirement 16: Storage Service Migration

**User Story:** As a developer, I want the object-storage-service to use Supabase Storage, so that file uploads work correctly.

#### Acceptance Criteria

1. THE object-storage-service SHALL replace InsForge storage API calls with Supabase Storage API calls
2. THE object-storage-service SHALL configure Supabase Storage buckets
3. WHEN a file is uploaded, THE object-storage-service SHALL store it in Supabase Storage
4. WHEN a file is downloaded, THE object-storage-service SHALL retrieve it from Supabase Storage
5. THE object-storage-service SHALL maintain the same file access permissions as before migration
6. THE object-storage-service SHALL migrate existing files from InsForge to Supabase Storage

### Requirement 17: Real-time Features Migration

**User Story:** As a developer, I want real-time features to use Supabase Realtime, so that live updates continue to work.

#### Acceptance Criteria

1. IF the System uses InsForge real-time features, THEN THE System SHALL replace them with Supabase Realtime
2. THE System SHALL configure Supabase Realtime channels for live data updates
3. WHEN database changes occur, THE System SHALL broadcast them via Supabase Realtime
4. THE System SHALL maintain the same real-time event structure as before migration
5. THE System SHALL handle Supabase Realtime connection errors gracefully

### Requirement 18: Deployment and CI/CD Updates

**User Story:** As a DevOps engineer, I want deployment scripts updated for Supabase, so that automated deployments work correctly.

#### Acceptance Criteria

1. THE System SHALL update deployment scripts to use Supabase environment variables
2. THE System SHALL configure CI/CD pipelines with Supabase credentials
3. THE System SHALL update Docker configurations with Supabase settings
4. THE System SHALL document deployment process with Supabase
5. THE System SHALL provide rollback procedures in case of migration issues
