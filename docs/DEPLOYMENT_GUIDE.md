# Deployment Guide - Supabase Migration

This guide provides instructions for deploying the Tayosa banking ecosystem with Supabase integration.

## Prerequisites

- Supabase project configured (see SUPABASE_PROJECT_CONFIGURATION.md)
- All environment variables documented
- Database migration tested in staging
- All services updated with Supabase integration

## Environment Variables

### Backend Services

All backend services require these environment variables:

```bash
# Supabase Configuration
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SERVICE_ROLE_KEY]

# Database Configuration
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:6543/postgres?pgbouncer=true

# Service-specific ports
PORT=8080  # api-gateway-service
# PORT=8081  # user-service
# PORT=8016  # affiliate-service
# PORT=8014  # audit-log-service
# PORT=8004  # fee-service
# PORT=8086  # kibiina-service
# PORT=8013  # loan-credit-service
# PORT=8010  # notification-service
# PORT=8015  # object-storage-service

# Service URLs (for API gateway)
USER_SERVICE_URL=http://localhost:8081
OBJECT_STORAGE_SERVICE_URL=http://localhost:8015
AFFILIATE_SERVICE_URL=http://localhost:8016
NOTIFICATION_SERVICE_URL=http://localhost:8010
AUDIT_SERVICE_URL=http://localhost:8014
LOAN_SERVICE_URL=http://localhost:8013
FEE_SERVICE_URL=http://localhost:8004
KIBIINA_SERVICE_URL=http://localhost:8086

# Admin Configuration
ADMIN_API_KEY=<your-admin-secret>
```

### Frontend Application

```bash
# Supabase Configuration
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=[YOUR_ANON_KEY]

# API Gateway
VITE_API_BASE_URL=https://api.tayosa.com
```

### Mobile Application (Flutter)

Build with dart-define flags:

```bash
flutter build apk \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=[YOUR_ANON_KEY] \
  --dart-define=API_BASE_URL=https://api.tayosa.com
```

## Deployment Steps

### 1. Staging Environment Deployment

#### 1.1 Deploy Backend Services

```bash
# Navigate to project root
cd /path/to/tayosaecosystem

# Pull latest code
git pull origin main

# Build services
cd services/user-service && go build -o user-service
cd ../api-gateway-service && go build -o api-gateway-service
cd ../affiliate-service && go build -o affiliate-service
cd ../audit-log-service && go build -o audit-log-service
cd ../fee-service && go build -o fee-service
cd ../kibiina-service && go build -o kibiina-service
cd ../loan-credit-service && go build -o loan-credit-service
cd ../notification-service && go build -o notification-service
cd ../object-storage-service && go build -o object-storage-service

# Copy environment file
cp .env.staging .env

# Start services (using systemd or docker-compose)
systemctl restart user-service
systemctl restart api-gateway-service
# Restart other services...
```

#### 1.2 Deploy Frontend

```bash
# Navigate to frontend directory
cd /path/to/tayosaecosystem/frontend

# Install dependencies
npm install

# Build for staging
npm run build -- --mode staging

# Deploy to hosting (example: Vercel, Netlify, or custom server)
# Option 1: Vercel
vercel --prod

# Option 2: Custom server
rsync -avz dist/ user@staging-server:/var/www/tayosa/
```

#### 1.3 Deploy Mobile App

```bash
# Navigate to mobile app directory
cd /path/to/tayosaecosystem/app/mobile_app

# Get dependencies
flutter pub get

# Build for Android (staging)
flutter build apk \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=[YOUR_ANON_KEY] \
  --dart-define=API_BASE_URL=https://staging-api.tayosa.com \
  --flavor staging

# Build for iOS (staging)
flutter build ios \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=[YOUR_ANON_KEY] \
  --dart-define=API_BASE_URL=https://staging-api.tayosa.com \
  --flavor staging

# Distribute via TestFlight/Firebase App Distribution
```

#### 1.4 Run Database Migration (Staging)

```bash
# Connect to staging database
psql -h db.[YOUR-PROJECT-REF].supabase.co \
     -U postgres \
     -d postgres_staging

# Execute migration
\i db/migrations/011_rename_insforge_to_supabase.sql

# Verify migration
\d users
```

