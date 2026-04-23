# Design Document: InsForge to Supabase Migration

## Overview

This design document specifies the technical implementation for migrating the Tayosa banking ecosystem from InsForge Backend-as-a-Service to Supabase. The migration is driven by InsForge's inability to send verification emails, which blocks user registration.

### System Context

The Tayosa ecosystem consists of:
- **Frontend**: React + Vite web application (Supabase SDK already installed at `@supabase/supabase-js@^2.104.0`)
- **Backend**: 9 Go microservices with existing Supabase client implementation
- **Mobile**: Flutter application requiring Supabase Flutter SDK integration
- **Database**: PostgreSQL with 10 migration files compatible with Supabase

### Current State Analysis

**Frontend**: 
- Supabase client already configured in `src/lib/insforge.ts`
- Uses `platformApi` wrapper for authentication calls
- Stores tokens in `sessionStorage`

**Backend**:
- `user-service` has complete Supabase client implementation (`supabase_client.go`)
- Authentication validation uses Supabase token verification
- Environment variables already support Supabase configuration
- InsForge references remain in variable names and comments

**Mobile**:
- Currently uses custom API client (`api_client.dart`)
- No Supabase SDK integration
- Requires full authentication flow migration

### Migration Strategy

This is a **replacement migration**, not a dual-backend migration. The system will:
1. Remove all InsForge SDK dependencies
2. Update variable names and references from "InsForge" to "Supabase"
3. Ensure all authentication flows use Supabase
4. Maintain backward compatibility with existing API contracts
5. Preserve all user data and session management behavior

## Architecture

### Authentication Flow

```mermaid
sequenceDiagram
    participant Client as Frontend/Mobile
    participant Gateway as API Gateway
    participant UserSvc as User Service
    participant Supabase as Supabase Auth
    participant DB as PostgreSQL

    Client->>UserSvc: POST /api/v1/auth/register
    UserSvc->>Supabase: POST /auth/v1/signup
    Supabase-->>UserSvc: {user, access_token}
    UserSvc->>DB: Create local profile
    UserSvc-->>Client: {session, user}
    
    Client->>UserSvc: POST /api/v1/auth/login
    UserSvc->>Supabase: POST /auth/v1/token
    Supabase-->>UserSvc: {access_token, refresh_token}
    UserSvc->>DB: Sync Supabase user ID
    UserSvc-->>Client: {session, user}
    
    Client->>Gateway: API request + Bearer token
    Gateway->>Supabase: GET /auth/v1/user (validate token)
    Supabase-->>Gateway: {user}
    Gateway->>UserSvc: Forward request
    UserSvc-->>Client: Response
```

### Service Architecture

```mermaid
graph TB
    Frontend[React Frontend]
    Mobile[Flutter Mobile]
    Gateway[API Gateway Service]
    UserSvc[User Service]
    Supabase[Supabase Auth]
    DB[(PostgreSQL)]
    
    Frontend -->|Supabase SDK| Supabase
    Mobile -->|Supabase SDK| Supabase
    Frontend -->|API Calls| Gateway
    Mobile -->|API Calls| Gateway
    Gateway -->|Route| UserSvc
    UserSvc -->|Validate Token| Supabase
    UserSvc -->|Store Profile| DB
    
    style Supabase fill:#3ecf8e
    style DB fill:#336791
```

### Data Flow

1. **Registration**: Client → Supabase signup → User Service creates local profile → Return session
2. **Login**: Client → Supabase token endpoint → User Service syncs user ID → Return session
3. **Token Validation**: Service → Supabase `/auth/v1/user` → Validate → Continue
4. **Session Refresh**: Client → Supabase `/auth/v1/token?grant_type=refresh_token` → New tokens

## Components and Interfaces

### Frontend Components

#### Supabase Client Configuration
**File**: `src/lib/insforge.ts` (rename to `src/lib/supabase.ts`)

```typescript
import { createClient } from '@supabase/supabase-js';

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL;
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY;

if (!supabaseUrl || !supabaseAnonKey) {
    throw new Error('Missing Supabase environment variables');
}

export const supabase = createClient(supabaseUrl, supabaseAnonKey);
```

