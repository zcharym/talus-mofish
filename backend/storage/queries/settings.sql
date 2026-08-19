-- name: GetSetting :one
SELECT key, value, updated_at
FROM settings
WHERE key = ?;

-- name: UpsertSetting :exec
INSERT INTO settings (key, value, updated_at)
VALUES (?, ?, datetime('now'))
ON CONFLICT(key) DO UPDATE SET
    value = excluded.value,
    updated_at = datetime('now');

-- name: ListSettings :many
SELECT key, value, updated_at
FROM settings
ORDER BY key;

-- name: DeleteSetting :exec
DELETE FROM settings
WHERE key = ?;
