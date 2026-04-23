# Supabase Authentication Settings Configuration

This document provides step-by-step instructions for configuring authentication settings in your Supabase project.

## Access Authentication Settings

1. Navigate to: https://supabase.com/dashboard/project/[YOUR-PROJECT-REF]
2. Click on **Authentication** in the left sidebar
3. Click on **Settings** tab

## General Settings

### Site URL

The Site URL is the base URL of your application where users will be redirected after authentication.

**Configuration:**
1. Navigate to **Authentication > URL Configuration**
2. Set **Site URL**:
   - **Production:** `https://app.tayosa.com`
   - **Staging:** `https://staging.tayosa.com`
   - **Development:** `http://localhost:3000`

**Purpose:** This is the default redirect URL for email confirmations and password resets.

### Redirect URLs

Whitelist all URLs where users can be redirected after authentication actions.

**Configuration:**
1. Navigate to **Authentication > URL Configuration**
2. Add **Redirect URLs** (one per line):

```
https://app.tayosa.com/auth/callback
https://app.tayosa.com/auth/verify-email
https://app.tayosa.com/auth/reset-password
https://app.tayosa.com/auth/magic-link
https://staging.tayosa.com/auth/callback
https://staging.tayosa.com/auth/verify-email
https://staging.tayosa.com/auth/reset-password
http://localhost:3000/auth/callback
http://localhost:3000/auth/verify-email
http://localhost:3000/auth/reset-password
tayosa://auth/callback
tayosa://auth/verify-email
tayosa://auth/reset-password
```

**Mobile Deep Links:**
- `tayosa://` - Custom URL scheme for mobile app
- Configure in your Flutter app's `AndroidManifest.xml` and `Info.plist`

### JWT Settings

Configure JSON Web Token settings for session management.

**Configuration:**
1. Navigate to **Authentication > Settings**
2. Configure **JWT Settings**:

**JWT Expiry:**
- **Value:** `3600` seconds (1 hour)
- **Purpose:** How long access tokens remain valid
- **Recommendation:** 1 hour for security, use refresh tokens for longer sessions

**Refresh Token Expiry:**
- **Value:** `2592000` seconds (30 days)
- **Purpose:** How long refresh tokens remain valid
- **Recommendation:** 30 days for good balance of security and UX

**Enable Refresh Token Rotation:**
- **Value:** ✅ Enabled
- **Purpose:** Issues new refresh token on each refresh, invalidating old one
- **Recommendation:** Enable for enhanced security

**Reuse Interval:**
- **Value:** `10` seconds
- **Purpose:** Grace period for refresh token reuse (prevents race conditions)
- **Recommendation:** 10 seconds is sufficient

## Email Provider Configuration

### Enable Email Authentication

**Configuration:**
1. Navigate to **Authentication > Providers**
2. Find **Email** provider
3. Toggle **Enable Email provider** to ON

### Email Confirmation Settings

**Configuration:**
1. Navigate to **Authentication > Providers > Email**
2. Configure **Email Confirmation**:

**Enable email confirmations:**
- **Value:** ✅ Enabled
- **Purpose:** Require users to verify email before accessing protected features
- **Recommendation:** Always enable for production

**Secure email change:**
- **Value:** ✅ Enabled
- **Purpose:** Require confirmation for email address changes
- **Recommendation:** Enable to prevent account takeover

**Double confirm email changes:**
- **Value:** ✅ Enabled
- **Purpose:** Send confirmation to both old and new email addresses
- **Recommendation:** Enable for maximum security

### Password Requirements

**Configuration:**
1. Navigate to **Authentication > Providers > Email**
2. Configure **Password Requirements**:

**Minimum password length:**
- **Value:** `8` characters
- **Recommendation:** 8-12 characters minimum

**Password strength:**
- **Value:** Medium
- **Options:** Weak, Medium, Strong
- **Recommendation:** Medium for balance of security and usability

**Allow password resets:**
- **Value:** ✅ Enabled
- **Purpose:** Enable password reset functionality
- **Recommendation:** Always enable

## Session Management

### Session Timeout

**Configuration:**
1. Navigate to **Authentication > Settings**
2. Configure **Session Settings**:

**Inactivity timeout:**
- **Value:** `0` (disabled) or `3600` seconds (1 hour)
- **Purpose:** Automatically log out inactive users
- **Recommendation:** Enable for sensitive applications

