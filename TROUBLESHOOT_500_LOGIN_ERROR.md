# Troubleshooting 500 Login Error

## Error
`POST https://tayosaecosystem.onrender.com/api/v1/auth/login 500 (Internal Server Error)`

## How to Find Runtime Logs in Render

1. Go to https://dashboard.render.com
2. Click on your "tayosaecosystem" service
3. Click on the **"Logs"** tab (should be next to "Events", "Settings", etc.)
4. You should see real-time logs from the running application
5. Try to login from the frontend
6. Watch the logs for error messages

## Possible Causes

### 1. Database Migration Not Run
**Symptoms:** Database queries failing because columns don't exist

**Fix:** Run these migrations in Supabase SQL Editor (in order):
```sql
-- First run: db/migrations/021_auto_approve_kyc.sql
-- Then run: db/migrations/022_remove_admin_system.sql
```

### 2. Missing Environment Variables
**Symptoms:** "auth backend not configured" or database connection errors

**Check in Render:**
- `SUPABASE_URL`
- `SUPABASE_ANON_KEY`
- `SUPABASE_SERVICE_ROLE_KEY`
- `DATABASE_URL` (PostgreSQL connection string)

### 3. Database Connection Issue
**Symptoms:** "connection refused" or "timeout" errors

**Fix:** Verify your Supabase database is accessible from Render

### 4. Code Error from Admin Removal
**Symptoms:** "undefined method" or "nil pointer" errors

**Likely culprits:**
- Store methods that were removed but still being called
- Database queries referencing dropped columns

## Quick Diagnostic Steps

### Step 1: Check if Backend is Running
```bash
curl https://tayosaecosystem.onrender.com/health
```

If this returns 404 or times out, the backend isn't running properly.

### Step 2: Check Database Connection
The backend should log database connection status on startup. Look for:
- "Connected to PostgreSQL"
- "Using in-memory store" (fallback if DB fails)

### Step 3: Test a Simple Endpoint
Try accessing a public endpoint:
```bash
curl https://tayosaecosystem.onrender.com/api/v1/auth/public-config
```

If this works but login doesn't, the issue is specific to the login handler.

### Step 4: Check Supabase Connection
The login handler calls Supabase. Verify:
- Supabase project is running
- API keys are correct
- No rate limiting

## Common Fixes

### Fix 1: Restart the Service
Sometimes a simple restart fixes transient issues:
1. Go to Render dashboard
2. Click "Manual Deploy" → "Clear build cache & deploy"

### Fix 2: Check Recent Deployments
Look at the deployment that's currently running:
- Is it the latest commit?
- Did it build successfully?
- Are there any warnings in the build logs?

### Fix 3: Rollback to Previous Version
If the issue started after our admin removal:
1. Go to Render dashboard
2. Find the last working deployment (before commit `21c4242`)
3. Click "Rollback to this version"

## What to Share for Debugging

If you can access the logs, please share:
1. **The exact error message** from the backend logs when login fails
2. **Stack trace** if available
3. **Any database errors** or connection issues
4. **Environment variable status** (without sharing the actual values)

## Temporary Workaround

If you need to get users logged in immediately while we debug:

1. **Rollback the backend** to the last working version
2. **Keep the frontend changes** (they're backward compatible)
3. Users can log in with the old backend
4. We fix the 500 error separately

## Next Steps

1. Access Render logs and find the actual error
2. Share the error message
3. I'll provide a targeted fix based on the actual error

The 500 error means the backend is crashing, so the logs will tell us exactly what's failing.
