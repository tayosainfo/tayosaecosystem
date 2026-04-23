# Pre-Deployment Verification Checklist

This document provides a comprehensive checklist to verify that the Tayosa banking ecosystem is ready for deployment with Supabase integration.

**Date:** _________________  
**Verified by:** _________________  
**Environment:** [ ] Staging [ ] Production

## 1. Code Changes Verification

### Backend Services

- [ ] **api-gateway-service**
  - [ ] Supabase token validation implemented
  - [ ] `requireAuth` middleware applied to protected routes
  - [ ] Authorization header forwarded to downstream services
  - [ ] Environment variables referenced correctly

- [ ] **user-service**
  - [ ] Struct fields renamed (SupabaseUserID, SupabaseLoginEmail)
  - [ ] Database queries updated with new column names
  - [ ] Comments updated to reference Supabase
  - [ ] Supabase client integration working

- [ ] **affiliate-service**
  - [ ] Supabase token validation implemented
  - [ ] Protected endpoints secured with `requireAuth`
  - [ ] Environment variables configured

- [ ] **audit-log-service**
  - [ ] Supabase token validation implemented
  - [ ] Protected endpoints secured with `requireAuth`
  - [ ] Environment variables configured

- [ ] **fee-service**
  - [ ] Supabase token validation implemented
  - [ ] Protected endpoints secured with `requireAuth`
  - [ ] Environment variables configured

- [ ] **kibiina-service**
  - [ ] Supabase token validation implemented
  - [ ] Protected endpoints secured with `requireAuth`
  - [ ] Environment variables configured

- [ ] **loan-credit-service**
  - [ ] Supabase token validation implemented
  - [ ] Protected endpoints secured with `requireAuth`
  - [ ] Environment variables configured

- [ ] **notification-service**
  - [ ] Supabase token validation implemented
  - [ ] Protected endpoints secured with `requireAuth`
  - [ ] Environment variables configured

- [ ] **object-storage-service**
  - [ ] Updated from InsForge to Supabase references
  - [ ] Supabase token validation implemented
  - [ ] Storage bucket configuration updated
  - [ ] Environment variables configured

### Frontend Application

- [ ] **Supabase client file**
  - [ ] File renamed from `insforge.ts` to `supabase.ts`
  - [ ] All imports updated
  - [ ] Comments updated to reference Supabase

- [ ] **Environment variables**
  - [ ] `VITE_SUPABASE_URL` configured
  - [ ] `VITE_SUPABASE_ANON_KEY` configured
  - [ ] API base URL configured

### Mobile Application

- [ ] **Supabase SDK integration**
  - [ ] `supabase_flutter` dependency added
  - [ ] Supabase client configuration created
  - [ ] Main.dart updated to initialize Supabase
  - [ ] API client updated with Supabase authentication

- [ ] **Authentication screens**
  - [ ] Login screen updated
  - [ ] Register screen updated
  - [ ] Forgot password screen updated
  - [ ] Reset password screen updated

- [ ] **Build configuration**
  - [ ] Dart-define flags documented
  - [ ] Build scripts updated

## 2. Database Migration

- [ ] **Migration file created**
  - [ ] File: `db/migrations/011_rename_insforge_to_supabase.sql`
  - [ ] Renames `insforge_user_id` to `supabase_user_id`
  - [ ] Renames `insforge_login_email` to `supabase_login_email`
  - [ ] Drops old index `idx_users_insforge_login_email`
  - [ ] Creates new index `idx_users_supabase_login_email`

- [ ] **Migration tested in staging**
  - [ ] Migration executed successfully
  - [ ] Data integrity verified
  - [ ] No data loss
  - [ ] Rollback tested

- [ ] **Backup plan ready**
  - [ ] Backup procedure documented
  - [ ] Backup script tested
  - [ ] Restore procedure documented
  - [ ] Rollback SQL prepared

## 3. Environment Configuration

### Backend Services Environment Variables

- [ ] **SUPABASE_URL**
  - Value: `https://[YOUR-PROJECT-REF].supabase.co`
  - Configured in: [ ] .env [ ] Docker [ ] CI/CD

- [ ] **SUPABASE_ANON_KEY**
  - Value: `[YOUR_ANON_KEY]`
  - Configured in: [ ] .env [ ] Docker [ ] CI/CD

- [ ] **SUPABASE_SERVICE_ROLE_KEY**
  - Value: `[YOUR_SERVICE_ROLE_KEY]`
  - Configured in: [ ] .env [ ] Docker [ ] CI/CD
  - ⚠️ **Security:** Never exposed in client-side code

