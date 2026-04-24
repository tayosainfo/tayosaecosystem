# Phone Number Capture Fix - Step-by-Step Guide

## Problem Summary
Users are registering successfully in Supabase Auth, but phone numbers are not being captured in the local database because the email verification redirect URL is incorrect.

## Root Cause
The verification email links are redirecting users to the wrong URL (Vercel preview URL), preventing them from completing the verification flow and creating their local profile with phone number.

## Solution: Fix Supabase Redirect URL Configuration

### Step 1: Open Supabase Dashboard
1. Go to https://supabase.com/dashboard
2. Click on your project: `ablvrbnbsdqshrorhmjf`

### Step 2: Navigate to Authentication Settings
1. In the left sidebar, click **"Authentication"**
2. Click **"URL Configuration"** tab

### Step 3: Update Site URL
1. Find the field labeled **"Site URL"**
2. **CRITICAL:** Set it to your PRODUCTION URL:
   ```
   https://tayosaecosystem.vercel.app
   ```
3. **DO NOT** use preview URLs like:
   - ❌ `https://tayosaecosystem-git-main-*.vercel.app`
   - ❌ `https://tayosaecosystem-*.vercel.app`
   - ✅ Use: `https://tayosaecosystem.vercel.app` (production only)

### Step 4: Update Redirect URLs
1. Find the field labeled **"Redirect URLs"**
2. Make sure it includes:
   ```
   https://tayosaecosystem.vercel.app/**
   http://localhost:5173/**
   ```
3. You can have multiple URLs (one per line), but the Site URL should be your main production URL

### Step 5: Verify Email Template
1. In the left sidebar, click **"Email Templates"**
2. Click on **"Confirm signup"** template
3. Verify the template uses the variable: `{{ .ConfirmationURL }}`
4. **DO NOT** hardcode any URL in the template
5. The template should look something like:
   ```html
   <p>Click this link to confirm your email:</p>
   <p><a href="{{ .ConfirmationURL }}">Confirm your email</a></p>
   ```

### Step 6: Save Changes
1. Click the **"Save"** button at the bottom of each page
2. Wait for the success message

### Step 7: Test with New User
1. Register a NEW user (use a different email)
2. Check the verification email
3. Click the verification link
4. Verify you're redirected to: `https://tayosaecosystem.vercel.app/verify?email=...`
5. Enter the 6-digit OTP code
6. Complete verification

### Step 8: Verify Phone Number is Captured
After a new user completes verification, check the database:

```sql
-- Check if user exists in Supabase Auth
SELECT id, email, phone, raw_user_meta_data 
FROM auth.users 
WHERE email = 'new-user-email@example.com';

-- Check if local profile was created with phone number
SELECT user_id, full_name, phone_e164, contact_email 
FROM public.users_identity 
WHERE contact_email = 'new-user-email@example.com';
```

The `public.users_identity` table should now have a record with the phone number.

## What About Existing Users?

The 2 existing users (Edson Kariyo and Tayosa Academy) who already verified their emails but don't have local profiles will need to:

### Option 1: Re-register (Recommended)
1. Use a different email address
2. Complete the full registration flow again
3. This time the verification redirect will work correctly

### Option 2: Manual Database Entry (Advanced)
If you want to keep the existing Supabase Auth users, you'll need to manually create their local profiles:

```sql
-- Get the Supabase user ID
SELECT id, email, raw_user_meta_data 
FROM auth.users 
WHERE email = 'edisonkariyo@gmail.com';

-- Manually insert into users_identity (replace values with actual data)
INSERT INTO public.users_identity (
    user_id, 
    full_name, 
    phone_e164, 
    auth_email, 
    contact_email, 
    supabase_login_email, 
    supabase_user_id, 
    created_at
) VALUES (
    '3bb3e6ee-a1f4-47b7-a12c-cc28e7d7bd40',  -- Supabase user ID
    'Edson Kariyo',                           -- Full name
    '+256700123456',                          -- Phone number (you'll need to ask the user)
    'edisonkariyo@gmail.com',
    'edisonkariyo@gmail.com',
    'edisonkariyo@gmail.com',
    '3bb3e6ee-a1f4-47b7-a12c-cc28e7d7bd40',
    NOW()
);

-- Create onboarding profile
INSERT INTO public.onboarding_profiles (
    user_id,
    phase,
    trust_score_seed,
    last_updated_at
) VALUES (
    '3bb3e6ee-a1f4-47b7-a12c-cc28e7d7bd40',
    1,
    10,
    NOW()
);
```

## Expected Behavior After Fix

### Registration Flow:
1. User fills registration form (email, phone, password, name)
2. Backend creates Supabase Auth user
3. User receives verification email
4. User clicks link → redirected to `https://tayosaecosystem.vercel.app/verify?email=...`
5. User enters 6-digit OTP code
6. Backend verifies OTP and creates local profile with phone number
7. User redirected to `/onboarding`

### Database State:
- ✅ User in `auth.users` (Supabase Auth)
- ✅ User in `public.users_identity` (local database with phone number)
- ✅ User in `public.onboarding_profiles` (onboarding state)

## Verification Checklist

- [ ] Site URL set to production URL in Supabase dashboard
- [ ] Redirect URLs include production and localhost
- [ ] Email template uses `{{ .ConfirmationURL }}` variable
- [ ] All changes saved in Supabase dashboard
- [ ] Test registration with new user
- [ ] Verification email redirects to correct URL
- [ ] OTP verification completes successfully
- [ ] Phone number appears in `public.users_identity` table

## Need Help?

If you're still having issues after following these steps, check:
1. Browser console for JavaScript errors
2. Network tab to see API requests/responses
3. Supabase Auth logs for verification errors
4. Backend logs for database insertion errors
