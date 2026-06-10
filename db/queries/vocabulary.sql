-- name: CreateVocabulary :exec
INSERT INTO vocabulary (
    id, word, phonetic, pos, definition, definition_en, examples, level, source, anki_note_guid, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'));

-- name: GetVocabulary :one
SELECT id, word, phonetic, pos, definition, definition_en, examples, level, source, anki_note_guid, created_at
FROM vocabulary
WHERE id = ?;

-- name: GetVocabularyByAnkiNoteGuid :one
SELECT id, word, phonetic, pos, definition, definition_en, examples, level, source, anki_note_guid, created_at
FROM vocabulary
WHERE anki_note_guid = ?;

-- name: ListVocabulary :many
SELECT id, word, phonetic, pos, definition, definition_en, examples, level, source, anki_note_guid, created_at
FROM vocabulary
ORDER BY word;

-- name: SearchVocabulary :many
SELECT id, word, phonetic, pos, definition, definition_en, examples, level, source, anki_note_guid, created_at
FROM vocabulary
WHERE word LIKE '%' || ? || '%'
   OR definition LIKE '%' || ? || '%'
ORDER BY word
LIMIT ?;

-- name: ListVocabularyByLevel :many
SELECT id, word, phonetic, pos, definition, definition_en, examples, level, source, anki_note_guid, created_at
FROM vocabulary
WHERE level = ?
ORDER BY word;

-- name: UpdateVocabulary :exec
UPDATE vocabulary
SET word = ?, phonetic = ?, pos = ?, definition = ?, definition_en = ?, examples = ?, level = ?, source = ?
WHERE id = ?;

-- name: DeleteVocabulary :exec
DELETE FROM vocabulary
WHERE id = ?;

-- name: LinkCardVocab :exec
INSERT INTO card_vocab (card_id, vocab_id)
VALUES (?, ?)
ON CONFLICT(card_id, vocab_id) DO NOTHING;

-- name: ListVocabForCard :many
SELECT v.id, v.word, v.phonetic, v.pos, v.definition, v.definition_en, v.examples, v.level, v.source, v.anki_note_guid, v.created_at
FROM vocabulary v
INNER JOIN card_vocab cv ON cv.vocab_id = v.id
WHERE cv.card_id = ?
ORDER BY v.word;

-- name: ListCardsForVocab :many
SELECT
    c.id, c.deck_id, c.front, c.back, c.example_sentence, c.hints, c.card_type,
    c.ease, c.interval, c.repetitions, c.due_at, c.lapsed, c.suspended,
    c.source, c.anki_card_id, c.anki_note_guid, c.template_ord, c.model_css,
    c.created_at, c.updated_at
FROM cards c
INNER JOIN card_vocab cv ON cv.card_id = c.id
WHERE cv.vocab_id = ?
ORDER BY c.created_at;
