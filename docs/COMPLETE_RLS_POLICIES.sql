-- Complete RLS Policies for All Tables
-- Run this in Supabase SQL Editor to set up security for all tables

-- =============================================================================
-- 1. ONBOARDING PROFILES - Users can manage their own onboarding data
-- =============================================================================
ALTER TABLE onboarding_profiles ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view own onboarding profile" ON onboarding_profiles
FOR SELECT USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = onboarding_profiles.user_id
  )
);

CREATE POLICY "Users can update own onboarding profile" ON onboarding_profiles
FOR UPDATE USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = onboarding_profiles.user_id
  )
);

CREATE POLICY "Users can insert own onboarding profile" ON onboarding_profiles
FOR INSERT WITH CHECK (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = onboarding_profiles.user_id
  )
);

CREATE POLICY "Service role full access onboarding" ON onboarding_profiles
USING (auth.role() = 'service_role');

-- =============================================================================
-- 2. USER CONSENTS - Users can manage their own consent records
-- =============================================================================
ALTER TABLE user_consents ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view own consents" ON user_consents
FOR SELECT USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = user_consents.user_id
  )
);

CREATE POLICY "Users can update own consents" ON user_consents
FOR UPDATE USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = user_consents.user_id
  )
);

CREATE POLICY "Users can insert own consents" ON user_consents
FOR INSERT WITH CHECK (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = user_consents.user_id
  )
);

CREATE POLICY "Service role full access consents" ON user_consents
USING (auth.role() = 'service_role');

-- =============================================================================
-- 3. KYC PROFILES - Users can manage their own KYC data
-- =============================================================================
ALTER TABLE kyc_profiles ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view own kyc profile" ON kyc_profiles
FOR SELECT USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = kyc_profiles.user_id
  )
);

CREATE POLICY "Users can update own kyc profile" ON kyc_profiles
FOR UPDATE USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = kyc_profiles.user_id
  )
);

CREATE POLICY "Users can insert own kyc profile" ON kyc_profiles
FOR INSERT WITH CHECK (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = kyc_profiles.user_id
  )
);

CREATE POLICY "Service role full access kyc" ON kyc_profiles
USING (auth.role() = 'service_role');

-- =============================================================================
-- 4. SACCO MEMBERSHIPS - Users can manage their own membership data
-- =============================================================================
ALTER TABLE sacco_memberships ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view own sacco membership" ON sacco_memberships
FOR SELECT USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = sacco_memberships.user_id
  )
);

CREATE POLICY "Users can update own sacco membership" ON sacco_memberships
FOR UPDATE USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = sacco_memberships.user_id
  )
);

CREATE POLICY "Users can insert own sacco membership" ON sacco_memberships
FOR INSERT WITH CHECK (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = sacco_memberships.user_id
  )
);

CREATE POLICY "Service role full access sacco memberships" ON sacco_memberships
USING (auth.role() = 'service_role');

-- =============================================================================
-- 5. KIBIINA PREFERENCES - Users can manage their own kibiina preferences
-- =============================================================================
ALTER TABLE kibiina_preferences ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view own kibiina preferences" ON kibiina_preferences
FOR SELECT USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = kibiina_preferences.user_id
  )
);

CREATE POLICY "Users can update own kibiina preferences" ON kibiina_preferences
FOR UPDATE USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = kibiina_preferences.user_id
  )
);

CREATE POLICY "Users can insert own kibiina preferences" ON kibiina_preferences
FOR INSERT WITH CHECK (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = kibiina_preferences.user_id
  )
);

CREATE POLICY "Service role full access kibiina preferences" ON kibiina_preferences
USING (auth.role() = 'service_role');

-- =============================================================================
-- 6. SHARES LEDGER - Users can view their own shares transactions
-- =============================================================================
ALTER TABLE shares_ledger ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view own shares ledger" ON shares_ledger
FOR SELECT USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = shares_ledger.user_id
  )
);

CREATE POLICY "Users can insert own shares transactions" ON shares_ledger
FOR INSERT WITH CHECK (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = shares_ledger.user_id
  )
);

CREATE POLICY "Service role full access shares ledger" ON shares_ledger
USING (auth.role() = 'service_role');

-- =============================================================================
-- 7. AFFILIATE REFERRALS - Users can view referrals they made or received
-- =============================================================================
ALTER TABLE affiliate_referrals ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view referrals they made" ON affiliate_referrals
FOR SELECT USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = affiliate_referrals.referrer_user_id
  )
);

CREATE POLICY "Users can view referrals they received" ON affiliate_referrals
FOR SELECT USING (
  auth.uid()::text = (
    SELECT supabase_user_id FROM users_identity WHERE user_id = affiliate_referrals.referee_user_id
  )
);

CREATE POLICY "Service role full access affiliate referrals" ON affiliate_referrals
USING (auth.role() = 'service_role');

-- =============================================================================
-- 8. UGANDA GEO UNITS - Read-only for all authenticated users (reference data)
-- =============================================================================
ALTER TABLE uganda_geo_units ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Authenticated users can view geo data" ON uganda_geo_units
FOR SELECT TO authenticated USING (true);

CREATE POLICY "Service role full access geo units" ON uganda_geo_units
USING (auth.role() = 'service_role');

-- =============================================================================
-- 9. PARISH SACCOS - Read-only for all authenticated users (reference data)
-- =============================================================================
ALTER TABLE parish_saccos ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Authenticated users can view parish saccos" ON parish_saccos
FOR SELECT TO authenticated USING (true);

CREATE POLICY "Service role full access parish saccos" ON parish_saccos
USING (auth.role() = 'service_role');

-- =============================================================================
-- 10. VILLAGE KIBIINA GROUPS - Read-only for all authenticated users (reference data)
-- =============================================================================
ALTER TABLE village_kibiina_groups ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Authenticated users can view kibiina groups" ON village_kibiina_groups
FOR SELECT TO authenticated USING (true);

CREATE POLICY "Service role full access kibiina groups" ON village_kibiina_groups
USING (auth.role() = 'service_role');

-- =============================================================================
-- 11. ADMIN SETTINGS - Read-only for all authenticated users (app configuration)
-- =============================================================================
ALTER TABLE admin_settings ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Authenticated users can view admin settings" ON admin_settings
FOR SELECT TO authenticated USING (true);

CREATE POLICY "Service role full access admin settings" ON admin_settings
USING (auth.role() = 'service_role');

-- =============================================================================
-- SUMMARY
-- =============================================================================
-- This script creates RLS policies for all tables in your database:
-- 
-- USER-SPECIFIC DATA (users can only see/modify their own):
-- - onboarding_profiles
-- - user_consents  
-- - kyc_profiles
-- - sacco_memberships
-- - kibiina_preferences
-- - shares_ledger
-- - affiliate_referrals (can see referrals they made or received)
--
-- REFERENCE DATA (all authenticated users can read):
-- - uganda_geo_units
-- - parish_saccos
-- - village_kibiina_groups
-- - admin_settings
--
-- SERVICE ROLE (backend) has full access to everything
-- =============================================================================