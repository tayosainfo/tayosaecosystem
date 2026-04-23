# 🌐 Deployment URLs Configuration Guide

## Your Current Setup

**Frontend (Vercel):** `https://tayosaecosystem.vercel.app`  
**Backend (Render):** `https://tayosaecosystem.onrender.com`  
**Production Domain:** `https://app.tayosa.com` (planned for later)  
**Mobile App:** `tayosa://` (planned for later)

---

## 📋 Supabase URL Configuration

### Current Configuration (For Development)

Use these URLs in your Supabase dashboard **right now**:

#### Site URL:
```
https://tayosaecosystem.vercel.app
```

#### Redirect URLs:
```
https://tayosaecosystem.vercel.app/auth/callback
https://tayosaecosystem.vercel.app/auth/verify-email
https://tayosaecosystem.vercel.app/auth/reset-password
http://localhost:3000/auth/callback
http://localhost:3000/auth/verify-email
http://localhost:3000/auth/reset-password
```

**Why these URLs?**
- `https://tayosaecosystem.vercel.app/*` - Your Vercel deployment (live development)
- `http://localhost:3000/*` - Local development on your machine

---

## 🚀 When to Add Production URLs

### When You're Ready for Production

When you deploy to `app.tayosa.com`, come back and **add** these URLs (don't remove the development ones):

#### Update Site URL to:
```
https://app.tayosa.com
```

#### Add These Redirect URLs:
```
https://app.tayosa.com/auth/callback
https://app.tayosa.com/auth/verify-email
https://app.tayosa.com/auth/reset-password
```

**Keep the development URLs!** This lets you test on both environments.

---

## 📱 When to Add Mobile App URLs

### When You're Ready to Test Mobile

When you start testing the mobile app, **add** these URLs:

```
tayosa://auth/callback
tayosa://auth/verify-email
tayosa://auth/reset-password
```

**What are these?**
- `tayosa://` is a custom URL scheme for your mobile app
- It allows the app to handle authentication redirects
- Only needed when testing the mobile app

---

## 🔧 Environment Variables for Each Environment

### Development (Current)

**Frontend (.env or Vercel Environment Variables):**
```bash
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
VITE_API_BASE_URL=https://tayosaecosystem.onrender.com
```

**Backend (Render Environment Variables):**
```bash
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SERVICE_ROLE_KEY]
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
```

**Mobile App (Flutter --dart-define):**
```bash
flutter build apk \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=[YOUR_ANON_KEY] \
  --dart-define=API_BASE_URL=https://tayosaecosystem.onrender.com
```

### Production (When Ready)

**Frontend (.env.production or Vercel Production Environment Variables):**
```bash
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
VITE_API_BASE_URL=https://api.tayosa.com  # Your production backend
```

**Backend (Production Environment Variables):**
```bash
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SERVICE_ROLE_KEY]
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
```

---

## 📝 Step-by-Step: Configuring Vercel

### 1. Set Environment Variables in Vercel

1. Go to your Vercel dashboard: https://vercel.com/dashboard
2. Select your project: `tayosaecosystem`
3. Click "Settings" tab
4. Click "Environment Variables" in the left sidebar
5. Add these variables:

| Name | Value | Environment |
|------|-------|-------------|
| `VITE_SUPABASE_URL` | `https://[YOUR-PROJECT-REF].supabase.co` | Production, Preview, Development |
| `VITE_SUPABASE_ANON_KEY` | `[YOUR_ANON_KEY]` | Production, Preview, Development |
| `VITE_API_BASE_URL` | `https://tayosaecosystem.onrender.com` | Production, Preview, Development |

6. Click "Save"
7. Redeploy your app for changes to take effect

---

## 📝 Step-by-Step: Configuring Render

### 1. Set Environment Variables in Render

1. Go to your Render dashboard: https://dashboard.render.com
2. Select your service: `tayosaecosystem`
3. Click "Environment" tab in the left sidebar
4. Add these environment variables:

| Key | Value |
|-----|-------|
| `SUPABASE_URL` | `https://[YOUR-PROJECT-REF].supabase.co` |
| `SUPABASE_ANON_KEY` | `[YOUR_ANON_KEY]` |
| `SUPABASE_SERVICE_ROLE_KEY` | `[YOUR_SERVICE_ROLE_KEY]` |
| `DATABASE_URL` | `postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres` |

5. Click "Save Changes"
6. Render will automatically redeploy your service

**Note:** Replace `[YOUR-PASSWORD]` with your actual Supabase database password.

---

## 🔄 Migration Timeline

### Phase 1: Development (Current) ✅

**Status:** Active  
**URLs:** Vercel + Render  
**Users:** Development team only

**Supabase Configuration:**
- Site URL: `https://tayosaecosystem.vercel.app`
- Redirect URLs: Vercel + localhost

### Phase 2: Production Preparation (Later)

**Status:** Planned  
**URLs:** `app.tayosa.com` + production backend  
**Users:** Beta testers

**Supabase Configuration:**
- Site URL: `https://app.tayosa.com`
- Redirect URLs: Add production URLs (keep development URLs)

### Phase 3: Mobile App Testing (Later)

**Status:** Planned  
**URLs:** Mobile app with `tayosa://` scheme  
**Users:** Mobile beta testers

**Supabase Configuration:**
- Add mobile redirect URLs: `tayosa://auth/*`

### Phase 4: Full Production (Later)

**Status:** Planned  
**URLs:** All production domains  
**Users:** All users

**Supabase Configuration:**
- All URLs configured
- Remove development URLs if desired

---

## ✅ Current Action Items

### For You (Right Now):

1. **Configure Supabase with development URLs:**
   - Site URL: `https://tayosaecosystem.vercel.app`
   - Redirect URLs: The 6 development URLs listed above

2. **Set environment variables in Vercel:**
   - Add the 3 Supabase variables
   - Redeploy

3. **Set environment variables in Render:**
   - Add the 4 Supabase variables
   - Let it auto-redeploy

### For Later (When Ready):

4. **Add production URLs to Supabase:**
   - When you deploy to `app.tayosa.com`
   - Add the production redirect URLs

5. **Add mobile URLs to Supabase:**
   - When you start testing mobile app
   - Add the `tayosa://` redirect URLs

---

## 🧪 Testing Your Setup

### Test Development URLs:

1. **Open your Vercel deployment:** https://tayosaecosystem.vercel.app
2. **Try to register a new user**
3. **Check your email for verification link**
4. **Click the verification link**
5. **You should be redirected to:** `https://tayosaecosystem.vercel.app/auth/verify-email`

**If it works:** ✅ Configuration is correct!  
**If it fails:** Check that the redirect URL is added in Supabase

### Test Local Development:

1. **Run your app locally:** `npm run dev` (should start on http://localhost:3000)
2. **Try to register a new user**
3. **Check your email for verification link**
4. **Click the verification link**
5. **You should be redirected to:** `http://localhost:3000/auth/verify-email`

**If it works:** ✅ Local configuration is correct!  
**If it fails:** Check that localhost URLs are added in Supabase

---

## 🔒 Security Notes

### Public vs Secret Keys:

**Anon Key (Public):** `[YOUR_ANON_KEY]`
- ✅ Safe to use in frontend code
- ✅ Safe to commit to GitHub (in .env.example)
- ✅ Safe to expose in browser

**Service Role Key (Secret):** `[YOUR_SERVICE_ROLE_KEY]`
- ❌ NEVER use in frontend code
- ❌ NEVER commit to GitHub
- ❌ NEVER expose in browser
- ✅ Only use in backend services (Render)

### GitHub Repository:

**Safe to commit:**
- `.env.example` with placeholder values
- Documentation with anon key

**NEVER commit:**
- `.env` files with actual values
- Service role key anywhere in frontend code

**Add to .gitignore:**
```
.env
.env.local
.env.production
.env.development
```

---

## 📋 Checklist

### Current Setup (Do Now):

- [ ] Configure Supabase Site URL: `https://tayosaecosystem.vercel.app`
- [ ] Add 6 development redirect URLs to Supabase
- [ ] Set 3 environment variables in Vercel
- [ ] Set 4 environment variables in Render
- [ ] Test registration on Vercel deployment
- [ ] Test registration on localhost
- [ ] Verify email redirects work correctly

### Production Setup (Do Later):

- [ ] Purchase/configure `app.tayosa.com` domain
- [ ] Update Supabase Site URL to production
- [ ] Add production redirect URLs to Supabase
- [ ] Update Vercel environment variables for production
- [ ] Update Render environment variables for production
- [ ] Test production deployment

### Mobile Setup (Do Later):

- [ ] Configure mobile app with `tayosa://` URL scheme
- [ ] Add mobile redirect URLs to Supabase
- [ ] Test mobile authentication flow
- [ ] Verify mobile redirects work correctly

---

## 📞 Need Help?

### Common Issues:

**Issue:** "Redirect URL not allowed"  
**Solution:** Make sure the exact URL is added in Supabase → Authentication → URL Configuration

**Issue:** "Invalid redirect URL"  
**Solution:** Check for typos in the URL. It must match exactly (including https:// or http://)

**Issue:** "Email verification link doesn't work"  
**Solution:** Check that the redirect URL in the email matches one of your configured URLs

**Issue:** "Environment variables not working"  
**Solution:** Make sure you redeployed after adding environment variables

---

## 🎯 Summary

**For Development (Now):**
- Use Vercel URL: `https://tayosaecosystem.vercel.app`
- Use Render URL: `https://tayosaecosystem.onrender.com`
- Configure 6 redirect URLs in Supabase
- Set environment variables in Vercel and Render

**For Production (Later):**
- Add `app.tayosa.com` URLs when ready
- Keep development URLs for testing
- Update environment variables as needed

**For Mobile (Later):**
- Add `tayosa://` URLs when testing mobile
- Configure mobile app with Supabase credentials

---

**Last Updated:** April 22, 2026  
**Current Environment:** Development (Vercel + Render)  
**Next Environment:** Production (app.tayosa.com)
