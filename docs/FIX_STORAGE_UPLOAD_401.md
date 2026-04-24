# Fix: Storage Upload 401 Unauthorized Error

## Problem

Users getting "invalid or expired session" error when uploading KYC documents:
- Browser console shows: `POST https://tayosaecosystem.onrender.com/api/v1/storage/upload 401 (Unauthorized)`
- Error message: "invalid or expired session"

## Root Cause Analysis

The storage upload endpoint now requires Supabase JWT token validation (added in commit 2f5855f). The 401 error can happen for several reasons:

1. **Session token expired** - Supabase JWT tokens expire after a certain time
2. **Missing SUPABASE_ANON_KEY** - Backend can't validate tokens without this key
3. **User logged in before token validation was added** - Old session tokens might not work with new validation

## Solution Steps

### Step 1: Verify Render Environment Variables

1. Go to Render Dashboard: https://dashboard.render.com/
2. Click on "tayosaecosystem" service (or your object-storage-service)
3. Click "Environment" tab
4. **Verify these variables exist:**
   - `SUPABASE_URL` = `https://ablvrbnbsdqshrorhmjf.supabase.co`
   - `SUPABASE_ANON_KEY` = Your Supabase anon key (starts with `eyJ...`)
   - `SUPABASE_SERVICE_ROLE_KEY` = Your Supabase service role key

5. **If any are missing, add them:**
   - Click "Add Environment Variable"
   - Enter the key name and value
   - Click "Save Changes"
   - Wait for automatic redeployment

### Step 2: Get Fresh Supabase Keys (if needed)

If you don't have the keys or they're incorrect:

1. Go to Supabase Dashboard: https://supabase.com/dashboard/project/ablvrbnbsdqshrorhmjf
2. Click "Settings" (gear icon in left sidebar)
3. Click "API" under Project Settings
4. Copy these values:
   - **Project URL** → Use for `SUPABASE_URL`
   - **anon public** key → Use for `SUPABASE_ANON_KEY`
   - **service_role** key (click "Reveal" button) → Use for `SUPABASE_SERVICE_ROLE_KEY`

### Step 3: Force User to Log In Again

The simplest fix is to have users log out and log back in to get a fresh session token:

1. **Option A: Clear browser storage (recommended)**
   - Open browser console (F12)
   - Go to "Application" tab (Chrome) or "Storage" tab (Firefox)
   - Click "Session Storage" → `https://tayosaecosystem.vercel.app`
   - Right-click → "Clear"
   - Refresh the page
   - User will be logged out automatically

2. **Option B: Add logout button**
   - User clicks logout in the app
   - Then logs back in with email/password
   - This generates a fresh Supabase JWT token

### Step 4: Test Upload Again

After logging in with fresh credentials:
1. Go to KYC upload page
2. Select ID front, ID back, and selfie photos
3. Click "Continue"
4. Upload should now work

## Verification Checklist

- [ ] `SUPABASE_URL` is set in Render environment
- [ ] `SUPABASE_ANON_KEY` is set in Render environment
- [ ] `SUPABASE_SERVICE_ROLE_KEY` is set in Render environment
- [ ] Render service has redeployed after environment variable changes
- [ ] User has logged out and logged back in (fresh session token)
- [ ] Browser console shows no 401 errors
- [ ] Files upload successfully

## Technical Details

### Token Flow

1. User logs in → Backend calls Supabase `/auth/v1/token?grant_type=password`
2. Supabase returns `access_token` (JWT)
3. Frontend stores token in `sessionStorage.auth_token`
4. Frontend sends token in `Authorization: Bearer <token>` header
5. Backend validates token by calling Supabase `/auth/v1/user` endpoint
6. If valid, backend extracts user ID and allows upload

### Token Validation Code

```go
// services/object-storage-service/main.go
func validateSupabaseToken(token string) (userID string, err error) {
    url := supabaseBaseURL() + "/auth/v1/user"
    req, err := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
    
    resp, err := client.Do(req)
    // ... parse response and extract user ID
}
```

### Why Token Expires

Supabase JWT tokens have a default expiration time (usually 1 hour). After expiration:
- Token validation fails with 401 Unauthorized
- User must log in again to get a new token
- OR implement token refresh using `refresh_token`

## Next Steps (Optional Improvements)

1. **Implement token refresh** - Automatically refresh expired tokens using `refresh_token`
2. **Add better error messages** - Show "Session expired, please log in again" instead of generic error
3. **Add token expiration check** - Check token expiration before upload and prompt login if needed

## Files Modified

- `services/object-storage-service/main.go` - Added `requireAuth` middleware and token validation
- Commit: 2f5855f "fix: Add token validation to storage upload endpoint"
