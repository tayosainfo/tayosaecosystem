# Complete Registration Flow Test & Fix

## Current State Analysis (Using Supabase MCP)

### ✅ What's Working:
1. **Supabase Auth Registration**: 2 users successfully created
   - Edson Kariyo (`edisonkariyo@gmail.com`)
   - Tayosa Academy (`tayosainfo@gmail.com`)
2. **Email Verification**: Both users verified their emails
3. **Metadata Storage**: User names stored in `raw_user_meta_data`

### ❌ What's Broken:
1. **Phone Numbers Missing**: 
   - `phone` field in `auth.users` is `null` for both users
   - `raw_user_meta_data` does NOT contain phone field
   - This means phone was never sent to Supabase during registration

2. **Local Database Empty**:
   - `public.users_identity` table has 0 rows
   - No local profiles created
   - Phone numbers never stored anywhere

### 🎯 Root Cause:
The registration flow is incomplete. Users are:
1. ✅ Creating Supabase Auth accounts
2. ✅ Verifying their emails
3. ❌ **NOT completing the verification flow that creates local profiles**

## Why Phone Numbers Are Missing

Looking at the backend code (`services/user-service/handlers.go`), the phone number SHOULD be sent in the signup request:

```go
signupResp, _, err := supabasePostWithQuery("/auth/v1/signup", clientTypeQuery(r), map[string]any{
    "email":    contactEmail,
    "password": req.Password,
    "data": map[string]any{
        "name":  req.FullName,
        "phone": phoneE164,  // ← Phone should be here
    },
})
```

But the Supabase logs show the phone is NOT in the metadata. This means either:
1. The registration request never reached the backend properly, OR
2. The phone field was empty/null when the request was made

## The Complete Flow (How It Should Work)

### Step 1: Registration (`POST /api/v1/register`)
**Frontend sends:**
```json
{
  "fullName": "John Doe",
  "phone": "+256700123456",
  "email": "john@example.com",
  "password": "password123",
  "nationality": "UG",
  "termsAccepted": true,
  "privacyAccepted": true
}
```

**Backend does:**
1. Validates phone number format
2. Creates Supabase Auth user with email + password
3. Stores phone in Supabase metadata: `data.phone`
4. Sends OTP verification email
5. Returns response with `requireEmailVerification: true`

**Frontend does:**
1. Stores pending profile in sessionStorage:
   ```json
   {
     "fullName": "John Doe",
     "phone": "+256700123456",
     "nationality": "UG"
   }
   ```
2. Redirects to `/verify?email=john@example.com`

### Step 2: Email Verification (`POST /api/v1/verify-email`)
**Frontend sends:**
```json
{
  "email": "john@example.com",
  "otp": "123456",
  "fullName": "John Doe",        // ← From sessionStorage
  "phone": "+256700123456",       // ← From sessionStorage
  "nationality": "UG"             // ← From sessionStorage
}
```

**Backend does:**
1. Verifies OTP with Supabase
2. Gets Supabase user ID from verification response
3. **Creates local database record** in `public.users_identity`:
   ```sql
   INSERT INTO public.users_identity (
       user_id, full_name, phone_e164, 
       auth_email, contact_email, 
       supabase_user_id, created_at
   ) VALUES (...)
   ```
4. Creates onboarding profile
5. Returns session + user data

**Frontend does:**
1. Clears sessionStorage
2. Applies session
3. Redirects to `/onboarding`

## What Went Wrong for Existing Users

The existing 2 users completed Step 1 (registration + email verification) but **never completed Step 2** (the verification callback that creates local profiles).

**Evidence:**
- ✅ Users exist in `auth.users`
- ✅ Emails are verified (`email_confirmed_at` is set)
- ❌ `public.users_identity` is empty
- ❌ Phone numbers not in metadata

**Most Likely Cause:**
The verification email redirect URL is pointing to the wrong domain (Vercel preview URL), so users couldn't complete the OTP verification flow on the correct frontend.

## Fix Strategy

### Option 1: Fix Redirect URL + Re-register (Recommended)
1. Fix Supabase Site URL configuration
2. Have users register again with new email addresses
3. Complete full flow end-to-end

### Option 2: Manually Complete Registration for Existing Users
Since we have MCP access, we can manually create the missing local profiles for the 2 existing users.

**Required Information:**
- User ID: `3bb3e6ee-a1f4-47b7-a12c-cc28e7d7bd40` (Edson Kariyo)
- User ID: `429e70ae-5d2c-4ac3-8295-94549a7959ac` (Tayosa Academy)
- **Missing**: Phone numbers (need to ask users)

## Testing the Complete Flow

### Test Plan:
1. ✅ Check current database state (DONE)
2. ⏳ Fix Supabase redirect URL configuration
3. ⏳ Register new test user
4. ⏳ Verify email with OTP
5. ⏳ Confirm local profile created with phone number
6. ⏳ Verify user can login
7. ⏳ Verify user can access onboarding

### Expected Database State After Fix:
```sql
-- Supabase Auth
SELECT id, email, phone, raw_user_meta_data->>'phone' as metadata_phone
FROM auth.users
WHERE email = 'test@example.com';
-- Should show: phone in metadata

-- Local Database
SELECT user_id, full_name, phone_e164, contact_email
FROM public.users_identity
WHERE contact_email = 'test@example.com';
-- Should show: complete profile with phone number
```

## Next Steps

1. **Immediate**: Fix Supabase Site URL to production URL
2. **Test**: Register new user and complete full flow
3. **Decide**: What to do about existing 2 users
   - Option A: Ask them to re-register
   - Option B: Manually create their profiles (need phone numbers)
