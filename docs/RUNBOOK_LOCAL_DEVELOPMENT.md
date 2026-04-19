# Local Development Runbook

## 1) Start Core Services

Run each service from its folder:

- `services/user-service` on `:8081`
- `services/api-gateway-service` on `:8080`
- `services/object-storage-service` on `:8015`
- `services/notification-service` on `:8010`
- `services/audit-log-service` on `:8014`
- `services/loan-credit-service` on `:8013`
- `services/affiliate-service` on `:8016`
- `services/fee-service` on `:8004` (override with `PORT`)
- `services/kibiina-service` on `:8086` (override with `PORT`)

Quick compile check before startup:

- PowerShell: `./scripts/build-backend.ps1`
- Bash: `./scripts/build-backend.sh`
- npm shortcut: `npm run backend:build`

### Environment file

Keep a **repo root** `.env` (gitignored) with at least `DATABASE_URL` (for `user-service` Postgres), `INSFORGE_BASE_URL`, `INSFORGE_ANON_KEY`, and optional service URLs. Both **`user-service`** and **`api-gateway-service`** auto-load `../../.env` then `./.env` when you run `go run .` from their service folders, so you do not need to manually `export` variables in PowerShell for normal local runs.

## 2) Start Web App

From repo root:

- `npm run dev`

Optional env:

- `VITE_API_BASE_URL=http://localhost:8080`

## 3) Start Mobile App

From `app/mobile_app`:

- `flutter pub get`
- `flutter run`

The mobile app already points to gateway local defaults in `lib/core/network/api_client.dart`.

## 4) Smoke Test Sequence

1. `POST /api/v1/auth/register`
2. `POST /api/v1/auth/login`
3. `GET /api/v1/geo?level=district`
4. `POST /api/v1/onboarding/phase`
5. `POST /api/v1/affiliate/referrals`
6. `POST /api/v1/notifications/send`
7. `POST /api/v1/audit/events`

## 5) Failure Handling

- 401 from gateway: verify `Authorization: Bearer <token>`.
- 502 from gateway: verify downstream service process is running.
- Geo empty result: ensure CSV exists at repo root and user-service restarted.

