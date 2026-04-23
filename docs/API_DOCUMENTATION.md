# Tayosa Banking Ecosystem - API Documentation

## Authentication

The Tayosa banking ecosystem uses Supabase for authentication. All protected endpoints require a valid JWT token obtained through Supabase authentication.

### Authentication Flow

1. User registers or logs in through the frontend/mobile app
2. Supabase issues a JWT access token
3. Client includes the token in the `Authorization` header for all API requests
4. API Gateway validates the token with Supabase before routing to backend services

### Token Format

All authenticated requests must include a Bearer token in the Authorization header:

```
Authorization: Bearer <supabase_jwt_token>
```

### Token Validation

The API Gateway validates tokens by calling Supabase's `/auth/v1/user` endpoint. Valid tokens return user information, while invalid or expired tokens return 401 Unauthorized.

**Token Validation Endpoint (Supabase):**
- **URL:** `https://<project-ref>.supabase.co/auth/v1/user`
- **Method:** GET
- **Headers:**
  - `Authorization: Bearer <token>`
  - `apikey: <supabase_anon_key>`

**Successful Response (200 OK):**
```json
{
  "id": "uuid-string",
  "email": "user@example.com",
  "email_confirmed_at": "2024-01-01T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  ...
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "invalid or expired session"
}
```

### Session Management

- **Token Expiry:** Supabase JWT tokens expire after 1 hour by default
- **Refresh Tokens:** Use Supabase's refresh token mechanism to obtain new access tokens
- **Session Storage:** Store tokens securely (HttpOnly cookies for web, secure storage for mobile)

### Authentication Endpoints

All authentication endpoints are handled by the `user-service` and proxied through the API Gateway.

#### Register

Create a new user account.

- **URL:** `/api/v1/auth/register`
- **Method:** POST
- **Auth Required:** No
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Success Response (201 Created):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "supabaseUserId": "supabase-uuid",
    "emailVerified": false
  },
  "message": "Registration successful. Please verify your email."
}
```

#### Login

Authenticate an existing user.

- **URL:** `/api/v1/auth/login`
- **Method:** POST
- **Auth Required:** No
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Success Response (200 OK):**
```json
{
  "accessToken": "eyJhbGc...",
  "refreshToken": "refresh_token_string",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "supabaseUserId": "supabase-uuid"
  }
}
```

#### Logout

Invalidate the current session.

- **URL:** `/api/v1/auth/logout`
- **Method:** POST
- **Auth Required:** Yes
- **Request Body:** None

**Success Response (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

#### Refresh Token

Obtain a new access token using a refresh token.

- **URL:** `/api/v1/auth/refresh`
- **Method:** POST
- **Auth Required:** No
- **Request Body:**
```json
{
  "refreshToken": "refresh_token_string"
}
```

**Success Response (200 OK):**
```json
{
  "accessToken": "new_jwt_token",
  "refreshToken": "new_refresh_token"
}
```

#### Get User Profile

Retrieve the authenticated user's profile.

- **URL:** `/api/v1/auth/profile`
- **Method:** GET
- **Auth Required:** Yes

**Success Response (200 OK):**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "supabaseUserId": "supabase-uuid",
  "firstName": "John",
  "lastName": "Doe",
  "emailVerified": true,
  "createdAt": "2024-01-01T00:00:00Z"
}
```

### Password Reset Flow

#### Request Password Reset

Send a password reset email.

- **URL:** `/api/v1/auth/send-reset-password-email`
- **Method:** POST
- **Auth Required:** No
- **Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200 OK):**
```json
{
  "message": "Password reset email sent"
}
```

#### Reset Password

Complete the password reset using the token from the email.

- **URL:** `/api/v1/auth/reset-password`
- **Method:** POST
- **Auth Required:** No
- **Request Body:**
```json
{
  "token": "reset_token_from_email",
  "newPassword": "newSecurePassword123"
}
```

**Success Response (200 OK):**
```json
{
  "message": "Password reset successful"
}
```

### Email Verification

#### Resend Verification Email

- **URL:** `/api/v1/auth/resend-verification-email`
- **Method:** POST
- **Auth Required:** No
- **Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200 OK):**
```json
{
  "message": "Verification email sent"
}
```

#### Verify Email

- **URL:** `/api/v1/auth/verify-email`
- **Method:** POST
- **Auth Required:** No
- **Request Body:**
```json
{
  "token": "verification_token_from_email"
}
```

**Success Response (200 OK):**
```json
{
  "message": "Email verified successfully"
}
```

## Error Responses

All endpoints follow a consistent error response format:

### 400 Bad Request
```json
{
  "error": "Invalid request parameters",
  "details": "Email is required"
}
```

### 401 Unauthorized
```json
{
  "error": "missing bearer token"
}
```

or

```json
{
  "error": "invalid or expired session"
}
```

### 403 Forbidden
```json
{
  "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "message": "An unexpected error occurred"
}
```

## Protected Endpoints

All endpoints except authentication endpoints require a valid Bearer token. The API Gateway validates tokens before routing requests to backend services.

### Common Headers

**Required for all protected endpoints:**
- `Authorization: Bearer <token>`
- `Content-Type: application/json`

**Optional:**
- `X-Request-Id: <unique-request-id>` - For request tracing

### Backend Services

The following services require authentication:

- **User Service** - User management and profiles
- **Affiliate Service** - Referral and rewards management
- **Audit Log Service** - Audit event tracking
- **Fee Service** - Transaction fee management
- **Kibiina Service** - Group savings management
- **Loan/Credit Service** - Credit scoring and loan eligibility
- **Notification Service** - Notification delivery
- **Object Storage Service** - File upload and storage

## Rate Limiting

Rate limiting is applied at the API Gateway level:

- **Authenticated requests:** 1000 requests per hour per user
- **Unauthenticated requests:** 100 requests per hour per IP

Rate limit headers are included in all responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## CORS

The API Gateway allows cross-origin requests from all origins with the following configuration:

- **Allowed Origins:** `*`
- **Allowed Methods:** GET, POST, PUT, PATCH, DELETE, OPTIONS
- **Allowed Headers:** Content-Type, Authorization, X-User-Id, X-Request-Id, X-CSRF-Token, X-Admin-Secret

## Versioning

The API uses URL-based versioning:

- Current version: `v1`
- Base URL: `https://api.tayosa.com/api/v1`

## Support

For API support and questions:
- Email: api-support@tayosa.com
- Documentation: https://docs.tayosa.com
