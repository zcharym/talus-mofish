-- name: CreateAnkiImport :exec
INSERT INTO anki_imports (id, filename, imported_at, status, error_message, stats_json)
VALUES (?, ?, datetime('now'), ?, ?, ?);

-- name: GetAnkiImport :one
SELECT id, filename, imported_at, status, error_message, stats_json
FROM anki_imports
WHERE id = ?;

-- name: ListAnkiImports :many
SELECT id, filename, imported_at, status, error_message, stats_json
FROM anki_imports
ORDER BY imported_at DESC;

-- name: UpdateAnkiImportStatus :exec
UPDATE anki_imports
SET status = ?, error_message = ?, stats_json = ?
WHERE id = ?;

-- name: CreateAnkiImportDeck :exec
INSERT INTO anki_import_decks (
    import_id, anki_deck_id, anki_deck_name, target_type, anki_model_id, field_mapping, app_deck_id
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: ListAnkiImportDecks :many
SELECT import_id, anki_deck_id, anki_deck_name, target_type, anki_model_id, field_mapping, app_deck_id
FROM anki_import_decks
WHERE import_id = ?
ORDER BY anki_deck_name;

-- name: GetAnkiImportDeck :one
SELECT import_id, anki_deck_id, anki_deck_name, target_type, anki_model_id, field_mapping, app_deck_id
FROM anki_import_decks
WHERE import_id = ? AND anki_deck_id = ?;
