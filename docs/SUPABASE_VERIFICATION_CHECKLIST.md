# Supabase Project Verification Checklist

This checklist ensures your Supabase project is properly configured before deployment.

## Project Information

- **Project URL:** `https://[YOUR-PROJECT-REF].supabase.co`
- **Project Reference:** `[YOUR-PROJECT-REF]`
- **Region:** US East
- **Database Host:** `db.[YOUR-PROJECT-REF].supabase.co`

## ✅ Verification Steps

### 1. API Keys Configuration

- [ ] **Anon Key (Publishable)** is available and documented
  - Key: `[YOUR_ANON_KEY]`
  - Used by: Frontend and mobile applications
  - Status: ✅ Available

- [ ] **Service Role Key (Secret)** is available and secured
  - Key: `[YOUR_SERVICE_ROLE_KEY]`
  - Used by: Backend services
  - Status: ✅ Available
  - ⚠️ **Security:** Never expose in client-side code

### 2. Database Connection

- [ ] **Direct Connection String** is accessible
  ```
  postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
  ```
  - Port: 5432
  - Database: postgres
  - User: postgres
  - Status: ✅ Configured

- [ ] **Connection Pooler** is available (port 6543)
  ```
  postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:6543/postgres?pgbouncer=true
  ```
  - Recommended for application connections
  - Status: ✅ Available

### 3. Row Level Security (RLS)

- [ ] **Automatic RLS** is enabled
  - Status: ✅ Enabled (as noted in connection details)
  - All tables should have appropriate RLS policies

- [ ] **Verify RLS Policies** on critical tables:
  - [ ] `users` table
  - [ ] `transactions` table
  - [ ] `accounts` table
  - [ ] Other sensitive tables

### 4. Authentication Configuration

- [ ] **Email Provider** is enabled
  - Navigate to: Authentication > Providers > Email
  - Status: ⏳ Verify in dashboard

- [ ] **Email Confirmations** are enabled
  - Required for user verification
  - Status: ⏳ Verify in dashboard

- [ ] **Password Requirements** are configured
  - Minimum length: 8 characters
  - Status: ⏳ Verify in dashboard

### 5. Email Templates

- [ ] **Email Verification Template** is customized
  - Template: Confirm signup
  - Status: ⏳ Customize in dashboard

- [ ] **Password Reset Template** is customized
  - Template: Reset password
  - Status: ⏳ Customize in dashboard

- [ ] **Test Email Delivery**
  - Send test verification email
  - Send test password reset email
  - Status: ⏳ Test required

### 6. Redirect URLs

- [ ] **Site URL** is configured
  - Production: `https://app.tayosa.com`
  - Status: ⏳ Configure in dashboard

