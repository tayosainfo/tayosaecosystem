-- Backfill referral codes + onboarding data into vertical-slice tables.

INSERT INTO user_referral_codes (user_id, referral_code)
SELECT u.user_id,
       'TAY-' || UPPER(SUBSTRING(MD5('tayosa:' || u.user_id) FROM 1 FOR 8))
FROM users_identity u
LEFT JOIN user_referral_codes rc ON rc.user_id = u.user_id
WHERE rc.user_id IS NULL;

INSERT INTO sacco_memberships (
  user_id, status, district, county, sub_county, parish, village, updated_at
)
SELECT o.user_id,
       CASE WHEN o.parish IS NOT NULL AND o.village IS NOT NULL THEN 'enrolled' ELSE 'not_started' END,
       o.district, o.county, o.sub_county, o.parish, o.village,
       COALESCE(o.updated_at, now())
FROM onboarding_profiles o
LEFT JOIN sacco_memberships s ON s.user_id = o.user_id
WHERE s.user_id IS NULL;

INSERT INTO kyc_profiles (user_id, status, updated_at)
SELECT o.user_id,
       CASE WHEN (o.phase_payload ? 'kyc') THEN 'pending' ELSE 'not_started' END,
       COALESCE(o.updated_at, now())
FROM onboarding_profiles o
LEFT JOIN kyc_profiles k ON k.user_id = o.user_id
WHERE k.user_id IS NULL;

