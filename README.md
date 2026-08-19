# Talus Mofish

Cross-platform desktop agent app (Windows / macOS), built with [Wails v3](https://v3.wails.io/) and [sqlc](https://sqlc.dev/) over SQLite.

Talus Echo is **chat-oriented**: the Agent window is the primary surface for interactive sessions. **English Learning** is the first domain (vocabulary, reading, Anki import, SRS); additional domains will plug into the same Management and Agent architecture.

## Prerequisites

- Go 1.24+ (this project uses `toolchain go1.25.4` for Wails v3 beta.4)
- [Wails v3 CLI](https://v3.wails.io/quick-start/installation/): `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.4`
- Pin `@wailsio/runtime` to `3.0.0-beta.1` in `frontend/package.json` (latest published runtime; keep near the Go Wails module version)
- Node.js (for the React frontend)
- Optional: [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) for regenerating query code

Run `wails3 doctor` to verify your environment.

## Quick start

```bash
cd frontend && npm install && cd ..
wails3 dev
```

Production build:

```bash
wails3 build
```

Cross-compile (from macOS or Windows with toolchains installed):

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
```

## Project layout

| Path | Purpose |
|------|---------|
| `main.go` | Wails app entry, registers five backend services |
| `frontend/` | React + TypeScript UI (Vite) |
| `backend/services/` | Wails-bound API facades (`System`, `Config`, `Auth`, `Chat`, `English`) |
| `backend/storage/` | SQLite open, schema, config.json, sqlc queries + store |
| `backend/types/` | Shared Wails/JSON DTOs |
| `backend/utils/` | Env loader, LLM client, autostart |
| `backend/consts/` | Domain IDs and shared constants |
| `backend/agent/`, `backend/auth/` | Chat orchestration and identity kernel |
| `backend/english/content/` | English Learning importers (Anki APKG) |
| `backend/watch/`, `backend/vdiupload/` | Sidecar domain packages |
| `cmd/echo-watch/`, `cmd/vdi-upload/` | Sidecar CLIs |
| `cloud/echo-watch/` | Cloudflare Worker + iOS PWA for watch alerts |
| `sqlc.yaml` | sqlc configuration |
| `docs/domains/` | DDD domain map and per-domain design docs |

## Domains (DDD)

Talus Echo is multi-domain. See **[docs/domains/README.md](docs/domains/README.md)** for the bounded-context map, shared-kernel rules, and migration plan.

| Domain | Kind | Entry |
|--------|------|-------|
| English Learning | desktop-agent | Wails app (`english.*`) |
| Echo Watch | sidecar | `task watch:build` / `cloud/echo-watch` |
| VDI Upload | sidecar | `task vdiupload:build` — design: [DESIGN.md](docs/domains/vdiupload/DESIGN.md) |

## Database

- Engine: SQLite via [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) (pure Go; works with Wails Windows builds where `CGO_ENABLED=0`)
- Default file: `talus-mofish/talus-mofish.db` under the OS user data directory:
  - **Windows**: `%LOCALAPPDATA%\talus-mofish\` (`backend/storage/paths_windows.go`)
  - **macOS**: `~/Library/Application Support/talus-mofish/`
  - **Linux**: `$XDG_CONFIG_HOME/talus-mofish/` or `~/.config/talus-mofish/` (`paths_unix.go`)
- Schema: idempotent SQL in `backend/storage/schema.sql` (embedded in `storage/database.go` at build time)

Print the default DB path:

```bash
task db:path
# or
go run ./cmd/dbpath
```

Regenerate store code after changing SQL:

```bash
task sqlc
```

## Adding schema / queries

1. Edit `backend/storage/schema.sql` (or add numbered migration files later).
2. Add queries under `backend/storage/queries/*.sql`.
3. Run `task sqlc` and commit `backend/storage/store/` changes.

## Wails services

Five services are bound from `backend/services/` (Tiny RDM-style split):

| Service | Role |
|---------|------|
| `SystemService` | DB path, windows, file picker, media paths, autostart status |
| `ConfigService` | `config.json`, settings KV |
| `AuthService` | Sign-in / sign-out |
| `ChatService` | Chat sessions and streaming turns |
| `EnglishService` | Anki import, articles, vocabulary, SRS |

Bindings are generated under `frontend/bindings/` when running `wails3 dev` or `wails3 build`.

## LLM Control (MCP)

Wails v3 can compile a loopback [MCP server](https://v3.wails.io/guides/mcp-service/) into the desktop app. Agents then list windows, inspect the DOM, click/type, and call bound Go services. The server is **opt-in**: without the `mcp` build tag it is absent from the binary.

```bash
task dev:mcp
# PowerShell equivalent:
# $env:WAILS_MCP=1; wails3 dev
```

On startup the app logs:

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

This repo's `.cursor/mcp.json` points Cursor at that URL. The endpoint is only up while the app is running with MCP enabled. Production `task build` / `task package` do not set `WAILS_MCP`, so release binaries do not include the server.

## Documentation

| Document | Description |
|----------|-------------|
| [docs/domains/README.md](docs/domains/README.md) | DDD domain map and dependency rules |
| [docs/domains/vdiupload/DESIGN.md](docs/domains/vdiupload/DESIGN.md) | H3C / VDI automated file upload engineering design |
| [docs/design-and-plan.md](docs/design-and-plan.md) | Product vision, architecture, data model, implementation phases |
| [docs/chat-learning-flows-plan.md](docs/chat-learning-flows-plan.md) | English Learning agent flows (IELTS, Anki recite, article reading, tools, UI) |
| [docs/system-tray-menu.md](docs/system-tray-menu.md) | System tray behavior |
| [cloud/echo-watch/README.md](cloud/echo-watch/README.md) | Echo Watch deploy and PWA setup |
