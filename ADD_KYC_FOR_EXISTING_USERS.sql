-- Add KYC profiles for existing users
-- This will create pending KYC submissions for users who already exist

-- First, let's see your existing users
SELECT user_id, full_name, contact_email, phone_e164, created_at 
FROM public.users_identity 
ORDER BY created_at DESC;

-- Now add KYC profiles for existing users
-- Replace 'USER_ID_HERE' with actual user IDs from the query above

-- Example for first user (replace with actual user_id)
INSERT INTO public.kyc_profiles (
    user_id, 
    status, 
    date_of_birth, 
    gender, 
    nationality, 
    occupation_status, 
    id_type, 
    id_number,
    nok_full_name,
    nok_relationship,
    nok_phone,
    source_of_funds,
    pep_status,
    submitted_at
)
VALUES (
    'YOUR_USER_ID_HERE',  -- Replace with actual user_id from the first query
    'pending',            -- Status: pending, approved, or rejected
    '1990-01-01',        -- Date of birth
    'male',              -- Gender: male, female, other
    'Uganda',            -- Nationality
    'employed',          -- Occupation: employed, self_employed, student, unemployed
    'national_id',       -- ID type: national_id, passport, driving_license
    'CM12345678',        -- ID number
    'Jane Doe',          -- Next of kin name
    'spouse',            -- Relationship: spouse, parent, sibling, child, friend
    '+256700000001',     -- Next of kin phone
    'salary',            -- Source of funds: salary, business, savings, investment
    false,               -- PEP status (Politically Exposed Person)
    NOW()                -- Submission timestamp
)
ON CONFLICT (user_id) DO UPDATE 
SET status = 'pending', submitted_at = NOW();

-- If you want to add KYC for ALL existing users at once:
-- (This creates pending KYC for all users who don't have one yet)
INSERT INTO public.kyc_profiles (
    user_id, 
    status, 
    date_of_birth, 
    gender, 
    nationality, 
    occupation_status, 
    id_type, 
    id_number,
    nok_full_name,
    nok_relationship,
    nok_phone,
    source_of_funds,
    pep_status,
    submitted_at
)
SELECT 
    user_id,
    'pending',
    '1990-01-01',
    'male',
    'Uganda',
    'employed',
    'national_id',
    'TEST' || SUBSTRING(user_id, 1, 8),
    'Test Next of Kin',
    'spouse',
    phone_e164,
    'salary',
    false,
    NOW()
FROM public.users_identity
WHERE user_id NOT IN (SELECT user_id FROM public.kyc_profiles)
ON CONFLICT (user_id) DO NOTHING;

-- Verify the KYC profiles were created
SELECT 
    u.user_id,
    u.full_name,
    u.contact_email,
    k.status,
    k.id_type,
    k.id_number,
    k.submitted_at
FROM public.users_identity u
LEFT JOIN public.kyc_profiles k ON u.user_id = k.user_id
ORDER BY k.submitted_at DESC;
