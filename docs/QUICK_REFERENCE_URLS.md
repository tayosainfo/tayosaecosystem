# 🎯 Quick Reference - URLs for Supabase Configuration

## Copy and Paste These URLs

### Site URL (Main Website):
```
https://tayosaecosystem.vercel.app
```

**Why Vercel and not localhost?**
- Site URL is your "default" redirect
- Email verification links will use this URL
- If you use localhost, emails won't work for other people
- You can still develop locally by adding localhost to "Redirect URLs" below!

### Redirect URLs (Copy All 6):
```
https://tayosaecosystem.vercel.app/auth/callback
https://tayosaecosystem.vercel.app/auth/verify-email
https://tayosaecosystem.vercel.app/auth/reset-password
http://localhost:3000/auth/callback
http://localhost:3000/auth/verify-email
http://localhost:3000/auth/reset-password
```

**What this means:**
- ✅ You can develop locally on `http://localhost:3000`
- ✅ You can test on Vercel at `https://tayosaecosystem.vercel.app`
- ✅ Email links work for everyone (they go to Vercel)
- ✅ Both environments work perfectly!

---

## Where to Use These

**In Supabase Dashboard:**
1. Go to: Authentication → URL Configuration
2. Paste Site URL in "Site URL" field
3. Paste all 6 Redirect URLs in "Redirect URLs" field
4. Click "Save"

---

## For Later (Production)

When you deploy to `app.tayosa.com`, add these:

### Production Site URL:
```
https://app.tayosa.com
```

### Production Redirect URLs:
```
https://app.tayosa.com/auth/callback
https://app.tayosa.com/auth/verify-email
https://app.tayosa.com/auth/reset-password
```

### Mobile App URLs:
```
tayosa://auth/callback
tayosa://auth/verify-email
tayosa://auth/reset-password
```

**Note:** Add these later, don't remove the development URLs!

---

## Environment Variables

### For Vercel (Frontend):
```
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
VITE_API_BASE_URL=https://tayosaecosystem.onrender.com
```

### For Render (Backend):
```
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SERVICE_ROLE_KEY]
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
```

**Replace `[YOUR-PASSWORD]` with your actual Supabase database password.**

---

## Quick Checklist

- [ ] Copy Site URL above
- [ ] Paste in Supabase → Authentication → URL Configuration
- [ ] Copy all 6 Redirect URLs above
- [ ] Paste in Supabase → Authentication → URL Configuration → Redirect URLs
- [ ] Click "Save"
- [ ] Add environment variables to Vercel
- [ ] Add environment variables to Render
- [ ] Test registration on https://tayosaecosystem.vercel.app

---

**Need more details?** See `DEPLOYMENT_URLS_GUIDE.md`