- [ ] **DATABASE_URL**
  - Connection pooler URL configured
  - Password secured

### Frontend Environment Variables

- [ ] **VITE_SUPABASE_URL** configured
- [ ] **VITE_SUPABASE_ANON_KEY** configured
- [ ] **VITE_API_BASE_URL** configured

### Mobile Environment Variables

- [ ] **SUPABASE_URL** in dart-define
- [ ] **SUPABASE_ANON_KEY** in dart-define
- [ ] **API_BASE_URL** in dart-define

## 4. Supabase Project Configuration

### Authentication Settings

- [ ] **Email provider enabled**
- [ ] **Email confirmations enabled**
- [ ] **Password requirements configured**
  - Minimum length: 8 characters
- [ ] **JWT expiry configured**
  - Access token: 3600 seconds (1 hour)
  - Refresh token: 2592000 seconds (30 days)
- [ ] **Refresh token rotation enabled**

### URL Configuration

- [ ] **Site URL configured**
  - Production: `https://app.tayosa.com`
- [ ] **Redirect URLs whitelisted**
  - [ ] `https://app.tayosa.com/auth/callback`
  - [ ] `https://app.tayosa.com/auth/verify-email`
  - [ ] `https://app.tayosa.com/auth/reset-password`
  - [ ] Development URLs
  - [ ] Mobile deep links

### Email Templates

- [ ] **Email verification template customized**
- [ ] **Password reset template customized**
- [ ] **Magic link template customized** (if used)
- [ ] **Test emails sent and received**

### SMTP Configuration (Production)

- [ ] **Custom SMTP provider configured**
  - Provider: _________________
- [ ] **Sender email verified**
  - Email: `noreply@tayosa.com`
- [ ] **Test email delivery working**

### Storage Configuration

- [ ] **Storage bucket created**
  - Bucket name: `collateral_docs`
  - Public: No (private)
- [ ] **Storage policies configured**
  - Upload policy for authenticated users
  - Read policy for file owners

### Security Settings

- [ ] **Rate limiting configured**
  - Anonymous: 100 requests/hour
  - Authenticated: 1000 requests/hour
- [ ] **Account lockout enabled**
  - Failed attempts: 5
  - Lockout duration: 10 minutes
- [ ] **RLS enabled on tables**
  - [ ] users table
  - [ ] Other sensitive tables

## 5. Documentation

- [ ] **README.md updated**
  - Supabase setup instructions added
  - Migration instructions included

- [ ] **API Documentation created**
  - File: `docs/API_DOCUMENTATION.md`
  - Authentication endpoints documented
  - Token format documented
  - Error responses documented

- [ ] **Supabase Project Configuration documented**
  - File: `docs/SUPABASE_PROJECT_CONFIGURATION.md`
  - Email templates documented
  - RLS policies documented
  - Storage configuration documented

- [ ] **Migration Runbook created**
  - File: `docs/SUPABASE_MIGRATION_RUNBOOK.md`
  - Pre-migration checklist
  - Step-by-step instructions
  - Rollback procedures

- [ ] **Database Migration Plan created**
  - File: `docs/DATABASE_MIGRATION_PLAN.md`
  - Execution timeline
  - Verification steps
  - Rollback plan

- [ ] **Deployment Guide created**
  - File: `docs/DEPLOYMENT_GUIDE.md`
  - Environment variables documented
  - Deployment steps
  - Docker configuration

- [ ] **Email Templates documented**
  - File: `docs/SUPABASE_EMAIL_TEMPLATES.md`
  - All templates provided
  - Configuration instructions
  - Testing procedures

- [ ] **Auth Settings documented**
  - File: `docs/SUPABASE_AUTH_SETTINGS.md`
  - All settings explained
  - Configuration steps
  - Troubleshooting guide

- [ ] **Verification Checklist created**
  - File: `docs/SUPABASE_VERIFICATION_CHECKLIST.md`
  - Comprehensive verification steps
  - Testing procedures

## 6. Testing

### Unit Tests

- [ ] **Backend services**
  - [ ] Token validation tests
  - [ ] Authentication middleware tests
  - [ ] Database query tests

- [ ] **Frontend**
  - [ ] Supabase client tests
  - [ ] Authentication flow tests

- [ ] **Mobile**
  - [ ] Supabase integration tests
  - [ ] Authentication flow tests

### Integration Tests

- [ ] **Authentication flow**
  - [ ] User registration
  - [ ] Email verification
  - [ ] User login
  - [ ] Password reset
  - [ ] Token refresh