#### Platform API Wrapper
**File**: `src/lib/platformApi.ts`

**Current**: Uses fetch to call backend endpoints
**Change**: No changes required - backend handles Supabase integration
**Rationale**: Maintains separation of concerns; backend abstracts auth provider

#### Authentication Context
**File**: `src/contexts/AuthContext.tsx`

**Current**: Stores user in `sessionStorage` as `auth_user` and token as `auth_token`
**Change**: No changes required - session format remains identical
**Rationale**: Backend returns same session structure regardless of auth provider

### Backend Components

#### User Service - Supabase Client
**File**: `services/user-service/supabase_client.go`

**Current State**: Fully implemented with all Supabase operations
**Required Changes**:
1. Rename function `sessionPayloadFromInsForge` → `sessionPayloadFromSupabase` (already done)
2. Update comments referencing InsForge
3. Verify all error messages reference Supabase

**Key Functions**:
- `supabaseDoJSON`: Core HTTP dispatcher for Supabase API calls
- `supabasePost`: POST requests with anon key
- `supabaseUserGet`: GET requests with user token
- `supabaseAdminPost`: POST requests with service role key
- `supabaseRefreshForward`: Token refresh forwarding
- `supabaseLogoutForward`: Logout forwarding

#### User Service - Authentication Handlers
**File**: `services/user-service/handlers.go`

**Registration Flow**:
```go
// POST /api/v1/auth/register
1. Validate input (email, password, phone, fullName)
2. Call Supabase: POST /auth/v1/signup
3. Extract user ID from response
4. Create local User record in PostgreSQL
5. Create OnboardingProfile record
6. Return session + user profile
```

**Login Flow**:
```go
// POST /api/v1/auth/login
1. Resolve identifier (email or phone) to Supabase email
2. Call Supabase: POST /auth/v1/token?grant_type=password
3. Sync Supabase user ID to local profile
4. Return session + user profile
```

**Token Validation**:
```go
// requireAuth middleware
1. Extract Bearer token from Authorization header
2. Call Supabase: GET /auth/v1/user with token
3. Extract user ID from response
4. Add user ID to request context
```

#### API Gateway Service
**File**: `services/api-gateway-service/main.go`

**Current**: Routes requests to backend services
**Required Changes**:
1. Add token validation before routing (call user-service or validate directly)
2. Forward `Authorization` header to downstream services
3. Handle Supabase authentication errors at gateway level

**Routing Logic**:
```go
/api/v1/auth/*     → user-service
/api/v1/users/*    → user-service (requires auth)
/api/v1/storage/*  → object-storage-service (requires auth)
/api/v1/affiliate/* → affiliate-service (requires auth)
// ... other routes
```

#### Other Backend Services

**Services**: affiliate, audit-log, fee, kibiina, loan-credit, notification, object-storage

**Required Changes**:
1. Update token validation to use Supabase
2. Use `SUPABASE_URL` and `SUPABASE_ANON_KEY` environment variables
3. Validate tokens by calling Supabase `/auth/v1/user` endpoint
4. Handle Supabase error responses

**Token Validation Pattern**:
```go
func validateSupabaseToken(token string) (userID string, err error) {
    req, _ := http.NewRequest("GET", 
        os.Getenv("SUPABASE_URL") + "/auth/v1/user", nil)
    req.Header.Set("Authorization", "Bearer " + token)
    req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
    
    resp, err := http.DefaultClient.Do(req)
    // Parse response and extract user ID
}
```

### Mobile Components

#### Supabase Client Setup
**File**: `app/mobile_app/lib/core/network/supabase_client.dart` (new file)

