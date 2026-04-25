-- Migration: Auto-approve KYC submissions
-- Description: Removes admin approval requirement - KYC is automatically approved upon submission
-- Date: 2026-04-25

-- Create a trigger function to auto-approve KYC submissions
CREATE OR REPLACE FUNCTION auto_approve_kyc()
RETURNS TRIGGER AS $$
BEGIN
  -- If KYC status is being set to 'pending', automatically change it to 'approved'
  IF NEW.status = 'pending' THEN
    NEW.status := 'approved';
    NEW.reviewed_at := CURRENT_TIMESTAMP;
    NEW.reviewed_by := 'system_auto_approval';
    NEW.review_note := 'Automatically approved upon submission';
  END IF;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger on kyc_profiles table
DROP TRIGGER IF EXISTS trigger_auto_approve_kyc ON kyc_profiles;
CREATE TRIGGER trigger_auto_approve_kyc
  BEFORE INSERT OR UPDATE ON kyc_profiles
  FOR EACH ROW
  EXECUTE FUNCTION auto_approve_kyc();

-- Update any existing pending KYC submissions to approved
UPDATE kyc_profiles
SET 
  status = 'approved',
  reviewed_at = CURRENT_TIMESTAMP,
  reviewed_by = 'system_auto_approval',
  review_note = 'Automatically approved during migration'
WHERE status = 'pending';

-- Add comment
COMMENT ON FUNCTION auto_approve_kyc IS 'Automatically approves KYC submissions without admin review';
