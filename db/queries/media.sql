-- name: CreateMediaAsset :exec
INSERT INTO media_assets (id, original_filename, stored_path, mime_type, sha256, source, created_at)
VALUES (?, ?, ?, ?, ?, ?, datetime('now'));

-- name: GetMediaAsset :one
SELECT id, original_filename, stored_path, mime_type, sha256, source, created_at
FROM media_assets
WHERE id = ?;

-- name: GetMediaAssetByFilename :one
SELECT id, original_filename, stored_path, mime_type, sha256, source, created_at
FROM media_assets
WHERE original_filename = ? AND source = ?;

-- name: ListMediaAssets :many
SELECT id, original_filename, stored_path, mime_type, sha256, source, created_at
FROM media_assets
WHERE source = ?
ORDER BY original_filename;