```dart
import 'package:supabase_flutter/supabase_flutter.dart';

class SupabaseConfig {
  static const String supabaseUrl = String.fromEnvironment(
    'SUPABASE_URL',
    defaultValue: 'https://ablvrbnbsdqshrorhmjf.supabase.co',
  );
  
  static const String supabaseAnonKey = String.fromEnvironment(
    'SUPABASE_ANON_KEY',
    defaultValue: '', // Set via --dart-define
  );
  
  static Future<void> initialize() async {
    await Supabase.initialize(
      url: supabaseUrl,
      anonKey: supabaseAnonKey,
    );
  }
  
  static SupabaseClient get client => Supabase.instance.client;
}
```

#### API Client Migration
**File**: `app/mobile_app/lib/core/network/api_client.dart`

**Current**: Custom HTTP client using Dio
**Change**: Add Supabase authentication integration

```dart
class ApiClient {
  final Dio _dio;
  final SupabaseClient _supabase;
  
  ApiClient() : 
    _dio = Dio(BaseOptions(baseUrl: apiBaseUrl)),
    _supabase = SupabaseConfig.client {
    
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        // Get current session from Supabase
        final session = _supabase.auth.currentSession;
        if (session != null) {
          options.headers['Authorization'] = 
            'Bearer ${session.accessToken}';
        }
        handler.next(options);
      },
    ));
  }
  
  Future<void> register({
    required String email,
    required String password,
    required String fullName,
    required String phone,
  }) async {
    // Call backend /api/v1/auth/register
    // Backend handles Supabase signup
  }
  
  Future<void> login({
    required String identifier,
    required String password,
  }) async {
    // Call backend /api/v1/auth/login
    // Backend handles Supabase authentication
  }
}
```

#### Authentication Screens
**Files**: 
- `app/mobile_app/lib/features/auth/presentation/login_screen.dart`
- `app/mobile_app/lib/features/auth/presentation/register_screen.dart`

**Current**: Call custom API endpoints
**Change**: Update to use Supabase-aware API client
**Note**: UI remains unchanged; only backend integration changes

## Data Models

### User Identity Model
**Table**: `users_identity`

```sql
CREATE TABLE users_identity (
  user_id TEXT PRIMARY KEY,              -- Supabase user ID
  full_name TEXT NOT NULL,
  phone_e164 TEXT NOT NULL UNIQUE,
  auth_email TEXT NOT NULL UNIQUE,       -- Email used for Supabase auth
  contact_email TEXT UNIQUE,
  insforge_user_id TEXT UNIQUE,          -- Rename to supabase_user_id
  phone_verified_at TIMESTAMP NULL,
  contact_email_verified_at TIMESTAMP NULL,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  password_hash TEXT,                    -- For local fallback (not used with Supabase)
  insforge_login_email TEXT              -- Rename to supabase_login_email
);
```

**Migration Required**:
```sql
-- Rename columns to reflect Supabase
ALTER TABLE users_identity 
  RENAME COLUMN insforge_user_id TO supabase_user_id;

ALTER TABLE users_identity 
  RENAME COLUMN insforge_login_email TO supabase_login_email;

-- Update index
DROP INDEX IF EXISTS idx_users_insforge_login_email;
CREATE INDEX idx_users_supabase_login_email 
  ON users_identity(supabase_login_email);
```

### Session Model
**Storage**: Client-side (sessionStorage for web, secure storage for mobile)

```typescript
interface Session {
  accessToken: string;      // Supabase JWT
  userId: string;           // Supabase user ID
  refreshToken?: string;    // Supabase refresh token
}

interface User {
  id: string;               // Supabase user ID
  fullName: string;
  phoneE164: string;
  contactEmail: string;
  contactEmailVerified: boolean;
  supabaseUserId: string;   // Same as id
  dateOfBirth?: string;
  nationality?: string;
  createdAt: string;
}
```

### Supabase Auth Response Models

