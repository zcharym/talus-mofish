-- name: CreateArticle :exec
INSERT INTO articles (
    id, title, url, content, translation, level, source, word_count,
    anki_note_guid, model_css, raw_fields, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'));

-- name: GetArticle :one
SELECT
    id, title, url, content, translation, level, source, word_count,
    anki_note_guid, model_css, raw_fields, created_at, updated_at
FROM articles
WHERE id = ?;

-- name: GetArticleByAnkiNoteGuid :one
SELECT
    id, title, url, content, translation, level, source, word_count,
    anki_note_guid, model_css, raw_fields, created_at, updated_at
FROM articles
WHERE anki_note_guid = ?;

-- name: ListArticles :many
SELECT
    id, title, url, content, translation, level, source, word_count,
    anki_note_guid, model_css, raw_fields, created_at, updated_at
FROM articles
ORDER BY created_at DESC;

-- name: ListArticlesBySource :many
SELECT
    id, title, url, content, translation, level, source, word_count,
    anki_note_guid, model_css, raw_fields, created_at, updated_at
FROM articles
WHERE source = ?
ORDER BY created_at DESC;

-- name: UpdateArticle :exec
UPDATE articles
SET title = ?, url = ?, content = ?, translation = ?, level = ?, word_count = ?, model_css = ?, raw_fields = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteArticle :exec
DELETE FROM articles
WHERE id = ?;

-- name: CreateClipping :exec
INSERT INTO clippings (id, article_id, text, translation, explanation, note, saved_to_card, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'));

-- name: GetClipping :one
SELECT id, article_id, text, translation, explanation, note, saved_to_card, created_at
FROM clippings
WHERE id = ?;

-- name: ListClippingsByArticle :many
SELECT id, article_id, text, translation, explanation, note, saved_to_card, created_at
FROM clippings
WHERE article_id = ?
ORDER BY created_at DESC;

-- name: UpdateClipping :exec
UPDATE clippings
SET text = ?, translation = ?, explanation = ?, note = ?, saved_to_card = ?
WHERE id = ?;

-- name: DeleteClipping :exec
DELETE FROM clippings
WHERE id = ?;
