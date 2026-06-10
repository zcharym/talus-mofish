-- name: CreateDeck :exec
INSERT INTO decks (id, name, description, parent_id, sort_order, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'));

-- name: GetDeck :one
SELECT id, name, description, parent_id, sort_order, created_at, updated_at
FROM decks
WHERE id = ?;

-- name: ListDecks :many
SELECT id, name, description, parent_id, sort_order, created_at, updated_at
FROM decks
ORDER BY sort_order, name;

-- name: ListChildDecks :many
SELECT id, name, description, parent_id, sort_order, created_at, updated_at
FROM decks
WHERE parent_id = ?
ORDER BY sort_order, name;

-- name: UpdateDeck :exec
UPDATE decks
SET name = ?, description = ?, parent_id = ?, sort_order = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteDeck :exec
DELETE FROM decks
WHERE id = ?;

-- name: CreateCard :exec
INSERT INTO cards (
    id, deck_id, front, back, example_sentence, hints, card_type,
    ease, interval, repetitions, due_at, lapsed, suspended,
    source, anki_card_id, anki_note_guid, template_ord, model_css,
    created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?,
    datetime('now'), datetime('now')
);

-- name: GetCard :one
SELECT
    id, deck_id, front, back, example_sentence, hints, card_type,
    ease, interval, repetitions, due_at, lapsed, suspended,
    source, anki_card_id, anki_note_guid, template_ord, model_css,
    created_at, updated_at
FROM cards
WHERE id = ?;

-- name: GetCardByAnkiCardID :one
SELECT
    id, deck_id, front, back, example_sentence, hints, card_type,
    ease, interval, repetitions, due_at, lapsed, suspended,
    source, anki_card_id, anki_note_guid, template_ord, model_css,
    created_at, updated_at
FROM cards
WHERE anki_card_id = ?;

-- name: ListCardsByDeck :many
SELECT
    id, deck_id, front, back, example_sentence, hints, card_type,
    ease, interval, repetitions, due_at, lapsed, suspended,
    source, anki_card_id, anki_note_guid, template_ord, model_css,
    created_at, updated_at
FROM cards
WHERE deck_id = ?
ORDER BY created_at;

-- name: ListDueCards :many
SELECT
    id, deck_id, front, back, example_sentence, hints, card_type,
    ease, interval, repetitions, due_at, lapsed, suspended,
    source, anki_card_id, anki_note_guid, template_ord, model_css,
    created_at, updated_at
FROM cards
WHERE deck_id = ?
  AND suspended = 0
  AND due_at IS NOT NULL
  AND due_at <= ?
ORDER BY due_at
LIMIT ?;

-- name: ListNewCards :many
SELECT
    id, deck_id, front, back, example_sentence, hints, card_type,
    ease, interval, repetitions, due_at, lapsed, suspended,
    source, anki_card_id, anki_note_guid, template_ord, model_css,
    created_at, updated_at
FROM cards
WHERE deck_id = ?
  AND suspended = 0
  AND due_at IS NULL
ORDER BY created_at
LIMIT ?;

-- name: UpdateCardSRS :exec
UPDATE cards
SET ease = ?, interval = ?, repetitions = ?, due_at = ?, lapsed = ?, suspended = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateCardContent :exec
UPDATE cards
SET front = ?, back = ?, example_sentence = ?, hints = ?, card_type = ?, model_css = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteCard :exec
DELETE FROM cards
WHERE id = ?;

-- name: CountCardsByDeck :one
SELECT COUNT(*) AS count
FROM cards
WHERE deck_id = ?;

-- name: CreateReviewLog :exec
INSERT INTO review_log (id, card_id, rating, elapsed_ms, review_at, day_seq, ease_before, interval_before)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListReviewLogByCard :many
SELECT id, card_id, rating, elapsed_ms, review_at, day_seq, ease_before, interval_before
FROM review_log
WHERE card_id = ?
ORDER BY review_at DESC;
