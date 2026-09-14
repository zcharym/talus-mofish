-- name: ListFeedSources :many
SELECT
    id,
    kind,
    title,
    url,
    remote_id,
    enabled,
    last_error,
    last_fetched_at,
    created_at,
    updated_at
FROM feed_sources
ORDER BY kind ASC, title ASC;

-- name: GetFeedSource :one
SELECT
    id,
    kind,
    title,
    url,
    remote_id,
    enabled,
    last_error,
    last_fetched_at,
    created_at,
    updated_at
FROM feed_sources
WHERE id = ?;

-- name: CreateFeedSource :exec
INSERT INTO feed_sources (
    id,
    kind,
    title,
    url,
    remote_id,
    enabled,
    last_error,
    last_fetched_at,
    created_at,
    updated_at
)
VALUES (?, ?, ?, ?, ?, 1, '', '', datetime('now'), datetime('now'));

-- name: DeleteFeedSource :exec
DELETE FROM feed_sources
WHERE id = ?;

-- name: UpdateFeedSourceFetch :exec
UPDATE feed_sources
SET
    title = ?,
    last_error = ?,
    last_fetched_at = datetime('now'),
    updated_at = datetime('now')
WHERE id = ?;

-- name: ListFeedItems :many
SELECT
    i.id,
    i.source_id,
    i.external_id,
    i.title,
    i.url,
    i.author,
    i.summary,
    i.thumbnail_url,
    i.published_at,
    i.saved,
    i.read,
    i.created_at,
    s.kind AS source_kind,
    s.title AS source_title
FROM feed_items i
INNER JOIN feed_sources s ON s.id = i.source_id
ORDER BY
    CASE WHEN i.published_at = '' THEN i.created_at ELSE i.published_at END DESC,
    i.id DESC
LIMIT 300;

-- name: ListSavedFeedItems :many
SELECT
    i.id,
    i.source_id,
    i.external_id,
    i.title,
    i.url,
    i.author,
    i.summary,
    i.thumbnail_url,
    i.published_at,
    i.saved,
    i.read,
    i.created_at,
    s.kind AS source_kind,
    s.title AS source_title
FROM feed_items i
INNER JOIN feed_sources s ON s.id = i.source_id
WHERE i.saved = 1
ORDER BY
    CASE WHEN i.published_at = '' THEN i.created_at ELSE i.published_at END DESC,
    i.id DESC
LIMIT 300;

-- name: ListFeedItemsByKind :many
SELECT
    i.id,
    i.source_id,
    i.external_id,
    i.title,
    i.url,
    i.author,
    i.summary,
    i.thumbnail_url,
    i.published_at,
    i.saved,
    i.read,
    i.created_at,
    s.kind AS source_kind,
    s.title AS source_title
FROM feed_items i
INNER JOIN feed_sources s ON s.id = i.source_id
WHERE s.kind = ?
ORDER BY
    CASE WHEN i.published_at = '' THEN i.created_at ELSE i.published_at END DESC,
    i.id DESC
LIMIT 300;

-- name: GetFeedItem :one
SELECT
    id,
    source_id,
    external_id,
    title,
    url,
    author,
    summary,
    thumbnail_url,
    published_at,
    saved,
    read,
    created_at
FROM feed_items
WHERE id = ?;

-- name: GetFeedItemBySourceExternal :one
SELECT
    id,
    source_id,
    external_id,
    title,
    url,
    author,
    summary,
    thumbnail_url,
    published_at,
    saved,
    read,
    created_at
FROM feed_items
WHERE source_id = ?
    AND external_id = ?;

-- name: CreateFeedItem :exec
INSERT INTO feed_items (
    id,
    source_id,
    external_id,
    title,
    url,
    author,
    summary,
    thumbnail_url,
    published_at,
    saved,
    read,
    created_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, datetime('now'));

-- name: UpdateFeedItemMetadata :exec
UPDATE feed_items
SET
    title = ?,
    url = ?,
    author = ?,
    summary = ?,
    thumbnail_url = ?,
    published_at = ?
WHERE id = ?;

-- name: SetFeedItemSaved :exec
UPDATE feed_items
SET saved = ?
WHERE id = ?;

-- name: SetFeedItemRead :exec
UPDATE feed_items
SET read = ?
WHERE id = ?;

-- name: DeleteFeedItem :exec
DELETE FROM feed_items
WHERE id = ?;
