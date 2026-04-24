# Email Verification Token Fix - Complete Summary

## ✅ What Was Fixed

### Problem 1: "Token has expired or is invalid"
**Root Cause:** The backend was calling `/auth/v1/otp` to send a second OTP after signup, creating two different tokens:
1. Signup confirmation token (from `/auth/v1/signup`)
2. OTP token (from `/auth/v1/otp`)

Users received the OTP token in email, but the backend was trying to verify it as a signup confirmation token.

**Solution:** Removed the duplicate `/auth/v1/otp` call. Now only the signup confirmation token is sent.

### Problem 2: Users had to manually copy token from email
**Root Cause:** When users clicked the verification link in the email, they were redirected to the verify page but had to manually copy the token from the email.

**Solution:** Added auto-extraction of token from URL parameters. Now when users click the email link, the token is automatically filled in.

## 🔧 Changes Made

### 1. Backend (`services/user-service/handlers.go`)
**Removed:**
```go
// Send OTP for email verification
needsVerify := !supabaseSignupHasAccessToken(signupResp)
if needsVerify && supabaseConfigured() {
    // Send OTP using the dedicated OTP endpoint
    _, _, otpErr := supabasePostWithQuery("/auth/v1/otp", clientTypeQuery(r), map[string]any{
        "email": contactEmail,
        "type":  "signup",
    })
    // ... error handling
}
```

**Added:**
```go
"options": map[string]any{
    "email_redirect_to": "https://tayosaecosystem.vercel.app/verify",
},
```

### 2. Frontend (`src/pages/auth/VerifyEmail.tsx`)
**Added auto-token extraction:**
```typescript
useEffect(() => {
  const params = new URLSearchParams(window.location.search);
  const emailParam = params.get('email');
  const tokenParam = params.get('token');
  const tokenHashParam = params.get('token_hash');
  
  if (emailParam) setEmail(emailParam);
  
  // Auto-extract token from URL if user clicked email link
  if (tokenParam) {
    setOtp(tokenParam);
  } else if (tokenHashParam) {
    // Supabase sometimes uses token_hash instead
    setOtp(tokenHashParam);
  }
}, []);
```

## 🎯 How It Works Now

### Registration Flow:
1. **User registers** at `/register`
2. **Backend calls** `/auth/v1/signup` with email + password
3. **Supabase sends ONE email** with confirmation token
4. **User clicks link** in email
5. **Redirected to** `/verify?token=XXXXXX&email=user@example.com`
6. **Frontend auto-fills** the token field
7. **User clicks** "Verify & Continue" (or can manually enter token)
8. **Backend verifies** token with `/auth/v1/verify`
9. **Local profile created** with phone number
10. **User redirected** to `/onboarding`

### Email Link Format:
```
https://ablvrbnbsdqshrorhmjf.supabase.co/auth/v1/verify?
  token=7fca06fd9245cfee0f0fc441150fb1d357bbc995dfe075d34edb9ce9
  &type=signup
  &redirect_to=https://tayosaecosystem.vercel.app/verify
```

After Supabase processes it, user lands on:
```
https://tayosaecosystem.vercel.app/verify?
  token=7fca06fd9245cfee0f0fc441150fb1d357bbc995dfe075d34edb9ce9
  &email=user@example.com
```

## 🧪 Testing Instructions

### Test 1: Click Email Link (Recommended)
1. Register a new user
2. Check email
3. **Click the verification link** in the email
4. Token should auto-fill
5. Click "Verify & Continue"
6. Should redirect to `/onboarding`

### Test 2: Manual Token Entry
1. Register a new user
2. Check email
3. **Copy the token** from the email (the long string in the link)
4. Go to `/verify?email=your-email@example.com`
5. **Paste the token** manually
6. Click "Verify & Continue"
7. Should redirect to `/onboarding`

## 📊 Expected Database State After Successful Verification

### Supabase Auth (`auth.users`):
```sql
SELECT id, email, phone, raw_user_meta_data, email_confirmed_at
FROM auth.users
WHERE email = 'test@example.com';
```

**Expected:**
- `email_confirmed_at`: Should have a timestamp
- `raw_user_meta_data`: Should contain `{"name": "...", "phone": "+256..."}`

### Local Database (`public.users_identity`):
```sql
SELECT user_id, full_name, phone_e164, contact_email, supabase_user_id
FROM public.users_identity
WHERE contact_email = 'test@example.com';
```

**Expected:**
- Row should exist
- `phone_e164`: Should have the phone number (e.g., `+256700123456`)
- `supabase_user_id`: Should match the Supabase user ID

### Onboarding Profile (`public.onboarding_profiles`):
```sql
SELECT user_id, phase, trust_score_seed
FROM public.onboarding_profiles
WHERE user_id = (SELECT user_id FROM public.users_identity WHERE contact_email = 'test@example.com');
```

**Expected:**
- Row should exist
- `phase`: Should be `1`
- `trust_score_seed`: Should be `10`

## ⚠️ Important Notes

1. **Token Expiry:** Supabase tokens expire after 24 hours by default
2. **One-Time Use:** Each token can only be used once
3. **Email Template:** Make sure the email template shows `{{ .Token }}` for the 6-digit code display
4. **Redirect URL:** Must be configured in Supabase dashboard

## 🔄 Deployment

Changes have been pushed to GitHub and will auto-deploy to:
- **Vercel:** https://tayosaecosystem.vercel.app (frontend)
- **Render:** https://tayosaecosystem.onrender.com (backend)

Wait 2-3 minutes for deployment to complete, then test with a new user registration.

## 🐛 Troubleshooting

### If token still shows as expired:
1. Check Supabase logs for verification attempts
2. Verify the token in the email matches the token being sent to backend
3. Check that backend is calling `/auth/v1/verify` with correct parameters

### If token doesn't auto-fill:
1. Check browser console for JavaScript errors
2. Verify the URL has `?token=...` parameter
3. Try manually copying the token from the email

### If local profile isn't created:
1. Check backend logs for database errors
2. Verify phone number is in sessionStorage
3. Check that verification response includes user ID

## ✅ Success Criteria

- ✅ User can register
- ✅ User receives email with token
- ✅ User can click link and token auto-fills
- ✅ User can verify and see "success" message
- ✅ User is redirected to `/onboarding`
- ✅ Phone number is stored in database
- ✅ User can login with email + password
