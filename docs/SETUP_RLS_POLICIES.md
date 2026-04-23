# 🔒 Setting Up Security Policies (RLS) - Simple Guide

**What is RLS?** It's like a security guard for your database - it makes sure users can only see and change their own data, not other people's data.

**You have automatic RLS enabled** - This means we MUST set up these security rules, or your app won't work properly!

⏱️ **Time needed:** 10-15 minutes

---

## 🚨 IMPORTANT: Do This BEFORE Testing Your App!

Without these security policies, users won't be able to:
- See their own profile
- Update their information
- Access their data

Let's fix that now! 👇

---

## 📝 Step-by-Step: Adding Security Policies

### Step 1: Open the SQL Editor (1 minute)

1. Go to your Supabase dashboard: https://supabase.com/dashboard/project/[YOUR-PROJECT-REF]
2. **Click on "SQL Editor"** in the left sidebar (it has a `</>` icon)
3. You'll see a big text box where you can type SQL commands

---

### Step 2: Add Security Policy for Users Table (5 minutes)

**What you're doing:** Letting users see and update their own profile information.

1. **Copy this entire code block:**

```sql
-- Enable RLS on users_identity table (if not already enabled)
ALTER TABLE users_identity ENABLE ROW LEVEL SECURITY;

-- Policy 1: Users can view their own profile
CREATE POLICY "Users can view own profile"
ON users_identity FOR SELECT
USING (auth.uid()::text = supabase_user_id);

-- Policy 2: Users can update their own profile
CREATE POLICY "Users can update own profile"
ON users_identity FOR UPDATE
USING (auth.uid()::text = supabase_user_id);

-- Policy 3: Service role (backend) can do everything
CREATE POLICY "Service role has full access"
ON users_identity
USING (auth.role() = 'service_role');
```

2. **Paste it into the SQL Editor** (the big text box)
3. **Click the "Run" button** (usually green, at the bottom right)
4. You should see a message saying **"Success"** or **"Completed"**

✅ **Done!** Users can now access their own profiles.

---

### Step 3: Add Security Policy for Onboarding Profiles Table (5 minutes)

**What you're doing:** Letting users see and update their own onboarding information.

1. **Copy this entire code block:**

```sql
-- Enable RLS on onboarding_profiles table
ALTER TABLE onboarding_profiles ENABLE ROW LEVEL SECURITY;

-- Policy 1: Users can view their own onboarding profile
CREATE POLICY "Users can view own onboarding profile"
ON onboarding_profiles FOR SELECT
USING (
  auth.uid()::text = (
    SELECT supabase_user_id 
    FROM users_identity 
    WHERE user_id = onboarding_profiles.user_id
  )
);

-- Policy 2: Users can update their own onboarding profile
CREATE POLICY "Users can update own onboarding profile"
ON onboarding_profiles FOR UPDATE
USING (
  auth.uid()::text = (
    SELECT supabase_user_id 
    FROM users_identity 
    WHERE user_id = onboarding_profiles.user_id
  )
);

-- Policy 3: Users can insert their own onboarding profile
CREATE POLICY "Users can create own onboarding profile"
ON onboarding_profiles FOR INSERT
WITH CHECK (
  auth.uid()::text = (
    SELECT supabase_user_id 
    FROM users_identity 
    WHERE user_id = onboarding_profiles.user_id
  )
);

-- Policy 4: Service role (backend) can do everything
CREATE POLICY "Service role has full access to onboarding profiles"
ON onboarding_profiles
USING (auth.role() = 'service_role');
```

2. **Paste it into the SQL Editor**
3. **Click the "Run" button**
4. You should see **"Success"**

✅ **Done!** Users can now manage their onboarding information.

---

### Step 4: Add Security Policy for User Consents Table (5 minutes)

**What you're doing:** Letting users see and manage their own consent records.

1. **Copy this entire code block:**

```sql
-- Enable RLS on user_consents table
ALTER TABLE user_consents ENABLE ROW LEVEL SECURITY;

-- Policy 1: Users can view their own consents
CREATE POLICY "Users can view own consents"
ON user_consents FOR SELECT
USING (
  auth.uid()::text = (
    SELECT supabase_user_id 
    FROM users_identity 
    WHERE user_id = user_consents.user_id
  )
);

-- Policy 2: Users can update their own consents
CREATE POLICY "Users can update own consents"
ON user_consents FOR UPDATE
USING (
  auth.uid()::text = (
    SELECT supabase_user_id 
    FROM users_identity 
    WHERE user_id = user_consents.user_id
  )
);

-- Policy 3: Users can insert their own consents
CREATE POLICY "Users can create own consents"
ON user_consents FOR INSERT
WITH CHECK (
  auth.uid()::text = (
    SELECT supabase_user_id 
    FROM users_identity 
    WHERE user_id = user_consents.user_id
  )
);

-- Policy 4: Service role (backend) can do everything
CREATE POLICY "Service role has full access to user consents"
ON user_consents
USING (auth.role() = 'service_role');
```

