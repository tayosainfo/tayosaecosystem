# Role-Based Admin Authentication - Deployment Guide

## Prerequisites Checklist

Before starting, ensure you have:

- [ ] Access to Supabase Dashboard
- [ ] Database connection string (from `.env` file)
- [ ] Admin user email/phone to assign initial admin role
- [ ] Backend services can be restarted
- [ ] Frontend can be rebuilt and deployed

## Step-by-Step Deployment

---

## Phase 1: Database Setup (15-20 minutes)

### Step 1.1: Connect to Database

**Option A: Using psql command line**
```bash
# Get your DATABASE_URL from .env file
# It should look like: postgresql://postgres:PASSWORD@db.ablvrbnbsdqshrorhmjf.supabase.co:5432/postgres

# Connect to database
psql "postgresql://postgres:YOUR_PASSWORD@db.ablvrbnbsdqshrorhmjf.supabase.co:5432/postgres"
```

**Option B: Using Supabase SQL Editor**
1. Go to Supabase Dashboard
2. Navigate to SQL Editor
3. Create a new query

### Step 1.2: Run Migration 013 (User Roles)

**Using psql:**
```bash
psql "postgresql://postgres:YOUR_PASSWORD@db.ablvrbnbsdqshrorhmjf.supabase.co:5432/postgres" -f db/migrations/013_add_user_roles.sql
```

**Using Supabase SQL Editor:**
1. Open `db/migrations/013_add_user_roles.sql` in your code editor
2. Copy the entire contents
3. Paste into Supabase SQL Editor
4. Click "Run" button

**Expected Output:**
```
CREATE TYPE
ALTER TABLE
CREATE INDEX
CREATE TABLE
CREATE INDEX
CREATE INDEX
CREATE FUNCTION
CREATE TRIGGER
COMMENT
COMMENT
COMMENT
COMMENT
```

**Verify Migration 013:**
```sql
-- Check if role column exists
SELECT column_name, data_type, column_default
FROM information_schema.columns 
WHERE table_name = 'users_identity' 
  AND column_name IN ('role', 'status', 'role_assigned_at', 'role_assigned_by', 'last_login');

-- Expected: 5 rows showing the new columns

-- Check if audit table exists
SELECT table_name 
FROM information_schema.tables 
WHERE table_name = 'admin_role_audit';

-- Expected: 1 row with 'admin_role_audit'
```

### Step 1.3: Run Migration 014 (Custom Claims Hook)

**Using psql:**
```bash
psql "postgresql://postgres:YOUR_PASSWORD@db.ablvrbnbsdqshrorhmjf.supabase.co:5432/postgres" -f db/migrations/014_configure_custom_claims.sql
```

**Using Supabase SQL Editor:**
1. Open `db/migrations/014_configure_custom_claims.sql`
2. Copy the entire contents
3. Paste into Supabase SQL Editor
4. Click "Run" button

**Expected Output:**
```
CREATE FUNCTION
GRANT
GRANT
GRANT
COMMENT
```

**Verify Migration 014:**
```sql
-- Check if custom claims hook function exists
SELECT routine_name, routine_type
FROM information_schema.routines 
WHERE routine_name = 'custom_access_token_hook';

-- Expected: 1 row with 'custom_access_token_hook' and 'FUNCTION'

-- Test the function (should not error)
SELECT public.custom_access_token_hook('{"user_id": "test", "claims": {}}'::jsonb);

-- Expected: JSON object with user_role claim
```

✅ **Checkpoint:** Both migrations should complete without errors.

---

## Phase 2: Configure Supabase Custom Claims Hook (5 minutes)

### Step 2.1: Access Supabase Dashboard

1. Go to https://supabase.com/dashboard
2. Select your project: `ablvrbnbsdqshrorhmjf`
3. You should see your project dashboard

### Step 2.2: Navigate to Authentication Hooks

1. In the left sidebar, click **"Authentication"**
2. Click on the **"Hooks"** tab
3. You should see a list of available hooks

### Step 2.3: Enable Custom Access Token Hook

1. Find **"Custom Access Token"** in the hooks list
2. Click the **"Enable Hook"** button or toggle switch
3. A configuration panel should appear

### Step 2.4: Configure Hook Function

1. In the **"Hook Function"** field, enter:
   ```
   public.custom_access_token_hook
   ```
2. Make sure there are no extra spaces
3. Click **"Save"** or **"Update"** button

### Step 2.5: Verify Hook Configuration

1. The hook should now show as **"Enabled"**
2. The function should be: `public.custom_access_token_hook`
3. Status indicator should be green/active

**Test Hook (Optional):**
```sql
-- Create a test user and check JWT claims
-- This will be tested in Phase 4
```

✅ **Checkpoint:** Custom Access Token hook is enabled and configured.

---

