# Talus Mofish

A cross-platform desktop agent application built with [Wails v3](https://v3.wails.io/) and [sqlc](https://sqlc.dev/) over SQLite. Talus Echo provides a **chat-oriented interface** where the Agent window serves as the primary surface for interactive sessions. **English Learning** is the first domain, featuring vocabulary management, reading tools, Anki import, and SRS (Spaced Repetition System); additional domains can be integrated using the same Management and Agent architecture.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Project Layout](#project-layout)
- [Domains (DDD)](#domains-ddd)
- [Database](#database)
- [Wails Services](#wails-services)
- [LLM Control (MCP)](#llm-control-mcp)
- [Documentation](#documentation)
- [Development](#development)

## Prerequisites

- **Go 1.27+** (this project uses `toolchain go1.27.0` for Wails v3 beta.4)
- **[Wails v3 CLI](https://v3.wails.io/quick-start/installation/)**: Install with `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.4`
- **Node.js** (for the React frontend)
- **[sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)** (optional, for regenerating query code)

> **Note**: Pin `@wailsio/runtime` to `3.0.0-beta.1` in `frontend/package.json` (latest published runtime; keep near the Go Wails module version).

Run `wails3 doctor` to verify your development environment.

## Quick Start

### Development Mode

Start the development server with hot reload:

```bash
cd frontend && npm install && cd ..
wails3 dev
```

The application will launch with live reload enabled for both Go and React code changes.

### Production Build

Create an optimized production build:

```bash
wails3 build
```

The compiled application will be available in the `build/bin` directory.

### Cross-Compilation

Build for different platforms (requires appropriate toolchains):

```bash
# Build for Windows
wails3 build GOOS=windows

# Build for macOS
wails3 build GOOS=darwin
```

## Project Layout

The codebase follows a domain-driven design with clear separation between frontend, backend services, and domain-specific logic:

| Path | Purpose |
|------|---------|
| `main.go` | Application entry point; registers backend services with Wails |
| `frontend/` | React + TypeScript UI with Vite build system |
| `backend/services/` | Wails-bound API facades (`System`, `Config`, `Auth`, `Chat`, `English`, `Sudoku`, `Obsidian`) |
| `backend/storage/` | SQLite connection management, schema, config.json, sqlc queries + store |
| `backend/types/` | Shared data transfer objects for Wails/JSON communication |
| `backend/utils/` | Environment loader, LLM client, autostart utilities |
| `backend/consts/` | Domain identifiers and shared constants |
| `backend/agent/` | Chat orchestration and streaming turn management |
| `backend/auth/` | Authentication and identity kernel |
| `backend/english/content/` | English Learning domain content importers (Anki APKG) |
| `backend/obsidian/` | Obsidian Local REST API client integration |
| `backend/watch/`, `backend/vdiupload/` | Sidecar domain packages |
| `cmd/echo-watch/`, `cmd/vdi-upload/` | Command-line interfaces for sidecar services |
| `cloud/echo-watch/` | Cloudflare Worker + iOS PWA for watch alerts |
| `sqlc.yaml` | sqlc configuration for type-safe SQL query generation |
| `docs/domains/` | DDD domain map and per-domain design documentation |

## Domains (DDD)

Talus Echo is architected as a multi-domain application following Domain-Driven Design principles. See **[docs/domains/README.md](docs/domains/README.md)** for the complete bounded-context map, shared-kernel rules, and domain integration patterns.

### Available Domains

| Domain | Type | Entry Point | Documentation |
|--------|------|-------------|---------------|
| **English Learning** | desktop-agent | Wails app (`english.*` routes) | [Design & Plan](docs/design-and-plan.md) |
| **Obsidian** | desktop-agent | Wails app (`obsidian.*` routes) | [README](docs/domains/obsidian/README.md) |
| **Sudoku** | desktop-agent | Wails app (Agent window games) | [README](docs/domains/sudoku/README.md) |
| **Echo Watch** | sidecar | `task watch:build` / `cloud/echo-watch` | [README](docs/domains/watch/README.md) |
| **VDI Upload** | sidecar | `task vdiupload:build` | [DESIGN](docs/domains/vdiupload/DESIGN.md) |

Each domain is self-contained and communicates through well-defined interfaces in the shared kernel.

## Database

### Technology Stack

- **Engine**: SQLite via [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) (pure Go implementation)
  - Works seamlessly with Wails Windows builds where `CGO_ENABLED=0`
  - No external dependencies or C compiler required
- **Schema Management**: Idempotent SQL in `backend/storage/schema.sql`
- **Query Generation**: Type-safe Go code via [sqlc](https://sqlc.dev/)

### Database Location

The database file is stored in the OS-specific user data directory:

- **Windows**: `%LOCALAPPDATA%\talus-mofish\talus-mofish.db` (via `backend/storage/paths_windows.go`)
- **macOS**: `~/Library/Application Support/talus-mofish/talus-mofish.db`
- **Linux**: `$XDG_CONFIG_HOME/talus-mofish/talus-mofish.db` or `~/.config/talus-mofish/talus-mofish.db` (via `paths_unix.go`)

### Common Database Tasks

**Print the database path:**

```bash
task db:path
# or
go run ./cmd/dbpath
```

**Regenerate store code after modifying SQL:**

```bash
task sqlc
```

### Working with Schema and Queries

1. **Edit schema**: Modify `backend/storage/schema.sql` (or add numbered migration files)
2. **Add queries**: Create/modify `*.sql` files under `backend/storage/queries/`
3. **Regenerate code**: Run `task sqlc` to generate type-safe Go code
4. **Commit changes**: Include both SQL and generated Go files in version control

## Wails Services

Services are registered in `main.go` and exposed to the frontend via automatic bindings. The architecture follows a clean separation of concerns with dedicated services for each major functional area:

| Service | Responsibilities |
|---------|------------------|
| `SystemService` | Database path inspection, window management, file picker dialogs, media paths, autostart configuration |
| `ConfigService` | Application settings via `config.json`, key-value settings store |
| `AuthService` | User authentication, sign-in/sign-out flows, session management |
| `ChatService` | Chat session management, message history, streaming LLM responses |
| `EnglishService` | Anki APKG import, article management, vocabulary, SRS card operations |
| `SudokuService` | YouDoSudoku puzzle generation and gameplay in the Agent window |
| `ObsidianService` | Vault browsing, note editing, search via Obsidian's Local REST API |

TypeScript bindings are automatically generated under `frontend/bindings/` when running `wails3 dev` or `wails3 build`.

## LLM Control (MCP)

Wails v3 can compile a loopback [MCP server](https://v3.wails.io/guides/mcp-service/) directly into the desktop application, enabling AI agents to interact with the app's features programmatically. This repo enables MCP by default via `Taskfile.yml` configuration (`WAILS_MCP=1`, `EXTRA_TAGS=mcp`).

### Starting with MCP Support

```bash
task dev
```

### MCP Server Information

On successful startup, you'll see:

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

The `.cursor/mcp.json` configuration in this repository points Cursor at this URL. **Important**: The MCP endpoint is only available while the application is running.

### Building Without MCP

To build a version without MCP support:

```bash
task build EXTRA_TAGS=
```

## Documentation

Comprehensive documentation is available in the `docs/` directory:

| Document | Description |
|----------|-------------|
| [docs/domains/README.md](docs/domains/README.md) | Complete DDD domain map, bounded-context diagram, and dependency rules |
| [docs/design-and-plan.md](docs/design-and-plan.md) | Product vision, architecture overview, data model, and implementation phases |
| [docs/chat-learning-flows-plan.md](docs/chat-learning-flows-plan.md) | English Learning agent flows: IELTS, Anki recitation, article reading, tools, and UI |
| [docs/domains/vdiupload/DESIGN.md](docs/domains/vdiupload/DESIGN.md) | Engineering design for H3C / VDI automated file upload |
| [docs/system-tray-menu.md](docs/system-tray-menu.md) | System tray behavior and menu structure |
| [cloud/echo-watch/README.md](cloud/echo-watch/README.md) | Echo Watch deployment guide and PWA setup instructions |

## Development

### Code Quality

The codebase follows these principles:

- **Type Safety**: Uses `sqlc` for type-safe database queries and TypeScript for frontend code
- **Error Handling**: Comprehensive error wrapping with context-specific messages
- **Documentation**: All exported functions and types include descriptive comments
- **Testing**: Unit tests for critical business logic (run with `go test ./...`)

### Contributing

When adding new features:

1. Follow the existing service patterns in `backend/services/`
2. Add database queries to `backend/storage/queries/` and regenerate with `task sqlc`
3. Update the relevant domain documentation in `docs/domains/`
4. Ensure TypeScript types align with Go types in `backend/types/`

### Useful Commands

```bash
# Run development server
task dev

# Build production binary
task build

# Regenerate database code
task sqlc

# Show database path
task db:path

# Build watch sidecar
task watch:build

# Build VDI upload sidecar
task vdiupload:build
```

## License

See the project repository for license information.