**Signup Response**:
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "email_confirmed_at": null,
  "confirmation_sent_at": "2024-01-01T00:00:00Z",
  "access_token": "jwt-token",  // Only if email confirmation disabled
  "refresh_token": "refresh-token"
}
```

**Login Response**:
```json
{
  "access_token": "jwt-token",
  "token_type": "bearer",
  "expires_in": 3600,
  "refresh_token": "refresh-token",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "email_confirmed_at": "2024-01-01T00:00:00Z"
  }
}
```

## Error Handling

### Supabase Error Mapping

**Error Response Structure**:
```json
{
  "error": "error_code",
  "error_description": "Human readable message",
  "message": "Alternative message field"
}
```

**Common Errors**:

| Supabase Error | HTTP Status | Application Message |
|----------------|-------------|---------------------|
| `email_not_confirmed` | 400 | "Please verify your email address" |
| `invalid_credentials` | 400 | "Invalid email or password" |
| `user_already_exists` | 422 | "Email already registered" |
| `invalid_grant` | 400 | "Invalid or expired token" |
| `too_many_requests` | 429 | "Too many attempts. Please try again later" |

**Error Handling Pattern**:
```go
func supabaseUpstreamHTTP(err error) (status int, msg string) {
    var sbErr *SupabaseRequestError
    if errors.As(err, &sbErr) {
        // Map Supabase status codes to application codes
        if sbErr.Status >= 400 && sbErr.Status <= 599 {
            return sbErr.Status, sbErr.Message
        }
        return http.StatusBadGateway, sbErr.Message
    }
    return http.StatusBadGateway, err.Error()
}
```

### Frontend Error Handling

**PlatformApiError Class**:
```typescript
export class PlatformApiError extends Error {
  readonly status: number;
  readonly body: Record<string, unknown>;
  
  constructor(message: string, status: number, body: Record<string, unknown>) {
    super(message);
    this.name = 'PlatformApiError';
    this.status = status;
    this.body = body;
  }
}
```

**Usage**:
```typescript
try {
  await platformApi.login({ identifier, password });
} catch (err) {
  if (isPlatformApiError(err)) {
    if (err.body.requireEmailVerification) {
      // Show email verification prompt
    } else {
      // Show error message
    }
  }
}
```

### Mobile Error Handling

```dart
class SupabaseException implements Exception {
  final String message;
  final int? statusCode;
  final Map<String, dynamic>? details;
  
  SupabaseException(this.message, {this.statusCode, this.details});
  
  bool get requiresEmailVerification =>
    details?['requireEmailVerification'] == true;
}
```

## Testing Strategy

### Unit Testing

**Frontend Unit Tests**:
- Test `platformApi` methods with mocked fetch responses
- Test `AuthContext` state management
- Test error handling for various Supabase error responses
- Test session storage and retrieval

**Backend Unit Tests**:
- Test Supabase client functions with mocked HTTP responses
- Test authentication handlers with various input scenarios
- Test token validation logic
- Test error mapping functions

**Mobile Unit Tests**:
- Test API client with mocked Dio responses
- Test authentication state management
- Test secure storage operations

### Integration Testing

**Authentication Flow Tests**:
1. **Registration Flow**:
   - Register new user → Verify Supabase user created → Verify local profile created
   - Register with existing email → Verify conflict error
   - Register with email confirmation → Verify email sent → Verify email → Complete profile

2. **Login Flow**:
   - Login with email → Verify token returned → Verify user profile returned
   - Login with phone → Verify phone resolved to email → Verify login succeeds
   - Login with unverified email → Verify error with `requireEmailVerification` flag

3. **Token Validation**:
   - Make authenticated request → Verify token validated → Verify request succeeds
   - Make request with expired token → Verify 401 error
   - Make request with invalid token → Verify 401 error

4. **Session Management**:
   - Refresh token → Verify new access token returned
   - Logout → Verify session invalidated

**Test Files to Update**:
- `test_unified_auth_flow.js`: Update to use Supabase SDK
- `test_unverified_login.js`: Update to use Supabase SDK

**Example Test Update**:
```javascript
// Before (InsForge)
const { createClient } = require('@insforge/sdk');
const client = createClient({
  baseUrl: process.env.INSFORGE_URL,
  anonKey: process.env.INSFORGE_ANON_KEY
});