#### 1.5 Test Staging Environment

```bash
# Run integration tests
npm run test:integration

# Test authentication flow
curl -X POST https://staging-api.tayosa.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","firstName":"Test","lastName":"User"}'

# Test protected endpoints
TOKEN="<access_token>"
curl https://staging-api.tayosa.com/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 2. Production Environment Deployment

#### 2.1 Pre-Deployment Checklist

- [ ] Staging environment tested successfully
- [ ] All tests passing
- [ ] Database backup created
- [ ] Rollback plan documented
- [ ] Maintenance window scheduled
- [ ] Users notified
- [ ] Team on standby

#### 2.2 Enable Maintenance Mode

```bash
# Enable maintenance page
touch /var/www/tayosa/maintenance.flag

# Or update load balancer to show maintenance page
```

#### 2.3 Stop Services

```bash
# Stop all backend services
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

#### 2.4 Backup Database

```bash
# Create full backup
pg_dump -h db.[YOUR-PROJECT-REF].supabase.co \
        -U postgres \
        -d postgres \
        -F c \
        -f backup_pre_migration_$(date +%Y%m%d_%H%M%S).dump

# Verify backup
ls -lh backup_pre_migration_*.dump
```

#### 2.5 Execute Database Migration

```bash
# Connect to production database
psql -h db.[YOUR-PROJECT-REF].supabase.co \
     -U postgres \
     -d postgres

# Execute migration in transaction
BEGIN;
\i db/migrations/011_rename_insforge_to_supabase.sql

# Verify migration
\d users
SELECT COUNT(*) FROM users;

# Commit if successful
COMMIT;
```

#### 2.6 Deploy Updated Code

```bash
# Pull latest code
git pull origin main

# Build services
cd services/user-service && go build -o user-service
cd ../api-gateway-service && go build -o api-gateway-service
# Build other services...

# Copy production environment file
cp .env.production .env

# Deploy frontend
cd frontend
npm run build -- --mode production
# Deploy to production hosting
```

#### 2.7 Start Services

```bash
# Start services in order
systemctl start user-service
sleep 5
systemctl start api-gateway-service
sleep 5
systemctl start affiliate-service
systemctl start audit-log-service
systemctl start fee-service
systemctl start kibiina-service
systemctl start loan-credit-service
systemctl start notification-service
systemctl start object-storage-service

# Verify all services are running
systemctl status user-service
systemctl status api-gateway-service
# Check other services...
```

#### 2.8 Verify Deployment

```bash
# Test health endpoints
curl https://api.tayosa.com/health
curl https://api.tayosa.com/health/ready

# Test authentication
curl -X POST https://api.tayosa.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"existing@user.com","password":"password123"}'

# Test protected endpoint
TOKEN="<access_token>"
curl https://api.tayosa.com/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

#### 2.9 Disable Maintenance Mode

```bash
# Remove maintenance flag
rm /var/www/tayosa/maintenance.flag

# Or update load balancer
```

#### 2.10 Monitor

```bash
# Monitor logs
tail -f /var/log/user-service/app.log
tail -f /var/log/api-gateway/app.log

# Monitor database
psql -h db.[YOUR-PROJECT-REF].supabase.co \
     -U postgres \
     -d postgres \
     -c "SELECT COUNT(*) FROM pg_stat_activity;"

