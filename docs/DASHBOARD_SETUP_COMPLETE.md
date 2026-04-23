# ✅ Supabase Dashboard Setup - COMPLETE!

**Date Completed:** April 22, 2026  
**Completed By:** Project Owner  
**Status:** Ready for Developer Team

---

## 🎉 What You Completed

Congratulations! You've successfully configured the Supabase dashboard. Here's what you did:

### ✅ Step 1: JWT Settings
- Verified JWT token expiry settings
- Confirmed default settings are correct

### ✅ Step 2: URL Configuration
- **Site URL:** `https://tayosaecosystem.vercel.app`
- **Redirect URLs (6 total):**
  - `https://tayosaecosystem.vercel.app/auth/callback`
  - `https://tayosaecosystem.vercel.app/auth/verify-email`
  - `https://tayosaecosystem.vercel.app/auth/reset-password`
  - `http://localhost:3000/auth/callback`
  - `http://localhost:3000/auth/verify-email`
  - `http://localhost:3000/auth/reset-password`

### ✅ Step 3: Email Login
- Enabled email authentication
- Enabled email confirmations
- Set minimum password length to 8 characters

### ✅ Step 4: Email Templates
- Customized email verification template
- Customized password reset template
- Added Tayosa branding

### ✅ Step 5: Storage Bucket
- Created `collateral_docs` bucket
- Set to private (not public)
- Set file size limit to 50 MB

### ⏳ Step 6: Test Email (Optional - Skipped)
- Test email feature not available or skipped
- Will test after deployment instead

---

## 🚨 IMPORTANT: RLS Policies Not Set Up Yet

**Why:** The database tables don't exist in Supabase yet.

**What this means:** 
- The tables are still in your old database
- The migration needs to be run first
- RLS policies will be set up after migration

**Who will do it:** Your developer team (see instructions below)

---

## 📋 For Your Developer Team - Next Steps

### Step 1: Set Environment Variables

#### Vercel (Frontend):
```bash
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
VITE_API_BASE_URL=https://tayosaecosystem.onrender.com
```

**How to add:**
1. Go to Vercel dashboard
2. Select project: `tayosaecosystem`
3. Settings → Environment Variables
4. Add the 3 variables above
5. Select: Production, Preview, Development
6. Save and redeploy

#### Render (Backend):
```bash
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SERVICE_ROLE_KEY]
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
```

**How to add:**
1. Go to Render dashboard
2. Select service: `tayosaecosystem`
3. Environment tab
4. Add the 4 variables above
5. Replace `[YOUR-PASSWORD]` with actual Supabase database password
6. Save (auto-redeploys)

---

### Step 2: Run Database Migration

**File:** `db/migrations/011_rename_insforge_to_supabase.sql`

**What it does:**
- Renames `insforge_user_id` → `supabase_user_id`
- Renames `insforge_login_email` → `supabase_login_email`
- Updates indexes

**How to run:**

**Option A: Using psql (Recommended)**
```bash
# Get the DATABASE_URL from Supabase dashboard or SUPABASEDATABASECONNECTION.txt
export DATABASE_URL="postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres"

# Run the migration
psql $DATABASE_URL -f db/migrations/011_rename_insforge_to_supabase.sql
```

**Option B: Using Supabase SQL Editor**
1. Open Supabase dashboard → SQL Editor
2. Open the migration file: `db/migrations/011_rename_insforge_to_supabase.sql`
3. Copy the entire contents
4. Paste into SQL Editor
5. Click "Run"

**Verify migration:**
```sql
-- Check that columns were renamed
SELECT column_name 
FROM information_schema.columns 
WHERE table_name = 'users' 
  AND column_name IN ('supabase_user_id', 'supabase_login_email');
```

You should see both columns listed.

---

### Step 3: Set Up RLS Policies

**File:** `docs/SETUP_RLS_POLICIES.md`

**CRITICAL:** This MUST be done because automatic RLS is enabled!

**How to do it:**
1. Open Supabase dashboard → SQL Editor
2. Follow the guide: `docs/SETUP_RLS_POLICIES.md`
3. Run all 4 SQL code blocks (one at a time)
4. Verify ~13 policies are created

**Quick verification:**
1. Go to: Database → Policies
2. You should see policies for:
   - users table (3 policies)
   - transactions table (3 policies)
   - accounts table (3 policies)
   - storage.objects (4 policies)

**If you get "table does not exist" errors:**
- That table hasn't been created yet
- Skip it for now
- You can add policies for that table later when it's created

---

### Step 4: Deploy Application

**Frontend (Vercel):**
- Should auto-deploy when you push to GitHub
- Or manually trigger deployment from Vercel dashboard
- Verify environment variables are set

