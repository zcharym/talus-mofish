package sudoku

import (
	"fmt"
	"strings"
)

const (
	GridSize = 81

	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"

	StatusPlaying = "playing"
	StatusSolved  = "solved"

	SessionKindChat   = "chat"
	SessionKindSudoku = "sudoku"
)

// NormalizeDifficulty returns a valid YouDoSudoku difficulty.
func NormalizeDifficulty(difficulty string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(difficulty)) {
	case "", DifficultyEasy:
		return DifficultyEasy, nil
	case DifficultyMedium:
		return DifficultyMedium, nil
	case DifficultyHard:
		return DifficultyHard, nil
	default:
		return "", fmt.Errorf("difficulty must be easy, medium, or hard")
	}
}

// TitleForDifficulty returns the sidebar session title for a new game.
func TitleForDifficulty(difficulty string) string {
	switch difficulty {
	case DifficultyMedium:
		return "Sudoku · Medium"
	case DifficultyHard:
		return "Sudoku · Hard"
	default:
		return "Sudoku · Easy"
	}
}

// NormalizeGrid accepts an 81-character puzzle/solution/board and maps '.' to '0'.
func NormalizeGrid(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if len(raw) != GridSize {
		return "", fmt.Errorf("grid must be %d characters, got %d", GridSize, len(raw))
	}
	out := make([]byte, GridSize)
	for i := 0; i < GridSize; i++ {
		ch := raw[i]
		switch {
		case ch == '.' || ch == '0':
			out[i] = '0'
		case ch >= '1' && ch <= '9':
			out[i] = ch
		default:
			return "", fmt.Errorf("invalid grid character %q at index %d", ch, i)
		}
	}
	return string(out), nil
}

// IsGiven reports whether the cell at index is a clue from the original puzzle.
func IsGiven(puzzle string, index int) bool {
	if index < 0 || index >= len(puzzle) {
		return false
	}
	return puzzle[index] >= '1' && puzzle[index] <= '9'
}

// SetCell writes value (0-9) into board at index, rejecting given clues.
func SetCell(puzzle, board string, index, value int) (string, error) {
	if index < 0 || index >= GridSize {
		return "", fmt.Errorf("cell index must be 0-80")
	}
	if value < 0 || value > 9 {
		return "", fmt.Errorf("cell value must be 0-9")
	}
	if len(puzzle) != GridSize || len(board) != GridSize {
		return "", fmt.Errorf("puzzle and board must be %d characters", GridSize)
	}
	if IsGiven(puzzle, index) {
		return "", fmt.Errorf("cannot edit a given clue")
	}
	out := []byte(board)
	if value == 0 {
		out[index] = '0'
	} else {
		out[index] = byte('0' + value)
	}
	return string(out), nil
}

// Conflicts returns indexes that are filled and do not match the solution.
func Conflicts(board, solution string) []int {
	conflicts := make([]int, 0)
	n := min(len(board), len(solution), GridSize)
	for i := 0; i < n; i++ {
		if board[i] == '0' {
			continue
		}
		if board[i] != solution[i] {
			conflicts = append(conflicts, i)
		}
	}
	return conflicts
}

// IsComplete reports whether every cell matches the solution.
func IsComplete(board, solution string) bool {
	return board == solution && len(board) == GridSize
}

// StatusForBoard returns solved when the board matches the solution.
func StatusForBoard(board, solution string) string {
	if IsComplete(board, solution) {
		return StatusSolved
	}
	return StatusPlaying
}
