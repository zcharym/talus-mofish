-- name: ListChatSessions :many
SELECT id, title, created_at, updated_at
FROM chat_sessions
ORDER BY updated_at DESC;

-- name: GetChatSession :one
SELECT id, title, created_at, updated_at
FROM chat_sessions
WHERE id = ?;

-- name: CreateChatSession :exec
INSERT INTO chat_sessions (id, title, created_at, updated_at)
VALUES (?, ?, datetime('now'), datetime('now'));

-- name: RenameChatSession :exec
UPDATE chat_sessions
SET title = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: TouchChatSession :exec
UPDATE chat_sessions
SET updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteChatSession :exec
DELETE FROM chat_sessions
WHERE id = ?;

-- name: ListChatMessages :many
SELECT id, session_id, role, content, created_at
FROM chat_messages
WHERE session_id = ?
ORDER BY created_at ASC;

-- name: CreateChatMessage :exec
INSERT INTO chat_messages (id, session_id, role, content, created_at)
VALUES (?, ?, ?, ?, datetime('now'));

-- name: GetChatMessage :one
SELECT id, session_id, role, content, created_at
FROM chat_messages
WHERE id = ?;

-- name: UpdateChatMessageContent :exec
UPDATE chat_messages
SET content = ?
WHERE id = ?;