# Check error rates
# Monitor application metrics dashboard
```

## Docker Deployment (Alternative)

### Docker Compose Configuration

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  api-gateway:
    build: ./services/api-gateway-service
    ports:
      - "8080:8080"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
      - USER_SERVICE_URL=http://user-service:8081
      - OBJECT_STORAGE_SERVICE_URL=http://object-storage:8015
      - AFFILIATE_SERVICE_URL=http://affiliate:8016
      - NOTIFICATION_SERVICE_URL=http://notification:8010
      - AUDIT_SERVICE_URL=http://audit-log:8014
      - LOAN_SERVICE_URL=http://loan-credit:8013
      - FEE_SERVICE_URL=http://fee:8004
      - KIBIINA_SERVICE_URL=http://kibiina:8086
    depends_on:
      - user-service
    restart: unless-stopped

  user-service:
    build: ./services/user-service
    ports:
      - "8081:8081"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
      - SUPABASE_SERVICE_ROLE_KEY=${SUPABASE_SERVICE_ROLE_KEY}
      - DATABASE_URL=${DATABASE_URL}
    restart: unless-stopped

  affiliate-service:
    build: ./services/affiliate-service
    ports:
      - "8016:8016"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    restart: unless-stopped

  audit-log-service:
    build: ./services/audit-log-service
    ports:
      - "8014:8014"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    restart: unless-stopped

  fee-service:
    build: ./services/fee-service
    ports:
      - "8004:8004"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    restart: unless-stopped

  kibiina-service:
    build: ./services/kibiina-service
    ports:
      - "8086:8086"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    restart: unless-stopped

  loan-credit-service:
    build: ./services/loan-credit-service
    ports:
      - "8013:8013"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    restart: unless-stopped

  notification-service:
    build: ./services/notification-service
    ports:
      - "8010:8010"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    restart: unless-stopped

  object-storage-service:
    build: ./services/object-storage-service
    ports:
      - "8015:8015"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    restart: unless-stopped
```

### Deploy with Docker

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## CI/CD Pipeline Configuration

### GitHub Actions Example

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Production

on:
  push:
    branches:
      - main

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: |
          cd services/user-service
          go test ./...
      
  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to production
        env:
          SUPABASE_URL: ${{ secrets.SUPABASE_URL }}
          SUPABASE_ANON_KEY: ${{ secrets.SUPABASE_ANON_KEY }}
          SUPABASE_SERVICE_ROLE_KEY: ${{ secrets.SUPABASE_SERVICE_ROLE_KEY }}
        run: |
          # Deploy script here
          ./scripts/deploy.sh production
```

## Rollback Procedure

If issues occur after deployment:

```bash
# 1. Enable maintenance mode
touch /var/www/tayosa/maintenance.flag

# 2. Stop services
systemctl stop api-gateway-service
systemctl stop user-service
# Stop other services...

# 3. Restore database backup
pg_restore -h db.[YOUR-PROJECT-REF].supabase.co \
           -U postgres \
           -d postgres \
           -c \
           backup_pre_migration_YYYYMMDD_HHMMSS.dump

# 4. Deploy previous version
git checkout <previous-commit>
# Build and deploy old version

# 5. Start services
systemctl start user-service
systemctl start api-gateway-service
# Start other services...

# 6. Verify rollback
curl https://api.tayosa.com/health

# 7. Disable maintenance mode
rm /var/www/tayosa/maintenance.flag
```

## Post-Deployment Monitoring

### First Hour

- Monitor error rates
- Check authentication success rate
- Verify email delivery
- Monitor database connections
- Check application logs

### First 24 Hours

- Review user feedback
- Monitor performance metrics
- Check for any anomalies
- Verify all features working

### First Week

- Analyze usage patterns
- Review security logs
- Check system performance
- Monitor costs

## Support Contacts

**Deployment Team:**
- DevOps Lead: [contact]
- Backend Lead: [contact]
- Database Admin: [contact]

**Escalation:**
- On-call Engineer: [contact]
- Technical Lead: [contact]

## Deployment Checklist

- [ ] Environment variables configured
- [ ] Database backup created
- [ ] Migration tested in staging
- [ ] All services built successfully
- [ ] Health checks passing
- [ ] Authentication flow tested
- [ ] Protected endpoints accessible
- [ ] Frontend deployed
- [ ] Mobile app deployed
- [ ] Monitoring configured
- [ ] Rollback plan ready
- [ ] Team notified
- [ ] Users notified

## Sign-Off

**Deployment Approved by:**
- Technical Lead: _________________ Date: _______
- DevOps Lead: _________________ Date: _______

**Deployment Executed by:**
- Engineer: _________________ Date: _______
- Start Time: _______
- End Time: _______
- Status: [ ] Success [ ] Rollback

**Notes:**
_________________________________________________________________
_________________________________________________________________
