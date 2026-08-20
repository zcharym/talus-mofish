# Obsidian domain

**ID:** `obsidian` · **Kind:** desktop-agent · **Constant:** `consts.Obsidian`

## Responsibility

Browse, edit, and search an Obsidian vault from the Management window through [Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api). Notes live in Obsidian; Talus Echo does not copy the vault into SQLite.

## Setup

1. Install and enable **Local REST API** in Obsidian.
2. Copy the API key from **Settings → Local REST API**.
3. In Talus Echo **Config → Obsidian**, paste the key and save, then **Test connection**.
4. Keep Obsidian running while using Notes and Search.

Default URL is `https://127.0.0.1:27124` (plugin HTTPS with a self-signed certificate). If you enable the plugin's HTTP server, use `http://127.0.0.1:27123` instead.

## Code ownership

| Layer | Path |
|-------|------|
| HTTP client + path rules | `backend/obsidian` |
| Wails API | `backend/services/obsidian.go` |
| Config | `config.json` (`obsidian.baseUrl`, `obsidian.apiKey`) |
| Management UI | `frontend/src/pages/ObsidianNotesPage.tsx`, `ObsidianSearchPage.tsx` |

## Boundaries

- **Owns:** Local REST API client, vault path validation, management Notes/Search UI.
- **Uses:** shared kernel (`storage` config, `types`).
- **Does not own:** tags, command palette, create/delete/rename, agent/MCP tools (later slices).

## Related

- Plugin API: [coddingtonbear/obsidian-local-rest-api](https://github.com/coddingtonbear/obsidian-local-rest-api)
- Domain map: [../README.md](../README.md)