// After (Supabase)
const { createClient } = require('@supabase/supabase-js');
const client = createClient(
  process.env.VITE_SUPABASE_URL,
  process.env.VITE_SUPABASE_ANON_KEY
);
```

### End-to-End Testing

**Test Scenarios**:
1. Complete user journey: Register → Verify email → Login → Access protected resource
2. Password reset flow: Request reset → Verify email → Reset password → Login
3. OAuth flow: Initiate OAuth → Complete OAuth → Sync user data
4. Mobile app flow: Register on mobile → Login on web → Verify session sync

**Test Environment**:
- Use Supabase test project with separate database
- Configure test environment variables
- Use test email addresses that can be verified programmatically

## Implementation Notes

### Environment Variables

**Frontend (.env)**:
```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co
VITE_SUPABASE_ANON_KEY=your-supabase-anon-key
```

**Backend (.env)**:
```bash
# Supabase Configuration
SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co
SUPABASE_ANON_KEY=your-supabase-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-supabase-service-role-key

# Database
DATABASE_URL=postgresql://postgres:PASSWORD@db.ablvrbnbsdqshrorhmjf.supabase.co:5432/postgres?sslmode=require

# Service Ports
PORT=8081  # user-service
```

**Mobile (--dart-define)**:
```bash
flutter run \
  --dart-define=SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-key \
  --dart-define=API_BASE_URL=http://localhost:8080
```

### Supabase Project Configuration

**Email Templates**:
1. **Confirmation Email**: Customize template in Supabase Dashboard → Authentication → Email Templates
2. **Password Reset Email**: Customize template for password recovery
3. **Magic Link Email**: Configure if using magic link authentication

**Email Settings**:
- Enable email confirmations: Dashboard → Authentication → Settings → Enable email confirmations
- Configure SMTP (optional): Use custom SMTP for production emails
- Set redirect URLs: Configure allowed redirect URLs for email links

**RLS Policies**:
```sql
-- Enable RLS on auth.users (if needed)
ALTER TABLE auth.users ENABLE ROW LEVEL SECURITY;

-- Policy: Users can read their own data
CREATE POLICY "Users can read own data"
  ON auth.users
  FOR SELECT
  USING (auth.uid() = id);
```

**Note**: Supabase RLS is enabled by default. The application uses service role key for admin operations, so RLS policies primarily affect direct database access.

### Database Migration Strategy

**Step 1**: Rename columns in `users_identity` table
```sql
ALTER TABLE users_identity 
  RENAME COLUMN insforge_user_id TO supabase_user_id;
ALTER TABLE users_identity 
  RENAME COLUMN insforge_login_email TO supabase_login_email;
```

**Step 2**: Update indexes
```sql
DROP INDEX IF EXISTS idx_users_insforge_login_email;
CREATE INDEX idx_users_supabase_login_email 
  ON users_identity(supabase_login_email);
```

**Step 3**: Verify data integrity
```sql
-- Check for null supabase_user_id
SELECT COUNT(*) FROM users_identity WHERE supabase_user_id IS NULL;

