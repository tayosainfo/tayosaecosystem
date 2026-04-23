#!/bin/bash
# Build Android Debug APK for Development
# Usage: ./scripts/build-android-dev.sh

set -e

echo "Building Android Debug APK..."

# Check if SUPABASE_ANON_KEY is set
if [ -z "$SUPABASE_ANON_KEY" ]; then
  echo "Error: SUPABASE_ANON_KEY environment variable is not set"
  echo "Please set it with: export SUPABASE_ANON_KEY=your-key-here"
  exit 1
fi

flutter build apk --debug \
  --dart-define=SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=http://10.0.2.2:8080

echo "✓ Build complete: build/app/outputs/flutter-apk/app-debug.apk"
