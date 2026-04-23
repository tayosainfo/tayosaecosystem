# Supabase Project Configuration Guide

This document provides detailed instructions for configuring your Supabase project for the Tayosa banking ecosystem.

## Project Information

- **Project URL:** `https://[YOUR-PROJECT-REF].supabase.co`
- **Project Reference:** `[YOUR-PROJECT-REF]`
- **Region:** US East

## Authentication Configuration

### Email Authentication Settings

1. Navigate to **Authentication > Providers** in the Supabase dashboard
2. Enable **Email** provider
3. Configure the following settings:

#### Email Confirmation

- **Enable email confirmations:** ✅ Enabled
- **Confirmation URL:** `https://app.tayosa.com/auth/verify-email`
- **Email confirmation expiry:** 24 hours (default)

**Purpose:** Users must verify their email address before accessing protected features.

#### Password Requirements

- **Minimum password length:** 8 characters
- **Password strength:** Medium (recommended)
- **Allow password resets:** ✅ Enabled

### Email Templates

Customize email templates to match your brand identity.

#### Email Verification Template

Navigate to **Authentication > Email Templates > Confirm signup**

**Subject:** Verify your Tayosa account

**Template:**
```html
<h2>Welcome to Tayosa!</h2>
<p>Thank you for signing up. Please verify your email address by clicking the link below:</p>
<p><a href="{{ .ConfirmationURL }}">Verify Email Address</a></p>
<p>This link will expire in 24 hours.</p>
<p>If you didn't create an account with Tayosa, you can safely ignore this email.</p>
```

**Variables available:**
- `{{ .ConfirmationURL }}` - Email verification link
- `{{ .Email }}` - User's email address
- `{{ .Token }}` - Verification token

#### Password Reset Template

Navigate to **Authentication > Email Templates > Reset password**

**Subject:** Reset your Tayosa password

**Template:**
```html
<h2>Password Reset Request</h2>
<p>We received a request to reset your password for your Tayosa account.</p>
<p>Click the link below to reset your password:</p>
<p><a href="{{ .ConfirmationURL }}">Reset Password</a></p>
<p>This link will expire in 1 hour.</p>
<p>If you didn't request a password reset, you can safely ignore this email.</p>
```

**Variables available:**
- `{{ .ConfirmationURL }}` - Password reset link
- `{{ .Email }}` - User's email address
- `{{ .Token }}` - Reset token

#### Magic Link Template (Optional)

Navigate to **Authentication > Email Templates > Magic Link**

**Subject:** Sign in to Tayosa

**Template:**
```html
<h2>Sign in to Tayosa</h2>
<p>Click the link below to sign in to your account:</p>
<p><a href="{{ .ConfirmationURL }}">Sign In</a></p>
<p>This link will expire in 1 hour.</p>
<p>If you didn't request this link, you can safely ignore this email.</p>
```

### Redirect URLs

Configure allowed redirect URLs for email links and OAuth callbacks.

Navigate to **Authentication > URL Configuration**

**Site URL:** `https://app.tayosa.com`

**Redirect URLs (whitelist):**
```
https://app.tayosa.com/auth/callback
https://app.tayosa.com/auth/verify-email
https://app.tayosa.com/auth/reset-password
http://localhost:3000/auth/callback
http://localhost:3000/auth/verify-email
http://localhost:3000/auth/reset-password
```

**Mobile Deep Links (for Flutter app):**
```
tayosa://auth/callback
tayosa://auth/verify-email
tayosa://auth/reset-password
```

### Session Configuration

Navigate to **Authentication > Settings**

**JWT expiry:** 3600 seconds (1 hour)
**Refresh token expiry:** 2592000 seconds (30 days)
**Enable refresh token rotation:** ✅ Enabled

## OAuth Providers (Optional)

If you plan to support social login, configure OAuth providers:

### Google OAuth

1. Navigate to **Authentication > Providers**
2. Enable **Google** provider
3. Enter your Google OAuth credentials:
   - **Client ID:** `<your-google-client-id>`
   - **Client Secret:** `<your-google-client-secret>`
4. Add authorized redirect URI in Google Cloud Console:
   - `https://[YOUR-PROJECT-REF].supabase.co/auth/v1/callback`

### GitHub OAuth

1. Navigate to **Authentication > Providers**
2. Enable **GitHub** provider
3. Enter your GitHub OAuth credentials:
   - **Client ID:** `<your-github-client-id>`
   - **Client Secret:** `<your-github-client-secret>`
4. Add authorized callback URL in GitHub OAuth App settings:
   - `https://[YOUR-PROJECT-REF].supabase.co/auth/v1/callback`

## Row Level Security (RLS) Policies

Configure RLS policies to secure your database tables.

### Users Table

```sql
-- Enable RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Policy: Users can read their own data
CREATE POLICY "Users can view own profile"
ON users FOR SELECT
USING (auth.uid() = supabase_user_id);

-- Policy: Users can update their own data
CREATE POLICY "Users can update own profile"
ON users FOR UPDATE
USING (auth.uid() = supabase_user_id);

-- Policy: Service role can do anything
CREATE POLICY "Service role has full access"
ON users
USING (auth.role() = 'service_role');
```

### Other Tables

Apply similar RLS policies to other tables based on your security requirements:

```sql
-- Example: Transactions table
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view own transactions"
ON transactions FOR SELECT
USING (auth.uid() = (SELECT supabase_user_id FROM users WHERE id = user_id));
```

## Storage Configuration

### Create Storage Buckets

Navigate to **Storage** in the Supabase dashboard

