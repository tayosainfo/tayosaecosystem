#!/bin/bash
# Build Android App Bundle for Google Play Store
# Usage: ./scripts/build-appbundle-prod.sh

set -e

echo "Building Android App Bundle for Play Store..."

# Check if SUPABASE_ANON_KEY is set
if [ -z "$SUPABASE_ANON_KEY" ]; then
  echo "Error: SUPABASE_ANON_KEY environment variable is not set"
  echo "Please set it with: export SUPABASE_ANON_KEY=your-production-key-here"
  exit 1
fi

flutter build appbundle --release \
  --dart-define=SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=https://api.tayosa.com

echo "✓ Build complete: build/app/outputs/bundle/release/app-release.aab"
echo "Ready to upload to Google Play Console"
