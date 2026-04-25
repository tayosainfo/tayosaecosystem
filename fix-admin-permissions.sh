#!/bin/bash

# Admin Permission Fix Deployment Script
# This script helps deploy the admin permission fixes

echo "========================================="
echo "Admin Permission Fix Deployment"
echo "========================================="
echo ""

# Step 1: Show the SQL to run
echo "STEP 1: Run this SQL in your Supabase SQL Editor"
echo "========================================="
echo ""
cat db/migrations/020_fix_custom_claims_app_metadata.sql
echo ""
echo "========================================="
echo ""
read -p "Have you run the SQL in Supabase? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Please run the SQL first, then re-run this script."
    exit 1
fi

# Step 2: Commit and push changes
echo ""
echo "STEP 2: Committing and pushing changes..."
echo "========================================="
git add services/api-gateway-service/auth.go
git add db/migrations/018_add_admin_rls_policies.sql
git add db/migrations/020_fix_custom_claims_app_metadata.sql
git add ADMIN_PERMISSION_FIX.md
git commit -m "Fix admin permission issue - update JWT claims and RLS policies"
git push

echo ""
echo "========================================="
echo "STEP 3: Redeploy on Render"
echo "========================================="
echo ""
echo "If auto-deploy is enabled, your api-gateway-service should redeploy automatically."
echo "Otherwise, manually trigger a deploy in the Render dashboard."
echo ""
echo "========================================="
echo "STEP 4: Admin User Re-login"
echo "========================================="
echo ""
echo "IMPORTANT: The admin user (baylesinfo@gmail.com) must:"
echo "1. Log out from the admin dashboard"
echo "2. Log back in"
echo "3. Try accessing the KYC endpoint again"
echo ""
echo "========================================="
echo "Deployment complete!"
echo "========================================="
