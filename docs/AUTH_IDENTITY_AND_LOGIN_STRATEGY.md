# Auth, Identity, And Login Strategy

This document defines how authentication and identity work across web and mobile, with phone-first UX and optional email login.

## Goals

- Allow users to sign in with either phone number or email.
- Keep one credential set usable across both apps.
- Keep InsForge-compatible auth implementation without exposing technical internals to users.
- Avoid login confusion that can damage trust and adoption.

## Core Principle

Every member has:
- a technical auth identifier used internally (`auth_email`)
- an optional user-facing email for contact (`contact_email`)
- a primary normalized phone identifier (`phone_e164`)

Users never see `auth_email`.

## Identity Data Model (User Service)

Recommended fields:
- `user_id` (UUID)
- `phone_e164` (required, unique, indexed)
- `auth_email` (required, unique, internal only)
- `contact_email` (optional, unique when present)
- `insforge_user_id` (required, unique)
- `phone_verified_at` (nullable timestamp)
- `contact_email_verified_at` (nullable timestamp)
- `status` (active/suspended/pending_kyc)
- `created_at`, `updated_at`

## Input Normalization Rules

### Phone normalization

- Convert to E.164 format (example Uganda: `+2567XXXXXXXX`).
- Remove spaces/dashes before parse.
- Reject invalid length/prefix early.
- Store and compare only normalized `phone_e164`.

### Email normalization

- Trim and lowercase before store/compare.
- Validate format at API boundary.
- Enforce uniqueness on normalized value.

## Registration Flow

1. Client sends phase-1 signup payload with phone, password/PIN, optional email.
2. Backend normalizes phone and optional email.
3. Backend creates deterministic internal `auth_email` from phone (example pattern: `<normalized_phone_digits>@tayosa.local`).
4. Backend creates InsForge account using `auth_email` + password/PIN.
5. Backend stores:
   - `phone_e164`
   - `auth_email`
   - `contact_email` (if provided)
   - `insforge_user_id`
6. OTP verification proceeds as configured by auth policy.

### Onboarding data additions

- During phase-2 and phase-3 onboarding, address fields must use Uganda geo master data from `uganda_geo_data_2025-11-26.csv`.
- Store geographic IDs (or normalized names) for: district, county, sub-county, parish, village.
- Capture optional `referralCode` during signup/membership setup to support affiliate rewards.
- Kibiina profile must include parish and village references to support village-level group creation.

## Login Flow (Dual Identifier)

Single login field label: `Phone number or email`.

Backend decision:
- If identifier looks like phone:
  - normalize to `phone_e164`
  - resolve account by `phone_e164`
  - authenticate via InsForge using stored `auth_email`
- If identifier looks like email:
  - resolve by `contact_email` (preferred)
  - fallback match by `auth_email` only for internal/admin cases
  - authenticate via InsForge using stored `auth_email`

On success, return normal app session/token response.

## UX Requirements (Non-Negotiable)

- Login placeholder must read: `Phone number or email`.
- Register helper text should say: `You can log in with either your phone number or email.`
- Never expose internal email format in UI, logs, or user-facing APIs.
- Use generic auth failure text: `Invalid credentials. Use your registered phone number or email.`

## Security And Compliance Controls

- Rate limit login and OTP endpoints at API Gateway.
- Lockout/risk-score repeated failures per IP + identifier fingerprint.
- Do not leak whether phone/email exists (anti-enumeration).
- Verify phone ownership before enabling transactions.
- Require KYC completion before high-risk actions.

## API Contract Guidance

### Register request (phase 1)

```json
{
  "fullName": "Jane Auma",
  "phone": "+2567XXXXXXXX",
  "email": "jane@example.com",
  "password": "******",
  "dateOfBirth": "1998-01-01",
  "nationality": "UG"
}
```

### Membership/onboarding extension request (phase 3/4)

```json
{
  "membershipType": "individual",
  "district": "Abim",
  "county": "Labwor County",
  "subCounty": "Abim",
  "parish": "Abongepach",
  "village": "VillageAbongepach",
  "referralCode": "FRIEND123",
  "kibiina": {
    "cycleFrequency": "weekly",
    "contributionAmount": 20000,
    "payoutMethod": "mobile_money"
  }
}
```

### Login request

```json
{
  "identifier": "+2567XXXXXXXX",
  "password": "******"
}
```

or

```json
{
  "identifier": "jane@example.com",
  "password": "******"
}
```

## Cross-App Consistency

Web (`src/`) and Mobile (`app/mobile_app/`) must both:
- call the same auth/login endpoints
- send the same `identifier` field
- share the same backend identity store
- use same token validation gateway rules

This guarantees one credential set works on both apps.

## Rollout Checklist

- [ ] Add DB columns and unique indexes (`phone_e164`, `auth_email`, `contact_email`).
- [ ] Implement phone/email normalization utilities.
- [ ] Implement dual-identifier login endpoint.
- [ ] Load and version Uganda geo master data from `uganda_geo_data_2025-11-26.csv`.
- [ ] Build cascading address selectors (`District -> County -> Sub-County -> Parish -> Village`).
- [ ] Add affiliate referral capture and reward trigger integration.
- [ ] Enforce parish-level SACCO and village-level Kibiina mapping constraints.
- [ ] Update web login UI copy and validation.
- [ ] Update mobile login UI copy and validation.
- [ ] Add audit logging for registration/login events.
- [ ] Add integration tests for both phone and email login.
- [ ] Add migration script for existing users.

