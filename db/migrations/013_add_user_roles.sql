-- Migration: Add user roles and admin audit logging
-- Description: Adds role-based access control to users_identity table with audit trail
-- Date: 2026-04-24

-- Create enum type for user roles (extensible for future roles)
CREATE TYPE user_role AS ENUM ('user', 'admin');

-- Add role columns to users_identity table
ALTER TABLE users_identity 
ADD COLUMN role user_role NOT NULL DEFAULT 'user',
ADD COLUMN role_assigned_at TIMESTAMP NULL,
ADD COLUMN role_assigned_by TEXT NULL;

-- Create index for efficient role-based queries
CREATE INDEX idx_users_identity_role ON users_identity(role);

-- Create admin role audit table for tracking role changes
CREATE TABLE admin_role_audit (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users_identity(user_id),
  action TEXT NOT NULL, -- 'granted', 'revoked'
  previous_role user_role NULL,
  new_role user_role NOT NULL,
  assigned_by TEXT NULL, -- user_id of admin who made change
  reason TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for audit table
CREATE INDEX idx_admin_audit_user ON admin_role_audit(user_id);
CREATE INDEX idx_admin_audit_created ON admin_role_audit(created_at DESC);

-- Function to automatically log role changes
CREATE OR REPLACE FUNCTION log_role_change()
RETURNS TRIGGER AS $$
BEGIN
  IF (TG_OP = 'UPDATE' AND OLD.role IS DISTINCT FROM NEW.role) THEN
    INSERT INTO admin_role_audit (user_id, action, previous_role, new_role, assigned_by)
    VALUES (
      NEW.user_id,
      CASE WHEN NEW.role = 'admin' THEN 'granted' ELSE 'revoked' END,
      OLD.role,
      NEW.role,
      NEW.role_assigned_by
    );
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for automatic audit logging on role changes
CREATE TRIGGER trigger_log_role_change
AFTER UPDATE ON users_identity
FOR EACH ROW
EXECUTE FUNCTION log_role_change();

-- Add comment for documentation
COMMENT ON TABLE admin_role_audit IS 'Audit trail for admin role assignments and revocations';
COMMENT ON COLUMN users_identity.role IS 'User role for access control (user, admin)';
COMMENT ON COLUMN users_identity.role_assigned_at IS 'Timestamp when role was last assigned';
COMMENT ON COLUMN users_identity.role_assigned_by IS 'User ID of admin who assigned the role';
