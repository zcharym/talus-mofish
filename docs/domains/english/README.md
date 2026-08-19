# English Learning domain

**ID:** `english` · **Kind:** desktop-agent · **Constant:** `domain.English`

## Responsibility

Vocabulary, articles/reading, Anki `.apkg` import, SRS review, and chat Agent flows for English learning. Primary surface: Wails Management (`english.*` routes) + Agent window.

## Code ownership (current)

| Area | Location |
|------|----------|
| Anki / content import | `backend/content`, `backend/content/apkg` |
| Persistence | `backend/store` (sqlc), `backend/database/schema.sql`, `backend/storage/queries/` |
| Wails API | `backend/services` (articles, vocabulary, srs, import, media, chat) |
| Agent tools / flows | `backend/agent` + docs in `docs/chat-learning-flows-plan.md` |

## Boundaries

- **Owns:** learning content model, import pipelines, SRS scheduling semantics, english Agent flows.
- **Uses:** shared kernel (`database`, `store`, `aiclient`, `auth`, `config`).
- **Does not own:** Echo Watch, VDI upload, Cloudflare push.

## Migration note

Physical package consolidation under `backend/english/` is deferred to avoid churn in Wails bindings and sqlc paths. Prefer extracting application services first, then relocating packages.
