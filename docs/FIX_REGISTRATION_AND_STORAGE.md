# Fix Registration Flow and Storage Upload Issues

## Problem Summary

**Problem 1: Users verified but not in local database**
- 3 users successfully verified emails in Supabase Auth
- BUT they don't exist in `public.users_identity` table
- Root cause: Users clicked email verification link BEFORE the frontend fix was deployed
- The old email link didn't pass `fullName` and `phone` to the verification endpoint

**Problem 2: Storage uploads failing with RLS error**
- Error: "new row violates row-level security policy" (403 Unauthorized)
- Root cause: NO RLS policies exist on `storage.objects` table
- Supabase blocks all operations when RLS is enabled but no policies are configured

---

## Solution 1: Fix Storage RLS Policies with Folder-Based Security (CRITICAL - Do This First)

### Why Folder-Based?
- ✅ **Better security**: Each user has their own isolated folder
- ✅ **Easy debugging**: Just check if file path starts with user's ID
- ✅ **Clear ownership**: Files stored as `{user_id}/kyc/filename.jpg`
- ✅ **Scalable**: Works for millions of users without policy changes

### Step 1: Open Supabase SQL Editor

1. Go to https://supabase.com/dashboard/project/ablvrbnbsdqshrorhmjf
2. Click "SQL Editor" in the left sidebar
3. Click "New query"

### Step 2: Run This SQL Script (Folder-Based Policies)

Copy and paste this ENTIRE script into the SQL editor and click "Run":

```sql
-- Storage RLS Policies for collateral_docs bucket (FOLDER-BASED)
-- Each user can only access files in their own folder: {user_id}/category/filename

-- Enable RLS on storage.objects (if not already enabled)
ALTER TABLE storage.objects ENABLE ROW LEVEL SECURITY;

-- Policy: Users can upload to their own folder only
CREATE POLICY "Users can upload to own folder" ON storage.objects
FOR INSERT TO authenticated
WITH CHECK (
  bucket_id = 'collateral_docs' 
  AND (storage.foldername(name))[1] = auth.uid()::text
);

-- Policy: Users can view files in their own folder only
CREATE POLICY "Users can view own folder files" ON storage.objects
FOR SELECT TO authenticated
USING (
  bucket_id = 'collateral_docs' 
  AND (storage.foldername(name))[1] = auth.uid()::text
);

-- Policy: Users can update files in their own folder only
CREATE POLICY "Users can update own folder files" ON storage.objects
FOR UPDATE TO authenticated
USING (
  bucket_id = 'collateral_docs' 
  AND (storage.foldername(name))[1] = auth.uid()::text
)
WITH CHECK (
  bucket_id = 'collateral_docs' 
  AND (storage.foldername(name))[1] = auth.uid()::text
);

-- Policy: Users can delete files in their own folder only
CREATE POLICY "Users can delete own folder files" ON storage.objects
FOR DELETE TO authenticated
USING (
  bucket_id = 'collateral_docs' 
  AND (storage.foldername(name))[1] = auth.uid()::text
);

-- Policy: Service role has full access to storage
CREATE POLICY "Service role full access storage" ON storage.objects
USING (auth.role() = 'service_role');
```

**How it works:**
- `storage.foldername(name)` splits the file path into folders
- `[1]` gets the first folder (user ID)
- `auth.uid()::text` is the authenticated user's Supabase ID
- Files are stored as: `{user_id}/kyc/20260423144838-id_front.jpg`

### Step 3: Verify Policies Were Created

Run this query to check:

```sql
SELECT policyname, cmd, roles 
FROM pg_policies 
WHERE schemaname = 'storage' AND tablename = 'objects';
```

You should see 5 policies listed.

### Step 4: Update Backend to Use Folder-Based Paths

The backend needs a small change to include the user ID in the storage path.

**File to update:** `services/object-storage-service/main.go`

**Find this code** (around line 280):

```go
// Sanitise filename and build storage path.
safeFilename := strings.ReplaceAll(strings.TrimSpace(header.Filename), " ", "_")
if safeFilename == "" {
	safeFilename = "upload.bin"
}
objectPath := fmt.Sprintf("%s/%s-%s",
	category,
	time.Now().UTC().Format("20060102150405"),
	safeFilename,
)
```

**Replace with:**

```go
// Sanitise filename and build storage path.
safeFilename := strings.ReplaceAll(strings.TrimSpace(header.Filename), " ", "_")
if safeFilename == "" {
	safeFilename = "upload.bin"
}

// Get user ID from validated token (set by requireAuth middleware)
userID := r.Header.Get("X-User-Id")
if userID == "" {
	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "user ID not found in request"})
	return
}

// Build folder-based path: {user_id}/{category}/{timestamp}-{filename}
objectPath := fmt.Sprintf("%s/%s/%s-%s",
	userID,    // User's folder (Supabase user ID)
	category,  // kyc, documents, etc.
	time.Now().UTC().Format("20060102150405"),
	safeFilename,
)
```

