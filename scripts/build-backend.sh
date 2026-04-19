#!/usr/bin/env bash
set -euo pipefail

SERVICES=(
  "api-gateway-service"
  "user-service"
  "fee-service"
  "kibiina-service"
  "affiliate-service"
  "notification-service"
  "audit-log-service"
  "loan-credit-service"
  "object-storage-service"
)

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "Starting backend build check..."
for service in "${SERVICES[@]}"; do
  echo "Building ${service}"
  go -C "services/${service}" build .
done

echo "Backend build check passed for all services."

