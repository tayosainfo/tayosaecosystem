# InsForge MCP Connection And Integration Setup

This document captures the verified MCP connection state and the exact integration wiring to follow.

## InsForge as the primary platform (your intent)

This project is **InsForge-first**: the hosted InsForge project is the **system of record** for authentication, storage, and PostgreSQL (including domain tables such as `members`, `wallets`, `loans`, and the Tayosa-specific identity/onboarding tables we added via migrations). **MCP** is the right tool for day-to-day **infrastructure** work against that project (raw SQL, schema checks, buckets, and later **deployment** via InsForge’s deployment flows). **Web and mobile apps** use the **InsForge SDK** (and your **API gateway** for platform-specific routes) for **runtime** behaviour—sign-in, verification codes, storage, and database access through InsForge’s HTTP APIs where applicable.

After development, **production** should use the same InsForge project: configure environment variables from the InsForge dashboard, run migrations or schema updates through MCP or your release process, and deploy the frontend through InsForge hosting or your chosen host while keeping **API base URL + keys** pointed at this backend.

## MCP Verification Status

Verified with MCP tool calls:
- `fetch-docs` with `docType: instructions` succeeded.
- `get-backend-metadata` succeeded and returned live backend metadata.
- `fetch-sdk-docs` succeeded for:
  - `auth` + `typescript`
  - `storage` + `typescript`

Current backend metadata highlights:
- `requireEmailVerification: true`
- `verifyEmailMethod: code`
- `resetPasswordMethod: code`
- Storage bucket available: `collateral_docs`

## Project Wiring Rules (Per InsForge Docs)

- Use SDK for application logic (auth, profile, storage usage patterns).
- Use MCP tools for infrastructure tasks (schema tools, bucket lifecycle, deployment).

## Environment Configuration

Use `.env.example` as the baseline and provide real values in local env files:

- Web:
  - `VITE_INSFORGE_BASE_URL`
  - `VITE_INSFORGE_ANON_KEY`
  - `VITE_API_BASE_URL`
- Backend services:
  - `INSFORGE_BASE_URL`
  - `INSFORGE_ANON_KEY`
  - `INSFORGE_ADMIN_API_KEY`
  - `INSFORGE_STORAGE_BUCKET`
  - `DATABASE_URL` (or `CONNECTION_STRING`) — Postgres URL for `user-service` when persisting identity/onboarding. Use `sslmode=require` for InsForge’s managed DB. **Only `KEY=value` lines are valid** in `.env` (not `HOST: …`); see `.env.example`.
- Mobile:
  - pass `API_BASE_URL`, `INSFORGE_BASE_URL`, `INSFORGE_ANON_KEY` via `--dart-define`

## Integration Points Already Wired

- Web InsForge client: `src/lib/insforge.ts`
- Web gateway API client: `src/lib/platformApi.ts`
- Mobile API config: `app/mobile_app/lib/core/network/api_client.dart`
- Object storage service env wiring: `services/object-storage-service/main.go`

## Mobile Run Example

```bash
flutter run \
  --dart-define=API_BASE_URL=http://10.0.2.2:8080 \
  --dart-define=INSFORGE_BASE_URL=https://74qj9u5z.us-east.insforge.app \
  --dart-define=INSFORGE_ANON_KEY=your-insforge-anon-key
```

## InsForge MCP in Cursor (why automation may be missing)

This repo’s `.cursor` **MCP** folder may only include built-in servers (`cursor-ide-browser`, etc.). The **`user-insforge`** MCP server must be **added in Cursor Settings → MCP** (with your InsForge token) before tools like `get-backend-metadata` or SQL helpers appear. Until then, the AI cannot call InsForge APIs from chat.

**Option 3 (migrations) without MCP:** open **InsForge Dashboard → SQL**, paste and run the bundled script:

- `db/insforge_run_on_dashboard.sql`

That file merges `001`, `002`, and `003` from `db/migrations/`. Geo rows are still loaded by `user-service` from the repo CSV at runtime unless you import `uganda_geo_data_2025-11-26.csv` separately into `uganda_geo_units` (optional; the service can seed from CSV when it runs).

## Storage bucket (Option 2)

`object-storage-service` uses:

