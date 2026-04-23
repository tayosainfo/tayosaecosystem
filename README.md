# Tayosa Ecosystem

**Tayosa** is a SACCO-oriented banking and savings platform for Uganda. This monorepo hosts the **Tayosa Ecosystem**: a **Go** microservice backend behind an **API gateway**, a **React + Vite + TypeScript** web client, a **Flutter** mobile app, SQL **migrations** for platform identity and ledgers, and integration with **[Supabase](https://supabase.com)** for hosted authentication, PostgreSQL database, and object storage.

**Upstream repository:** [github.com/tayosainfo/tayosaecosystem](https://github.com/tayosainfo/tayosaecosystem)

---

## Contents

- [Supabase setup](#supabase-setup)
- [Architecture](#architecture)
- [Repository layout](#repository-layout)
- [Backend services](#backend-services)
- [Web application](#web-application)
- [Mobile application](#mobile-application)
- [Data and migrations](#data-and-migrations)
- [Environment configuration](#environment-configuration)
- [Prerequisites](#prerequisites)
- [Quick start](#quick-start)
- [Authentication and email verification](#authentication-and-email-verification)
- [API surface (gateway)](#api-surface-gateway)
- [Scripts](#scripts)
- [Documentation](#documentation)
- [License and support](#license-and-support)

---

## Supabase setup

The Tayosa ecosystem uses **Supabase** as its Backend-as-a-Service provider for authentication, PostgreSQL database, and object storage. Follow these steps to set up your Supabase project and configure the system.

### 1. Create Supabase project

1. **Sign up** at [supabase.com](https://supabase.com) if you don't have an account
2. **Create a new project**:
   - Choose your organization
   - Enter project name (e.g., "tayosa-banking")
   - Set a strong database password
   - Select your preferred region
3. **Wait for project initialization** (usually takes 1-2 minutes)

### 2. Obtain Supabase credentials

Once your project is ready, navigate to **Project Settings > API** to find:

#### Project URL
- **Format**: `https://[project-ref].supabase.co`
- **Example**: `https://ablvrbnbsdqshrorhmjf.supabase.co`
- **Usage**: Used by both frontend and backend services

#### API Keys
- **Anon/Public Key**: Safe to use in frontend code, provides read access with Row Level Security
- **Service Role Key**: ⚠️ **Secret key** - never expose in frontend, used for admin operations

### 3. Configure environment variables

1. **Copy** `.env.example` to `.env` at the repository root
2. **Update** the following Supabase variables in your `.env` file:

```bash
# Frontend Supabase Configuration
VITE_SUPABASE_URL=https://your-project-ref.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key-here

# Backend Supabase Configuration  
SUPABASE_URL=https://your-project-ref.supabase.co
SUPABASE_ANON_KEY=your-anon-key-here
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key-here

# Database Connection (Supabase PostgreSQL)
DATABASE_URL=postgresql://postgres:YOUR_PASSWORD@db.your-project-ref.supabase.co:5432/postgres?sslmode=require
```

**Replace**:
- `your-project-ref` with your actual Supabase project reference
- `your-anon-key-here` with your actual anon key
- `your-service-role-key-here` with your actual service role key
- `YOUR_PASSWORD` with your database password

### 4. Configure Supabase authentication

#### Enable email confirmations
1. Go to **Authentication > Settings** in your Supabase dashboard
2. **Enable** "Enable email confirmations"
3. **Configure** redirect URLs for email verification:
   - Add your frontend URL (e.g., `http://localhost:5173/verify`)
   - Add production URLs when deploying

#### Customize email templates (optional)
1. Navigate to **Authentication > Email Templates**
2. **Customize** the "Confirm signup" template
3. **Customize** the "Reset password" template
4. **Test** email delivery using the "Send test email" feature

#### Configure OAuth providers (optional)
1. Go to **Authentication > Providers**
2. **Enable** desired providers (Google, GitHub, etc.)
3. **Configure** OAuth credentials for each provider
4. **Set** redirect URLs for OAuth flows

### 5. Set up database schema

#### Apply migrations
1. **Navigate** to your Supabase project dashboard
2. **Go to** SQL Editor
3. **Run** each migration file from `db/migrations/` in order:
   ```sql
   -- Copy and paste content from db/migrations/001_*.sql
   -- Then 002_*.sql, 003_*.sql, etc.
   ```

#### Alternative: Use local migration tools
```bash
# Using psql with Supabase connection string
psql "postgresql://postgres:PASSWORD@db.your-project-ref.supabase.co:5432/postgres?sslmode=require" -f db/migrations/001_initial_schema.sql
```

### 6. Configure Row Level Security (RLS)

Supabase enables RLS by default. For the Tayosa ecosystem:

1. **Most operations** use the service role key, which bypasses RLS
2. **Client-side operations** use the anon key with RLS policies
3. **Custom policies** can be added in the Supabase dashboard under Authentication > Policies

### 7. Mobile app configuration

For the Flutter mobile app, configure Supabase using build-time variables:

```bash
# Development build
flutter run \
  --dart-define=SUPABASE_URL=https://your-project-ref.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-anon-key-here \
  --dart-define=API_BASE_URL=http://localhost:8080

# Production build
flutter build apk \
  --dart-define=SUPABASE_URL=https://your-project-ref.supabase.co \
  --dart-define=SUPABASE_ANON_KEY=your-anon-key-here \
  --dart-define=API_BASE_URL=https://your-production-api.com
```

### 8. Verify setup

#### Test authentication flow
1. **Start** the backend services (see [Quick start](#quick-start))
2. **Start** the web application (`npm run dev`)
3. **Register** a new user account
4. **Check** your email for verification link
5. **Verify** email and complete login

#### Test database connection
```bash
# Test user service database connection
curl http://localhost:8081/health

# Test geo data endpoint
curl http://localhost:8080/api/v1/geo?level=district
```

### Troubleshooting

#### Common issues and solutions

**Email verification not working**
- ✅ Check that "Enable email confirmations" is enabled in Supabase
- ✅ Verify SMTP configuration in Authentication > Settings
- ✅ Check spam folder for verification emails
- ✅ Ensure redirect URLs are correctly configured

**Database connection failed**
- ✅ Verify `DATABASE_URL` format and credentials
- ✅ Check that your IP is allowed (Supabase allows all by default)
- ✅ Ensure SSL mode is set to `require` for Supabase connections
- ✅ Test connection using `psql` or database client

**Authentication errors**
- ✅ Verify `SUPABASE_URL` and `SUPABASE_ANON_KEY` are correct
- ✅ Check that service role key is not exposed in frontend code
- ✅ Ensure environment variables are loaded correctly
- ✅ Check browser network tab for 401/403 errors

**Mobile app build issues**
- ✅ Verify `--dart-define` flags are correctly formatted
- ✅ Check that Supabase Flutter SDK is properly installed
- ✅ Ensure API base URL is reachable from mobile device
- ✅ Test on both iOS and Android platforms

**CORS errors**
- ✅ Configure allowed origins in Supabase dashboard
- ✅ Check that frontend URL is included in allowed origins
- ✅ Verify API gateway CORS configuration

#### Getting help

- **Supabase Documentation**: [docs.supabase.com](https://docs.supabase.com)
- **Supabase Community**: [github.com/supabase/supabase/discussions](https://github.com/supabase/supabase/discussions)
- **Tayosa Support**: [tech@tayosa.ug](mailto:tech@tayosa.ug)

---

## Architecture

Clients talk to **`api-gateway-service`** (default port **8080**). The gateway forwards requests to small **Go** services by path prefix. **`user-service`** owns identity, onboarding phases, Uganda geo lookups, and Supabase-backed auth (register, login, email verification, password reset, OAuth helpers, profile).

```mermaid
flowchart LR
  subgraph clients [Clients]
    Web[Vite React]
    Mobile[Flutter]
  end
  GW[api-gateway :8080]
  US[user-service :8081]
  ST[object-storage :8015]
  AF[affiliate :8016]
  NT[notification :8010]
  AU[audit :8014]
  LN[loan-credit :8013]
  FE[fee :8004]
  KB[kibiina :8086]
  SB[(Supabase)]
  PG[(PostgreSQL)]

  Web --> GW
  Mobile --> GW
  GW --> US
  GW --> ST
  GW --> AF
  GW --> NT
  GW --> AU
  GW --> LN
  GW --> FE
  GW --> KB
  US --> SB
  US --> PG
  ST --> SB
```

---

## Repository layout

| Path | Description |
|------|-------------|
| [`services/`](services/) | Independent Go modules per service (`go.mod` in each folder). |
| [`src/`](src/) | Web SPA: auth pages, `AuthContext`, `platformApi` client to the gateway. |
| [`app/mobile_app/`](app/mobile_app/) | Flutter app (Android, iOS, web targets). |
| [`db/migrations/`](db/migrations/) | Ordered SQL migrations for `users_identity`, onboarding, shares ledger, etc. |
| [`db/supabase_run_on_dashboard.sql`](db/supabase_run_on_dashboard.sql) | Optional SQL to run from Supabase SQL editor when applicable. |
| [`docs/`](docs/) | Runbooks, HTTP contracts, auth strategy, Supabase wiring. |
| [`scripts/`](scripts/) | `build-backend.sh` / `build-backend.ps1` to compile all Go services. |
| [`docker-compose.yml`](docker-compose.yml) | Local **Postgres 15** for development (optional). |
| [`.env.example`](.env.example) | Template for root `.env` (copy to `.env`; **never commit** `.env`). |
| [`uganda_geo_data_2025-11-26.csv`](uganda_geo_data_2025-11-26.csv) | Geo hierarchy seed data consumed by `user-service`. |

---

## Backend services

| Service | Default port | Role |
|---------|--------------|------|
| **api-gateway-service** | 8080 | Single HTTP entry; CORS; proxies to services below. |
| **user-service** | 8081 | Auth, `users/me`, onboarding phase, `/geo`, group policy; Supabase + Postgres. |
| **object-storage-service** | 8015 | Upload strategy / storage proxy toward Supabase. |
| **affiliate-service** | 8016 | Referrals and rewards. |
| **notification-service** | 8010 | Outbound notifications. |
| **audit-log-service** | 8014 | Append-only audit API. |
| **loan-credit-service** | 8013 | Loan scoring placeholder / API surface. |
| **fee-service** | 8004 | Fee-related API (`PORT` overridable). |
| **kibiina-service** | 8086 | Group / merry-go-round flows (`PORT` overridable). |

Each service is started with `go run .` from its directory (or run compiled binaries). **`user-service`** and **`api-gateway-service`** load the **repo root** `.env` via `godotenv` (see runbook).

---

## Web application

- **Stack:** React 18, TypeScript, Vite 6, Tailwind CSS, `@supabase/supabase-js` (for Supabase integration), `lucide-react`.
- **Dev server:** `npm run dev` (default Vite port **5173**).
- **Production build:** `npm run build` → output in `dist/`.
- **Env:** `VITE_API_BASE_URL` (gateway, default `http://localhost:8080`), `VITE_SUPABASE_URL` / `VITE_SUPABASE_ANON_KEY` for Supabase authentication and database access.

**Routes (path-based in [`src/App.tsx`](src/App.tsx)):** login (`/`), email verification (`/verify`), registration (`/register`), forgot/reset password. Authenticated shell is minimal while the product UI is rebuilt around Supabase-backed flows.

---

## Mobile application

- **Path:** [`app/mobile_app/`](app/mobile_app/)
- **Commands:** `flutter pub get` then `flutter run` (device or emulator).
- **Networking:** Base URL for the gateway is configured in [`lib/core/network/api_client.dart`](app/mobile_app/lib/core/network/api_client.dart) (local defaults documented in the runbook).

---

## Data and migrations

1. Choose **Postgres**: local Docker (`docker compose up -d postgres`) or Supabase-hosted DB (connection string from your Supabase project dashboard).
2. Apply migrations in order under [`db/migrations/`](db/migrations/) (`001` … `010` and any newer files).
3. Keep the **Uganda geo CSV** at the **repository root**; `user-service` uses it when seeding geo units.

Details: [`docs/RUNBOOK_LOCAL_DEVELOPMENT.md`](docs/RUNBOOK_LOCAL_DEVELOPMENT.md), [`docs/USER_SERVICE_DATA_AND_AUTH.md`](docs/USER_SERVICE_DATA_AND_AUTH.md).

---

## Environment configuration

1. Copy **`.env.example`** to **`.env`** at the repo root.
2. Set at least:
   - **`SUPABASE_URL`** / **`SUPABASE_ANON_KEY`** / **`SUPABASE_SERVICE_ROLE_KEY`** — Supabase authentication and database access (see [Supabase setup](#supabase-setup) for details).
   - **`DATABASE_URL`** — Supabase PostgreSQL connection string for identities and onboarding; optional **`USER_SERVICE_ALLOW_MEMORY_FALLBACK=1`** if you must run without DB.
   - Gateway service URLs (**`USER_SERVICE_URL`**, etc.) if not using localhost defaults.
3. For the **Vite** app, mirror Supabase URL/anon key: **`VITE_SUPABASE_URL`**, **`VITE_SUPABASE_ANON_KEY`**, **`VITE_API_BASE_URL`**.


**Security:** `.env` is **gitignored**. Do not commit secrets, database passwords, or GitHub tokens. Use your platform’s secret store for CI/CD.

---

## Prerequisites

| Tool | Notes |
|------|--------|
| **Go** | 1.22+ (see each service’s `go.mod`). |
| **Node.js** | 18+ and npm (web). |
| **Flutter** | 3.24+ (mobile). |
| **Docker** | Optional; for local Postgres only. |

---

## Quick start

1. **Clone** this repository and `cd` into it.
2. **`.env`** — Copy `.env.example` → `.env` and fill values (see above).
3. **Database** — Start Postgres if needed (`docker compose up -d postgres`), then apply [`db/migrations/`](db/migrations/).
4. **Backend** — From repo root, either:
   - `npm run backend:build` (Windows PowerShell), or  
   - `./scripts/build-backend.sh` (Unix), then run each service binary **or** `go run .` per service (see runbook for order; start **user-service** then **api-gateway** for auth smoke tests).
5. **Web** — `npm install` then `npm run dev`.
6. **Mobile** — `cd app/mobile_app` → `flutter pub get` → `flutter run`.

**Smoke test ideas:** register → verify email (if Supabase email is configured) → login → `GET /api/v1/geo?level=district` → authenticated onboarding phase update. See [`docs/RUNBOOK_LOCAL_DEVELOPMENT.md`](docs/RUNBOOK_LOCAL_DEVELOPMENT.md) §4.

---

## Authentication and email verification

- With **Supabase** enabled, `user-service` delegates sign-up, sessions, verification, and password reset to InsForge’s REST API (same patterns as `@supabase/supabase-js`).
- **Email verification:** Outbound mail must be configured in the **Supabase project** (SMTP / provider). Otherwise verification resend/create-token calls will fail on InsForge’s side.
- **Pending local profile:** If sign-up completes in Supabase before a full Tayosa row exists, the web flow can send users to **`/verify`** with profile completion; see [`docs/AUTH_IDENTITY_AND_LOGIN_STRATEGY.md`](docs/AUTH_IDENTITY_AND_LOGIN_STRATEGY.md).

---

## API surface (gateway)

All paths below are relative to the gateway (e.g. `http://localhost:8080`). **`public: false`** routes require `Authorization: Bearer <access_token>`.

**Public (auth / geo):**  
`/api/v1/auth/register`, `login`, `resend-verification-email`, `verify-email`, `send-reset-password-email`, `exchange-reset-password-token`, `reset-password`, `oauth/start`, `oauth/exchange`, `refresh`, `logout`, `public-config`, **`GET`** `/api/v1/geo`, `/api/v1/groups/policy`.

**Authenticated:**  
`/api/v1/auth/profile` (**PATCH**), `/api/v1/users/me`, `/api/v1/onboarding/phase`, storage, affiliate, notifications, audit, loans, fees, kibiina — see [`docs/SERVICE_CONTRACTS.md`](docs/SERVICE_CONTRACTS.md) for shapes and methods.

---

## Scripts

| Script | Purpose |
|--------|---------|
| [`scripts/build-backend.ps1`](scripts/build-backend.ps1) | Build all Go services on Windows (`npm run backend:build`). |
| [`scripts/build-backend.sh`](scripts/build-backend.sh) | Same on macOS/Linux. |
| [`test_unified_auth_flow.js`](test_unified_auth_flow.js), [`test_unverified_login.js`](test_unverified_login.js) | Ad-hoc Node scripts for local auth experiments (run with Node as needed). |

---

## Documentation

| Document | Topic |
|----------|--------|
| [`docs/RUNBOOK_LOCAL_DEVELOPMENT.md`](docs/RUNBOOK_LOCAL_DEVELOPMENT.md) | Ports, start order, smoke tests, troubleshooting. |
| [`docs/SERVICE_CONTRACTS.md`](docs/SERVICE_CONTRACTS.md) | Gateway routing and HTTP contracts. |
| [`docs/ARCHITECTURE_WORKFLOW.md`](docs/ARCHITECTURE_WORKFLOW.md) | Target workflows (onboarding, Kibiina). |
| [`docs/AUTH_IDENTITY_AND_LOGIN_STRATEGY.md`](docs/AUTH_IDENTITY_AND_LOGIN_STRATEGY.md) | Identity model and auth strategy. |
| [`docs/SUPABASE_CONNECTION_SETUP.md`](docs/SUPABASE_CONNECTION_SETUP.md) | Supabase env vars, admin key, setup guide. |
| [`docs/USER_SERVICE_DATA_AND_AUTH.md`](docs/USER_SERVICE_DATA_AND_AUTH.md) | User-service data model and auth behavior. |
| [`AGENTS.md`](AGENTS.md) | Notes for AI coding agents working in this repo. |

---

## License and support

- **License:** Proprietary software owned by **TAYOSA GROUP (U) LTD.**
- **Support:** [tech@tayosa.ug](mailto:tech@tayosa.ug)

---

*README version aligned with the Tayosa Ecosystem monorepo; for the latest file tree and scripts, browse [the repository on GitHub](https://github.com/tayosainfo/tayosaecosystem).*
