#!/bin/bash
# Run Flutter app on iOS simulator for Development
# Usage: ./scripts/run-ios-dev.sh

set -e

echo "Running Flutter app on iOS..."

# Check if SUPABASE_ANON_KEY is set
if [ -z "$SUPABASE_ANON_KEY" ]; then
  echo "Error: SUPABASE_ANON_KEY environment variable is not set"
  echo "Please set it with: export SUPABASE_ANON_KEY=your-key-here"
  exit 1
fi

flutter run \
  --dart-define=SUPABASE_URL=https://ablvrbnbsdqshrorhmjf.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=$SUPABASE_ANON_KEY \
  --dart-define=API_BASE_URL=http://localhost:8080

echo "✓ App running on iOS"