#### Collateral Documents Bucket

1. Click **New bucket**
2. **Name:** `collateral_docs`
3. **Public:** ❌ Disabled (private bucket)
4. **File size limit:** 50 MB
5. **Allowed MIME types:** 
   - `image/jpeg`
   - `image/png`
   - `image/gif`
   - `application/pdf`
   - `application/msword`
   - `application/vnd.openxmlformats-officedocument.wordprocessingml.document`

#### Storage Policies

Configure storage policies for the bucket:

```sql
-- Policy: Authenticated users can upload files
CREATE POLICY "Authenticated users can upload"
ON storage.objects FOR INSERT
TO authenticated
WITH CHECK (bucket_id = 'collateral_docs');

-- Policy: Users can view their own files
CREATE POLICY "Users can view own files"
ON storage.objects FOR SELECT
TO authenticated
USING (bucket_id = 'collateral_docs' AND auth.uid()::text = (storage.foldername(name))[1]);

-- Policy: Service role has full access
CREATE POLICY "Service role has full access"
ON storage.objects
TO service_role
USING (bucket_id = 'collateral_docs');
```

## SMTP Configuration (Optional)

For production environments, configure a custom SMTP provider for better email deliverability.

Navigate to **Project Settings > Auth > SMTP Settings**

### Recommended Providers

- **SendGrid**
- **AWS SES**
- **Mailgun**
- **Postmark**

### Configuration Example (SendGrid)

```
SMTP Host: smtp.sendgrid.net
SMTP Port: 587
SMTP User: apikey
SMTP Password: <your-sendgrid-api-key>
Sender Email: noreply@tayosa.com
Sender Name: Tayosa Banking
```

**Enable SMTP:** ✅ Enabled

## API Keys and Secrets

### Anon Key (Public)

Used by client applications (web and mobile):

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Usage:**
- Frontend applications
- Mobile applications
- Public API calls

**Security:** Safe to expose in client-side code

### Service Role Key (Secret)

Used by backend services for privileged operations:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Usage:**
- Backend services
- Admin operations
- Bypassing RLS policies

**Security:** ⚠️ NEVER expose in client-side code or public repositories

### JWT Secret

Used to sign and verify JWT tokens:

```
<your-jwt-secret>
```

**Security:** ⚠️ Keep this secret secure and never expose it

## Database Connection

### Connection String

**Direct connection (for migrations):**
```
postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
```

**Pooler connection (for applications):**
```
postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:6543/postgres?pgbouncer=true
```

### Connection Pooling

Supabase uses PgBouncer for connection pooling:

- **Pool mode:** Transaction
- **Max connections:** 15 (Free tier), 200+ (Pro tier)
- **Default pool size:** 15

## Security Best Practices

### 1. Environment Variables

Store all sensitive credentials in environment variables:

```bash
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=<anon-key>
SUPABASE_SERVICE_ROLE_KEY=<service-role-key>
```

### 2. RLS Policies

Always enable RLS on tables containing user data:

```sql
ALTER TABLE <table_name> ENABLE ROW LEVEL SECURITY;
```

### 3. API Rate Limiting

Configure rate limiting in **Project Settings > API**:

- **Anonymous requests:** 100 per hour
- **Authenticated requests:** 1000 per hour

### 4. CORS Configuration

Configure allowed origins in **Project Settings > API > CORS**:

```
https://app.tayosa.com
http://localhost:3000
```

### 5. Database Backups

Enable automatic backups in **Project Settings > Database > Backups**:

- **Backup frequency:** Daily
- **Retention period:** 7 days (Free tier), 30 days (Pro tier)

## Monitoring and Logs

### Authentication Logs

View authentication events in **Authentication > Logs**:

- Sign-ups
- Sign-ins
- Password resets
- Email verifications

### Database Logs

View database queries in **Database > Logs**:

- Slow queries
- Error logs
- Connection logs

### API Logs

View API requests in **Logs > API**:

- Request count
- Response times
- Error rates

## Testing Configuration

### Test Email Addresses

For development and testing, use Supabase's test email feature:

Navigate to **Authentication > Settings > Test Email**

**Enable test emails:** ✅ Enabled

This allows you to see verification and reset emails in the dashboard without sending actual emails.

### Test Users

Create test users for development:

```sql
-- Insert test user
INSERT INTO auth.users (
  id,
  email,
  encrypted_password,
  email_confirmed_at,
  created_at,
  updated_at
) VALUES (
  gen_random_uuid(),
  'test@tayosa.com',
  crypt('testpassword123', gen_salt('bf')),
  now(),
  now(),
  now()
);
```

## Troubleshooting

### Email Delivery Issues

1. Check SMTP configuration in **Project Settings > Auth > SMTP**
2. Verify sender email is verified with your SMTP provider
3. Check spam folder for test emails
4. Review authentication logs for email sending errors

### Token Validation Failures

1. Verify `SUPABASE_URL` and `SUPABASE_ANON_KEY` are correct
2. Check token expiry (default 1 hour)
3. Ensure client is sending token in `Authorization: Bearer <token>` format
4. Verify RLS policies are not blocking requests

### Connection Issues

1. Check database connection string
2. Verify IP allowlist in **Project Settings > Database > Connection Pooling**
3. Test connection using `psql` or database client
4. Check connection pool limits

## Support

For Supabase-specific issues:
- Documentation: https://supabase.com/docs
- Community: https://github.com/supabase/supabase/discussions
- Support: https://supabase.com/support

For Tayosa-specific issues:
- Email: support@tayosa.com
- Documentation: https://docs.tayosa.com
