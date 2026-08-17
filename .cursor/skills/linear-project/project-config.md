# Linear project bindings — talus_echo_loop

Fill in these values after the first Linear MCP discovery call. The agent reads this file to avoid re-querying IDs every session.

## Workspace

| Field | Value |
|-------|-------|
| Team name | Taluship |
| Team ID | TAL |
| Team UUID | 9014d09c-a528-4e93-8cb2-b186a179c4b7 |
| Project name | talus-mofish |
| Project ID | 0f90a64a-bba4-42d7-982f-c69fe0046fb6 |
| Project URL | https://linear.app/taluship/project/talus-mofish-edffbfec79d9 |
| Initiative name (optional) | _TODO_ |
| Initiative ID (optional) | _TODO_ |

## Defaults for new issues

| Field | Value |
|-------|-------|
| Default team ID | 9014d09c-a528-4e93-8cb2-b186a179c4b7 |
| Default project ID | 0f90a64a-bba4-42d7-982f-c69fe0046fb6 |
| Default labels | ""|
| Default priority | 0 |

## Repo context (auto-filled)

| Field | Value |
|-------|-------|
| Repository | `talus-mofish` |
| Product | Talus Mofish — chat-oriented desktop agent (English Learning is the first domain) |
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
