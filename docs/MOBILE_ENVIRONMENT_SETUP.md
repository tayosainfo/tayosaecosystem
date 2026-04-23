# Mobile App Environment Setup (Flutter)

This document provides instructions for configuring environment variables for the Tayosa Flutter mobile application using Supabase.

## Overview

Flutter applications use `--dart-define` flags to pass environment variables at build/run time. These variables are compiled into the application and accessed via `String.fromEnvironment()`.

## Required Environment Variables

### `SUPABASE_URL`
- **Purpose**: Supabase project URL for mobile app authentication and API access
- **Format**: `https://[project-ref].supabase.co`
- **Example**: `https://[YOUR-PROJECT-REF].supabase.co`
- **How to obtain**: Supabase Dashboard → Project Settings → API → Project URL

### `SUPABASE_ANON_KEY`
- **Purpose**: Public anonymous key for client-side Supabase operations
- **Format**: JWT token (starts with `eyJ`)
- **Security**: Safe to include in mobile app (RLS policies enforced)
- **How to obtain**: Supabase Dashboard → Project Settings → API → anon public key

### `API_BASE_URL`
- **Purpose**: Base URL for backend API gateway
- **Format**: `http://[host]:[port]` or `https://api.yourdomain.com`
- **Example**: `http://localhost:8080` (development) or `https://api.tayosa.com` (production)

## Build Commands

### Development Build (Debug)

```bash
flutter run \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-supabase-anon-key-here \
  --dart-define=API_BASE_URL=http://localhost:8080
```

### Production Build (Android APK)

```bash
flutter build apk --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com
```

### Production Build (Android App Bundle)

```bash
flutter build appbundle --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com
```

### iOS Build

```bash
flutter build ios --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-production-anon-key \
  --dart-define=API_BASE_URL=https://api.tayosa.com
```

## Configuration in Code

The Supabase configuration is defined in `lib/core/network/supabase_client.dart`:

```dart
import 'package:supabase_flutter/supabase_flutter.dart';

class SupabaseConfig {
  static const String supabaseUrl = String.fromEnvironment(
    'SUPABASE_URL',
    defaultValue: 'https://[YOUR-PROJECT-REF].supabase.co',
  );

  static const String supabaseAnonKey = String.fromEnvironment(
    'SUPABASE_ANON_KEY',
    defaultValue: '', // Must be provided via --dart-define
  );

  static Future<void> initialize() async {
    await Supabase.initialize(
      url: supabaseUrl,
      anonKey: supabaseAnonKey,
    );
  }

  static SupabaseClient get client => Supabase.instance.client;
}
```

## Initialization

The Supabase client must be initialized before the app runs. This is done in `lib/main.dart`:

```dart
import 'package:flutter/material.dart';
import 'core/network/supabase_client.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize Supabase
  await SupabaseConfig.initialize();
  
  runApp(MyApp());
}
```

## Environment-Specific Configurations

### Development
- Use local backend: `API_BASE_URL=http://localhost:8080`
- Use development Supabase project (if separate from production)
- Anon key from development Supabase project

### Staging
- Use staging backend: `API_BASE_URL=https://staging-api.tayosa.com`
- Use staging Supabase project
- Anon key from staging Supabase project

### Production
- Use production backend: `API_BASE_URL=https://api.tayosa.com`
- Use production Supabase project: `https://[YOUR-PROJECT-REF].supabase.co`
- Anon key from production Supabase project

## Best Practices

### 1. Never Hardcode Credentials
- Always use `--dart-define` flags
- Never commit API keys to version control
- Use different keys for development and production

### 2. Use Build Scripts
Create shell scripts for common build configurations:

**build-dev.sh**:
```bash
#!/bin/bash
flutter run \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$DEV_SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=http://localhost:8080
```

**build-prod.sh**:
```bash
#!/bin/bash
flutter build apk --release \
  --dart-define=SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$PROD_SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=https://api.tayosa.com
```

### 3. CI/CD Integration
Store environment variables as secrets in your CI/CD platform:
- GitHub Actions: Repository Secrets
- GitLab CI: CI/CD Variables
- Bitrise: Secrets

Example GitHub Actions workflow:
```yaml
- name: Build APK
  run: |
    flutter build apk --release \
      --dart-define=SUPABASE_URL=${{ secrets.SUPABASE_URL }} \
      --dart-define=SUPABASE_ANON_KEY=${{ secrets.SUPABASE_ANON_KEY }} \
      --dart-define=API_BASE_URL=${{ secrets.API_BASE_URL }}
```

## Troubleshooting

### Issue: "Missing SUPABASE_ANON_KEY"
**Solution**: Ensure you're passing the `--dart-define=SUPABASE_ANON_KEY=...` flag when running or building

### Issue: Authentication fails in mobile app
**Solutions**:
1. Verify `SUPABASE_URL` matches your Supabase project
2. Confirm `SUPABASE_ANON_KEY` is correct (check Supabase dashboard)
3. Ensure you've run `flutter pub get` after adding `supabase_flutter` dependency
4. Rebuild the app after changing environment variables

### Issue: Cannot connect to backend API
**Solutions**:
1. Verify `API_BASE_URL` is accessible from your device/emulator
2. For local development on Android emulator, use `http://10.0.2.2:8080` instead of `localhost`
3. For local development on iOS simulator, use `http://localhost:8080`
4. Ensure backend services are running

## Verification

To verify environment variables are correctly configured:

1. **Check Supabase initialization**:
   - App should start without errors
   - Authentication screens should load

2. **Test authentication**:
   - Try registering a new user
   - Try logging in with existing credentials
   - Verify tokens are being issued

3. **Check API connectivity**:
   - Test API calls to backend
   - Verify requests include proper authentication headers

## Additional Resources

- [Flutter Environment Variables](https://docs.flutter.dev/deployment/flavors)
- [Supabase Flutter Documentation](https://supabase.com/docs/reference/dart/introduction)
- [Supabase Flutter Auth](https://supabase.com/docs/reference/dart/auth-signup)
- Main environment variables documentation: `docs/ENVIRONMENT_VARIABLES.md`
