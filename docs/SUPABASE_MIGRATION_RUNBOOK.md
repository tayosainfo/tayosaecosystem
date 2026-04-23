# Supabase Migration Runbook

## Overview

This runbook provides comprehensive procedures for migrating the Tayosa banking ecosystem from InsForge Backend-as-a-Service to Supabase. The migration addresses InsForge's inability to send verification emails while maintaining system functionality and user experience.

**Migration Scope:**
- Frontend: React + Vite web application
- Backend: 9 Go microservices
- Mobile: Flutter application
- Database: PostgreSQL with 10 migration files
- Tests: Integration test files

**Target Supabase Project:** `https://[YOUR-PROJECT-REF].supabase.co`

## Pre-Migration Checklist

### 1. Environment Preparation

#### 1.1 Supabase Project Verification
- [ ] Confirm Supabase project URL: `https://[YOUR-PROJECT-REF].supabase.co`
- [ ] Verify access to Supabase dashboard
- [ ] Confirm anon key is available and valid
- [ ] Confirm service role key is available and valid
- [ ] Test database connection using provided connection string

#### 1.2 Email Configuration
- [ ] Configure email templates in Supabase Dashboard → Authentication → Email Templates
  - [ ] Confirmation email template customized
  - [ ] Password reset email template customized
  - [ ] Magic link template (if used)
- [ ] Test email delivery from Supabase dashboard
- [ ] Configure SMTP settings (if using custom email provider)
- [ ] Verify allowed redirect URLs are configured

#### 1.3 Authentication Settings
- [ ] Enable email confirmations: Dashboard → Authentication → Settings
- [ ] Configure session timeout settings
- [ ] Set up OAuth providers (if required):
  - [ ] Google OAuth configured
  - [ ] GitHub OAuth configured
  - [ ] Other providers as needed
- [ ] Configure allowed redirect URLs for email links

#### 1.4 Database Preparation
- [ ] Verify PostgreSQL connection to Supabase database
- [ ] Confirm all 10 existing migration files are compatible
- [ ] Create backup of current database state
- [ ] Prepare migration file: `db/migrations/011_rename_insforge_to_supabase.sql`
- [ ] Test migration on staging/development environment

#### 1.5 Environment Variables
- [ ] Prepare environment variables for all environments:
  ```bash
  # Frontend
  VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
  VITE_SUPABASE_ANON_KEY=your-anon-key
  
  # Backend Services
  SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
  SUPABASE_ANON_KEY=your-anon-key
  SUPABASE_SERVICE_ROLE_KEY=your-service-role-key
  DATABASE_URL=postgresql://postgres:PASSWORD@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres?sslmode=require
  
  # Mobile (--dart-define)
  SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
  SUPABASE_ANON_KEY=your-anon-key
  API_BASE_URL=your-api-base-url
  ```

#### 1.6 Code Readiness
- [ ] All migration code changes are committed and tested
- [ ] Frontend refactoring completed (InsForge → Supabase references)
- [ ] Backend services updated with Supabase token validation
- [ ] Mobile app Supabase SDK integration completed
- [ ] Test files updated for Supabase
- [ ] Documentation updated

#### 1.7 Testing Environment
- [ ] Set up Supabase test project for integration testing
- [ ] Configure test environment variables
- [ ] Run integration tests against Supabase test environment
- [ ] Verify all authentication flows work correctly
- [ ] Test email verification and password reset flows

## Step-by-Step Deployment Instructions

### Phase 1: Database Migration

#### Step 1.1: Pre-Migration Database Backup
```bash
# Create full database backup
pg_dump $CURRENT_DATABASE_URL > backup_pre_supabase_migration_$(date +%Y%m%d_%H%M%S).sql

# Verify backup integrity
psql $CURRENT_DATABASE_URL -c "SELECT COUNT(*) FROM users_identity;"
```

#### Step 1.2: Execute Database Migration
```bash
# Connect to Supabase database
export DATABASE_URL="postgresql://postgres:PASSWORD@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres?sslmode=require"

# Execute migration file
psql $DATABASE_URL -f db/migrations/011_rename_insforge_to_supabase.sql

# Verify migration success
psql $DATABASE_URL -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'users_identity' AND column_name LIKE '%supabase%';"
```

