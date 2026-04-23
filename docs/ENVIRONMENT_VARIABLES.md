# Environment Variables Documentation

This document provides comprehensive information about all required environment variables for the Tayosa banking ecosystem using Supabase as the Backend-as-a-Service provider.

## Overview

The Tayosa system consists of multiple components that require proper environment configuration:
- **Frontend**: React + Vite web application
- **Backend**: 9 Go microservices
- **Mobile**: Flutter application
- **Database**: PostgreSQL (hosted on Supabase)

## Supabase Project Information

**Project URL**: `https://[YOUR-PROJECT-REF].supabase.co`

This is the base URL for all Supabase API endpoints including authentication, database, and storage services.

## Obtaining Supabase Credentials

### Step 1: Access Supabase Dashboard
1. Navigate to [https://supabase.com](https://supabase.com)
2. Sign in to your account
3. Select your project (or create a new one)

### Step 2: Locate API Keys
1. In the Supabase dashboard, go to **Project Settings** (gear icon in sidebar)
2. Click on **API** in the settings menu
3. You'll find the following credentials:
   - **Project URL**: Your unique Supabase project URL
   - **anon/public key**: Public key for client-side operations
   - **service_role key**: Secret key for server-side admin operations

### Step 3: Copy Credentials
- Copy the **Project URL** (format: `https://[project-ref].supabase.co`)
- Copy the **anon public** key (starts with `eyJ...`)
- Copy the **service_role** key (starts with `eyJ...`)

⚠️ **Security Warning**: The `service_role` key bypasses Row Level Security (RLS) and should NEVER be exposed in frontend code or committed to version control.

## Environment Variables Reference

### Frontend Variables (React + Vite)

These variables are used by the web application and must be prefixed with `VITE_` to be accessible in the browser.

#### `VITE_API_BASE_URL`
- **Purpose**: Base URL for backend API gateway
- **Format**: `http://localhost:8080` (development) or `https://api.yourdomain.com` (production)
- **Example**: `VITE_API_BASE_URL=http://localhost:8080`
- **Required**: Yes
- **Exposed to browser**: Yes

#### `VITE_SUPABASE_URL`
- **Purpose**: Supabase project URL for frontend authentication and database access
- **Format**: `https://[project-ref].supabase.co`
- **Example**: `VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co`
- **Required**: Yes
- **Exposed to browser**: Yes
- **How to obtain**: Supabase Dashboard → Project Settings → API → Project URL

#### `VITE_SUPABASE_ANON_KEY`
- **Purpose**: Public anonymous key for client-side Supabase operations
- **Format**: JWT token (starts with `eyJ`)
- **Example**: `VITE_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
- **Required**: Yes
- **Exposed to browser**: Yes (safe to expose)
- **Security**: Provides read access with Row Level Security (RLS) policies enforced
- **How to obtain**: Supabase Dashboard → Project Settings → API → Project API keys → anon public

#### `VITE_ADMIN_API_KEY`
- **Purpose**: Admin API key for privileged frontend operations
- **Format**: Custom secret string
- **Example**: `VITE_ADMIN_API_KEY=change-me-admin-secret`
- **Required**: Yes
- **Exposed to browser**: Yes
- **Security**: Change default value in production

### Backend Variables (Go Microservices)

These variables are used by backend services and should be kept secure.

#### `SUPABASE_URL`
- **Purpose**: Supabase project URL for backend services
- **Format**: `https://[project-ref].supabase.co`
- **Example**: `SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co`
- **Required**: Yes
- **Used by**: All backend services for token validation
- **How to obtain**: Same as `VITE_SUPABASE_URL`

#### `SUPABASE_ANON_KEY`
- **Purpose**: Public anonymous key for backend client-side operations
- **Format**: JWT token (starts with `eyJ`)
- **Example**: `SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
- **Required**: Yes
- **Used by**: Backend services for public API calls to Supabase
- **How to obtain**: Same as `VITE_SUPABASE_ANON_KEY`

#### `SUPABASE_SERVICE_ROLE_KEY`
- **Purpose**: Secret key for privileged backend operations
- **Format**: JWT token (starts with `eyJ`)
- **Example**: `SUPABASE_SERVICE_ROLE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
- **Required**: Yes
- **Used by**: Backend services for admin operations (user lookup, bypassing RLS)
- **Security**: ⚠️ **CRITICAL** - Never expose in frontend or commit to version control
- **How to obtain**: Supabase Dashboard → Project Settings → API → Project API keys → service_role

#### `ADMIN_API_KEY`
- **Purpose**: Admin API key for backend service authentication
- **Format**: Custom secret string
- **Example**: `ADMIN_API_KEY=change-me-admin-secret`
- **Required**: Yes
- **Security**: Change default value in production

#### `DATABASE_URL`
- **Purpose**: PostgreSQL connection string for direct database access
- **Format**: `postgresql://[user]:[password]@[host]:[port]/[database]?sslmode=require`
- **Example (Supabase)**: `postgresql://postgres:PASSWORD@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres?sslmode=require`
- **Example (Local)**: `postgres://tayosa:password@localhost:5432/tayosa_banking?sslmode=disable`
- **Required**: Optional (services can fall back to in-memory storage)
- **How to obtain**: Supabase Dashboard → Project Settings → Database → Connection string → URI
- **Note**: Replace `[YOUR-PASSWORD]` with your actual database password

#### Service Routing Variables

These variables configure the API gateway to route requests to backend microservices.

- `USER_SERVICE_URL` - Default: `http://localhost:8081`
- `OBJECT_STORAGE_SERVICE_URL` - Default: `http://localhost:8015`
- `AFFILIATE_SERVICE_URL` - Default: `http://localhost:8016`
- `NOTIFICATION_SERVICE_URL` - Default: `http://localhost:8010`
- `AUDIT_SERVICE_URL` - Default: `http://localhost:8014`
- `LOAN_SERVICE_URL` - Default: `http://localhost:8013`
- `FEE_SERVICE_URL` - Default: `http://localhost:8004`
- `KIBIINA_SERVICE_URL` - Default: `http://localhost:8086`

**Format**: `http://[host]:[port]`
**Required**: Yes for production, defaults work for local development

### Mobile Variables (Flutter)

Flutter applications use `--dart-define` flags to pass environment variables at build time. These flags must be provided every time you run, build, or test the Flutter application.

#### `SUPABASE_URL`
- **Purpose**: Supabase project URL for mobile app
- **Format**: `https://[project-ref].supabase.co`
- **Example**: `--dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co`
- **Required**: Yes
- **How to obtain**: Same as frontend `VITE_SUPABASE_URL`

#### `SUPABASE_ANON_KEY`
- **Purpose**: Public anonymous key for mobile Supabase operations
- **Format**: JWT token
- **Example**: `--dart-define=SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
- **Required**: Yes
- **How to obtain**: Same as frontend `VITE_SUPABASE_ANON_KEY`

#### `API_BASE_URL`
- **Purpose**: Base URL for backend API gateway
- **Format**: `http://[host]:[port]` or `https://api.yourdomain.com`
- **Example**: `--dart-define=API_BASE_URL=http://localhost:8080`
- **Required**: Yes

#### Flutter Build Commands Reference

**Development - Run on Device/Emulator**
```bash
# Android
flutter run \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-supabase-anon-key-here \
  --dart-define=API_BASE_URL=http://10.0.2.2:8080

# iOS
flutter run \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-supabase-anon-key-here \
  --dart-define=API_BASE_URL=http://localhost:8080
```

**Note**: Android emulator uses `10.0.2.2` to access host machine's localhost. iOS simulator can use `localhost` directly.

**Development - Debug Build**
```bash
# Android APK (debug)
flutter build apk --debug \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-supabase-anon-key-here \
  --dart-define=API_BASE_URL=http://10.0.2.2:8080

# iOS (debug)
flutter build ios --debug \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-supabase-anon-key-here \
  --dart-define=API_BASE_URL=http://localhost:8080
```

**Production - Release Build**
```bash
# Android APK (release)
flutter build apk --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com

# Android App Bundle (for Play Store)
flutter build appbundle --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com

# iOS (release)
flutter build ios --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com

# iOS Archive (for App Store)
flutter build ipa --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com
```

**Testing**
```bash
# Run tests with environment variables
flutter test \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-test-anon-key \
  --dart-define=API_BASE_URL=http://localhost:8080
```

#### Accessing Environment Variables in Flutter Code

In your Dart code, access these variables using `String.fromEnvironment`:

```dart
// lib/core/network/supabase_client.dart
class SupabaseConfig {
  static const String supabaseUrl = String.fromEnvironment(
    'SUPABASE_URL',
    defaultValue: 'https://[YOUR-PROJECT-REF].supabase.co',
  );
  
  static const String supabaseAnonKey = String.fromEnvironment(
    'SUPABASE_ANON_KEY',
    defaultValue: '', // No default for security
  );
  
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );
}
```

#### Flutter Environment Configuration Best Practices

1. **Never hardcode credentials**: Always use `--dart-define` flags
2. **Use build scripts**: Create shell scripts to avoid typing long commands
3. **Separate environments**: Use different Supabase projects for dev/staging/prod
4. **Validate on startup**: Check that required variables are set when app initializes
5. **Use CI/CD secrets**: Store production keys in your CI/CD pipeline secrets

**Example Build Script** (`scripts/build-android-dev.sh`):
```bash
#!/bin/bash
flutter build apk --debug \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=http://10.0.2.2:8080
```

**Example Build Script** (`scripts/build-ios-prod.sh`):
```bash
#!/bin/bash
flutter build ipa --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=https://api.tayosa.com
```

## Configuration Examples

### Development Environment

#### Frontend (.env)
```bash
# API Gateway
VITE_API_BASE_URL=http://localhost:8080

# Supabase Configuration
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=your-supabase-anon-key-here

# Admin
VITE_ADMIN_API_KEY=dev-admin-secret
```

#### Backend (.env)
```bash
# Supabase Configuration
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=your-supabase-anon-key-here
SUPABASE_SERVICE_ROLE_KEY=your-supabase-service-role-key-here

# Admin
ADMIN_API_KEY=dev-admin-secret

# Database (optional for development)
DATABASE_URL=postgres://tayosa:password@localhost:5432/tayosa_banking?sslmode=disable

# Service URLs (defaults work for local development)
USER_SERVICE_URL=http://localhost:8081
OBJECT_STORAGE_SERVICE_URL=http://localhost:8015
AFFILIATE_SERVICE_URL=http://localhost:8016
NOTIFICATION_SERVICE_URL=http://localhost:8010
AUDIT_SERVICE_URL=http://localhost:8014
LOAN_SERVICE_URL=http://localhost:8013
FEE_SERVICE_URL=http://localhost:8004
KIBIINA_SERVICE_URL=http://localhost:8086
```

#### Mobile (Flutter Development)
```bash
# Android Emulator (use 10.0.2.2 for host machine)
flutter run \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-supabase-anon-key-here \
  --dart-define=API_BASE_URL=http://10.0.2.2:8080

# iOS Simulator (use localhost)
flutter run \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-supabase-anon-key-here \
  --dart-define=API_BASE_URL=http://localhost:8080

# Physical Device (use your machine's IP address)
flutter run \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-supabase-anon-key-here \
  --dart-define=API_BASE_URL=http://192.168.1.100:8080
```

### Production Environment

#### Frontend (.env.production)
```bash
# API Gateway
VITE_API_BASE_URL=https://api.tayosa.com

# Supabase Configuration
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=your-production-anon-key

# Admin
VITE_ADMIN_API_KEY=secure-production-admin-key
```

#### Backend (Production environment variables)
```bash
# Supabase Configuration
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=your-production-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-production-service-role-key

# Admin
ADMIN_API_KEY=secure-production-admin-key

# Database
DATABASE_URL=postgresql://postgres:SECURE_PASSWORD@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres?sslmode=require

# Service URLs (production internal network)
USER_SERVICE_URL=http://user-service:8081
OBJECT_STORAGE_SERVICE_URL=http://storage-service:8015
AFFILIATE_SERVICE_URL=http://affiliate-service:8016
NOTIFICATION_SERVICE_URL=http://notification-service:8010
AUDIT_SERVICE_URL=http://audit-service:8014
LOAN_SERVICE_URL=http://loan-service:8013
FEE_SERVICE_URL=http://fee-service:8004
KIBIINA_SERVICE_URL=http://kibiina-service:8086
```

#### Mobile (Production Build)
```bash
# Android Release APK
flutter build apk --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com

# Android App Bundle (Google Play Store)
flutter build appbundle --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com

# iOS Release (App Store)
flutter build ipa --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com
```

## Security Best Practices

### 1. Never Commit Secrets
- Add `.env` files to `.gitignore`
- Use `.env.example` as a template with placeholder values
- Never commit actual API keys or passwords

### 2. Rotate Keys Regularly
- Change `ADMIN_API_KEY` periodically
- Rotate Supabase keys if compromised
- Update database passwords on a schedule

### 3. Use Different Keys Per Environment
- Development, staging, and production should use separate Supabase projects
- Each environment should have unique admin keys
- Never use production keys in development

### 4. Protect Service Role Key
- The `SUPABASE_SERVICE_ROLE_KEY` bypasses all Row Level Security (RLS)
- Only use in backend services, never in frontend or mobile apps
- Store securely in environment variables or secrets management systems
- Limit access to this key to essential personnel only

### 5. Environment-Specific Configuration
- Use `.env.local` for local overrides (add to `.gitignore`)
- Use CI/CD secrets management for production deployments
- Validate all required variables on application startup

## Troubleshooting

### Frontend Cannot Connect to Supabase
**Symptoms**: Authentication fails, "Invalid API key" errors

**Solutions**:
1. Verify `VITE_SUPABASE_URL` matches your Supabase project URL
2. Confirm `VITE_SUPABASE_ANON_KEY` is correct (check Supabase dashboard)
3. Ensure variables are prefixed with `VITE_` for Vite to expose them
4. Restart development server after changing `.env` files

### Backend Token Validation Fails
**Symptoms**: 401 Unauthorized errors, "Invalid token" messages

**Solutions**:
1. Verify `SUPABASE_URL` and `SUPABASE_ANON_KEY` match frontend values
2. Check that `SUPABASE_SERVICE_ROLE_KEY` is set correctly
3. Ensure tokens are being passed in `Authorization: Bearer <token>` header
4. Verify Supabase project is active and not paused

### Mobile App Authentication Issues
**Symptoms**: Login fails, "Network error" messages

**Solutions**:
1. Verify `--dart-define` flags are passed correctly during build
2. Check that `SUPABASE_URL` and `SUPABASE_ANON_KEY` match other environments
3. Ensure `API_BASE_URL` points to accessible backend:
   - Android emulator: Use `http://10.0.2.2:8080` for localhost
   - iOS simulator: Use `http://localhost:8080`
   - Physical device: Use your machine's IP (e.g., `http://192.168.1.100:8080`)
4. Rebuild app after changing environment variables (hot reload won't pick up `--dart-define` changes)
5. Check that Supabase Flutter SDK is properly initialized in `main.dart`
6. Verify network permissions in `AndroidManifest.xml` and `Info.plist`
7. Test API connectivity using a tool like Postman before testing in app

### Database Connection Fails
**Symptoms**: "Connection refused", "Database unavailable" errors

**Solutions**:
1. Verify `DATABASE_URL` format is correct
2. Check database password is correct (no special characters causing issues)
3. Ensure `sslmode=require` for Supabase, `sslmode=disable` for local
4. Verify database is accessible from your network
5. Check Supabase project is not paused or over quota

### Missing Environment Variables
**Symptoms**: Application crashes on startup, "Missing required variable" errors

**Solutions**:
1. Copy `.env.example` to `.env` and fill in all values
2. Verify all required variables are set (see reference above)
3. Check for typos in variable names
4. Ensure `.env` file is in the correct directory (project root)

## Validation Checklist

Before deploying, verify:

- [ ] All Supabase credentials are obtained from dashboard
- [ ] `VITE_SUPABASE_URL` and `SUPABASE_URL` match
- [ ] `VITE_SUPABASE_ANON_KEY` and `SUPABASE_ANON_KEY` match
- [ ] `SUPABASE_SERVICE_ROLE_KEY` is set in backend only
- [ ] `ADMIN_API_KEY` is changed from default value
- [ ] `DATABASE_URL` is correct for target environment
- [ ] All service URLs are accessible
- [ ] `.env` files are in `.gitignore`
- [ ] Production uses different keys than development
- [ ] Mobile build commands include all `--dart-define` flags

## Additional Resources

- [Supabase Documentation](https://supabase.com/docs)
- [Supabase API Reference](https://supabase.com/docs/reference)
- [Vite Environment Variables](https://vitejs.dev/guide/env-and-mode.html)
- [Flutter Environment Variables](https://docs.flutter.dev/deployment/flavors)

## Support

For issues related to:
- **Supabase configuration**: Check Supabase Dashboard or contact Supabase support
- **Application setup**: Refer to `docs/RUNBOOK_LOCAL_DEVELOPMENT.md`
- **Authentication flow**: See `docs/AUTH_IDENTITY_AND_LOGIN_STRATEGY.md`
- **Service architecture**: Review `docs/ARCHITECTURE_WORKFLOW.md`
