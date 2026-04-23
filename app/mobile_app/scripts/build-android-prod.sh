#!/bin/bash
# Build Android Release APK for Production
# Usage: ./scripts/build-android-prod.sh

set -e

echo "Building Android Release APK..."

# Check if SUPABASE_ANON_KEY is set
if [ -z "$SUPABASE_ANON_KEY" ]; then
  echo "Error: SUPABASE_ANON_KEY environment variable is not set"
  echo "Please set it with: export SUPABASE_ANON_KEY=your-production-key-here"
  exit 1
fi

flutter build apk --release \
  --dart-define=SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=https://api.tayosa.com

echo "✓ Build complete: build/app/outputs/flutter-apk/app-release.apk"
