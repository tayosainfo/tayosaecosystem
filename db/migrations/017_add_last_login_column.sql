-- Migration: Add last_login column to users_identity table
-- Description: Track when users last logged in
-- Date: 2026-04-25

ALTER TABLE users_identity 
ADD COLUMN IF NOT EXISTS last_login TIMESTAMP NULL;

-- Create index for efficient queries
CREATE INDEX IF NOT EXISTS idx_users_last_login ON users_identity(last_login);
