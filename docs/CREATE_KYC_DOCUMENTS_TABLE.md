# Create KYC Documents Table

## Problem

The `kyc_documents` table is missing from your Supabase database. This table is required to store references to uploaded KYC document files.

## Solution

Run the SQL migration to create the table.

---

## Step 1: Open Supabase SQL Editor

1. Go to https://supabase.com/dashboard/project/ablvrbnbsdqshrorhmjf
2. Click **"SQL Editor"** in the left sidebar
3. Click **"New query"** button

---

## Step 2: Copy and Run This SQL

Copy this entire SQL script and paste it into the SQL Editor:

```sql
-- Create kyc_documents table
-- This table stores references to KYC document files in Supabase storage

CREATE TABLE IF NOT EXISTS public.kyc_documents (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users_identity(user_id) ON DELETE CASCADE,
  doc_type TEXT NOT NULL,  -- 'id_document', 'selfie', 'proof_of_address', etc.
  doc_side TEXT NULL,      -- 'front', 'back' (for ID documents)
  storage_key TEXT NOT NULL,  -- Path to file in Supabase storage
  uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_kyc_documents_user ON public.kyc_documents(user_id);
CREATE INDEX IF NOT EXISTS idx_kyc_documents_type ON public.kyc_documents(doc_type);

-- Verify table was created
SELECT 
    table_name, 
    column_name, 
    data_type 
FROM information_schema.columns 
WHERE table_name = 'kyc_documents' 
ORDER BY ordinal_position;
```

---

## Step 3: Click "Run" Button

Click the **"Run"** button (or press Ctrl+Enter / Cmd+Enter)

You should see:
- Success message
- Table structure showing 6 columns (id, user_id, doc_type, doc_side, storage_key, uploaded_at)

---

## Step 4: Verify Table Exists

Run this query to confirm:

```sql
SELECT COUNT(*) as table_exists 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name = 'kyc_documents';
```

You should see: `table_exists: 1`

---

## What This Table Does

The `kyc_documents` table stores **references** to files uploaded to Supabase storage:

| Column | Type | Description |
|--------|------|-------------|
| `id` | BIGSERIAL | Auto-incrementing primary key |
| `user_id` | TEXT | Foreign key to users_identity table |
| `doc_type` | TEXT | Type of document (e.g., "id_document", "selfie") |
| `doc_side` | TEXT | Side of document (e.g., "front", "back") |
| `storage_key` | TEXT | Path to file in Supabase storage |
| `uploaded_at` | TIMESTAMP | When the file was uploaded |

### Example Data

```
id | user_id | doc_type    | doc_side | storage_key                                    | uploaded_at
---|---------|-------------|----------|------------------------------------------------|-------------------
1  | abc123  | id_document | front    | abc123/kyc/20260424120000-national_id_front.jpg | 2026-04-24 12:00:00
2  | abc123  | id_document | back     | abc123/kyc/20260424120001-national_id_back.jpg  | 2026-04-24 12:00:01
3  | abc123  | selfie      | NULL     | abc123/kyc/20260424120002-selfie.jpg            | 2026-04-24 12:00:02
```

---

## After Creating the Table

Once the table is created, you can:

1. **Submit KYC forms** - The backend will save document references to this table
2. **Run test queries** - Use `docs/TEST_KYC_DATA.sql` to verify data
3. **View documents** - Query this table to see all uploaded documents

---

## Troubleshooting

### Error: "relation already exists"
- This is fine! It means the table was already created
- You can proceed to the next step

### Error: "permission denied"
- Make sure you're logged in to Supabase dashboard
- Make sure you're in the correct project (ablvrbnbsdqshrorhmjf)

### Error: "foreign key constraint"
- Make sure the `users_identity` table exists first
- Run the main migration `db/migrations/011_create_supabase_schema.sql` first

---

## Next Steps

After creating the table:

1. **Test KYC submission** - Fill out the KYC form and upload documents
2. **Verify data** - Run queries from `docs/TEST_KYC_DATA.sql`
3. **Check storage** - Go to Storage → collateral_docs to see uploaded files

The KYC flow will now work completely end-to-end! ✅
