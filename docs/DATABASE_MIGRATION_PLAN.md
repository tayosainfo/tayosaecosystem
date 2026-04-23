# Database Migration Execution Plan

This document provides a detailed plan for executing the database migration from InsForge to Supabase naming conventions.

## Migration Overview

**Migration File:** `db/migrations/011_rename_insforge_to_supabase.sql`

**Purpose:** Rename database columns from InsForge naming to Supabase naming conventions

**Affected Tables:**
- `users` table

**Changes:**
- Rename `insforge_user_id` → `supabase_user_id`
- Rename `insforge_login_email` → `supabase_login_email`
- Drop index `idx_users_insforge_login_email`
- Create index `idx_users_supabase_login_email`

**Estimated Duration:** 5-10 minutes (depends on table size)

**Downtime Required:** Yes (recommended 15-30 minutes maintenance window)

## Pre-Migration Checklist

### 1. Backup Preparation

- [ ] **Full Database Backup**
  ```bash
  # Create backup using pg_dump
  pg_dump -h db.[YOUR-PROJECT-REF].supabase.co \
          -U postgres \
          -d postgres \
          -F c \
          -f backup_pre_migration_$(date +%Y%m%d_%H%M%S).dump
  ```
  - Verify backup file is created
  - Test backup restoration on staging environment
  - Store backup in secure location

- [ ] **Table-Specific Backup**
  ```sql
  -- Backup users table
  CREATE TABLE users_backup_20240101 AS SELECT * FROM users;
  
  -- Verify backup
  SELECT COUNT(*) FROM users;
  SELECT COUNT(*) FROM users_backup_20240101;
  ```

### 2. Data Integrity Verification

- [ ] **Check for NULL values**
  ```sql
  -- Verify no NULL insforge_user_id values
  SELECT COUNT(*) FROM users WHERE insforge_user_id IS NULL;
  
  -- Verify no NULL insforge_login_email values
  SELECT COUNT(*) FROM users WHERE insforge_login_email IS NULL;
  ```
  - Expected result: 0 rows for both queries

- [ ] **Check for duplicate emails**
  ```sql
  -- Find duplicate emails
  SELECT insforge_login_email, COUNT(*) 
  FROM users 
  GROUP BY insforge_login_email 
  HAVING COUNT(*) > 1;
  ```
  - Expected result: 0 rows

- [ ] **Record current row counts**
  ```sql
  SELECT 
    COUNT(*) as total_users,
    COUNT(DISTINCT insforge_user_id) as unique_user_ids,
    COUNT(DISTINCT insforge_login_email) as unique_emails
  FROM users;
  ```
  - Document these numbers for post-migration verification

### 3. Application Preparation

- [ ] **Stop all backend services**
  ```bash
  # Stop services to prevent writes during migration
  systemctl stop api-gateway-service
  systemctl stop user-service
  systemctl stop affiliate-service
  systemctl stop audit-log-service
  systemctl stop fee-service
  systemctl stop kibiina-service
  systemctl stop loan-credit-service
  systemctl stop notification-service
  systemctl stop object-storage-service
  ```

- [ ] **Display maintenance page**
  - Enable maintenance mode on frontend
  - Display message: "System maintenance in progress. We'll be back shortly."

- [ ] **Notify users**
  - Send notification about scheduled maintenance
  - Provide estimated completion time

### 4. Environment Verification

- [ ] **Verify database connection**
  ```bash
  psql -h db.[YOUR-PROJECT-REF].supabase.co \
       -U postgres \
       -d postgres \
       -c "SELECT version();"
  ```

- [ ] **Check active connections**
  ```sql
  SELECT COUNT(*) FROM pg_stat_activity 
  WHERE datname = 'postgres' AND state = 'active';
  ```
  - Should be minimal (only admin connections)

- [ ] **Verify migration file exists**
  ```bash
  ls -lh db/migrations/011_rename_insforge_to_supabase.sql
  cat db/migrations/011_rename_insforge_to_supabase.sql
  ```

## Migration Execution

### Step 1: Connect to Database

```bash
psql -h db.[YOUR-PROJECT-REF].supabase.co \
     -U postgres \
     -d postgres
```

### Step 2: Begin Transaction

```sql
-- Start transaction for rollback capability
BEGIN;

-- Set statement timeout (10 minutes)
SET statement_timeout = '600s';
```

### Step 3: Execute Migration

```sql
-- Execute migration file
\i db/migrations/011_rename_insforge_to_supabase.sql
```

