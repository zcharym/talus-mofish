# English Learning domain

**ID:** `english` · **Kind:** desktop-agent · **Constant:** `domain.English`

## Responsibility

Vocabulary, articles/reading, Anki `.apkg` import, SRS review, and chat Agent flows for English learning. Primary surface: Wails Management (`english.*` routes) + Agent window.

## Code ownership (current)

| Area | Location |
|------|----------|
| Anki / content import | `internal/content`, `internal/content/apkg` |
| Persistence | `internal/store` (sqlc), `internal/database/schema.sql`, `db/queries/` |
| Wails API | `internal/appservice` (articles, vocabulary, srs, import, media, chat) |
| Agent tools / flows | `internal/agent` + docs in `docs/chat-learning-flows-plan.md` |

## Boundaries

- **Owns:** learning content model, import pipelines, SRS scheduling semantics, english Agent flows.
- **Uses:** shared kernel (`database`, `store`, `aiclient`, `auth`, `config`).
- **Does not own:** Echo Watch, VDI upload, Cloudflare push.

## Migration note

Physical package consolidation under `internal/english/` is deferred to avoid churn in Wails bindings and sqlc paths. Prefer extracting application services first, then relocating packages.