**Absolute timeout:**
- **Value:** `0` (disabled) or `86400` seconds (24 hours)
- **Purpose:** Force re-authentication after specified time
- **Recommendation:** Enable for high-security applications

### Multi-Factor Authentication (MFA)

**Configuration:**
1. Navigate to **Authentication > Settings**
2. Configure **MFA Settings**:

**Enable MFA:**
- **Value:** ⏳ Optional (can be enabled later)
- **Purpose:** Add second factor authentication
- **Recommendation:** Enable for admin accounts

**MFA Factors:**
- TOTP (Time-based One-Time Password)
- SMS (requires Twilio integration)

## OAuth Providers (Optional)

### Google OAuth

**Configuration:**
1. Navigate to **Authentication > Providers**
2. Find **Google** provider
3. Toggle to enable
4. Enter credentials:

**Client ID:**
```
<your-google-client-id>.apps.googleusercontent.com
```

**Client Secret:**
```
<your-google-client-secret>
```

**Authorized redirect URI (add to Google Cloud Console):**
```
https://[YOUR-PROJECT-REF].supabase.co/auth/v1/callback
```

**Scopes:**
- `email`
- `profile`

### GitHub OAuth

**Configuration:**
1. Navigate to **Authentication > Providers**
2. Find **GitHub** provider
3. Toggle to enable
4. Enter credentials:

**Client ID:**
```
<your-github-client-id>
```

**Client Secret:**
```
<your-github-client-secret>
```

**Authorization callback URL (add to GitHub OAuth App):**
```
https://[YOUR-PROJECT-REF].supabase.co/auth/v1/callback
```

## Rate Limiting

Protect your authentication endpoints from abuse.

**Configuration:**
1. Navigate to **Authentication > Settings**
2. Configure **Rate Limiting**:

**Anonymous requests:**
- **Value:** `100` requests per hour
- **Purpose:** Limit unauthenticated requests
- **Recommendation:** 100-200 for public endpoints

**Authenticated requests:**
- **Value:** `1000` requests per hour
- **Purpose:** Limit authenticated requests
- **Recommendation:** 1000-5000 based on usage

**Email sending rate:**
- **Value:** `4` emails per hour per user
- **Purpose:** Prevent email spam
- **Recommendation:** 4-6 emails per hour

## Security Settings

### Password Policy

**Configuration:**
1. Navigate to **Authentication > Settings**
2. Configure **Password Policy**:

**Minimum length:** 8 characters
**Require uppercase:** ❌ Optional
**Require lowercase:** ❌ Optional
**Require numbers:** ❌ Optional
**Require special characters:** ❌ Optional

**Recommendation:** Start with length requirement only, add complexity as needed.

### Account Lockout

**Configuration:**
1. Navigate to **Authentication > Settings**
2. Configure **Account Lockout**:

**Enable account lockout:**
- **Value:** ✅ Enabled
- **Purpose:** Lock account after failed login attempts
- **Recommendation:** Enable for security

**Failed attempts threshold:**
- **Value:** `5` attempts
- **Purpose:** Number of failed attempts before lockout
- **Recommendation:** 5-10 attempts

**Lockout duration:**
- **Value:** `600` seconds (10 minutes)
- **Purpose:** How long account remains locked
- **Recommendation:** 10-30 minutes

### CAPTCHA Protection

**Configuration:**
1. Navigate to **Authentication > Settings**
2. Configure **CAPTCHA**:

**Enable CAPTCHA:**
- **Value:** ⏳ Optional (recommended for production)
- **Purpose:** Prevent automated attacks
- **Recommendation:** Enable for production

**CAPTCHA Provider:**
- hCaptcha (recommended)
- reCAPTCHA v2
- reCAPTCHA v3

**Site Key:**
```
<your-captcha-site-key>
```

**Secret Key:**
```
<your-captcha-secret-key>
```

## Advanced Settings

### Custom Claims

Add custom data to JWT tokens.

**Configuration:**
1. Navigate to **Authentication > Settings**
2. Configure **Custom Claims**:

