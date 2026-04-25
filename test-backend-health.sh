#!/bin/bash

# Backend Health Check Script
# Run this to diagnose the 500 login error

BASE_URL="https://tayosaecosystem.onrender.com"

echo "========================================="
echo "Backend Health Check"
echo "========================================="
echo ""

# Test 1: Health endpoint
echo "Test 1: Checking /health endpoint..."
curl -s "$BASE_URL/health" | jq '.' || echo "FAILED: Health endpoint not responding"
echo ""

# Test 2: Health ready endpoint
echo "Test 2: Checking /health/ready endpoint..."
curl -s "$BASE_URL/health/ready" | jq '.' || echo "FAILED: Health ready endpoint not responding"
echo ""

# Test 3: Public config endpoint
echo "Test 3: Checking /api/v1/auth/public-config..."
curl -s "$BASE_URL/api/v1/auth/public-config" | jq '.' || echo "FAILED: Public config endpoint not responding"
echo ""

# Test 4: Try login (will fail but shows the error)
echo "Test 4: Attempting login to see error..."
curl -s -X POST "$BASE_URL/api/v1/auth/login?client_type=web" \
  -H "Content-Type: application/json" \
  -d '{"identifier":"test@example.com","password":"test123"}' \
  | jq '.' || echo "FAILED: Login endpoint not responding"
echo ""

echo "========================================="
echo "Health Check Complete"
echo "========================================="
