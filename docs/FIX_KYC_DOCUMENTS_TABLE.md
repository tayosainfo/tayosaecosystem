# Fix KYC Documents Table - Step-by-Step Guide

## Problem
You're seeing these errors:
1. `relation "public.kyc_documents" does not exist` - The table hasn't been created yet
2. `column ui.id does not exist` - SQL query was using wrong column name (already fixed)

## Solution

### Step 1: Create the kyc_documents Table

1. **Open Supabase Dashboard**
   - Go to: https://supabase.com/dashboard/project/ablvrbnbsdqshrorhmjf

2. **Open SQL Editor**
   - Click "SQL Editor" in the left sidebar
   - Click "New Query" button

3. **Copy and Paste This SQL**
   
```sql
-- Create kyc_documents table
-- This table stores references to KYC document files in Supabase storage

CREATE TABLE IF NOT EXISTS public.kyc_documents (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users_identity(user_id) ON DELETE CASCADE,
  doc_type TEXT NOT NULL,  -- 'id_document', 'selfie', 'proof_of_address', etc.
  doc_side TEXT NULL,      -- 'front', 'back' (for ID documents)
  storage_key TEXT NOT NULL,  -- Path to file in Supabase storage (e.g., 'user_id/kyc/timestamp-filename.jpg')
  uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups by user_id
CREATE INDEX IF NOT EXISTS idx_kyc_documents_user ON public.kyc_documents(user_id);

-- Create index for faster lookups by doc_type
CREATE INDEX IF NOT EXISTS idx_kyc_documents_type ON public.kyc_documents(doc_type);

-- Add comment to table
COMMENT ON TABLE public.kyc_documents IS 'Stores references to KYC document files uploaded to Supabase storage';

-- Add comments to columns
COMMENT ON COLUMN public.kyc_documents.doc_type IS 'Type of document: id_document, selfie, proof_of_address, etc.';
COMMENT ON COLUMN public.kyc_documents.doc_side IS 'Side of document (for ID cards): front, back';
COMMENT ON COLUMN public.kyc_documents.storage_key IS 'Path to file in Supabase storage bucket';

-- Enable Row Level Security
ALTER TABLE public.kyc_documents ENABLE ROW LEVEL SECURITY;

-- RLS Policy: Users can view their own KYC documents
CREATE POLICY "Users can view own KYC documents"
ON public.kyc_documents
FOR SELECT
TO authenticated
USING (user_id = auth.uid()::text);

-- RLS Policy: Users can insert their own KYC documents
CREATE POLICY "Users can insert own KYC documents"
ON public.kyc_documents
FOR INSERT
TO authenticated
WITH CHECK (user_id = auth.uid()::text);

-- RLS Policy: Users can update their own KYC documents
CREATE POLICY "Users can update own KYC documents"
ON public.kyc_documents
FOR UPDATE
TO authenticated
USING (user_id = auth.uid()::text)
WITH CHECK (user_id = auth.uid()::text);

-- RLS Policy: Users can delete their own KYC documents
CREATE POLICY "Users can delete own KYC documents"
ON public.kyc_documents
FOR DELETE
TO authenticated
USING (user_id = auth.uid()::text);

-- RLS Policy: Service role can do everything (for backend operations)
CREATE POLICY "Service role full access to KYC documents"
ON public.kyc_documents
FOR ALL
TO service_role
USING (true)
WITH CHECK (true);
```

4. **Run the Query**
   - Click the "Run" button (or press Ctrl+Enter)
   - You should see "Success. No rows returned"

### Step 2: Verify the Table Was Created

Run this query to confirm:

```sql
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name = 'kyc_documents';
```

You should see one row with `kyc_documents`.

### Step 3: Verify RLS Policies Were Created

Run this query:

```sql
SELECT schemaname, tablename, policyname, permissive, roles, cmd
FROM pg_policies
WHERE tablename = 'kyc_documents';
```

You should see 5 policies:
1. Users can view own KYC documents (SELECT)
2. Users can insert own KYC documents (INSERT)
3. Users can update own KYC documents (UPDATE)
4. Users can delete own KYC documents (DELETE)
5. Service role full access to KYC documents (ALL)

### Step 4: Test KYC Data Queries

Now you can run the verification queries from `docs/TEST_KYC_DATA.sql` without errors.

Start with this simple query:

```sql
-- Check if any KYC documents exist
SELECT COUNT(*) as document_count 
FROM public.kyc_documents;
```

Then try the full queries from `TEST_KYC_DATA.sql`.

## What This Does

### Table Structure
- **id**: Auto-incrementing primary key
- **user_id**: Links to users_identity table
- **doc_type**: Type of document (id_document, selfie, etc.)
- **doc_side**: For ID cards (front/back)
- **storage_key**: Path to file in Supabase storage
- **uploaded_at**: Timestamp when uploaded

### RLS Policies
- **Users**: Can only see/modify their own documents
- **Service role**: Full access for backend operations (using service_role_key)

### How It Works with Your KYC Flow
1. User uploads ID photos and selfie on frontend
2. Files go to Supabase storage (`collateral_docs` bucket)
3. Backend saves file references to `kyc_documents` table
4. Backend saves KYC form data to `kyc_profiles` table
5. You can query both tables to see complete KYC submission

## Next Steps

After creating the table:
1. Try submitting a KYC form again
2. Check if data appears in both `kyc_profiles` and `kyc_documents` tables
3. Verify files are in storage with matching `storage_key` values

## Troubleshooting

**If you see "permission denied":**
- Make sure you're logged in to Supabase dashboard
- The SQL Editor uses your admin credentials, so it should work

**If policies fail to create:**
- They might already exist from a previous attempt
- You can drop them first with: `DROP POLICY IF EXISTS "policy_name" ON public.kyc_documents;`
- Then re-run the CREATE POLICY statements

**If foreign key constraint fails:**
- Make sure `users_identity` table exists first
- Run the main migration `011_create_supabase_schema.sql` if needed