#### Step 1.3: Data Integrity Verification
```bash
# Check for null supabase_user_id values
psql $DATABASE_URL -c "SELECT COUNT(*) FROM users_identity WHERE supabase_user_id IS NULL;"

# Check for duplicate emails
psql $DATABASE_URL -c "SELECT auth_email, COUNT(*) FROM users_identity GROUP BY auth_email HAVING COUNT(*) > 1;"

# Verify index creation
psql $DATABASE_URL -c "SELECT indexname FROM pg_indexes WHERE tablename = 'users_identity' AND indexname LIKE '%supabase%';"
```

### Phase 2: Backend Services Deployment

#### Step 2.1: Deploy User Service
```bash
# Update environment variables
export SUPABASE_URL="https://[YOUR-PROJECT-REF].supabase.co"
export SUPABASE_ANON_KEY="your-anon-key"
export SUPABASE_SERVICE_ROLE_KEY="your-service-role-key"

# Deploy user-service
cd services/user-service
go build -o user-service
./user-service &

# Verify service health
curl http://localhost:8081/health
```

#### Step 2.2: Deploy API Gateway Service
```bash
# Deploy api-gateway-service with Supabase token validation
cd services/api-gateway-service
go build -o api-gateway-service
./api-gateway-service &

# Verify gateway health
curl http://localhost:8080/health
```

#### Step 2.3: Deploy Other Backend Services
```bash
# Deploy each service with Supabase environment variables
for service in affiliate audit-log fee kibiina loan-credit notification object-storage; do
  cd services/${service}-service
  export SUPABASE_URL="https://[YOUR-PROJECT-REF].supabase.co"
  export SUPABASE_ANON_KEY="your-anon-key"
  go build -o ${service}-service
  ./${service}-service &
  cd ../..
done
```

#### Step 2.4: Verify Backend Services
```bash
# Test authentication endpoint
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier": "test@example.com", "password": "testpassword"}'

# Test token validation
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/users/profile
```

### Phase 3: Frontend Deployment

#### Step 3.1: Update Environment Variables
```bash
# Update .env file
cat > .env << EOF
VITE_API_BASE_URL=http://localhost:8080
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key
EOF
```

#### Step 3.2: Build and Deploy Frontend
```bash
# Install dependencies and build
npm install
npm run build

# Deploy to hosting platform
npm run deploy
# OR copy dist/ to web server
```

#### Step 3.3: Verify Frontend
```bash
# Test frontend authentication
# 1. Open browser to application URL
# 2. Test user registration with email verification
# 3. Test user login
# 4. Test password reset flow
# 5. Verify session persistence across page refreshes
```

### Phase 4: Mobile App Deployment

#### Step 4.1: Build Mobile App
```bash
cd app/mobile_app

# Build for Android
flutter build apk \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-anon-key \
  --dart-define=API_BASE_URL=your-api-base-url

# Build for iOS
flutter build ios \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-anon-key \
  --dart-define=API_BASE_URL=your-api-base-url
```

#### Step 4.2: Deploy Mobile App
```bash
# Deploy to app stores or distribute internally
# Follow platform-specific deployment procedures
```

## Database Migration Execution Steps

### Migration File: `db/migrations/011_rename_insforge_to_supabase.sql`

```sql
-- Rename InsForge columns to Supabase
BEGIN;

-- Rename insforge_user_id to supabase_user_id
ALTER TABLE users_identity 
  RENAME COLUMN insforge_user_id TO supabase_user_id;

-- Rename insforge_login_email to supabase_login_email  
ALTER TABLE users_identity 
  RENAME COLUMN insforge_login_email TO supabase_login_email;

-- Drop old index
DROP INDEX IF EXISTS idx_users_insforge_login_email;

-- Create new index
CREATE INDEX idx_users_supabase_login_email 
  ON users_identity(supabase_login_email);

-- Verify changes
SELECT column_name 
FROM information_schema.columns 
WHERE table_name = 'users_identity' 
  AND column_name LIKE '%supabase%';

COMMIT;
```

### Execution Steps

1. **Pre-execution Verification**
   ```bash
   # Check current schema
   psql $DATABASE_URL -c "\d users_identity"
   
   # Count existing records
   psql $DATABASE_URL -c "SELECT COUNT(*) FROM users_identity;"
   ```

2. **Execute Migration**
   ```bash
   # Run migration with transaction safety
   psql $DATABASE_URL -f db/migrations/011_rename_insforge_to_supabase.sql
   ```