## Phase 3: Assign Initial Admin Roles (5 minutes)

### Step 3.1: Identify Your Admin User(s)

First, find the user_id of the account(s) you want to make admin:

```sql
-- Find your user by email
SELECT user_id, full_name, auth_email, phone_e164, role, status
FROM users_identity 
WHERE auth_email = 'your-email@example.com';

-- OR find by phone
SELECT user_id, full_name, auth_email, phone_e164, role, status
FROM users_identity 
WHERE phone_e164 = '+256700000000';

-- List all users to choose from
SELECT user_id, full_name, auth_email, phone_e164, role, status
FROM users_identity 
ORDER BY created_at DESC
LIMIT 10;
```

**Copy the `user_id` value(s)** - you'll need them in the next step.

### Step 3.2: Assign Admin Role

Replace `'USER_ID_HERE'` with the actual user_id from Step 3.1:

```sql
-- Assign admin role to one user
UPDATE users_identity 
SET 
  role = 'admin',
  status = 'active',
  role_assigned_at = NOW(),
  role_assigned_by = 'system',
  updated_at = NOW()
WHERE user_id = 'USER_ID_HERE';

-- For multiple users, use IN clause
UPDATE users_identity 
SET 
  role = 'admin',
  status = 'active',
  role_assigned_at = NOW(),
  role_assigned_by = 'system',
  updated_at = NOW()
WHERE user_id IN (
  'USER_ID_1',
  'USER_ID_2',
  'USER_ID_3'
);
```

**Expected Output:**
```
UPDATE 1  -- or UPDATE 3 if you assigned 3 users
```

### Step 3.3: Verify Admin Assignment

```sql
-- Check admin users
SELECT 
  user_id,
  full_name,
  auth_email,
  phone_e164,
  role,
  status,
  role_assigned_at,
  role_assigned_by
FROM users_identity 
WHERE role = 'admin';

-- Expected: Your admin user(s) should appear with role = 'admin'

-- Check audit log
SELECT 
  user_id,
  action,
  previous_role,
  new_role,
  assigned_by,
  created_at
FROM admin_role_audit 
ORDER BY created_at DESC 
LIMIT 5;

-- Expected: Role change entries for your admin assignments
```

✅ **Checkpoint:** At least one user has admin role assigned.

---

## Phase 4: Test Backend Services (10 minutes)

### Step 4.1: Set Migration Mode Environment Variable

Edit your `.env` file and add/update:

```bash
# Enable migration mode (supports both JWT and API key)
AUTH_MIGRATION_MODE=true

# Keep the API key temporarily for fallback
ADMIN_API_KEY=change-me-admin-secret
```

### Step 4.2: Start Backend Services

**Terminal 1 - User Service:**
```bash
cd services/user-service
go run .
```

**Expected Output:**
```
user-service: loaded env file ...
user-service: using PostgreSQL store
user-service: Supabase auth enabled (https://ablvrbnbsdqshrorhmjf.supabase.co)
user-service: SUPABASE_SERVICE_ROLE_KEY is set (admin operations available)
user-service listening on :8081
```

**Terminal 2 - API Gateway (if separate):**
```bash
cd services/api-gateway-service
go run .
```

### Step 4.3: Test Admin Endpoints with JWT

