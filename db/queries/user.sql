-- name: GetCurrentUser :one
SELECT id, provider, provider_user_id, display_name, email, avatar_url, created_at, last_login_at
FROM user_account
ORDER BY last_login_at DESC
LIMIT 1;

-- name: DeleteAllUserAccounts :exec
DELETE FROM user_account;

-- name: InsertUserAccount :exec
INSERT INTO user_account (
    id,
    provider,
    provider_user_id,
    display_name,
    email,
    avatar_url,
    created_at,
    last_login_at
)
VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'));