3. **Post-execution Verification**
   ```bash
   # Verify column rename
   psql $DATABASE_URL -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'users_identity' AND column_name LIKE '%supabase%';"
   
   # Verify index creation
   psql $DATABASE_URL -c "SELECT indexname FROM pg_indexes WHERE tablename = 'users_identity';"
   
   # Verify data integrity
   psql $DATABASE_URL -c "SELECT COUNT(*) FROM users_identity WHERE supabase_user_id IS NOT NULL;"
   ```

## Rollback Procedures

### Emergency Rollback Triggers
- Authentication failures > 50% of requests
- Database connection errors
- Email verification not working
- Critical application errors
- User data corruption detected

### Phase 1: Immediate Service Rollback

#### Step 1.1: Stop New Deployments
```bash
# Stop all new services
pkill -f user-service
pkill -f api-gateway-service
pkill -f affiliate-service
# ... stop other services
```

#### Step 1.2: Restore Previous Service Versions
```bash
# Deploy previous versions from backup
cd services/user-service
git checkout HEAD~1  # or specific commit
go build -o user-service
./user-service &

# Repeat for all services
```

#### Step 1.3: Restore Environment Variables
```bash
# Switch back to InsForge environment variables
export VITE_INSFORGE_URL="previous-insforge-url"
export VITE_INSFORGE_ANON_KEY="previous-insforge-key"
# ... restore other variables
```

### Phase 2: Database Rollback

#### Step 2.1: Create Rollback Migration
```sql
-- File: db/migrations/012_rollback_supabase_to_insforge.sql
BEGIN;

-- Rename supabase_user_id back to insforge_user_id
ALTER TABLE users_identity 
  RENAME COLUMN supabase_user_id TO insforge_user_id;

-- Rename supabase_login_email back to insforge_login_email
ALTER TABLE users_identity 
  RENAME COLUMN supabase_login_email TO insforge_login_email;

-- Drop Supabase index
DROP INDEX IF EXISTS idx_users_supabase_login_email;

-- Recreate InsForge index
CREATE INDEX idx_users_insforge_login_email 
  ON users_identity(insforge_login_email);

COMMIT;
```

#### Step 2.2: Execute Database Rollback
```bash
# Execute rollback migration
psql $DATABASE_URL -f db/migrations/012_rollback_supabase_to_insforge.sql

# Verify rollback
psql $DATABASE_URL -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'users_identity' AND column_name LIKE '%insforge%';"
```

### Phase 3: Frontend Rollback

#### Step 3.1: Restore Frontend Code
```bash
# Checkout previous version
git checkout HEAD~1  # or specific commit

# Restore environment variables
cat > .env << EOF
VITE_API_BASE_URL=http://localhost:8080
VITE_INSFORGE_URL=previous-insforge-url
VITE_INSFORGE_ANON_KEY=previous-insforge-key
EOF

# Rebuild and deploy
npm run build
npm run deploy
```

### Phase 4: Mobile App Rollback

#### Step 4.1: Restore Mobile App
```bash
# Checkout previous version
cd app/mobile_app
git checkout HEAD~1

# Rebuild with InsForge configuration
flutter build apk \
  --dart-define=INSFORGE_URL=previous-insforge-url \
  --dart-define=INSFORGE_ANON_KEY=previous-insforge-key

# Redeploy to app stores (emergency release)
```

### Rollback Verification Checklist
- [ ] All services are running on previous versions
- [ ] Database schema restored to pre-migration state
- [ ] Environment variables restored
- [ ] Authentication flows working
- [ ] User registration working
- [ ] Email verification working (if it was working before)
- [ ] All critical user journeys functional
- [ ] No data loss detected

## Post-Migration Verification Steps

### 1. System Health Checks

#### 1.1 Service Health Verification
```bash
# Check all services are running
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # User Service
# ... check other services

# Verify service logs for errors
tail -f logs/user-service.log
tail -f logs/api-gateway.log
```

#### 1.2 Database Connectivity
```bash
# Test database connection
psql $DATABASE_URL -c "SELECT version();"

# Verify migration applied correctly
psql $DATABASE_URL -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'users_identity' AND column_name LIKE '%supabase%';"

# Check data integrity
psql $DATABASE_URL -c "SELECT COUNT(*) FROM users_identity WHERE supabase_user_id IS NOT NULL;"
```

