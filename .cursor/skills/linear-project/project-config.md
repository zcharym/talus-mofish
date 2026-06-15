# Linear project bindings — talus_echo_loop

Fill in these values after the first Linear MCP discovery call. The agent reads this file to avoid re-querying IDs every session.

## Workspace

| Field | Value |
|-------|-------|
| Team name | Talus |
| Team ID | TAL |
| Team UUID | cb7a680d-472b-4bfc-8cc7-209f55acf3b1 |
| Project name | Talus Mofish |
| Project ID | 32601757-e459-4f86-9e8e-228ebc69fe15 |
| Project URL | https://linear.app/talus/project/talus-mofish-a8bb88fce0f5 |
| Initiative name (optional) | _TODO_ |
| Initiative ID (optional) | _TODO_ |

## Defaults for new issues

| Field | Value |
|-------|-------|
| Default team ID | _same as Team ID_ |
| Default project ID | 32601757-e459-4f86-9e8e-228ebc69fe15 |
| Default labels | ""|
| Default priority | 0 |

## Repo context (auto-filled)

| Field | Value |
|-------|-------|
| Repository | `talus-mofish` |
| Product | Talus Mofish — English learning desktop app |
| Stack | Wails v3, Go, React/TypeScript, SQLite/sqlc |
| Default branch | `main` |

## Naming conventions

- **Issue titles**: `[area] short imperative` — e.g. `[frontend] Add spaced-repetition review screen`
- **Areas**: `frontend`, `backend`, `db`, `wails`, `docs`, `infra`
- **Branch links**: include `Branch: feature/ABC-123-slug` in description when work starts
- **PR links**: add PR URL to issue when opened; move to In Review

## Discovery commands

Run once after plugin install to populate the TODO fields above:

1. List teams → find team name/ID
2. Search or list projects → find project linked to Talus Mofish
3. Update this file with resolved IDs
