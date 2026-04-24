# MCP-Powered Complete Registration Flow Fix

## What I Can Do With Supabase MCP Access

### ✅ Database Operations I Can Perform:
1. **Query Data**: Check users, profiles, and all tables
2. **Insert Data**: Create missing local profiles manually
3. **Update Data**: Fix existing records
4. **Execute Migrations**: Apply schema changes
5. **Check Logs**: View auth, API, and service logs
6. **Get Advisors**: Security and performance recommendations

### ❌ What I Cannot Do:
1. **Change Supabase Dashboard Settings**: Site URL, redirect URLs, email templates
2. **Send Emails**: Trigger verification emails
3. **Generate OTP Codes**: Create verification codes
4. **Access Frontend**: Can't interact with the web UI directly

## Action Plan to Complete the Flow

### Phase 1: Fix Configuration (REQUIRES YOUR ACTION)
**You must do this in Supabase Dashboard:**

1. Go to https://supabase.com/dashboard/project/ablvrbnbsdqshrorhmjf
2. Click "Authentication" → "URL Configuration"
3. Set **Site URL** to: `https://tayosaecosystem.vercel.app`
4. Add to **Redirect URLs**:
   ```
   https://tayosaecosystem.vercel.app/**
   http://localhost:5173/**
   ```
5. Click "Save"

### Phase 2: Fix Existing Users (I CAN DO THIS WITH MCP)
**For the 2 existing users who completed email verification but have no local profiles:**

#### User 1: Edson Kariyo
- Supabase ID: `3bb3e6ee-a1f4-47b7-a12c-cc28e7d7bd40`
- Email: `edisonkariyo@gmail.com`
- Name: `Edson Kariyo`
- **Missing**: Phone number

#### User 2: Tayosa Academy
- Supabase ID: `429e70ae-5d2c-4ac3-8295-94549a7959ac`
- Email: `tayosainfo@gmail.com`
- Name: `Tayosa Academy`
- **Missing**: Phone number

**What I need from you:**
Please provide the phone numbers for these 2 users, then I can:
1. Insert records into `public.users_identity` with phone numbers
2. Create `public.onboarding_profiles` records
3. Set up initial user data

**SQL I will execute:**
```sql
-- For Edson Kariyo
INSERT INTO public.users_identity (
    user_id, 
    full_name, 
    phone_e164, 
    auth_email, 
    contact_email, 
    supabase_login_email,
    supabase_user_id, 
    nationality,
    created_at
) VALUES (
    '3bb3e6ee-a1f4-47b7-a12c-cc28e7d7bd40',
    'Edson Kariyo',
    '+256XXXXXXXXX',  -- Need phone number
    'edisonkariyo@gmail.com',
    'edisonkariyo@gmail.com',
    'edisonkariyo@gmail.com',
    '3bb3e6ee-a1f4-47b7-a12c-cc28e7d7bd40',
    'UG',
    '2026-04-23 12:31:32'
);

INSERT INTO public.onboarding_profiles (
    user_id,
    phase,
    trust_score_seed,
    updated_at
) VALUES (
    '3bb3e6ee-a1f4-47b7-a12c-cc28e7d7bd40',
    1,
    10,
    NOW()
);

-- Same for Tayosa Academy with their phone number
```

### Phase 3: Test New User Registration (REQUIRES YOUR ACTION + MY VERIFICATION)
**After you fix the Site URL:**

1. **You register a new test user** on https://tayosaecosystem.vercel.app/register
2. **You verify the email** by entering the OTP code
3. **I verify in database** that:
   - User created in `auth.users`
   - Phone number in metadata
   - Local profile created in `public.users_identity`
   - Onboarding profile created

### Phase 4: Verify Complete Flow (I CAN DO THIS WITH MCP)
**After new user registration, I will check:**

```sql
-- Check Supabase Auth
SELECT 
    id, 
    email, 
    phone, 
    raw_user_meta_data->>'phone' as metadata_phone,
    raw_user_meta_data->>'name' as metadata_name,
    email_confirmed_at
FROM auth.users 
WHERE email = 'your-test-email@example.com';

-- Check Local Profile
SELECT 
    user_id, 
    full_name, 
    phone_e164, 
    contact_email,
    supabase_user_id,
    created_at
FROM public.users_identity 
WHERE contact_email = 'your-test-email@example.com';

-- Check Onboarding
SELECT 
    user_id, 
    phase, 
    trust_score_seed
FROM public.onboarding_profiles 
WHERE user_id = (
    SELECT user_id FROM public.users_identity 
    WHERE contact_email = 'your-test-email@example.com'
);
```

## What I Need From You Right Now

### Option A: Fix Existing Users (Fastest)
**Provide me with:**
1. Edson Kariyo's phone number: `+256XXXXXXXXX`
2. Tayosa Academy's phone number: `+256XXXXXXXXX`

**I will:**
1. Create their local profiles immediately
2. Set up onboarding records
3. They can login and use the system

### Option B: Test New User Flow (Most Thorough)
**You do:**
1. Fix Supabase Site URL (see Phase 1 above)
2. Register a new test user
3. Tell me the email address

**I will:**
1. Monitor the database in real-time
2. Verify each step completes correctly
3. Confirm phone number is captured
4. Validate the complete flow

### Option C: Both (Recommended)
1. Give me phone numbers → I fix existing users
2. You fix Site URL → We test new registration
3. Confirm everything works end-to-end

## Summary

**What's blocking the flow:**
- Redirect URL configuration (you must fix in dashboard)
- Missing local profiles for existing users (I can fix with phone numbers)

**What I can verify with MCP:**
- Database state at each step
- Phone number capture
- Profile creation
- Complete data flow

**What you need to do:**
1. Fix Supabase Site URL configuration
2. Provide phone numbers for existing users (optional)
3. Test new user registration

Tell me which option you want to proceed with!