### 2. Authentication Flow Testing

#### 2.1 User Registration Flow
```bash
# Test registration via API
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-migration@example.com",
    "password": "TestPassword123!",
    "fullName": "Test User",
    "phone": "+256700000000"
  }'

# Expected: Success response with user data
# Expected: Verification email sent via Supabase
```

#### 2.2 Email Verification Flow
```bash
# Manual verification:
# 1. Check email inbox for verification email
# 2. Click verification link
# 3. Verify user can now login
# 4. Check user status in database
psql $DATABASE_URL -c "SELECT auth_email, contact_email_verified_at FROM users_identity WHERE auth_email = 'test-migration@example.com';"
```

#### 2.3 User Login Flow
```bash
# Test login via API
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "test-migration@example.com",
    "password": "TestPassword123!"
  }'

# Expected: Success response with access token
# Expected: Token can be used for authenticated requests
```

#### 2.4 Token Validation
```bash
# Extract token from login response
TOKEN="your-access-token-here"

# Test authenticated endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/users/profile

# Expected: User profile data returned
# Expected: No authentication errors
```

#### 2.5 Password Reset Flow
```bash
# Test password reset request
curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "test-migration@example.com"}'

# Manual verification:
# 1. Check email inbox for reset email
# 2. Click reset link
# 3. Complete password reset
# 4. Login with new password
```

### 3. Frontend Verification

#### 3.1 Web Application Testing
- [ ] Navigate to application URL
- [ ] Test user registration with email verification
- [ ] Test user login with verified account
- [ ] Test password reset flow
- [ ] Verify session persistence across page refreshes
- [ ] Test logout functionality
- [ ] Verify protected routes require authentication
- [ ] Test error handling for invalid credentials

#### 3.2 Frontend Console Checks
```javascript
// Open browser developer console
// Check for JavaScript errors
console.log('Checking for errors...');

// Verify Supabase client initialization
console.log(window.supabase);  // Should show Supabase client

// Check authentication state
console.log(localStorage.getItem('auth_token'));
console.log(sessionStorage.getItem('auth_user'));
```

### 4. Mobile App Verification

#### 4.1 Mobile Authentication Testing
- [ ] Install updated mobile app
- [ ] Test user registration flow
- [ ] Test user login flow
- [ ] Test password reset flow
- [ ] Verify session persistence across app restarts
- [ ] Test logout functionality
- [ ] Verify API calls include authentication headers

#### 4.2 Mobile App Logs
```bash
# Android logs
adb logcat | grep -i supabase

# iOS logs (via Xcode or device console)
# Look for Supabase-related log entries
```

### 5. Integration Testing

#### 5.1 Run Automated Tests
```bash
# Run updated integration tests
node test_unified_auth_flow.js
node test_unverified_login.js

# Expected: All tests pass
# Expected: Tests connect to Supabase successfully
```

#### 5.2 End-to-End User Journey
1. **Registration Journey**:
   - Register new user via web app
   - Receive verification email
   - Verify email address
   - Login to web app
   - Login to mobile app with same credentials

2. **Password Reset Journey**:
   - Request password reset via web app
   - Receive reset email
   - Complete password reset
   - Login with new password on both web and mobile

3. **Cross-Platform Session**:
   - Login on web app
   - Verify session works across browser tabs
   - Login on mobile app
   - Verify independent session management

### 6. Performance and Monitoring

#### 6.1 Response Time Verification
```bash
# Test authentication endpoint response times
time curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier": "test@example.com", "password": "password"}'

# Expected: < 200ms for authentication endpoints
# Expected: < 100ms for token validation
```

#### 6.2 Error Rate Monitoring
```bash
# Monitor service logs for error patterns
grep -i error logs/user-service.log | wc -l
grep -i "supabase" logs/user-service.log | grep -i error

# Expected: No critical errors
# Expected: Supabase API calls succeeding
```

#### 6.3 Supabase Dashboard Monitoring
- [ ] Check Supabase Dashboard → Authentication → Users
- [ ] Verify new user registrations appear
- [ ] Check authentication metrics
- [ ] Monitor API usage and rate limits
- [ ] Verify email delivery statistics

### 7. Data Integrity Verification

