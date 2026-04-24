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
