# Talus Mofish

Cross-platform desktop app (Windows / macOS) for English learning, built with [Wails v3](https://v3.wails.io/) and [sqlc](https://sqlc.dev/) over SQLite.

## Prerequisites

- Go 1.24+ (this project uses `toolchain go1.25.4` for Wails v3 alpha.60)
- [Wails v3 CLI](https://v3.wails.io/quick-start/installation/): `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-alpha.60`
- Pin `@wailsio/runtime` to `3.0.0-alpha.76` in `frontend/package.json` (must match the Go Wails module version family)
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
| `main.go` | Wails app entry, service registration |
| `appservice.go` | Backend API exposed to the frontend |
| `frontend/` | React + TypeScript UI (Vite) |
| `internal/database/schema.sql` | SQLite schema (sqlc source of truth, embedded at runtime) |
| `db/queries/` | SQL queries consumed by sqlc |
| `internal/store/` | sqlc-generated Go data access (`task sqlc`) |
| `internal/database/` | DB open, embedded schema apply, default path |
| `sqlc.yaml` | sqlc configuration |

## Database

- Engine: SQLite via [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) (pure Go; works with Wails Windows builds where `CGO_ENABLED=0`)
- Default file: `talus-mofish/talus-mofish.db` under the OS user data directory:
  - **Windows**: `%LOCALAPPDATA%\talus-mofish\` (`internal/database/paths_windows.go`)
  - **macOS**: `~/Library/Application Support/talus-mofish/`
  - **Linux**: `$XDG_CONFIG_HOME/talus-mofish/` or `~/.config/talus-mofish/` (`paths_unix.go`)
- Schema: idempotent SQL in `internal/database/schema.sql` (embedded in `database.go` at build time)

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

1. Edit `internal/database/schema.sql` (or add numbered migration files later).
2. Add queries under `db/queries/*.sql`.
3. Run `task sqlc` and commit `internal/store/` changes.

## Wails services

- `AppService` — settings CRUD, config, autostart, and `DatabasePath()`

Bindings are generated under `frontend/bindings/` when running `wails3 dev` or `wails3 build`.
