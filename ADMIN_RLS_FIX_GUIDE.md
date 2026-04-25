# Admin RLS Permission Fix Guide

## Problem
The admin user (baylesinfo@gmail.com) is getting a **403 Forbidden** error when trying to access the KYC management endpoint:
```
GET https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending → 403 (Forbidden)
```

## Root Cause
The Row-Level Security (RLS) policies on the `kyc_documents` table only allow:
1. Users to view/manage their own KYC documents
2. Service role to have full access

However, the admin user needs to be able to view **all** KYC documents for management purposes. The current RLS policies don't have an admin exception.

## Solution
Update the RLS policies to allow admin users (those with `role = 'admin'` in the `users_identity` table) to view all KYC documents.

## How to Apply the Fix

### Option 1: Using Supabase Dashboard (Recommended)

1. **Open Supabase Dashboard**
   - Go to: https://app.supabase.com
   - Select your project: `ablvrbnbsdqshrorhmjf`

2. **Navigate to SQL Editor**
   - Click on "SQL Editor" in the left sidebar
   - Click "New Query"

3. **Copy and Paste the SQL**
   - Open the file: `SUPABASE_FIX_ADMIN_RLS.sql` in this repository
   - Copy the entire contents
   - Paste into the Supabase SQL Editor

4. **Execute the Query**
   - Click the "Run" button (or press Ctrl+Enter)
   - Wait for the query to complete
   - You should see a success message

5. **Verify the Fix**
   - Test the endpoint: `GET https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending`
   - The admin user should now be able to access the KYC queue

### Option 2: Using Supabase CLI (If Local Setup Works)

```bash
# Navigate to the project directory
cd tayosaecosystem

# Apply the migration
supabase db push

# This will apply all pending migrations including:
# - 018_add_admin_rls_policies.sql
# - 019_create_exec_sql_function.sql
```

## What Changed

### Before
```sql
-- Only users could view their own documents
CREATE POLICY "Users can view own KYC documents"
ON public.kyc_documents
FOR SELECT
TO authenticated
USING (user_id = auth.uid()::text);
```

### After
```sql
-- Users can view their own documents OR admins can view all
CREATE POLICY "Users can view own KYC documents"
ON public.kyc_documents
FOR SELECT
TO authenticated
USING (
  user_id = auth.uid()::text
  OR
  -- Allow admin users to view all KYC documents
  EXISTS (
    SELECT 1 FROM public.users_identity
    WHERE supabase_user_id = auth.uid()::text
    AND role = 'admin'
  )
);
```

## Verification

After applying the fix, verify that:

1. **Admin user can access KYC endpoint**
   ```bash
   curl -H "Authorization: Bearer <admin_token>" \
     https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending
   ```
   Should return: `200 OK` with KYC queue data

2. **Regular users still can't access admin endpoint**
   ```bash
   curl -H "Authorization: Bearer <user_token>" \
     https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending
   ```
   Should return: `403 Forbidden`

3. **Regular users can still access their own KYC documents**
   - This is handled by the backend API, not RLS

## Files Modified

- `db/migrations/018_add_admin_rls_policies.sql` - Updated RLS policies
- `db/migrations/019_create_exec_sql_function.sql` - Helper function for SQL execution
- `SUPABASE_FIX_ADMIN_RLS.sql` - Complete SQL script for manual application

## Troubleshooting

### Issue: "Could not find the function public.exec_sql"
**Solution**: The `exec_sql` function hasn't been created yet. Run the SQL script which includes the function creation.

### Issue: "Permission denied for schema public"
**Solution**: Make sure you're using the Service Role Key, not the Anon Key. The Service Role Key has elevated permissions.

### Issue: Admin user still getting 403 after applying fix
**Solution**: 
1. Verify the user has `role = 'admin'` in the `users_identity` table
2. Check that the custom JWT claims hook is enabled in Supabase Auth settings
3. Clear browser cache and re-login to get a fresh JWT token

## Additional Notes

- The fix maintains security by only allowing admins to view KYC documents
- Regular users can still only view their own documents
- The service role still has full access for backend operations
- The `users_identity` table RLS policy was also updated to allow reading the role field (needed for the JWT claims hook)

## Support

If you encounter any issues:
1. Check the Supabase logs: Dashboard → Logs → API
2. Verify the admin user's role: Dashboard → SQL Editor → `SELECT * FROM users_identity WHERE supabase_user_id = '<user_id>'`
3. Check the custom JWT claims hook is enabled: Dashboard → Authentication → Hooks