**Get JWT Token:**
1. Open your browser
2. Go to your application (e.g., http://localhost:5173)
3. Login with your admin user credentials
4. Open browser DevTools (F12)
5. Go to Application/Storage → Local Storage or Session Storage
6. Find and copy the JWT token (usually stored as `auth_token` or in Supabase session)

**Test with curl:**
```bash
# Replace YOUR_JWT_TOKEN with the actual token
TOKEN="YOUR_JWT_TOKEN"

# Test user list endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/admin/users

# Expected: JSON response with user list (status 200)

# Test user detail endpoint (replace USER_ID)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/admin/users/USER_ID

# Expected: JSON response with user details (status 200)
```

**Test with Non-Admin User (should fail):**
```bash
# Login as regular user and get their token
NON_ADMIN_TOKEN="REGULAR_USER_JWT_TOKEN"

curl -H "Authorization: Bearer $NON_ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/users

# Expected: {"error":"insufficient permissions"} (status 403)
```

**Test without Token (should fail):**
```bash
curl http://localhost:8081/api/v1/admin/users

# Expected: {"error":"authentication required"} (status 401)
```

✅ **Checkpoint:** Admin endpoints work with JWT, reject non-admin users.

---

## Phase 5: Test Frontend (10 minutes)

### Step 5.1: Start Frontend Development Server

```bash
# In project root
npm run dev
```

**Expected Output:**
```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
```

### Step 5.2: Test Admin Login

1. Open browser to http://localhost:5173
2. Login with your admin user credentials
3. After login, you should be redirected to home page

### Step 5.3: Access Admin Dashboard

1. Navigate to http://localhost:5173/admin/users
2. **Expected:** You should see the user list page
3. **If you see "Access Denied":** 
   - Logout and login again (to get fresh JWT with role claim)
   - Check browser console for errors
   - Verify custom claims hook is enabled in Supabase

### Step 5.4: Test User Management Features

**Test User List:**
- [ ] User list loads and displays users
- [ ] Search bar works (try searching by name/email)
- [ ] Status filter works (try filtering by Active/Suspended)
- [ ] KYC filter works (try filtering by Pending/Approved)
- [ ] Pagination works (if you have >20 users)

**Test User Detail:**
1. Click on any user in the list
2. **Expected:** User detail page opens
3. Verify you see:
   - [ ] User information (name, email, phone)
   - [ ] Role and status badges
   - [ ] Account details
   - [ ] KYC information
   - [ ] Activity log (may be empty)

**Test Status Change:**
1. On user detail page, click "Change Status"
2. Select a new status (e.g., "Suspended")
3. Enter a reason (required)
4. Click "Update Status"
5. **Expected:** Success message, status badge updates

**Test Role Change:**
1. Click "Change Role"
2. Select "Admin" role
3. **Expected:** Warning message appears
4. Click "Continue"
5. **Expected:** Confirmation dialog appears
6. Click "Yes, Grant Admin Access"
7. **Expected:** Success message, role badge updates

**Test Password Reset:**
1. Click "Reset Password"
2. **Expected:** Confirmation dialog appears
3. Click "Send Reset Email"
4. **Expected:** Success message
5. Check email inbox for password reset email

### Step 5.5: Test Non-Admin Access

1. Logout from admin account
2. Login with a regular (non-admin) user
3. Try to access http://localhost:5173/admin/users
4. **Expected:** "Access Denied" message or redirect

✅ **Checkpoint:** All admin features work correctly, non-admin users are blocked.

---

## Phase 6: Deploy to Production (30 minutes)

### Step 6.1: Build Frontend

```bash
# In project root
npm run build
```

**Expected Output:**
```
vite v5.x.x building for production...
✓ xxxx modules transformed.
dist/index.html                   x.xx kB
dist/assets/index-xxxxx.js       xxx.xx kB
✓ built in xxxs
```

### Step 6.2: Deploy Backend with Migration Mode

**Update Production Environment Variables:**

```bash
# On your production server/platform
AUTH_MIGRATION_MODE=true
ADMIN_API_KEY=your-secure-api-key-here
SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key
DATABASE_URL=your-database-url
```

**Deploy Backend:**
- If using Docker: `docker-compose up -d --build`
- If using cloud platform: Follow your platform's deployment process
- If using systemd: Restart services with `systemctl restart user-service`

**Verify Backend Deployment:**
```bash
# Check health endpoint
curl https://your-production-domain.com/health

# Expected: {"status":"active","service":"user-service",...}
```

### Step 6.3: Deploy Frontend

**Upload built files:**
- Upload contents of `dist/` folder to your web server
- Or deploy to Vercel/Netlify/etc.

**Verify Frontend Deployment:**
1. Visit your production URL
2. Login with admin credentials
3. Access admin dashboard
4. Test one feature to confirm it works

### Step 6.4: Monitor Production (3-7 days)

**Daily Checks:**

1. **Check Authentication Logs:**
   ```bash
   # Look for authentication method usage
   grep "AUTH_METHOD" /var/log/api-gateway.log | tail -20
   
   # You should see mostly "jwt" with occasional "api_key (fallback)"
   ```

2. **Check Error Logs:**
   ```bash
   # Look for authentication failures
   grep "Failed to extract role\|insufficient permissions" /var/log/api-gateway.log
   ```

3. **Monitor Admin Activity:**
   ```sql
   -- Check recent admin actions
   SELECT 
     action,
     COUNT(*) as count,
     MAX(created_at) as last_action
   FROM admin_role_audit
   WHERE created_at >= NOW() - INTERVAL '24 hours'
   GROUP BY action;
   ```

4. **Check User Complaints:**
   - Monitor support tickets
   - Check for access issues
   - Verify no admin users are locked out

**Success Criteria (after 3-7 days):**
- [ ] No authentication errors in logs
- [ ] All admin users can access admin features
- [ ] No API key fallback usage (or very minimal)
- [ ] No user complaints about access issues

---

## Phase 7: Remove API Key Fallback (Final Step)

⚠️ **Only proceed if Phase 6 monitoring shows no issues!**

### Step 7.1: Update Backend Code

**Option A: Remove Fallback Code (Recommended)**

Edit `services/api-gateway-service/main.go`:

```go
// BEFORE (with fallback):
mux.HandleFunc("/api/v1/admin/*", requireAdminWithFallback(adminHandler))

// AFTER (JWT only):
mux.HandleFunc("/api/v1/admin/*", requireAdmin(adminHandler))
```

**Option B: Just Disable Migration Mode**

Update environment variables:
```bash
AUTH_MIGRATION_MODE=false
# ADMIN_API_KEY can be removed or left (won't be used)
```

### Step 7.2: Remove API Key from Environment

Edit production `.env`:
```bash
# Remove or comment out:
# ADMIN_API_KEY=your-api-key
# AUTH_MIGRATION_MODE=true
```

### Step 7.3: Rebuild and Deploy

```bash
# Rebuild backend
cd services/api-gateway-service
go build

# Deploy updated backend
# (Follow your deployment process)
```

### Step 7.4: Final Verification

```bash
# Test that API key no longer works
curl -H "X-Admin-Secret: your-old-api-key" \
  https://your-production-domain.com/api/v1/admin/users

# Expected: {"error":"authentication required"} (status 401)

# Test that JWT still works
curl -H "Authorization: Bearer $ADMIN_JWT_TOKEN" \
  https://your-production-domain.com/api/v1/admin/users

# Expected: User list (status 200)
```

✅ **Checkpoint:** API key authentication is completely removed, JWT-only authentication is active.

---

## Rollback Procedures

### If Issues Occur During Testing (Phase 4-5)

**Rollback Database:**
```sql
-- Remove role columns
ALTER TABLE users_identity 
  DROP COLUMN IF EXISTS role,
  DROP COLUMN IF EXISTS status,
  DROP COLUMN IF EXISTS role_assigned_at,
  DROP COLUMN IF EXISTS role_assigned_by,
  DROP COLUMN IF EXISTS last_login;

-- Drop audit table
DROP TABLE IF EXISTS admin_role_audit;

-- Drop functions
DROP FUNCTION IF EXISTS public.custom_access_token_hook;
DROP FUNCTION IF EXISTS log_role_change;

-- Drop enum
DROP TYPE IF EXISTS user_role;
```

### If Issues Occur in Production (Phase 6)

**Emergency: Re-enable API Key:**
```bash
# Set environment variable
AUTH_MIGRATION_MODE=true
ADMIN_API_KEY=your-secure-api-key

# Restart services
systemctl restart user-service api-gateway-service
```

This gives you time to investigate while keeping admin access working.

---

## Troubleshooting

### Issue: "insufficient permissions" error for admin user

**Diagnosis:**
```sql
-- Check user role in database
SELECT user_id, role FROM users_identity WHERE auth_email = 'admin@example.com';
```

**Solutions:**
1. Verify role is 'admin' in database
2. User must logout and login again to get new JWT with role claim
3. Check Supabase hook is enabled
4. Decode JWT token to verify user_role claim exists

### Issue: Custom claims hook not working

**Diagnosis:**
```sql
-- Check hook function exists
SELECT routine_name FROM information_schema.routines 
WHERE routine_name = 'custom_access_token_hook';
```

**Solutions:**
1. Re-run migration 014
2. Verify hook is enabled in Supabase Dashboard
3. Check function permissions: `GRANT EXECUTE ON FUNCTION public.custom_access_token_hook TO supabase_auth_admin;`

### Issue: Frontend shows "Access Denied" for admin

**Solutions:**
1. Clear browser cache and cookies
2. Logout and login again
3. Check browser console for errors
4. Verify JWT token contains user_role claim (use jwt.io to decode)
5. Check backend logs for authentication errors

---

## Success Checklist

- [ ] Phase 1: Database migrations completed successfully
- [ ] Phase 2: Supabase custom claims hook enabled
- [ ] Phase 3: At least one admin user assigned
- [ ] Phase 4: Backend services running and responding correctly
- [ ] Phase 5: Frontend admin dashboard accessible and functional
- [ ] Phase 6: Production deployment successful
- [ ] Phase 7: API key fallback removed (after monitoring period)

---

## Support

If you encounter issues:

1. Check logs: `grep "AUTH_METHOD\|Failed to extract role" /var/log/*.log`
2. Review documentation: `docs/ADMIN_AUTH_MIGRATION.md`
3. Check audit trail: `SELECT * FROM admin_role_audit ORDER BY created_at DESC LIMIT 10;`
4. Verify environment variables are set correctly
5. Test with curl commands to isolate frontend vs backend issues

**Emergency Contact:** [Your team contact information]

---

## Completion

Once all phases are complete and Phase 7 is done:

🎉 **Congratulations!** Your admin authentication is now fully migrated to secure JWT-based authentication with role-based access control.

**Final Steps:**
- [ ] Document any custom configurations in your team wiki
- [ ] Train admin users on new user management features
- [ ] Set up monitoring alerts for authentication failures
- [ ] Schedule monthly admin role audits