2. **Paste it into the SQL Editor**
3. **Click the "Run" button**
4. You should see **"Success"**

✅ **Done!** Users can now manage their consent records.

---

### Step 5: Add Security Policy for Storage (File Uploads) (3 minutes)

**What you're doing:** Letting users upload and view their own documents (like ID photos).

1. **Copy this entire code block:**

```sql
-- Policy 1: Authenticated users can upload files
CREATE POLICY "Authenticated users can upload"
ON storage.objects FOR INSERT
TO authenticated
WITH CHECK (bucket_id = 'collateral_docs');

-- Policy 2: Users can view their own files
CREATE POLICY "Users can view own files"
ON storage.objects FOR SELECT
TO authenticated
USING (
  bucket_id = 'collateral_docs' 
  AND auth.uid()::text = (storage.foldername(name))[1]
);

-- Policy 3: Users can delete their own files
CREATE POLICY "Users can delete own files"
ON storage.objects FOR DELETE
TO authenticated
USING (
  bucket_id = 'collateral_docs' 
  AND auth.uid()::text = (storage.foldername(name))[1]
);

-- Policy 4: Service role (backend) can do everything
CREATE POLICY "Service role has full access to storage"
ON storage.objects
TO service_role
USING (bucket_id = 'collateral_docs');
```

2. **Paste it into the SQL Editor**
3. **Click the "Run" button**
4. You should see **"Success"**

✅ **Done!** Users can now upload and view their documents.

---

## ✅ Verification: Make Sure It Worked

After running all the SQL commands above:

1. **Click on "Database"** in the left sidebar
2. **Click on "Policies"** at the top
3. You should see a list of all the policies you just created

**You should see policies for:**
- ✅ users_identity table (3 policies)
- ✅ onboarding_profiles table (4 policies)
- ✅ user_consents table (4 policies)
- ✅ storage.objects (4 policies)

**Total: About 15 policies**

---

## 🎉 Congratulations!

Your database is now secure! Here's what you accomplished:

✅ Users can only see their own data  
✅ Users can only update their own information  
✅ Users can only upload and view their own files  
✅ Your backend services can still access everything (using the service role key)  
✅ Other users cannot see or change each other's data  

---

## 🤔 What If I Get an Error?

### Error: "Policy already exists"
**What it means:** You already ran this command before.  
**What to do:** This is fine! Just skip to the next step.

### Error: "Table does not exist"
**What it means:** That table hasn't been created yet in your database.  
**What to do:** Skip that table for now. Your developer team will add it later.

### Error: "Permission denied"
**What it means:** You might not have admin access.  
**What to do:** Make sure you're logged in as the project owner.

---

## 📞 Need Help?

**If you get stuck:**
1. Take a screenshot of the error message
2. Email: support@tayosa.com
3. Or ask your developer team

**Common Questions:**

**Q: What if I make a mistake?**  
A: You can delete a policy and recreate it. Go to Database > Policies, find the policy, and click the trash icon.

**Q: Can I test if the policies work?**  
A: Yes! After your app is deployed, try logging in as a user. You should only see your own data.

**Q: Do I need to do this for every table?**  
A: Only for tables that store user-specific data. Your developer team will handle any other tables.

---

## 📋 Quick Checklist

Before you're done, make sure:

- [ ] I ran all 4 SQL code blocks
- [ ] Each one showed "Success"
- [ ] I can see the policies in Database > Policies
- [ ] I have about 15 policies total

**All checked?** Your security is set up! 🔒

---

## 🚀 What's Next?

Now that RLS policies are set up:

1. ✅ Complete the Supabase dashboard setup (if you haven't already)
2. ✅ Your developer team can now deploy the app
3. ✅ Users will be able to access their data securely

---

**⏱️ Total Time:** 10-15 minutes  
**Difficulty:** Easy - Just copy and paste!  
**Status:** Your database is now secure! 🔒✅
