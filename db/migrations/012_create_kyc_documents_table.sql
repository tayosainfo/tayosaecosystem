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
