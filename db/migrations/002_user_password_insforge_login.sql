-- Local password (bcrypt) when running without InsForge; InsForge login email snapshot
ALTER TABLE users_identity ADD COLUMN IF NOT EXISTS password_hash TEXT;
ALTER TABLE users_identity ADD COLUMN IF NOT EXISTS insforge_login_email TEXT;

CREATE INDEX IF NOT EXISTS idx_users_insforge_login_email ON users_identity(insforge_login_email);
