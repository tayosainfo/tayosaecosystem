# TAYOSA Mobile App

Flutter 3.24+ application for iOS and Android.

## Build Instructions

### Android
```bash
flutter build apk --release
```

### iOS
```bash
flutter build ios --release
```

## Security Notes
- Certificate pinning enabled via Dio interceptors
- Biometric authentication via local_auth
- Secure storage for tokens via flutter_secure_storage
- Play Integrity API (Android) / App Attest (iOS) for device attestation