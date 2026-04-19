-- Phase-1 profile fields from registration
ALTER TABLE users_identity ADD COLUMN IF NOT EXISTS date_of_birth DATE NULL;
ALTER TABLE users_identity ADD COLUMN IF NOT EXISTS nationality TEXT NULL;