- [ ] **Allowed Redirect URLs** are whitelisted
  - [ ] `https://app.tayosa.com/auth/callback`
  - [ ] `https://app.tayosa.com/auth/verify-email`
  - [ ] `https://app.tayosa.com/auth/reset-password`
  - [ ] Development URLs (localhost)
  - [ ] Mobile deep links (tayosa://)
  - Status: ⏳ Configure in dashboard

### 7. Storage Configuration

- [ ] **Storage Bucket** exists
  - Bucket name: `collateral_docs`
  - Public: No (private)
  - Status: ⏳ Verify/create in dashboard

- [ ] **Storage Policies** are configured
  - Upload policy for authenticated users
  - Read policy for file owners
  - Status: ⏳ Configure in dashboard

### 8. SMTP Configuration (Production)

- [ ] **Custom SMTP Provider** is configured (optional but recommended)
  - Provider: SendGrid/AWS SES/Mailgun
  - Status: ⏳ Configure for production

- [ ] **Sender Email** is verified
  - Email: `noreply@tayosa.com`
  - Status: ⏳ Verify with SMTP provider

- [ ] **Test Email Delivery** from custom SMTP
  - Send test emails
  - Check spam folder
  - Status: ⏳ Test required

### 9. Environment Variables

- [ ] **Backend Services** have Supabase credentials
  - [ ] `SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co`
  - [ ] `SUPABASE_ANON_KEY=[YOUR_ANON_KEY]`
  - [ ] `SUPABASE_SERVICE_ROLE_KEY=[YOUR_SERVICE_ROLE_KEY]`
  - Status: ⏳ Configure in deployment

- [ ] **Frontend Application** has Supabase credentials
  - [ ] `VITE_SUPABASE_URL`
  - [ ] `VITE_SUPABASE_ANON_KEY`
  - Status: ⏳ Configure in deployment

- [ ] **Mobile Application** has Supabase credentials
  - [ ] `--dart-define=SUPABASE_URL=...`
  - [ ] `--dart-define=SUPABASE_ANON_KEY=...`
  - Status: ⏳ Configure in build scripts

### 10. Database Migration

- [ ] **Migration File** is ready
  - File: `db/migrations/011_rename_insforge_to_supabase.sql`
  - Status: ✅ Created

- [ ] **Backup Plan** is documented
  - Backup existing data before migration
  - Status: ⏳ See migration runbook

- [ ] **Rollback Plan** is prepared
  - SQL statements to revert changes
  - Status: ⏳ See migration runbook

### 11. Security Configuration

- [ ] **API Rate Limiting** is configured
  - Anonymous: 100 requests/hour
  - Authenticated: 1000 requests/hour
  - Status: ⏳ Configure in dashboard

- [ ] **CORS Settings** are configured
  - Allowed origins: Production and development URLs
  - Status: ⏳ Configure in dashboard

- [ ] **Database Backups** are enabled
  - Frequency: Daily
  - Retention: 7-30 days
  - Status: ⏳ Verify in dashboard

### 12. Monitoring and Logs

- [ ] **Authentication Logs** are accessible
  - Location: Authentication > Logs
  - Status: ⏳ Verify access

- [ ] **Database Logs** are accessible
  - Location: Database > Logs
  - Status: ⏳ Verify access

- [ ] **API Logs** are accessible
  - Location: Logs > API
  - Status: ⏳ Verify access

## Testing Checklist

### Authentication Flow Testing

- [ ] **User Registration**
  - [ ] Create new user account
  - [ ] Verify email is sent
  - [ ] Confirm email verification works
  - Status: ⏳ Test required

- [ ] **User Login**
  - [ ] Login with verified account
  - [ ] Receive valid JWT token
  - [ ] Token validates successfully
  - Status: ⏳ Test required

- [ ] **Password Reset**
  - [ ] Request password reset
  - [ ] Receive reset email
  - [ ] Complete password reset
  - [ ] Login with new password
  - Status: ⏳ Test required

- [ ] **Token Validation**
  - [ ] Valid token returns user data
  - [ ] Expired token returns 401
  - [ ] Invalid token returns 401
  - Status: ⏳ Test required

### API Gateway Testing

- [ ] **Protected Endpoints**
  - [ ] Request without token returns 401
  - [ ] Request with valid token succeeds
  - [ ] Request with expired token returns 401
  - Status: ⏳ Test required

- [ ] **Token Forwarding**
  - [ ] Authorization header forwarded to services
  - [ ] X-User-Id header set correctly
  - Status: ⏳ Test required

### Backend Services Testing

- [ ] **All Services Validate Tokens**
  - [ ] api-gateway-service
  - [ ] affiliate-service
  - [ ] audit-log-service
  - [ ] fee-service
  - [ ] kibiina-service
  - [ ] loan-credit-service
  - [ ] notification-service
  - [ ] object-storage-service
  - Status: ⏳ Test required

### Mobile App Testing

- [ ] **Flutter App Authentication**
  - [ ] Register new user
  - [ ] Login existing user
  - [ ] Password reset flow
  - [ ] Token refresh works
  - Status: ⏳ Test required

- [ ] **API Calls from Mobile**
  - [ ] Authenticated requests succeed
  - [ ] Token automatically attached
  - [ ] Error handling works
  - Status: ⏳ Test required

## Pre-Deployment Verification

### Critical Checks

- [ ] All environment variables are configured
- [ ] Database migration is tested in staging
- [ ] Email delivery is working
- [ ] Authentication flow is tested end-to-end
- [ ] All backend services are updated
- [ ] Mobile app is updated and tested
- [ ] Documentation is complete
- [ ] Rollback plan is ready

### Performance Checks

- [ ] Token validation latency is acceptable (<100ms)
- [ ] Database queries are optimized
- [ ] Connection pooling is configured
- [ ] Rate limiting is appropriate

### Security Checks

- [ ] Service role key is not exposed
- [ ] RLS policies are enabled on all tables
- [ ] CORS is properly configured
- [ ] HTTPS is enforced
- [ ] Sensitive data is encrypted

## Post-Deployment Verification

### Immediate Checks (First Hour)

- [ ] Monitor authentication logs for errors
- [ ] Check API error rates
- [ ] Verify email delivery
- [ ] Monitor database connections
- [ ] Check application logs

### 24-Hour Checks

- [ ] Review authentication success rate
- [ ] Check for token validation errors
- [ ] Monitor email delivery rate
- [ ] Review user feedback
- [ ] Check system performance

### Week 1 Checks

- [ ] Analyze authentication patterns
- [ ] Review security logs
- [ ] Check database performance
- [ ] Monitor storage usage
- [ ] Review error rates

## Troubleshooting Guide

### Common Issues

**Issue: Email not delivered**
- Check SMTP configuration
- Verify sender email is verified
- Check spam folder
- Review authentication logs

**Issue: Token validation fails**
- Verify SUPABASE_URL is correct
- Check SUPABASE_ANON_KEY is valid
- Ensure token is not expired
- Check RLS policies

**Issue: Database connection fails**
- Verify connection string
- Check IP allowlist
- Test connection pooler
- Review connection limits

**Issue: RLS policy blocks request**
- Review policy conditions
- Check user authentication
- Verify auth.uid() matches
- Test with service role key

## Support Contacts

**Supabase Support:**
- Documentation: https://supabase.com/docs
- Community: https://github.com/supabase/supabase/discussions
- Support: https://supabase.com/support

**Tayosa Team:**
- Email: support@tayosa.com
- Documentation: https://docs.tayosa.com

## Sign-Off

**Verified by:** _________________  
**Date:** _________________  
**Environment:** [ ] Staging [ ] Production  
**Ready for Deployment:** [ ] Yes [ ] No  

**Notes:**
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