-- Check for duplicate emails
SELECT auth_email, COUNT(*) 
FROM users_identity 
GROUP BY auth_email 
HAVING COUNT(*) > 1;
```

**Step 4**: Apply migration
```bash
# Run migration file
psql $DATABASE_URL -f db/migrations/011_rename_insforge_to_supabase.sql
```

### Code Refactoring Checklist

**Frontend**:
- [ ] Rename `src/lib/insforge.ts` → `src/lib/supabase.ts`
- [ ] Update imports in all files using the client
- [ ] Update comments referencing InsForge
- [ ] Verify environment variable names

**Backend - User Service**:
- [ ] Update comments in `supabase_client.go`
- [ ] Rename `InsforgeUserID` → `SupabaseUserID` in structs
- [ ] Rename `InsforgeEmail` → `SupabaseLoginEmail` in structs
- [ ] Update function `supabaseLoginEmail` logic
- [ ] Update error messages

**Backend - Other Services**:
- [ ] Add Supabase token validation to each service
- [ ] Update environment variable loading
- [ ] Add error handling for Supabase errors

**Mobile**:
- [ ] Add `supabase_flutter` dependency
- [ ] Create Supabase client configuration
- [ ] Update API client to use Supabase session
- [ ] Update authentication screens
- [ ] Update secure storage for tokens

**Tests**:
- [ ] Update `test_unified_auth_flow.js`
- [ ] Update `test_unverified_login.js`
- [ ] Add new integration tests for Supabase flows

**Documentation**:
- [ ] Update README.md with Supabase setup instructions
- [ ] Update API documentation
- [ ] Create migration runbook
- [ ] Document rollback procedures

### Deployment Considerations

**Pre-Deployment**:
1. Verify Supabase project is configured correctly
2. Test email delivery in Supabase
3. Verify database migrations are ready
4. Update environment variables in deployment environment
5. Run integration tests against Supabase test project

**Deployment Steps**:
1. Deploy database migrations
2. Deploy backend services with new environment variables
3. Deploy frontend with updated Supabase configuration
4. Deploy mobile app update
5. Monitor error logs for Supabase-related issues

**Rollback Plan**:
1. Revert database migrations (rename columns back)
2. Revert backend services to previous version
3. Revert frontend to previous version
4. Switch environment variables back to InsForge (if dual-backend was implemented)

**Monitoring**:
- Monitor Supabase dashboard for authentication metrics
- Monitor backend logs for Supabase API errors
- Monitor frontend error tracking for authentication failures
- Set up alerts for high error rates

### Security Considerations

**Token Security**:
- Access tokens are JWTs signed by Supabase
- Tokens are validated on every request
- Refresh tokens are stored securely (httpOnly cookies or secure storage)
- Service role key is never exposed to clients

**RLS (Row Level Security)**:
- Enabled by default in Supabase
- Application uses service role key for admin operations
- Client-side operations use anon key with RLS policies

**Environment Variables**:
- Never commit `.env` files to version control
- Use `.env.example` as template
- Rotate keys periodically
- Use different keys for development, staging, and production

**CORS Configuration**:
- Configure allowed origins in Supabase dashboard
- Restrict to known domains in production
- Allow localhost for development

### Performance Considerations

**Token Validation**:
- Cache Supabase user lookups (with short TTL)
- Use connection pooling for database queries
- Implement rate limiting for authentication endpoints

**Database Queries**:
- Use indexes on frequently queried columns
- Optimize onboarding profile queries
- Use prepared statements to prevent SQL injection

**API Response Times**:
- Target: < 200ms for authentication endpoints
- Target: < 100ms for token validation
- Monitor Supabase API latency

## Migration Phases

### Phase 1: Backend Preparation (Complete)
- ✅ Supabase client implementation in user-service
- ✅ Environment variable support
- ✅ Token validation logic
- ✅ Authentication handlers

### Phase 2: Code Refactoring (In Progress)
- Rename InsForge references to Supabase
- Update variable names and comments
- Update error messages
- Update documentation

### Phase 3: Mobile Integration (Pending)
- Add Supabase Flutter SDK
- Implement authentication flows
- Update API client
- Test on iOS and Android

### Phase 4: Testing (Pending)
- Update test files
- Run integration tests
- Perform end-to-end testing
- Load testing

### Phase 5: Deployment (Pending)
- Deploy database migrations
- Deploy backend services
- Deploy frontend
- Deploy mobile app
- Monitor and verify

## Conclusion

This design provides a comprehensive technical specification for migrating from InsForge to Supabase. The migration is straightforward because:

1. **Backend Already Implemented**: The user-service has complete Supabase integration
2. **Frontend Ready**: Supabase SDK is installed and client is configured
3. **Database Compatible**: PostgreSQL migrations work with Supabase
4. **API Contracts Maintained**: Backend abstracts auth provider from clients

The primary work involves:
- Renaming variables and references
- Mobile app Supabase SDK integration
- Testing and verification
- Documentation updates

The migration maintains backward compatibility and preserves all existing functionality while enabling email verification through Supabase.
