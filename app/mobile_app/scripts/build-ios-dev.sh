#!/bin/bash
# Build iOS Debug for Development
# Usage: ./scripts/build-ios-dev.sh

set -e

echo "Building iOS Debug..."

# Check if SUPABASE_ANON_KEY is set
if [ -z "$SUPABASE_ANON_KEY" ]; then
  echo "Error: SUPABASE_ANON_KEY environment variable is not set"
  echo "Please set it with: export SUPABASE_ANON_KEY=your-key-here"
  exit 1
fi

flutter build ios --debug \
  --dart-define=SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=http://localhost:8080

echo "✓ Build complete"
echo "Open ios/Runner.xcworkspace in Xcode to run on device"
