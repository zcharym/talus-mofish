package types

import "github.com/songwei.ma/talus-mofish/backend/storage/store"

// SudokuGame is the player-facing board state. The solution stays in Go.
type SudokuGame struct {
	ID         string `json:"id"`
	SessionID  string `json:"session_id"`
	Difficulty string `json:"difficulty"`
	Puzzle     string `json:"puzzle"`
	Board      string `json:"board"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// SudokuSession is a new or updated game plus its chat-like session row.
type SudokuSession struct {
	Session store.ChatSession `json:"session"`
	Game    SudokuGame        `json:"game"`
}

// SudokuCheckResult reports wrong cells after a Check action.
type SudokuCheckResult struct {
	Game      SudokuGame `json:"game"`
	Conflicts []int      `json:"conflicts"`
	Solved    bool       `json:"solved"`
}

// PublicSudokuGame copies a stored game without the solution.
func PublicSudokuGame(row store.SudokuGame) SudokuGame {
	return SudokuGame{
		ID:         row.ID,
		SessionID:  row.SessionID,
		Difficulty: row.Difficulty,
		Puzzle:     row.Puzzle,
		Board:      row.Board,
		Status:     row.Status,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
