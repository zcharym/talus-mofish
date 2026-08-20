# Sudoku domain

**ID:** `sudoku` · **Kind:** desktop-agent · **Constant:** `consts.Sudoku`

## Responsibility

YouDoSudoku-backed puzzle games in the Agent window. Each game is a sidebar session (`chat_sessions.kind = sudoku`) with a dedicated board UI — not a chat thread.

## Code ownership

| Layer | Path |
|-------|------|
| Domain logic + HTTP client | `backend/sudoku` |
| Wails API | `backend/services/sudoku.go` |
| Persistence | `backend/storage/schema.sql`, `backend/storage/queries/sudoku.sql` |
| Agent UI | `frontend/src/components/agent/SudokuBoard` |

## Boundaries

- **Owns:** puzzle fetch from [YouDoSudoku](https://www.youdosudoku.com/), given-cell rules, check/complete against the server-side solution.
- **Uses:** shared kernel (`storage`, `types`, config). API key lives in `config.json` (`sudoku.apiKey`) and never goes to the renderer.
- **Does not own:** LLM chat orchestration (`backend/agent`).

## Related

- Domain map: [../README.md](../README.md)