**Result:** Files will be stored as `{user_id}/kyc/20260423144838-id_front.jpg` instead of `kyc/20260423144838-id_front.jpg`

**After making this change:**
1. Commit and push to GitHub
2. Render will auto-deploy the updated service
3. Test file upload - it should now work with folder-based security

---

## Solution 2: Fix Existing Verified Users (Manual Database Insertion)

The 3 verified users need to be manually added to the local database. Here's the data from Supabase Auth:

### User 1: TIGER MUSIC UGA
- **Supabase ID**: `23a0e476-d7f7-4142-acc7-e56eb866de8a`
- **Email**: `tger5018@gmail.com`
- **Phone**: `+256775846264`
- **Name**: `TIGER  MUSIC UGA`
- **Verified**: ✅ Yes (2026-04-23 14:43:11)

### User 2: Alex Twahirwa
- **Supabase ID**: `fdcaf337-332e-4341-a218-748c597c00ac`
- **Email**: `twahirwaalex@gmail.com`
- **Phone**: `+256775846265`
- **Name**: `Alex  Twahirwa`
- **Verified**: ✅ Yes (2026-04-23 14:14:20)

### User 3: BAYLES INFO
- **Supabase ID**: `eb3f84db-707f-4c55-9106-22da0776af73`
- **Email**: `baylesinfo@gmail.com`
- **Phone**: `+256775846267`
- **Name**: `BAYLES INFO`
- **Verified**: ✅ Yes (2026-04-23 14:03:57)

### SQL Script to Add These Users

Run this in Supabase SQL Editor:

```sql
-- Insert User 1: TIGER MUSIC UGA
INSERT INTO public.users_identity (
  user_id, full_name, phone_e164, auth_email, contact_email, 
  supabase_login_email, supabase_user_id, nationality, created_at, contact_email_checked
) VALUES (
  '23a0e476-d7f7-4142-acc7-e56eb866de8a',
  'TIGER  MUSIC UGA',
  '+256775846264',
  'tger5018@gmail.com',
  'tger5018@gmail.com',
  'tger5018@gmail.com',
  '23a0e476-d7f7-4142-acc7-e56eb866de8a',
  'UG',
  '2026-04-23 14:24:11.941484+00',
  true
);

-- Insert User 2: Alex Twahirwa
INSERT INTO public.users_identity (
  user_id, full_name, phone_e164, auth_email, contact_email, 
  supabase_login_email, supabase_user_id, nationality, created_at, contact_email_checked
) VALUES (
  'fdcaf337-332e-4341-a218-748c597c00ac',
  'Alex  Twahirwa',
  '+256775846265',
  'twahirwaalex@gmail.com',
  'twahirwaalex@gmail.com',
  'twahirwaalex@gmail.com',
  'fdcaf337-332e-4341-a218-748c597c00ac',
  'UG',
  '2026-04-23 14:13:13.772964+00',
  true
);

-- Insert User 3: BAYLES INFO
INSERT INTO public.users_identity (
  user_id, full_name, phone_e164, auth_email, contact_email, 
  supabase_login_email, supabase_user_id, nationality, created_at, contact_email_checked
) VALUES (
  'eb3f84db-707f-4c55-9106-22da0776af73',
  'BAYLES INFO',
  '+256775846267',
  'baylesinfo@gmail.com',
  'baylesinfo@gmail.com',
  'baylesinfo@gmail.com',
  'eb3f84db-707f-4c55-9106-22da0776af73',
  'UG',
  '2026-04-23 14:03:11.752953+00',
  true
);

-- Create onboarding profiles for all 3 users
INSERT INTO public.onboarding_profiles (user_id, phase, trust_score_seed, last_updated_at)
VALUES 
  ('23a0e476-d7f7-4142-acc7-e56eb866de8a', 1, 10, NOW()),
  ('fdcaf337-332e-4341-a218-748c597c00ac', 1, 10, NOW()),
  ('eb3f84db-707f-4c55-9106-22da0776af73', 1, 10, NOW());
```

### Verify Users Were Added

Run this query:

```sql
SELECT user_id, full_name, phone_e164, contact_email, contact_email_checked 
FROM public.users_identity 
ORDER BY created_at DESC 
LIMIT 5;
```

You should see all 3 users listed.

---

## Storage Folder Structure

After implementing folder-based policies, files will be organized like this:

