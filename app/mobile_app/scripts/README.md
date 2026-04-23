# Flutter Build Scripts

This directory contains convenience scripts for building and running the Tayosa mobile app with proper environment configuration.

## Prerequisites

1. **Set the SUPABASE_ANON_KEY environment variable**:
   ```bash
   export SUPABASE_ANON_KEY=your-supabase-anon-key-here
   ```

2. **Make scripts executable** (first time only):
   ```bash
   chmod +x scripts/*.sh
   ```

## Development Scripts

### Run on Android Emulator
```bash
./scripts/run-android-dev.sh
```
- Uses `http://10.0.2.2:8080` to connect to backend on host machine
- Connects to development Supabase project

### Run on iOS Simulator
```bash
./scripts/run-ios-dev.sh
```
- Uses `http://localhost:8080` to connect to backend
- Connects to development Supabase project

### Build Android Debug APK
```bash
./scripts/build-android-dev.sh
```
- Output: `build/app/outputs/flutter-apk/app-debug.apk`
- For testing on physical devices

### Build iOS Debug
```bash
./scripts/build-ios-dev.sh
```
- Opens Xcode workspace for device deployment

## Production Scripts

### Build Android Release APK
```bash
export SUPABASE_ANON_KEY=your-production-key
./scripts/build-android-prod.sh
```
- Output: `build/app/outputs/flutter-apk/app-release.apk`
- For direct distribution

### Build Android App Bundle (Play Store)
```bash
export SUPABASE_ANON_KEY=your-production-key
./scripts/build-appbundle-prod.sh
```
- Output: `build/app/outputs/bundle/release/app-release.aab`
- Upload to Google Play Console

### Build iOS IPA (App Store)
```bash
export SUPABASE_ANON_KEY=your-production-key
./scripts/build-ios-prod.sh
```
- Output: `build/ios/ipa/*.ipa`
- Upload to App Store Connect

## Environment Variables

All scripts use the following environment variables:

| Variable | Development | Production |
|----------|-------------|------------|
| `SUPABASE_URL` | `https://ablvrbnbsdqshrorhmjf.supabase.co` | Same |
| `SUPABASE_ANON_KEY` | Development key | Production key |
| `API_BASE_URL` | `http://10.0.2.2:8080` (Android)<br>`http://localhost:8080` (iOS) | `https://api.tayosa.com` |

## Tips

1. **Use different Supabase keys for dev/prod**: Never use production keys in development
2. **Store keys securely**: Add `SUPABASE_ANON_KEY` to your shell profile or use a secrets manager
3. **Physical devices**: For testing on physical devices, replace `API_BASE_URL` with your machine's IP address
4. **CI/CD**: Store `SUPABASE_ANON_KEY` as a secret in your CI/CD pipeline

## Troubleshooting

### "SUPABASE_ANON_KEY environment variable is not set"
Set the environment variable before running the script:
```bash
export SUPABASE_ANON_KEY=your-key-here
./scripts/run-android-dev.sh
```

### "Permission denied"
Make the script executable:
```bash
chmod +x scripts/run-android-dev.sh
```

### Android emulator can't connect to backend
- Ensure backend is running on `http://localhost:8080`
- Android emulator uses `10.0.2.2` to access host machine's localhost
- Check firewall settings

### iOS simulator can't connect to backend
- Ensure backend is running on `http://localhost:8080`
- iOS simulator can access `localhost` directly
- Check that backend is listening on all interfaces (not just 127.0.0.1)

## Manual Build Commands

If you prefer to run commands manually without scripts:

**Android Development**:
```bash
flutter run \
  --dart-define=SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=http://10.0.2.2:8080
```

**iOS Development**:
```bash
flutter run \
  --dart-define=SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=http://localhost:8080
```

For more details, see `docs/ENVIRONMENT_VARIABLES.md`.