- `INSFORGE_BASE_URL`
- `INSFORGE_STORAGE_BUCKET` (defaults to `collateral_docs` if unset)

Create or confirm the bucket name in the **InsForge dashboard** (must match `INSFORGE_STORAGE_BUCKET` in `.env`). `INSFORGE_ADMIN_API_KEY` is only needed if you add admin-only automation later.

## Vercel + InsForge (deploy shape)

- **Vercel** hosts the **Vite** build. Set **build-time** env vars in the Vercel project (same names as local, `VITE_` prefix):

  - `VITE_INSFORGE_BASE_URL`
  - `VITE_INSFORGE_ANON_KEY`
  - `VITE_API_BASE_URL` → public URL of your **API gateway** (not localhost), e.g. `https://api.yourdomain.com`

- **Go services** (gateway, `user-service`, etc.) are **not** deployed by Vercel by default. Host them on **InsForge** (if you use their compute), **Fly.io**, **Railway**, **Cloud Run**, etc., and point `VITE_API_BASE_URL` at that gateway.

- **Secrets** (`DATABASE_URL`, `INSFORGE_ANON_KEY`, DB password): set on the **server** host (InsForge env / your container platform), **not** in the Vercel client unless you use serverless functions that need them (avoid putting DB URLs in the browser bundle).

## Troubleshooting: `relation "public.shares_ledger" does not exist`

Registration calls InsForge, then inserts into `users_identity`. If a **trigger** (created in the SQL editor or elsewhere) inserts into `public.shares_ledger` on user creation, that table must exist or Postgres returns this error and the API responds with a failure (often surfaced as **400** from InsForge or **500** from `user-service`).

**Fix:** run migration `db/migrations/004_shares_ledger.sql` (or the same block at the end of `db/insforge_run_on_dashboard.sql`) on your InsForge / `DATABASE_URL` database. If your trigger uses different column names, `ALTER TABLE public.shares_ledger` to match, or adjust the trigger.

### `column "shares_balance" of relation "shares_ledger" does not exist`

Your trigger expects **`shares_balance`**. Apply `db/migrations/005_shares_ledger_shares_balance.sql` (or the `005` block in `db/insforge_run_on_dashboard.sql`): it runs `ADD COLUMN IF NOT EXISTS shares_balance ...`.

### `InsForge sign-up response missing user id` (502 from `user-service`)

InsForge sometimes returns **`requireEmailVerification: true`** with **`user` omitted or without `id`**. Tayosa now:

1. Parses **`user.id`**, **`data.user`**, **`userId`**, **`user_id`**, **`sub`**, and **`user` as an array**.
2. If still empty and **`INSFORGE_ADMIN_API_KEY`** is set, calls **`GET /api/auth/users?search={email}`** to resolve the new user’s id.

Set **`INSFORGE_ADMIN_API_KEY`** in the same `.env` as `INSFORGE_ANON_KEY` (project admin key from InsForge). If the error persists, use the JSON field **`insforgeResponseKeys`** in the error body to compare with the live InsForge response.

### `insert or update on table "shares_ledger" violates foreign key constraint "shares_ledger_user_id_fkey"`

- **Same transaction / ordering:** Apply **`006_shares_ledger_deferrable_fk.sql`** (deferrable FK) and ensure `user-service` runs **`SET CONSTRAINTS ALL DEFERRED`** in the identity transaction (already done in `CreateIdentityWithOnboarding`). Prefer triggers defined as **`AFTER INSERT`** on `public.users_identity` so `NEW.user_id` is committed for child inserts.

- **Still failing:** A trigger may run in **another transaction** (e.g. InsForge `auth` user creation) before `public.users_identity` has the row—no FK can validate. Apply **`007_shares_ledger_drop_user_fk.sql`**, which **drops** `shares_ledger_user_id_fkey`. You can add the FK back later if triggers only touch `shares_ledger` after `users_identity` exists in the same DB.

## Backend Metadata Snapshot

The following table names were returned from InsForge metadata and should remain the source of truth for backend domain alignment:

- `members`
- `wallets`
- `loans`
- `loan_schedules`
- `credit_scores`
- `shares_ledger`
- `referrals`
- `village_group_pools`
- `group_members`
- `group_contributions`
- `notifications`
- `transactions`
- `audit_logs`

