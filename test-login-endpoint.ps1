# PowerShell script to test the login endpoint
# Run this to diagnose the 500 error

$baseUrl = "https://tayosaecosystem.onrender.com"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Login Endpoint Diagnostic Test" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: Health check
Write-Host "Test 1: Checking backend health..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/health/ready" -UseBasicParsing
    Write-Host "✅ Backend is ready" -ForegroundColor Green
    Write-Host ($health | ConvertTo-Json) -ForegroundColor Gray
} catch {
    Write-Host "❌ Backend health check failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 2: Public config
Write-Host "Test 2: Checking public config..." -ForegroundColor Yellow
try {
    $config = Invoke-RestMethod -Uri "$baseUrl/api/v1/auth/public-config" -UseBasicParsing
    Write-Host "✅ Public config accessible" -ForegroundColor Green
    Write-Host ($config | ConvertTo-Json) -ForegroundColor Gray
} catch {
    Write-Host "❌ Public config failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 3: Login with invalid credentials (should return 401, not 500)
Write-Host "Test 3: Testing login with invalid credentials..." -ForegroundColor Yellow
$testBody = @{
    identifier = "test@example.com"
    password = "wrongpassword"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/auth/login?client_type=web" `
        -Method POST `
        -ContentType "application/json" `
        -Body $testBody `
        -UseBasicParsing
    Write-Host "Response: $($response.StatusCode)" -ForegroundColor Gray
    Write-Host ($response.Content | ConvertFrom-Json | ConvertTo-Json) -ForegroundColor Gray
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 401 -or $statusCode -eq 400) {
        Write-Host "✅ Login endpoint working (returned $statusCode as expected)" -ForegroundColor Green
    } elseif ($statusCode -eq 500) {
        Write-Host "❌ 500 Internal Server Error detected!" -ForegroundColor Red
        Write-Host "This is the error we're trying to fix." -ForegroundColor Red
    } else {
        Write-Host "⚠️  Unexpected status code: $statusCode" -ForegroundColor Yellow
    }
    
    try {
        $errorBody = $_.ErrorDetails.Message | ConvertFrom-Json
        Write-Host "Error details:" -ForegroundColor Gray
        Write-Host ($errorBody | ConvertTo-Json) -ForegroundColor Gray
    } catch {
        Write-Host "Could not parse error details" -ForegroundColor Gray
    }
}
Write-Host ""

# Test 4: Login with malformed data (should return 400, not 500)
Write-Host "Test 4: Testing login with malformed data..." -ForegroundColor Yellow
$malformedBody = @{
    identifier = ""
    password = ""
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/auth/login?client_type=web" `
        -Method POST `
        -ContentType "application/json" `
        -Body $malformedBody `
        -UseBasicParsing
    Write-Host "Response: $($response.StatusCode)" -ForegroundColor Gray
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 400) {
        Write-Host "✅ Validation working (returned 400 as expected)" -ForegroundColor Green
    } elseif ($statusCode -eq 500) {
        Write-Host "❌ 500 Internal Server Error on empty credentials!" -ForegroundColor Red
    } else {
        Write-Host "⚠️  Unexpected status code: $statusCode" -ForegroundColor Yellow
    }
}
Write-Host ""

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Diagnostic Complete" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Summary:" -ForegroundColor Cyan
Write-Host "- If all tests passed, the login endpoint is working correctly" -ForegroundColor White
Write-Host "- If you see 500 errors, there's a backend issue that needs fixing" -ForegroundColor White
Write-Host "- Check Render logs for the actual error message" -ForegroundColor White
