# Admin Permission Fix - Step by Step Guide

## Problem
Admin user (baylesinfo@gmail.com) is getting a 403 Forbidden error when trying to access the KYC management endpoint at `tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending`.

## Root Cause
The custom JWT claims hook was adding `user_role` to the JWT token, but the API gateway was looking for it in the `app_metadata` field returned by the Supabase `/auth/v1/user` endpoint.

## Solution

### Step 1: Update RLS Policies (COMPLETED ✅)
You've already run the SQL to update the RLS policies on the `kyc_documents` table to allow admin users to view all KYC documents.

### Step 2: Fix Custom Claims Hook

Run this SQL in your Supabase SQL Editor:

```sql
-- Drop the existing function
DROP FUNCTION IF EXISTS public.custom_access_token_hook(jsonb);

-- Create updated custom access token hook function
CREATE OR REPLACE FUNCTION public.custom_access_token_hook(event jsonb)
RETURNS jsonb
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $
DECLARE
  claims jsonb;
  user_role text;
BEGIN
  -- Fetch the user role from users_identity table
  SELECT role INTO user_role
  FROM public.users_identity
  WHERE supabase_user_id = (event->>'user_id')::text;

  -- Get existing claims from the event
  claims := event->'claims';

  IF user_role IS NOT NULL THEN
    -- Set custom claim for user role in JWT claims
    claims := jsonb_set(claims, '{user_role}', to_jsonb(user_role));
    
    -- Also add to app_metadata for easier access via /auth/v1/user endpoint
    claims := jsonb_set(claims, '{app_metadata, user_role}', to_jsonb(user_role));
  ELSE
    -- Default to 'user' if no role found
    claims := jsonb_set(claims, '{user_role}', '"user"');
    claims := jsonb_set(claims, '{app_metadata, user_role}', '"user"');
  END IF;

  -- Update the 'claims' object in the original event
  event := jsonb_set(event, '{claims}', claims);

  RETURN event;
END;
$;

-- Grant necessary permissions
GRANT USAGE ON SCHEMA public TO supabase_auth_admin;
GRANT SELECT ON public.users_identity TO supabase_auth_admin;
GRANT EXECUTE ON FUNCTION public.custom_access_token_hook TO supabase_auth_admin;
```

### Step 3: Redeploy API Gateway Service

The `services/api-gateway-service/auth.go` file has been updated to check both `app_metadata` and `user_metadata` for the `user_role` field.

**To deploy:**

1. Commit the changes:
   ```bash
   git add services/api-gateway-service/auth.go
   git commit -m "Fix admin role extraction from JWT token"
   git push
   ```

2. Redeploy the api-gateway-service on Render (it should auto-deploy if you have auto-deploy enabled)

### Step 4: Admin User Re-login

**IMPORTANT:** The admin user (baylesinfo@gmail.com) must log out and log back in to get a new JWT token with the updated claims.

1. Log out from the admin dashboard
2. Log back in with the admin credentials
3. Try accessing the KYC endpoint again

## Verification

After completing all steps, test the admin endpoint:

```bash
# The admin user should now be able to access this endpoint
GET https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending
```

Expected result: 200 OK with a list of pending KYC submissions.

## Troubleshooting

If the issue persists:

1. **Check the admin user's role in the database:**
   ```sql
   SELECT user_id, supabase_login_email, role 
   FROM users_identity 
   WHERE supabase_login_email = 'baylesinfo@gmail.com';
   ```
   Expected: `role = 'admin'`

2. **Verify the custom claims hook is enabled:**
   - Go to Supabase Dashboard → Authentication → Hooks
   - Ensure "Custom Access Token" hook is enabled
   - Ensure it points to `public.custom_access_token_hook`

3. **Check the JWT token:**
   - After logging in, inspect the JWT token in the browser's developer tools
   - Decode it at https://jwt.io
   - Verify it contains `"user_role": "admin"` in the claims

4. **Check API Gateway logs:**
   - Look for log messages like `"AUTH_METHOD: jwt, user_role: admin"`
   - If you see `"user role 'user' is not admin"`, the role is not being extracted correctly

## Files Changed

- `services/api-gateway-service/auth.go` - Updated to check both app_metadata and user_metadata
- `db/migrations/018_add_admin_rls_policies.sql` - Added admin exception to RLS policies
- `db/migrations/020_fix_custom_claims_app_metadata.sql` - Fixed custom claims hook to add role to app_metadata
