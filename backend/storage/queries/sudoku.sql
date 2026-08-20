-- name: GetSudokuGameBySessionID :one
SELECT id, session_id, difficulty, puzzle, solution, board, status, created_at, updated_at
FROM sudoku_games
WHERE session_id = ?;

-- name: CreateSudokuGame :exec
INSERT INTO sudoku_games (
    id,
    session_id,
    difficulty,
    puzzle,
    solution,
    board,
    status,
    created_at,
    updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'));

-- name: UpdateSudokuGameBoard :exec
UPDATE sudoku_games
SET board = ?, status = ?, updated_at = datetime('now')
WHERE session_id = ?;

-- name: ReplaceSudokuPuzzle :exec
UPDATE sudoku_games
SET difficulty = ?,
    puzzle = ?,
    solution = ?,
    board = ?,
    status = 'playing',
    updated_at = datetime('now')
WHERE session_id = ?;
