param(
  [switch]$VerboseOutput
)

$ErrorActionPreference = "Stop"

$services = @(
  "api-gateway-service",
  "user-service",
  "fee-service",
  "kibiina-service",
  "affiliate-service",
  "notification-service",
  "audit-log-service",
  "loan-credit-service",
  "object-storage-service"
)

$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

Write-Host "Starting backend build check..." -ForegroundColor Cyan

foreach ($service in $services) {
  $servicePath = Join-Path "services" $service
  Write-Host "Building $service" -ForegroundColor Yellow
  if ($VerboseOutput) {
    go -C $servicePath build -v .
  } else {
    go -C $servicePath build .
  }
  if ($LASTEXITCODE -ne 0) {
    throw "Build failed for $service"
  }
}

Write-Host "Backend build check passed for all services." -ForegroundColor Green

