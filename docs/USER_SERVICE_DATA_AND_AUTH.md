# User service: data access and InsForge auth

This document closes the “parity” notes for `user-service`: what is implemented in Go, what is delegated to InsForge HTTP APIs, and how phone-first identities relate to InsForge’s email-centric auth.

## Database: Postgres via `pgx` (not InsForge REST)

- **Canonical access:** `user-service` reads and writes platform tables (`users_identity`, `onboarding_profiles`, `uganda_geo_units`, etc.) using **`DATABASE_URL` / `CONNECTION_STRING`** and **`pgx`**.
- **InsForge `GET /api/database/records/...`** is **not** proxied by this service. That REST surface is optional for other clients; using it from the backend would duplicate logic and add latency without benefit for the current design.
- **Schema alignment:** Migrations under `db/migrations/` and bundled SQL (`db/insforge_run_on_dashboard.sql`) should stay consistent with whatever is applied to the InsForge-hosted Postgres project.

## Phone-first sign-up vs InsForge email/password

- Tayosa still collects a **Uganda-normalized phone** and stores `phone_e164` as the primary human identifier in `users_identity`.
- InsForge auth is **email-based**. The service maps:
  - **`contact_email`** when the user supplies one, or
  - a synthetic **`{digits}@tayosa.local`** (`auth_email`) when there is no contact email, so InsForge always receives an email for `POST /api/auth/users` and `POST /api/auth/sessions`.
- **`insforge_login_email`** in Postgres records which address is used for InsForge login (contact vs synthetic), so resend/verify/reset flows target the correct mailbox.
- **OAuth** users are created under InsForge’s user id; new OAuth-only accounts get a **synthetic Uganda mobile** derived deterministically from the InsForge user id so `phone_e164` remains unique without a real SIM.

## InsForge features exposed through the gateway

| Capability | Route(s) | Notes |
|------------|-----------|--------|
| `client_type` on password flows | `POST .../register?client_type=`, `POST .../login?client_type=`, `POST .../verify-email?client_type=` | Forwards to InsForge query params. Non-web clients receive `refreshToken` in JSON when InsForge returns it. |
| OAuth PKCE start | `GET /api/v1/auth/oauth/start?provider=|customKey=&redirect_uri=&code_challenge=` | Returns `{ "authUrl": "..." }` from InsForge. |
| OAuth code exchange | `POST /api/v1/auth/oauth/exchange?client_type=` | Body: `{ "code", "code_verifier" }`. Syncs user into `users_identity` when possible. |
| Refresh | `POST /api/v1/auth/refresh?client_type=` | Proxies body; forwards `X-CSRF-Token` and `Cookie` for web refresh. |
| Logout | `POST /api/v1/auth/logout` | Proxies `Authorization` and `Cookie` when present. |
| Public auth config | `GET /api/v1/auth/public-config` | Pass-through from InsForge. |
| Profile | `PATCH /api/v1/auth/profile` (Bearer) | Proxies to InsForge `PATCH /api/auth/profiles/current`, then updates local `full_name` from `profile.name` when present. Requires a real InsForge access token (not `dev-token-*`). |

Configure **allowed redirect URLs** for OAuth and email links in the **InsForge dashboard** (`allowedRedirectUrls`).