- [ ] **API Gateway**
  - [ ] Token validation
  - [ ] Protected endpoints
  - [ ] Token forwarding

- [ ] **Backend services**
  - [ ] All services validate tokens
  - [ ] Authenticated requests succeed
  - [ ] Unauthenticated requests rejected

### End-to-End Tests

- [ ] **Web application**
  - [ ] Complete registration flow
  - [ ] Complete login flow
  - [ ] Access protected pages
  - [ ] Password reset flow

- [ ] **Mobile application**
  - [ ] Complete registration flow
  - [ ] Complete login flow
  - [ ] Access protected features
  - [ ] Password reset flow

## 7. Performance Testing

- [ ] **Token validation latency**
  - Target: < 100ms
  - Actual: _________ ms

- [ ] **Database query performance**
  - Queries optimized
  - Indexes created

- [ ] **API response times**
  - Authentication endpoints: < 500ms
  - Protected endpoints: < 1000ms

- [ ] **Load testing**
  - Concurrent users tested: _________
  - Success rate: _________ %

## 8. Security Verification

- [ ] **Credentials secured**
  - [ ] Service role key not exposed
  - [ ] Environment variables not in git
  - [ ] Secrets stored securely

- [ ] **RLS policies tested**
  - [ ] Users can only access own data
  - [ ] Service role bypasses RLS
  - [ ] Unauthorized access blocked

- [ ] **Token validation tested**
  - [ ] Valid tokens accepted
  - [ ] Invalid tokens rejected
  - [ ] Expired tokens rejected

- [ ] **HTTPS enforced**
  - [ ] All endpoints use HTTPS
  - [ ] HTTP redirects to HTTPS

## 9. Monitoring and Logging

- [ ] **Application logs configured**
  - [ ] Log level set appropriately
  - [ ] Logs accessible

- [ ] **Supabase logs accessible**
  - [ ] Authentication logs
  - [ ] Database logs
  - [ ] API logs

- [ ] **Monitoring dashboards**
  - [ ] Error rates
  - [ ] Response times
  - [ ] Authentication metrics

- [ ] **Alerts configured**
  - [ ] High error rate alerts
  - [ ] Performance degradation alerts
  - [ ] Security event alerts

## 10. Rollback Preparation

- [ ] **Rollback plan documented**
  - [ ] Step-by-step instructions
  - [ ] Database restore procedure
  - [ ] Code rollback procedure

- [ ] **Rollback tested in staging**
  - [ ] Database restore works
  - [ ] Application rollback works
  - [ ] Services restart correctly

- [ ] **Team trained on rollback**
  - [ ] All team members know procedure
  - [ ] Contact information available
  - [ ] Escalation path defined

## 11. Communication

- [ ] **Maintenance window scheduled**
  - Date: _________________
  - Time: _________________
  - Duration: _________________

- [ ] **Users notified**
  - [ ] T-7 days notification
  - [ ] T-24 hours reminder
  - [ ] T-1 hour final reminder

- [ ] **Team availability confirmed**
  - [ ] Database admin: _________________
  - [ ] Backend lead: _________________
  - [ ] DevOps lead: _________________
  - [ ] On-call engineer: _________________

- [ ] **Communication channels ready**
  - [ ] Status page
  - [ ] Email notifications
  - [ ] In-app notifications

## 12. Final Checks

- [ ] **All code committed**
  - [ ] No uncommitted changes
  - [ ] All branches merged
  - [ ] Tags created

- [ ] **All tests passing**
  - [ ] Unit tests: 100% pass
  - [ ] Integration tests: 100% pass
  - [ ] E2E tests: 100% pass

- [ ] **Staging environment verified**
  - [ ] All features working
  - [ ] No critical bugs
  - [ ] Performance acceptable

- [ ] **Production environment ready**
  - [ ] Infrastructure provisioned
  - [ ] DNS configured
  - [ ] SSL certificates valid

- [ ] **Team ready**
  - [ ] All team members briefed
  - [ ] Roles and responsibilities clear
  - [ ] Emergency contacts available

## Decision

Based on the verification above:

**Ready for deployment:** [ ] Yes [ ] No

**If No, blockers:**
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________

**Deployment approved by:**

- Technical Lead: _________________ Date: _______
- Database Admin: _________________ Date: _______
- DevOps Lead: _________________ Date: _______
- Product Owner: _________________ Date: _______

**Deployment scheduled for:**
- Date: _________________
- Time: _________________
- Duration: _________________

**Notes:**
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
