# Echo Watch

Local screen watcher for VDI build alerts, delivered to your phone via **Safari Web Push** and a **Cloudflare Worker**.

- **Runs on your physical PC** — never inside the isolated VDI guest
- Watches the **local VDI client window** (`mstsc.exe`, Horizon, Citrix, etc.)
- Skips capture when the client is **minimized** or not visible
- Sends **minimal JSON alerts** to Cloudflare; push is delivered to an iOS PWA on the Home Screen

## Architecture

```
VDI guest UI  →  local VDI client window  →  echo-watch (OCR + rules)
                                              ↓
                                    Cloudflare Worker
                                              ↓
                              iPhone PWA (web.push.apple.com)
```

## 1. Generate VAPID keys

```bash
npx @pushforge/builder vapid
```

Or with OpenSSL / any ES256 P-256 VAPID generator. You need:

- `VAPID_PUBLIC_KEY` (URL-safe base64)
- `VAPID_PRIVATE_KEY` (URL-safe base64)
- `VAPID_SUBJECT` — `mailto:you@example.com` or `https://your-worker-url`

## 2. Cloudflare setup

```bash
cd cloud/echo-watch/worker
npm install

# Create KV namespace
npx wrangler kv namespace create ECHO_WATCH_KV
npx wrangler kv namespace create ECHO_WATCH_KV --preview
```

Update `wrangler.toml` with the KV `id` and `preview_id` from the command output.

> **Note:** `local-dev-echo-watch-kv` placeholders only work with `wrangler dev --local`. Production deploy requires real namespace IDs from `wrangler kv namespace create`.

Set secrets:

```bash
npx wrangler secret put VAPID_PUBLIC_KEY
npx wrangler secret put VAPID_PRIVATE_KEY
npx wrangler secret put VAPID_SUBJECT
```

Deploy (Worker + PWA static assets):

```bash
task watch:deploy
# or: npx wrangler deploy
```

Local dev:

```bash
task watch:dev
```

Quick local route check (starts `wrangler dev`, hits `/health` and PWA):

```powershell
.\cloud\echo-watch\scripts\verify-local.ps1
```

Production deploy (requires `wrangler login` + KV + VAPID secrets):

```powershell
.\cloud\echo-watch\scripts\deploy.ps1
```

## 3. iPhone PWA setup (Safari)

1. Open the deployed Worker URL in **Safari** on iPhone.
2. **Share → Add to Home Screen**.
3. Launch **Echo Watch** from the Home Screen icon (required for push on iOS 16.4+).
4. Enter your **pairing token** and tap **Enable notifications**.
5. Tap **Send test push** to verify.

> Push does **not** work from a normal Safari tab on iOS — only from the installed PWA.

## 4. Desktop agent (`echo-watch`)

Build:

```bash
task watch:build
# → bin/echo-watch.exe
```

Copy and edit config:

```bash
copy cloud\echo-watch\echo-watch.example.yaml %APPDATA%\TalusEcho\echo-watch.yaml
```

Set `worker_url` and `secret` (same token as the PWA). Adjust:

- `process_name` — local VDI client executable (`mstsc.exe`, `vmware-view.exe`, `wfica32.exe`)
- `window_title_regex` — title of the client window on your PC
- `region` — crop inside the client where status text appears

Test alert (no screen capture):

```bash
bin\echo-watch.exe --trigger-test
# or: go run ./cmd/echo-watch --trigger-test
```

Run watcher:

```bash
task watch:run
# CONFIG=%APPDATA%\TalusEcho\echo-watch.yaml task watch:run
```

## API

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | — | Liveness |
| GET | `/api/vapid-public-key` | — | Public VAPID key for subscribe |
| POST | `/api/subscribe` | token in body | Store Web Push subscription |
| DELETE | `/api/subscribe` | token in body | Remove subscription |
| POST | `/api/alert` | `Bearer <token>` | Send notification to paired device |
| POST | `/api/test-push` | `Bearer <token>` | Send test notification |

### Example alert (`curl`)

```bash
curl -X POST "https://echo-watch.example.workers.dev/api/alert" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"rule_id":"success","title":"Build","body":"Successful"}'
```

## VDI window state

Before each capture, echo-watch checks:

| Check | On failure |
|-------|------------|
| VDI process running (`process_name`) | `vdi_process_not_found` — skip |
| Client HWND found | `vdi_window_not_found` — skip |
| Not minimized (`IsIconic`) | `vdi_window_minimized` — skip |
| Visible + non-zero size | `vdi_window_not_visible` — skip |

Keep the VDI client **open and not minimized** while a build is running.

## Security

- Treat the pairing token like an API key.
- Default alerts contain only rule id + short title/body — no screenshots.
- OCR runs locally; only matched alerts leave your PC.

## Repo layout

```
cmd/echo-watch/           CLI entry
internal/watch/           config, capture, OCR, rules, notifier
cloud/echo-watch/worker/  Cloudflare Worker (TypeScript)
cloud/echo-watch/pwa/     Safari-first PWA
```

Separate from the main Talus Echo Wails app — `task build` is unchanged.
