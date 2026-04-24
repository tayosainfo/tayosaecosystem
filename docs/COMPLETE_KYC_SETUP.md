# Complete KYC Setup Guide

## Current Status

✅ **Fixed Issues:**
1. Build error in object-storage-service (unused variable) - FIXED
2. Token validation for storage uploads - WORKING
3. Folder-based storage paths ({user_id}/kyc/filename) - WORKING
4. Storage RLS policies - CONFIGURED
5. Test queries column names - CORRECT

## Step 1: Create kyc_documents Table in Supabase

**IMPORTANT:** The kyc_documents table does NOT exist yet in your Supabase database. You need to create it.

### How to Create the Table:

1. **Open Supabase Dashboard**
   - Go to: https://supabase.com/dashboard/project/ablvrbnbsdqshrorhmjf

2. **Open SQL Editor**
   - Click "SQL Editor" in the left sidebar
   - Click "New query"

3. **Copy and Paste This SQL**
   - Open the file: `db/migrations/012_create_kyc_documents_table.sql`
   - Copy ALL the SQL code from that file
   - Paste it into the SQL Editor

4. **Run the SQL**
   - Click "Run" button (or press Ctrl+Enter)
   - You should see: "Success. No rows returned"

5. **Verify Table Was Created**
   - Click "Table Editor" in the left sidebar
   - Look for "kyc_documents" in the list of tables
   - You should see it with columns: id, user_id, doc_type, doc_side, storage_key, uploaded_at

## Step 2: Verify Deployment is Successful

1. **Check Render Dashboard**
   - Go to: https://dashboard.render.com
   - Check that all services are "Live" (green)
   - Object Storage Service should show no errors

2. **Check Service Logs**
   - Click on "Object Storage Service"
   - Click "Logs" tab
   - Look for: "Object Storage Service listening on :8015"
   - Should NOT see any build errors

## Step 3: Test KYC Flow End-to-End

### 3.1 Register a New User

1. **Go to Registration Page**
   - Open: https://tayosaecosystem.vercel.app/register

2. **Fill in Registration Form**
   - Use a NEW email address (not used before)
   - Use a valid Uganda phone number (e.g., 0700123456)
   - Fill in all required fields
   - Click "Register"

3. **Verify Email**
   - Check your email inbox
   - You should receive a 6-digit verification code
   - Go to: https://tayosaecosystem.vercel.app/verify
   - Enter the 6-digit code
   - Click "Verify"

4. **Confirm Login**
   - After verification, you should be logged in
   - You should see the dashboard

### 3.2 Submit KYC Form

1. **Navigate to KYC Page**
   - From the dashboard, click "Complete KYC" or go to the KYC section

2. **Fill in KYC Form**
   - **Personal Information:**
     - Date of Birth
     - Gender
     - Nationality
     - Occupation Status
   
   - **ID Information:**
     - ID Type (National ID, Passport, etc.)
     - ID Number
   
   - **Upload Documents:**
     - ID Front Photo (click "Upload" and select image)
     - ID Back Photo (click "Upload" and select image)
     - Selfie Photo (click "Upload" and select image)
   
   - **Next of Kin:**
     - Full Name
     - Relationship
     - Phone Number
     - Email (optional)
   
   - **Financial Information:**
     - Source of Funds
     - PEP Status (Yes/No)
     - SACCO Membership Disclosures

3. **Submit Form**
   - Click "Submit KYC"
   - You should see a success message

### 3.3 Verify Data Was Saved

1. **Open Supabase Dashboard**
   - Go to: https://supabase.com/dashboard/project/ablvrbnbsdqshrorhmjf

2. **Open SQL Editor**
   - Click "SQL Editor" in the left sidebar
   - Click "New query"

3. **Run Verification Queries**

   **Query 1: Check KYC Profile**
   ```sql
   SELECT * FROM public.kyc_profiles 
   ORDER BY submitted_at DESC 
   LIMIT 1;
   ```
   - You should see your KYC data with status = 'pending'

   **Query 2: Check KYC Documents**
   ```sql
   SELECT * FROM public.kyc_documents 
   ORDER BY uploaded_at DESC 
   LIMIT 3;
   ```
   - You should see 3 rows (ID front, ID back, selfie)
   - Each row should have a storage_key like: `{user_id}/kyc/20260424123456-filename.jpg`

   **Query 3: Check Storage Files**
   ```sql
   SELECT name, created_at, metadata 
   FROM storage.objects 
   WHERE bucket_id = 'collateral_docs' 
   ORDER BY created_at DESC 
   LIMIT 3;
   ```
   - You should see 3 files with paths matching the storage_keys from Query 2

   **Query 4: Complete View (Profile + Documents + User)**
   ```sql
   SELECT 
       ui.user_id,
       ui.full_name,
       ui.contact_email,
       kp.status AS kyc_status,
       kp.id_type,
       kp.id_number,
       kp.submitted_at,
       COUNT(kd.id) AS document_count
   FROM public.users_identity ui
   LEFT JOIN public.kyc_profiles kp ON ui.user_id = kp.user_id
   LEFT JOIN public.kyc_documents kd ON ui.user_id = kd.user_id
   WHERE kp.user_id IS NOT NULL
   GROUP BY ui.user_id, ui.full_name, ui.contact_email, 
            kp.status, kp.id_type, kp.id_number, kp.submitted_at
   ORDER BY kp.submitted_at DESC
   LIMIT 1;
   ```
   - You should see 1 row with document_count = 3

## Step 4: Troubleshooting

### If KYC Profile is NOT Saved:

1. **Check Backend Logs**
   - Go to Render dashboard
   - Click "User Service"
   - Click "Logs" tab
   - Look for errors when you submitted the KYC form

2. **Check Browser Console**
   - Open browser DevTools (F12)
   - Click "Console" tab
   - Look for errors when you submitted the form

### If Documents are NOT Saved:

1. **Check Storage Service Logs**
   - Go to Render dashboard
   - Click "Object Storage Service"
   - Click "Logs" tab
   - Look for upload errors

2. **Check Storage RLS Policies**
   - Go to Supabase Dashboard
   - Click "Authentication" → "Policies"
   - Click "storage" schema → "objects" table
   - Verify you have 6 policies (INSERT, SELECT, UPDATE, DELETE for users + 2 for service_role)

### If Queries Return "No Rows":

This means the data was NOT saved. Possible causes:

1. **kyc_documents table doesn't exist** → Run Step 1 again
2. **Backend error during submission** → Check backend logs
3. **RLS policies blocking the insert** → Check RLS policies
4. **Token validation failing** → Check that you're logged in

## Expected Results

After completing all steps, you should have:

✅ kyc_documents table exists in Supabase
✅ Render deployment is successful (no build errors)
✅ User can register and verify email
✅ User can submit KYC form with file uploads
✅ KYC profile is saved in kyc_profiles table
✅ 3 document references are saved in kyc_documents table
✅ 3 files are uploaded to storage.objects
✅ All verification queries return data

## Next Steps

Once KYC flow is working:

1. Test SACCO membership enrollment
2. Test Kibiina group creation/joining
3. Test shares purchase
4. Test mobile money integration

## Need Help?

If you encounter any errors:

1. Copy the EXACT error message
2. Note which step you were on
3. Check the relevant logs (backend or browser console)
4. Share the error message for debugging