**Backend (Render):**
- Should auto-deploy when you push to GitHub
- Or manually trigger deployment from Render dashboard
- Verify environment variables are set

**Mobile App (Flutter):**
- Build with Supabase credentials:
```bash
flutter build apk \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=[YOUR_ANON_KEY] \
  --dart-define=API_BASE_URL=https://tayosaecosystem.onrender.com
```

---

### Step 5: Test Everything

**Test Scenarios:**

1. **User Registration:**
   - Go to: https://tayosaecosystem.vercel.app
   - Click "Register"
   - Enter email and password
   - Check email for verification link
   - Click verification link
   - Should redirect to Vercel app
   - Should be able to log in

2. **User Login:**
   - Go to: https://tayosaecosystem.vercel.app
   - Click "Login"
   - Enter credentials
   - Should successfully log in
   - Should see user dashboard

3. **Password Reset:**
   - Go to: https://tayosaecosystem.vercel.app
   - Click "Forgot Password"
   - Enter email
   - Check email for reset link
   - Click reset link
   - Enter new password
   - Should be able to log in with new password

4. **RLS Verification:**
   - Register two different users
   - Log in as User A
   - Verify User A can only see their own data
   - Log in as User B
   - Verify User B can only see their own data
   - Verify User B cannot see User A's data

5. **Backend API:**
   - Test protected endpoints
   - Verify token validation works
   - Check that unauthorized requests are rejected

---

## 📊 Deployment Checklist

### Pre-Deployment:
- [ ] Environment variables set in Vercel
- [ ] Environment variables set in Render
- [ ] Database migration file reviewed
- [ ] RLS policies reviewed
- [ ] Code committed to GitHub

### Deployment:
- [ ] Database migration executed successfully
- [ ] RLS policies created (verify in Supabase dashboard)
- [ ] Frontend deployed to Vercel
- [ ] Backend deployed to Render
- [ ] Mobile app built with Supabase credentials

### Post-Deployment:
- [ ] User registration works
- [ ] Email verification works
- [ ] User login works
- [ ] Password reset works
- [ ] RLS policies work (users can only see own data)
- [ ] Backend API works
- [ ] Mobile app works

---

## 🔒 Security Reminders

### Public Keys (Safe to Use Anywhere):
- **Anon Key:** `[YOUR_ANON_KEY]`
  - ✅ Use in frontend
  - ✅ Use in mobile app
  - ✅ Safe to commit to GitHub (in .env.example)

### Secret Keys (NEVER Expose):
- **Service Role Key:** `[YOUR_SERVICE_ROLE_KEY]`
  - ❌ NEVER use in frontend
  - ❌ NEVER commit to GitHub
  - ✅ Only use in backend (Render)
  - ✅ Only in environment variables

---

## 📞 Support

**If you encounter issues:**

1. **Check the documentation:**
   - `docs/DEPLOYMENT_GUIDE.md`
   - `docs/DATABASE_MIGRATION_PLAN.md`
   - `docs/SETUP_RLS_POLICIES.md`
   - `docs/DEPLOYMENT_URLS_GUIDE.md`

2. **Common issues:**
   - "Redirect URL not allowed" → Check URL Configuration in Supabase
   - "Table does not exist" → Run database migration first
   - "Permission denied" → Check RLS policies are set up
   - "Invalid token" → Check environment variables

3. **Contact:**
   - Email: support@tayosa.com
   - Check Supabase logs: Dashboard → Logs
   - Check Vercel logs: Vercel Dashboard → Deployments → Logs
   - Check Render logs: Render Dashboard → Logs

---

## 🎉 Summary

**Dashboard Setup:** ✅ COMPLETE  
**Database Migration:** ⏳ Pending (Developer Team)  
**RLS Policies:** ⏳ Pending (Developer Team)  
**Deployment:** ⏳ Pending (Developer Team)  

**Next Action:** Developer team should follow the steps above to complete the migration!

---

## 📋 Quick Reference

**Supabase Dashboard:** https://supabase.com/dashboard/project/[YOUR-PROJECT-REF]  
**Vercel Dashboard:** https://vercel.com/dashboard  
**Render Dashboard:** https://dashboard.render.com  
**Frontend URL:** https://tayosaecosystem.vercel.app  
**Backend URL:** https://tayosaecosystem.onrender.com  

**Supabase Credentials:**
- URL: `https://[YOUR-PROJECT-REF].supabase.co`
- Anon Key: `[YOUR_ANON_KEY]`
- Service Key: `[YOUR_SERVICE_ROLE_KEY]`

---

**Great job completing the dashboard setup!** 🎉  
**Your developer team can now complete the migration!** 🚀
