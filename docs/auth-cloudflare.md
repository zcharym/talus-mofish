# Talus Auth — Cloudflare Magic Link Deployment

This guide deploys the email magic-link auth API used by the Talus desktop app.

Stack: **Cloudflare Workers** + **D1** + **Resend** (outbound email).

## Prerequisites

- Cloudflare account with Workers and D1 enabled
- [Wrangler CLI](https://developers.cloudflare.com/workers/wrangler/install-and-update/) installed
- Resend account with a verified sending domain
- Domain DNS managed in Cloudflare (recommended for Resend SPF/DKIM)

## 1. Create D1 database

```bash
cd workers/auth
npm install
wrangler d1 create talus-auth
```

Copy the `database_id` from the output into [`wrangler.toml`](wrangler.toml):

```toml
[[d1_databases]]
binding = "AUTH_DB"
database_name = "talus-auth"
database_id = "<your-database-id>"
```

## 2. Apply migrations

Local (for `wrangler dev`):

```bash
npm run db:migrate:local
```

Production:

```bash
npm run db:migrate:remote
```

## 3. Configure Resend

1. Add your domain in [Resend](https://resend.com/domains).
2. Add the DNS records Resend provides (SPF, DKIM) in Cloudflare DNS.
3. Create an API key with send permission.

Use a from-address on the verified domain, e.g. `login@yourdomain.com`.

## 4. Set Worker secrets and vars

Update `AUTH_PUBLIC_URL` in `wrangler.toml` to your Worker URL or custom domain:

```toml
[vars]
AUTH_APP_NAME = "Talus Agent"
AUTH_PUBLIC_URL = "https://auth.yourdomain.com"
```

Set secrets:

```bash
wrangler secret put AUTH_JWT_SECRET
wrangler secret put RESEND_API_KEY
wrangler secret put AUTH_FROM_EMAIL
```

- `AUTH_JWT_SECRET` — long random string (32+ bytes)
- `RESEND_API_KEY` — from Resend dashboard
- `AUTH_FROM_EMAIL` — e.g. `Talus Agent <login@yourdomain.com>`

## 5. Deploy

```bash
npm run deploy
```

For a custom domain, add a route in Cloudflare Workers → your worker → Triggers → Custom Domains → `auth.yourdomain.com`.

## 6. Configure the desktop app

Set the Auth server URL (no trailing slash), e.g. `https://auth.yourdomain.com`. Priority (highest first):

1. **OS env** — `TALUS_AUTH_SERVER_URL` already set in the shell
2. **Env file** — `.env.development` (default) or `.env.production` when `TALUS_ENV=production` (see `.env.*.example`)
3. **`config.json`** — `auth.authServerUrl` (path shown when hovering Save in Management → Config)

```bash
cp .env.development.example .env.development
# edit TALUS_AUTH_SERVER_URL, then:
wails3 dev
```

Then in **Agent**, use **Continue with email** to test the flow.

## Development loop

```bash
cd workers/auth
npm run dev
```

Point the desktop app at the Worker URL via `.env.development` (`TALUS_AUTH_SERVER_URL`), `auth.authServerUrl` in `config.json`, or an already-exported `TALUS_AUTH_SERVER_URL`. Apply local D1 migrations with `npm run db:migrate:local`.

Resend sandbox mode only delivers to your Resend account email until the domain is verified.

## API reference

| Method | Path                                    | Purpose                                                           |
| ------ | --------------------------------------- | ----------------------------------------------------------------- |
| `POST` | `/v1/auth/magic-link`                   | Start login `{ "email": "..." }` → `{ "requestId", "expiresAt" }` |
| `GET`  | `/v1/auth/magic-link/status/:requestId` | Poll until verified                                               |
| `GET`  | `/v1/auth/verify?token=...&request=...` | Browser link from email                                           |
| `GET`  | `/v1/auth/me`                           | Validate JWT (`Authorization: Bearer ...`)                        |
| `POST` | `/v1/auth/signout`                      | No-op for v1 (client clears local session)                        |

## Troubleshooting

| Issue                       | Check                                                                                                              |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| Email not received          | Resend dashboard logs; domain verified; `AUTH_FROM_EMAIL` matches domain                                           |
| Link invalid                | `AUTH_PUBLIC_URL` matches deployed URL exactly                                                                     |
| Poll never verifies         | User clicked link; D1 migration applied; check Worker logs                                                         |
| Desktop cannot reach server | `.env.development` / `TALUS_AUTH_SERVER_URL` / `auth.authServerUrl`; firewall/system proxy for desktop HTTP client |