```
collateral_docs/
├── 23a0e476-d7f7-4142-acc7-e56eb866de8a/  (User 1's folder)
│   ├── kyc/
│   │   ├── 20260423144838-id_front.jpg
│   │   ├── 20260423144840-id_back.jpg
│   │   └── 20260423144842-selfie.jpg
│   └── documents/
│       └── 20260423150000-contract.pdf
├── fdcaf337-332e-4341-a218-748c597c00ac/  (User 2's folder)
│   └── kyc/
│       ├── 20260423145000-id_front.jpg
│       └── ...
└── eb3f84db-707f-4c55-9106-22da0776af73/  (User 3's folder)
    └── kyc/
        └── ...
```

**Benefits:**
- Each user can ONLY access their own folder
- Easy to debug: Check if file path starts with user's Supabase ID
- Scalable: Works for millions of users
- Clear ownership: No confusion about who owns which file

---

## Solution 3: Test New User Registration (End-to-End)

After applying the storage RLS policies, test the complete flow:

### Test Steps:

1. **Register a new user** at https://tayosaecosystem.vercel.app/register
   - Use a NEW email address (not one of the 3 above)
   - Use a valid Uganda phone number
   - Accept terms and privacy policy

2. **Check email** for 6-digit verification code

3. **Verify email** at https://tayosaecosystem.vercel.app/verify
   - Enter the 6-digit code
   - Frontend will auto-send `fullName` and `phone` from sessionStorage

4. **Check database** - User should now exist in `public.users_identity`

5. **Complete onboarding** - Upload KYC documents
   - Storage uploads should now work (no more 403 errors)

---

## Why This Happened

### Timeline of Events:

1. **Initial deployment** - Registration worked but had duplicate OTP call issue
2. **Fix deployed** (commit f9fed70) - Removed duplicate OTP call, added auto-token extraction
3. **3 users registered** - They clicked OLD email links that didn't pass `fullName`/`phone`
4. **Verification succeeded** - Supabase Auth verified them
5. **Local profile creation failed** - Backend didn't receive `fullName`/`phone` parameters
6. **Storage uploads failed** - No RLS policies on `storage.objects` table

### Current State:

✅ **Frontend code is correct** - Register and Verify pages properly handle sessionStorage
✅ **Backend code is correct** - verifyEmailHandler creates local profiles when it receives data
❌ **3 users stuck** - Verified in Auth but not in local database (need manual fix)
❌ **Storage blocked** - No RLS policies (need SQL script + backend update)

### Why Folder-Based Storage?

The updated solution uses **folder-based RLS policies** instead of bucket-level policies:

**Old approach (bucket-level):**
- ❌ Anyone authenticated could access any file in the bucket
- ❌ No user isolation
- ❌ Hard to debug access issues

**New approach (folder-based):**
- ✅ Each user has their own folder: `{user_id}/kyc/filename.jpg`
- ✅ Users can ONLY access files in their own folder
- ✅ Easy to debug: Just check if path starts with user's ID
- ✅ Better security and scalability

---

## Next Steps

1. ✅ **Apply storage RLS policies** (Solution 1 - SQL script) - CRITICAL for KYC uploads
2. ✅ **Update backend code** (Solution 1 - Step 4) - Add folder-based paths
3. ✅ **Deploy backend changes** - Commit and push to GitHub (Render auto-deploys)
4. ✅ **Add 3 verified users to database** (Solution 2) - So they can log in
5. ✅ **Test new registration** (Solution 3) - Verify end-to-end flow works
6. ✅ **Monitor logs** - Check for any new errors

---

## Expected Results After Fixes

### For Existing 3 Users:
- Can log in with email + password
- Will see onboarding screen (phase 1)
- Can complete KYC and upload documents

### For New Users:
- Registration → Email verification → Local profile creation (all automatic)
- No manual database intervention needed
- Storage uploads work correctly

---

## Verification Queries

### Check Auth Users:
```sql
SELECT id, email, phone, email_confirmed_at, raw_user_meta_data->>'phone' as metadata_phone
FROM auth.users 
ORDER BY created_at DESC 
LIMIT 5;
```

### Check Local Users:
```sql
SELECT user_id, full_name, phone_e164, contact_email, contact_email_checked
FROM public.users_identity 
ORDER BY created_at DESC 
LIMIT 5;
```

### Check Storage Policies:
```sql
SELECT policyname, cmd, roles 
FROM pg_policies 
WHERE schemaname = 'storage' AND tablename = 'objects';
```

### Check Storage Bucket:
```sql
SELECT id, name, public, created_at 
FROM storage.buckets 
WHERE name = 'collateral_docs';
```

---

## Contact

If you encounter any issues after applying these fixes, check:
1. Supabase dashboard logs (Authentication & Storage sections)
2. Render backend logs (User Service & Object Storage Service)
3. Browser console (Frontend errors)
