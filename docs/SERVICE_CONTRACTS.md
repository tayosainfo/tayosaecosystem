# Service Contracts Baseline

This file locks the initial cross-service contracts for the SACCO rollout.

## API Gateway

- Base URL: `http://localhost:8080`
- Auth pass-through:
  - `POST /api/v1/auth/register` (optional `?client_type=web|mobile|desktop|server`)
  - `POST /api/v1/auth/login` (optional `client_type`)
  - `GET /api/v1/auth/oauth/start?provider=...&redirect_uri=...&code_challenge=...` (or `customKey=` instead of `provider`)
  - `POST /api/v1/auth/oauth/exchange?client_type=`
  - `POST /api/v1/auth/refresh?client_type=`
  - `POST /api/v1/auth/logout`
  - `GET /api/v1/auth/public-config`
  - `PATCH /api/v1/auth/profile` (Bearer)
- User/onboarding:
  - `POST /api/v1/onboarding/phase` (Bearer token required)
  - `GET /api/v1/geo?level={district|county|subcounty|parish|village}&parent={name}`
  - `GET /api/v1/groups/policy`
- Domain integrations:
  - `POST /api/v1/affiliate/referrals`
  - `POST /api/v1/notifications/send`
  - `GET|POST /api/v1/audit/events`
  - `POST /api/v1/loans/score`
  - `GET /api/v1/storage/upload-url`

## User Service

- `POST /api/v1/auth/register` creates identity with:
  - phone normalization (`phone_e164`)
  - internal auth key (`auth_email`)
  - optional `contact_email`
  - optional `dateOfBirth` (`YYYY-MM-DD`) and optional `nationality` (trimmed, max 64 characters), persisted when `DATABASE_URL` is set (see migration `003_user_dob_nationality.sql`)
- `POST /api/v1/auth/login` supports `identifier` as phone or email. The `user` object uses the same public profile fields as register and `GET /api/v1/users/me` (including `contactEmailVerified`, `dateOfBirth`, `nationality`, and `insforgeUserId` when present). With InsForge, `session` may include `refreshToken` and `csrfToken` when returned for the requested `client_type`.
- OAuth PKCE: `GET /api/v1/auth/oauth/start` → open `authUrl`; callback receives `insforge_code`; `POST /api/v1/auth/oauth/exchange` with `{ "code", "code_verifier" }` completes sign-in. InsForge user id is synced into `users_identity` (new OAuth users get a synthetic `phone_e164` until a real phone is collected in onboarding).
- `POST /api/v1/auth/refresh` and `POST /api/v1/auth/logout` proxy InsForge with optional `X-CSRF-Token` / cookies for web clients.
- `PATCH /api/v1/auth/profile` updates InsForge profile and mirrors `name` into local `full_name` when present.
- Platform data remains **`pgx` → Postgres**; InsForge database REST is not required for parity (see `docs/USER_SERVICE_DATA_AND_AUTH.md`).
- `POST /api/v1/onboarding/phase` stores phase 2-4 payload. **`geo` must include all of** `district`, `county`, `sub_county`, `parish`, `village` **and** the combination must match a row in `uganda_geo_units` (same hierarchy as `GET /api/v1/geo`).
- `GET /api/v1/geo` serves cascading location data from Uganda CSV preload.
- `GET /api/v1/groups/policy` returns policy text plus **`parishSaccosRegistered`** and **`villageKibiinaGroupsRegistered`** counts from Postgres (`parish_saccos`, `village_kibiina_groups`); in-memory dev store these counts are always `0`.

## Affiliate Service

- `POST /api/v1/affiliate/referrals` creates pending reward event.
- `GET /api/v1/affiliate/rewards` lists reward records.

## Notification Service

- `POST /api/v1/notifications/send` queues notification jobs.
- `GET /api/v1/notifications/outbox` lists queued jobs.

## Audit Service

- `POST /api/v1/audit/events` append-only event write.
- `GET /api/v1/audit/events` event list read.

## Loan/Credit Service

- `POST /api/v1/loans/score` returns initial score and max eligible amount.

## Object Storage Service

- `GET` or `POST /api/v1/storage/upload-url` — server calls InsForge `POST /api/storage/buckets/{bucket}/upload-strategy` and **returns the upstream JSON body** (`method`, `uploadUrl`, `fields`, `key`, `confirmRequired`, `confirmUrl`, `expiresAt`, etc.). Query params: `fileName`, `category` (default `kyc`), optional `contentType`, optional `size`. For `POST`, the same fields may be sent as JSON (`application/json`). Requires `Authorization: Bearer` (user JWT) or `INSFORGE_ANON_KEY` on the service when calling the storage service directly.
- `GET /api/v1/storage/object`

