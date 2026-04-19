# TAYOSA Banking Platform

SACCO-oriented banking stack for Uganda: **Go microservices** behind an **API gateway**, **React + Vite** web app, **Flutter** mobile app, and **InsForge** for hosted auth, database, and storage integration.

## Repository layout

| Path | Role |
|------|------|
| `services/*` | Go services (each folder has its own `go.mod`) |
| `src/` | Web client (Vite + React + TypeScript) |
| `app/mobile_app/` | Flutter iOS/Android app |
| `db/migrations/` | SQL migrations for platform Postgres (e.g. geo + identity) |
| `docs/` | Runbooks, contracts, and architecture notes for this stack |

## Go services (current)

- `api-gateway-service` — single entry for clients; proxies to downstream services
- `user-service` — identity onboarding, geo APIs, InsForge-backed auth flows
- `affiliate-service` — referrals and rewards
- `notification-service` — outbound notifications
- `audit-log-service` — append-only audit events
- `loan-credit-service` — loan scoring placeholder/API surface
- `object-storage-service` — proxies InsForge upload-strategy for uploads; object reads placeholder
- `kibiina-service` — group / merry-go-round flows
- `fee-service` — fee-related API surface

## Prerequisites

- **Go** 1.22+ (see each service’s `go.mod` for exact toolchain)
- **Node.js** 18+ and npm (web app)
- **Flutter** 3.24+ (mobile app)
- **Docker** (optional local Postgres via `docker-compose.yml`)

## Quick start

1. **Environment** — Copy `.env.example` to `.env` at the repo root and fill values (`INSFORGE_BASE_URL`, `INSFORGE_ANON_KEY`, optional `DATABASE_URL`, `VITE_API_BASE_URL`, etc.). See `docs/INSFORGE_CONNECTION_SETUP.md`. For **InsForge Postgres**, use the dashboard connection string as `DATABASE_URL=postgresql://…?sslmode=require` (only works while the InsForge project is **active**, not paused). If direct DB access still fails from your network, use **local Docker Postgres** or omit `DATABASE_URL` / set `USER_SERVICE_ALLOW_MEMORY_FALLBACK=1`.

2. **Database** — Apply SQL under `db/migrations/` to your Postgres instance (or follow `docs/RUNBOOK_LOCAL_DEVELOPMENT.md`).

3. **Uganda geo CSV** — Keep `uganda_geo_data_2025-11-26.csv` at the **repo root**. `user-service` loads it when seeding `uganda_geo_units`.

4. **Backend** — From each service directory: `go run .` (ports and env loading are described in `docs/RUNBOOK_LOCAL_DEVELOPMENT.md`). Compile all services:

   ```bash
   ./scripts/build-backend.sh
   ```

   On Windows: `.\scripts\build-backend.ps1` or `npm run backend:build`.

5. **Web** — `npm install` then `npm run dev` (default dev server: Vite).

6. **Mobile** — `cd app/mobile_app` → `flutter pub get` → `flutter run`.

7. **Local Postgres only** — `docker compose up -d postgres` (see `docker-compose.yml`).

## Documentation

- `docs/RUNBOOK_LOCAL_DEVELOPMENT.md` — start order, smoke tests, troubleshooting
- `docs/SERVICE_CONTRACTS.md` — gateway and service HTTP contracts
- `docs/ARCHITECTURE_WORKFLOW.md` — target workflows (onboarding, Kibiina)
- `docs/AUTH_IDENTITY_AND_LOGIN_STRATEGY.md` — auth and identity model
- `docs/INSFORGE_CONNECTION_SETUP.md` — InsForge wiring and MCP checks

## License

Proprietary software owned by TAYOSA GROUP (U) LTD.

## Support

- Email: tech@tayosa.ug