**Or execute commands directly:**

```sql
-- Rename insforge_user_id to supabase_user_id
ALTER TABLE users 
RENAME COLUMN insforge_user_id TO supabase_user_id;

-- Rename insforge_login_email to supabase_login_email
ALTER TABLE users 
RENAME COLUMN insforge_login_email TO supabase_login_email;

-- Drop old index
DROP INDEX IF EXISTS idx_users_insforge_login_email;

-- Create new index
CREATE INDEX idx_users_supabase_login_email 
ON users(supabase_login_email);
```

### Step 4: Verify Migration

```sql
-- Verify columns exist
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users' 
  AND column_name IN ('supabase_user_id', 'supabase_login_email');

-- Verify old columns are gone
SELECT column_name 
FROM information_schema.columns 
WHERE table_name = 'users' 
  AND column_name IN ('insforge_user_id', 'insforge_login_email');

-- Verify index exists
SELECT indexname 
FROM pg_indexes 
WHERE tablename = 'users' 
  AND indexname = 'idx_users_supabase_login_email';

-- Verify data integrity
SELECT 
  COUNT(*) as total_users,
  COUNT(DISTINCT supabase_user_id) as unique_user_ids,
  COUNT(DISTINCT supabase_login_email) as unique_emails
FROM users;
```

**Expected Results:**
- 2 columns found (supabase_user_id, supabase_login_email)
- 0 old columns found
- 1 index found (idx_users_supabase_login_email)
- Row counts match pre-migration numbers

### Step 5: Commit or Rollback

**If verification passes:**
```sql
-- Commit the transaction
COMMIT;
```

**If verification fails:**
```sql
-- Rollback the transaction
ROLLBACK;
```

## Post-Migration Verification

### 1. Database Verification

- [ ] **Verify column names**
  ```sql
  \d users
  ```
  - Should show `supabase_user_id` and `supabase_login_email`

- [ ] **Test queries with new column names**
  ```sql
  -- Test SELECT
  SELECT supabase_user_id, supabase_login_email 
  FROM users 
  LIMIT 5;
  
  -- Test WHERE clause
  SELECT * FROM users 
  WHERE supabase_login_email = 'test@example.com';
  
  -- Test index usage
  EXPLAIN SELECT * FROM users 
  WHERE supabase_login_email = 'test@example.com';
  ```

- [ ] **Verify data integrity**
  ```sql
  -- Check for NULL values
  SELECT COUNT(*) FROM users WHERE supabase_user_id IS NULL;
  SELECT COUNT(*) FROM users WHERE supabase_login_email IS NULL;
  
  -- Verify row count unchanged
  SELECT COUNT(*) FROM users;
  ```

### 2. Application Verification

- [ ] **Update backend services with new code**
  ```bash
  # Deploy updated services
  git pull origin main
  # Build and restart services
  ```

- [ ] **Start backend services**
  ```bash
  systemctl start user-service
  systemctl start api-gateway-service
  systemctl start affiliate-service
  systemctl start audit-log-service
  systemctl start fee-service
  systemctl start kibiina-service
  systemctl start loan-credit-service
  systemctl start notification-service
  systemctl start object-storage-service
  ```

- [ ] **Verify services are running**
  ```bash
  systemctl status user-service
  systemctl status api-gateway-service
  # Check other services...
  ```

- [ ] **Test health endpoints**
  ```bash
  curl http://localhost:8080/health
  curl http://localhost:8081/health
  # Test other services...
  ```

### 3. Integration Testing

- [ ] **Test user registration**
  ```bash
  curl -X POST http://localhost:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{
      "email": "test@example.com",
      "password": "testpass123",
      "firstName": "Test",
      "lastName": "User"
    }'
  ```

- [ ] **Test user login**
  ```bash
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{
      "email": "test@example.com",
      "password": "testpass123"
    }'
  ```

- [ ] **Test authenticated endpoints**
  ```bash
  # Get token from login response
  TOKEN="<access_token>"
  
  curl http://localhost:8080/api/v1/auth/profile \
    -H "Authorization: Bearer $TOKEN"
  ```

### 4. Frontend/Mobile Testing

- [ ] **Test web application**
  - Login with existing account
  - Register new account
  - Access protected pages
  - Verify no errors in console

- [ ] **Test mobile application**
  - Login with existing account
  - Register new account
  - Access protected features
  - Verify no crashes

### 5. Monitoring

