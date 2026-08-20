package sudoku

import (
	"strings"
	"testing"
)

func TestNormalizeGridAcceptsDotsAndZeros(t *testing.T) {
	t.Parallel()

	raw := strings.Repeat(".", 80) + "1"
	got, err := NormalizeGrid(raw)
	if err != nil {
		t.Fatalf("NormalizeGrid: %v", err)
	}
	if got[80] != '1' {
		t.Fatalf("clue not preserved: %q", got[80])
	}
	if strings.Count(got, "0") != 80 {
		t.Fatalf("expected 80 empties, got %q", got)
	}
}

func TestSetCellRejectsGivenClues(t *testing.T) {
	t.Parallel()

	puzzle := strings.Repeat("0", 80) + "5"
	board := puzzle
	if _, err := SetCell(puzzle, board, 80, 1); err == nil {
		t.Fatal("expected error editing given clue")
	}
	next, err := SetCell(puzzle, board, 0, 9)
	if err != nil {
		t.Fatalf("SetCell empty: %v", err)
	}
	if next[0] != '9' {
		t.Fatalf("cell 0 = %q, want 9", next[0])
	}
}

func TestConflictsAndComplete(t *testing.T) {
	t.Parallel()

	solution := strings.Repeat("123456789", 9)
	board := solution
	if !IsComplete(board, solution) {
		t.Fatal("expected complete board")
	}
	wrong := []byte(board)
	wrong[0] = '9'
	conflicts := Conflicts(string(wrong), solution)
	if len(conflicts) != 1 || conflicts[0] != 0 {
		t.Fatalf("conflicts = %v, want [0]", conflicts)
	}
	if StatusForBoard(string(wrong), solution) != StatusPlaying {
		t.Fatal("wrong board should still be playing")
	}
	if StatusForBoard(solution, solution) != StatusSolved {
		t.Fatal("matching board should be solved")
	}
}

func TestNormalizeDifficulty(t *testing.T) {
	t.Parallel()

	got, err := NormalizeDifficulty("HARD")
	if err != nil || got != DifficultyHard {
		t.Fatalf("got %q err %v", got, err)
	}
	if _, err := NormalizeDifficulty("expert"); err == nil {
		t.Fatal("expected invalid difficulty error")
	}
}
