# Feeds domain

**ID:** `feeds` · **Kind:** desktop-agent · **Constant:** `consts.Feeds`

## Responsibility

Pinned **Feeds** tab in the Agent window: a local inbox for RSS/Atom, YouTube channels and playlists, Bilibili watch later / uploads / favorites, and a Read later list.

The tab is always visible. Credentials are optional:

- RSS works with no keys (`github.com/mmcdole/gofeed`).
- YouTube `@handle` URLs resolve via [YouTube Data API v3](https://developers.google.com/youtube/v3) when `feeds.youtubeApiKey` is set, otherwise by reading the channel page RSS link. Items are pulled from official YouTube Atom feeds.
- Bilibili **稍后再看** uses `GET /x/v2/history/toview` with a `SESSDATA` cookie. Space uploads and public favorites use the documented web APIs.

YouTube’s private Watch Later playlist (`WL`) is not available without user OAuth; add a public playlist or use Bilibili watch later instead.

## Setup

1. Open the Agent window → **Feeds** → **Add**.
2. Paste an RSS URL, `https://www.youtube.com/@handle`, a playlist URL, `https://space.bilibili.com/<uid>`, a favorites URL with `fid=`, or type `稍后再看`.
3. Optional: **Management → Configuration → Feeds** for a YouTube Data API key and Bilibili `SESSDATA`. Save, then Test.

## Code ownership

| Layer | Path |
|-------|------|
| Fetch + URL resolve ([gofeed](https://github.com/mmcdole/gofeed), [youtube/v3](https://pkg.go.dev/google.golang.org/api/youtube/v3)) | `backend/feeds` |
| Wails API | `backend/services/feeds.go` |
| Persistence | `backend/storage/schema.sql`, `backend/storage/queries/feeds.sql` |
| Config | `config.json` (`feeds.youtubeApiKey`, `feeds.bilibiliSessdata`) in the OS keyring |
| Agent UI | `frontend/src/components/agent/FeedsInbox` |
| Config UI | `frontend/src/pages/config/FeedsTab.tsx` |

## Boundaries

- **Owns:** source/item SQLite tables, inbox DTOs, RSS/YouTube/Bilibili fetchers, Agent pinned tab, feed config fields.
- **Uses:** shared kernel (`storage`, `types`). Tokens stay in the keyring; the Agent UI never displays them.
- **Does not own:** English Learning articles, chat sessions, or opening URLs (uses `SystemService.OpenURL`).

## Related

- Domain map: [../README.md](../README.md)