- [ ] **Check application logs**
  ```bash
  tail -f /var/log/user-service/app.log
  tail -f /var/log/api-gateway/app.log
  ```

- [ ] **Monitor database connections**
  ```sql
  SELECT COUNT(*) FROM pg_stat_activity 
  WHERE datname = 'postgres';
  ```

- [ ] **Check for errors**
  ```bash
  # Check for database errors
  grep -i error /var/log/postgresql/postgresql.log
  
  # Check for application errors
  grep -i error /var/log/user-service/app.log
  ```

## Rollback Plan

### When to Rollback

Rollback if:
- Migration verification fails
- Post-migration tests fail
- Critical errors in application logs
- Data integrity issues detected
- User-facing errors occur

### Rollback Procedure

#### Option 1: Transaction Rollback (During Migration)

If still in transaction:
```sql
ROLLBACK;
```

#### Option 2: Restore from Backup (After Commit)

```bash
# Stop all services
systemctl stop api-gateway-service
systemctl stop user-service
# Stop other services...

# Restore from backup
pg_restore -h db.[YOUR-PROJECT-REF].supabase.co \
           -U postgres \
           -d postgres \
           -c \
           backup_pre_migration_YYYYMMDD_HHMMSS.dump

# Or restore from table backup
psql -h db.[YOUR-PROJECT-REF].supabase.co \
     -U postgres \
     -d postgres \
     -c "DROP TABLE users; 
         ALTER TABLE users_backup_20240101 RENAME TO users;"
```

#### Option 3: Reverse Migration (Manual)

```sql
BEGIN;

-- Rename columns back
ALTER TABLE users 
RENAME COLUMN supabase_user_id TO insforge_user_id;

ALTER TABLE users 
RENAME COLUMN supabase_login_email TO insforge_login_email;

-- Recreate old index
DROP INDEX IF EXISTS idx_users_supabase_login_email;
CREATE INDEX idx_users_insforge_login_email 
ON users(insforge_login_email);

-- Verify
\d users

COMMIT;
```

#### Post-Rollback Steps

- [ ] Redeploy old application code
- [ ] Restart all services
- [ ] Verify application functionality
- [ ] Notify users of issue
- [ ] Document rollback reason
- [ ] Plan corrective actions

## Migration Timeline

### Recommended Schedule

**Maintenance Window:** Saturday 2:00 AM - 4:00 AM (low traffic period)

**Timeline:**
- **1:45 AM** - Send user notification
- **2:00 AM** - Enable maintenance mode
- **2:05 AM** - Stop all services
- **2:10 AM** - Create backups
- **2:20 AM** - Execute migration
- **2:30 AM** - Verify migration
- **2:40 AM** - Deploy updated code
- **2:50 AM** - Start services
- **3:00 AM** - Run integration tests
- **3:15 AM** - Disable maintenance mode
- **3:30 AM** - Monitor for issues
- **4:00 AM** - Confirm success / Execute rollback if needed

## Communication Plan

### Pre-Migration

**T-7 days:**
- Announce scheduled maintenance
- Provide date, time, and expected duration

**T-24 hours:**
- Reminder notification
- Confirm maintenance window

**T-1 hour:**
- Final reminder
- Advise users to save work

### During Migration

- Display maintenance page
- Provide status updates if extended

### Post-Migration

**Success:**
- Announce completion
- Thank users for patience
- Confirm all systems operational

**Issues:**
- Communicate any problems
- Provide estimated resolution time
- Keep users updated

## Success Criteria

Migration is successful when:
- [ ] All database columns renamed correctly
- [ ] All indexes recreated successfully
- [ ] Data integrity verified (no data loss)
- [ ] All backend services running
- [ ] Authentication flow working
- [ ] No errors in application logs
- [ ] Users can login successfully
- [ ] Protected endpoints accessible
- [ ] Mobile app functioning correctly
- [ ] No increase in error rates

## Contact Information

**Migration Team:**
- Database Admin: [contact]
- Backend Lead: [contact]
- DevOps Lead: [contact]

**Escalation:**
- On-call Engineer: [contact]
- Technical Lead: [contact]

## Sign-Off

**Migration Plan Reviewed by:**
- Database Admin: _________________ Date: _______
- Backend Lead: _________________ Date: _______
- DevOps Lead: _________________ Date: _______

**Migration Executed by:**
- Engineer: _________________ Date: _______
- Start Time: _______
- End Time: _______
- Status: [ ] Success [ ] Rollback [ ] Partial

**Notes:**
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
