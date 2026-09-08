# Cloudflare domain

**ID:** `cloudflare` · **Kind:** desktop-agent · **Constant:** `consts.Cloudflare`

## Responsibility

Read-only operations dashboard in the Agent window for a saved Cloudflare account. v1 shows account health, Workers inventory, 24h invocation analytics, D1 databases, and KV namespaces. Talus Echo does not deploy, tail, or mutate Cloudflare resources from this UI.

A **system-pinned** sidebar tab labeled **Cloudflare** appears when:

- the signed-in user is the debug admin (`provider = debug`), or
- `config.json` has both `cloudflare.accountId` and `cloudflare.apiToken`

Debug users without credentials still see the tab and a setup empty state that points at **Management → Configuration → Cloudflare**.

## Setup

1. Create a Cloudflare API token with at least: **Account Settings Read**, **Workers Scripts Read**, **Account Analytics Read**, **D1 Read**, **Workers KV Storage Read**.
2. Copy the **Account ID** from the Cloudflare dashboard.
3. In Talus Echo **Config → Cloudflare**, paste both values and save, then **Test connection**.
4. Open the Agent window — the pinned **Cloudflare** tab loads a snapshot (Refresh to reload).

## Code ownership

| Layer | Path |
|-------|------|
| API client ([cloudflare-go/v7](https://github.com/cloudflare/cloudflare-go)) | `backend/cloudflare` |
| Wails API | `backend/services/cloudflare.go` |
| Config | `config.json` (`cloudflare.accountId`, `cloudflare.apiToken`) |
| Agent UI | `frontend/src/components/agent/CloudflareDashboard` |
| Config UI | `frontend/src/pages/config/CloudflareTab.tsx` |

## Boundaries

- **Owns:** Cloudflare API client (`cloudflare-go/v7` REST resources; GraphQL analytics via `Client.Post`), dashboard snapshot DTO, Agent pinned tab, config fields.
- **Uses:** shared kernel (`storage` config, `types`). The API token stays in `config.json`; the Agent UI calls `IsConfigured` / `GetDashboard` / `Ping` and does not parse the token for visibility.
- **Does not own:** Echo Watch push Worker (`cloud/echo-watch`), Talus Auth Worker (`workers/auth`), write APIs (deploy, rollback, tail), Management-window duplicate of the dashboard.

Known-service cards highlight `talus-auth` and `echo-watch` when those scripts exist on the account.

## Related

- Cloudflare API: [developers.cloudflare.com/api](https://developers.cloudflare.com/api/)
- Go SDK: [github.com/cloudflare/cloudflare-go](https://github.com/cloudflare/cloudflare-go) (`v7`)
- Domain map: [../README.md](../README.md)