**Example custom claims function:**
```sql
CREATE OR REPLACE FUNCTION auth.custom_access_token_hook(event jsonb)
RETURNS jsonb
LANGUAGE plpgsql
AS $$
DECLARE
  claims jsonb;
  user_role text;
BEGIN
  -- Get user role from users table
  SELECT role INTO user_role
  FROM public.users
  WHERE supabase_user_id = (event->>'user_id')::uuid;

  -- Add custom claims
  claims := event->'claims';
  claims := jsonb_set(claims, '{user_role}', to_jsonb(user_role));
  
  -- Return modified event
  RETURN jsonb_set(event, '{claims}', claims);
END;
$$;
```

### Hooks

Configure authentication hooks for custom logic.

**Available Hooks:**
- `auth.custom_access_token_hook` - Modify JWT claims
- `auth.send_email_hook` - Customize email sending
- `auth.send_sms_hook` - Customize SMS sending

**Configuration:**
1. Navigate to **Database > Functions**
2. Create hook function
3. Navigate to **Authentication > Hooks**
4. Enable and configure hook

## Testing Configuration

### Test Email Flow

```bash
# Test registration
curl -X POST https://[YOUR-PROJECT-REF].supabase.co/auth/v1/signup \
  -H "apikey: [YOUR_ANON_KEY]" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123"
  }'

# Test password reset
curl -X POST https://[YOUR-PROJECT-REF].supabase.co/auth/v1/recover \
  -H "apikey: [YOUR_ANON_KEY]" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'
```

### Test OAuth Flow

1. Navigate to your application
2. Click "Sign in with Google" (or other provider)
3. Complete OAuth flow
4. Verify redirect to correct URL
5. Check user is created in Supabase

### Test Rate Limiting

```bash
# Send multiple requests to test rate limiting
for i in {1..10}; do
  curl -X POST https://[YOUR-PROJECT-REF].supabase.co/auth/v1/signup \
    -H "apikey: [YOUR_ANON_KEY]" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"test$i@example.com\",\"password\":\"test123\"}"
done
```

## Configuration Checklist

### Essential Settings

- [ ] Site URL configured
- [ ] Redirect URLs whitelisted
- [ ] Email provider enabled
- [ ] Email confirmations enabled
- [ ] Password requirements set
- [ ] JWT expiry configured
- [ ] Refresh token rotation enabled

### Security Settings

- [ ] Rate limiting configured
- [ ] Account lockout enabled
- [ ] Password policy defined
- [ ] CAPTCHA enabled (production)
- [ ] Secure email change enabled

### Email Settings

- [ ] Email templates customized
- [ ] SMTP configured (production)
- [ ] Sender email verified
- [ ] Test emails sent successfully

### Optional Settings

- [ ] OAuth providers configured
- [ ] MFA enabled (if required)
- [ ] Custom claims configured
- [ ] Hooks configured (if needed)

## Monitoring

### Authentication Logs

View authentication events:
1. Navigate to **Authentication > Logs**
2. Filter by event type:
   - Sign ups
   - Sign ins
   - Password resets
   - Email verifications
   - Failed attempts

### Metrics

Monitor authentication metrics:
1. Navigate to **Authentication > Metrics**
2. View:
   - Daily active users
   - Sign-up rate
   - Sign-in success rate
   - Failed login attempts
   - Email delivery rate

## Troubleshooting

### Users Not Receiving Emails

1. Check SMTP configuration
2. Verify sender email is verified
3. Check spam folder
4. Review authentication logs
5. Test with different email providers

### Redirect URL Errors

1. Verify URL is whitelisted
2. Check Site URL configuration
3. Ensure URL format is correct
4. Test with different browsers

### OAuth Not Working

1. Verify OAuth credentials
2. Check redirect URI in provider settings
3. Ensure provider is enabled in Supabase
4. Test OAuth flow in incognito mode

### Rate Limiting Issues

1. Review rate limit settings
2. Check if legitimate traffic is blocked
3. Adjust limits based on usage patterns
4. Implement exponential backoff in client

## Support

For authentication configuration issues:
- Supabase Docs: https://supabase.com/docs/guides/auth
- Community: https://github.com/supabase/supabase/discussions
- Support: support@tayosa.com

## Sign-Off

**Configuration Completed by:**
- Engineer: _________________ Date: _______

**Settings Verified:**
- [ ] All essential settings configured
- [ ] Security settings enabled
- [ ] Email flow tested
- [ ] OAuth tested (if enabled)
- [ ] Rate limiting tested
- [ ] Documentation updated

**Notes:**
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