#### 7.1 User Data Consistency
```bash
# Verify user data migration
psql $DATABASE_URL -c "
  SELECT 
    COUNT(*) as total_users,
    COUNT(supabase_user_id) as users_with_supabase_id,
    COUNT(CASE WHEN supabase_user_id IS NULL THEN 1 END) as users_without_supabase_id
  FROM users_identity;
"

# Expected: All users have supabase_user_id
# Expected: No data loss from migration
```

#### 7.2 Authentication State Consistency
```bash
# Check for orphaned sessions or inconsistent states
psql $DATABASE_URL -c "
  SELECT auth_email, supabase_user_id, created_at 
  FROM users_identity 
  WHERE supabase_user_id IS NULL 
  LIMIT 10;
"

# Expected: No users without supabase_user_id
```

### 8. Security Verification

#### 8.1 Token Security
```bash
# Verify tokens are properly validated
curl -H "Authorization: Bearer invalid-token" \
  http://localhost:8080/api/v1/users/profile

# Expected: 401 Unauthorized response
```

#### 8.2 Environment Variable Security
```bash
# Verify sensitive keys are not exposed
grep -r "supabase.*key" . --exclude-dir=node_modules --exclude-dir=.git

# Expected: Keys only in .env files and secure configuration
# Expected: No keys in source code or logs
```

## Success Criteria

### Migration is considered successful when:

1. **Authentication Functionality**:
   - [ ] User registration works with email verification
   - [ ] User login works for existing and new users
   - [ ] Password reset flow works end-to-end
   - [ ] Token validation works across all services
   - [ ] Session management works on web and mobile

2. **System Stability**:
   - [ ] All backend services are running without errors
   - [ ] Database migration completed successfully
   - [ ] No data loss detected
   - [ ] Response times meet performance targets
   - [ ] Error rates are within acceptable limits

3. **Email Functionality**:
   - [ ] Email verification emails are delivered
   - [ ] Password reset emails are delivered
   - [ ] Email templates are properly formatted
   - [ ] Email delivery rates are satisfactory

4. **Cross-Platform Compatibility**:
   - [ ] Web application authentication works
   - [ ] Mobile application authentication works
   - [ ] API endpoints work for all clients
   - [ ] Session management is consistent

5. **Monitoring and Observability**:
   - [ ] Supabase dashboard shows user activity
   - [ ] Service logs show successful Supabase integration
   - [ ] No critical errors in application logs
   - [ ] Authentication metrics are being tracked

## Troubleshooting Guide

### Common Issues and Solutions

#### Issue: Email verification not working
**Symptoms**: Users not receiving verification emails
**Solutions**:
1. Check Supabase Dashboard → Authentication → Settings → Email confirmations enabled
2. Verify SMTP configuration in Supabase
3. Check email templates are configured
4. Verify allowed redirect URLs

#### Issue: Token validation failing
**Symptoms**: 401 errors on authenticated endpoints
**Solutions**:
1. Verify SUPABASE_ANON_KEY is correct in all services
2. Check token format and expiration
3. Verify Supabase project URL is correct
4. Check service logs for Supabase API errors

#### Issue: Database connection errors
**Symptoms**: Services cannot connect to database
**Solutions**:
1. Verify DATABASE_URL format and credentials
2. Check Supabase database is accessible
3. Verify SSL mode is set correctly
4. Check firewall and network connectivity

#### Issue: Mobile app authentication failing
**Symptoms**: Mobile app cannot authenticate users
**Solutions**:
1. Verify Supabase SDK is properly initialized
2. Check --dart-define variables are set correctly
3. Verify API client is sending correct headers
4. Check mobile app logs for Supabase errors

### Emergency Contacts

- **Technical Lead**: [Contact Information]
- **DevOps Engineer**: [Contact Information]
- **Database Administrator**: [Contact Information]
- **Supabase Support**: [Support Channel]

### Rollback Decision Matrix

| Issue Severity | User Impact | Rollback Decision |
|----------------|-------------|-------------------|
| Critical | >50% users affected | Immediate rollback |
| High | 10-50% users affected | Rollback within 1 hour |
| Medium | <10% users affected | Fix forward or rollback within 4 hours |
| Low | Minimal impact | Fix forward |

---

**Document Version**: 1.0  
**Last Updated**: [Current Date]  
**Next Review**: [Date + 3 months]
