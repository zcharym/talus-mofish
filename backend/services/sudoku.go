package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/storage/store"
	"github.com/songwei.ma/talus-mofish/backend/sudoku"
	"github.com/songwei.ma/talus-mofish/backend/types"
)

// SudokuService exposes YouDoSudoku-backed games as Agent window sessions.
type SudokuService struct {
	db     *storage.DB
	config *storage.ConfigStore
	client *sudoku.Client
}

// NewSudokuService creates the Sudoku Wails service.
func NewSudokuService(db *storage.DB, cfg *storage.ConfigStore) *SudokuService {
	return &SudokuService{
		db:     db,
		config: cfg,
	}
}

func (s *SudokuService) apiClient() *sudoku.Client {
	if s.client != nil {
		return s.client
	}
	return sudoku.NewClient(s.config.Get().Sudoku.APIKey)
}

func (s *SudokuService) fetchPuzzle(difficulty string) (sudoku.Puzzle, error) {
	ctx := context.Background()
	puzzle, err := s.apiClient().Generate(ctx, difficulty)
	if err != nil {
		return sudoku.Puzzle{}, err
	}
	return puzzle, nil
}

func (s *SudokuService) loadGame(ctx context.Context, sessionID string) (store.ChatSession, store.SudokuGame, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return store.ChatSession{}, store.SudokuGame{}, fmt.Errorf("session id is required")
	}
	session, err := s.db.Queries.GetChatSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ChatSession{}, store.SudokuGame{}, fmt.Errorf("session not found")
		}
		return store.ChatSession{}, store.SudokuGame{}, fmt.Errorf("get chat session: %w", err)
	}
	if session.Kind != sudoku.SessionKindSudoku {
		return store.ChatSession{}, store.SudokuGame{}, fmt.Errorf("session is not a Sudoku game")
	}
	game, err := s.db.Queries.GetSudokuGameBySessionID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ChatSession{}, store.SudokuGame{}, fmt.Errorf("sudoku game not found")
		}
		return store.ChatSession{}, store.SudokuGame{}, fmt.Errorf("get sudoku game: %w", err)
	}
	return session, game, nil
}

func (s *SudokuService) reloadPublic(ctx context.Context, sessionID string) (types.SudokuSession, error) {
	session, game, err := s.loadGame(ctx, sessionID)
	if err != nil {
		return types.SudokuSession{}, err
	}
	return types.SudokuSession{
		Session: session,
		Game:    types.PublicSudokuGame(game),
	}, nil
}

// NewSudokuGame fetches a puzzle and creates a sidebar session.
func (s *SudokuService) NewSudokuGame(difficulty string) (types.SudokuSession, error) {
	ctx := context.Background()
	puzzle, err := s.fetchPuzzle(difficulty)
	if err != nil {
		return types.SudokuSession{}, err
	}

	sessionID := uuid.NewString()
	if err := s.db.Queries.CreateChatSession(ctx, store.CreateChatSessionParams{
		ID:    sessionID,
		Title: sudoku.TitleForDifficulty(puzzle.Difficulty),
		Kind:  sudoku.SessionKindSudoku,
	}); err != nil {
		return types.SudokuSession{}, fmt.Errorf("create sudoku session: %w", err)
	}
	if err := s.db.Queries.CreateSudokuGame(ctx, store.CreateSudokuGameParams{
		ID:         uuid.NewString(),
		SessionID:  sessionID,
		Difficulty: puzzle.Difficulty,
		Puzzle:     puzzle.Puzzle,
		Solution:   puzzle.Solution,
		Board:      puzzle.Puzzle,
		Status:     sudoku.StatusPlaying,
	}); err != nil {
		return types.SudokuSession{}, fmt.Errorf("create sudoku game: %w", err)
	}

	return s.reloadPublic(ctx, sessionID)
}

// GetSudokuGame returns the current board for a Sudoku session.
func (s *SudokuService) GetSudokuGame(sessionID string) (types.SudokuGame, error) {
	ctx := context.Background()
	_, game, err := s.loadGame(ctx, sessionID)
	if err != nil {
		return types.SudokuGame{}, err
	}
	return types.PublicSudokuGame(game), nil
}

// SetSudokuCell writes a digit (1-9) or 0 to clear. Given clues cannot be edited.
func (s *SudokuService) SetSudokuCell(sessionID string, index, value int) (types.SudokuGame, error) {
	ctx := context.Background()
	_, game, err := s.loadGame(ctx, sessionID)
	if err != nil {
		return types.SudokuGame{}, err
	}
	if game.Status == sudoku.StatusSolved {
		return types.SudokuGame{}, fmt.Errorf("puzzle is already solved")
	}

	board, err := sudoku.SetCell(game.Puzzle, game.Board, index, value)
	if err != nil {
		return types.SudokuGame{}, err
	}
	status := sudoku.StatusForBoard(board, game.Solution)
	if err := s.db.Queries.UpdateSudokuGameBoard(ctx, store.UpdateSudokuGameBoardParams{
		Board:     board,
		Status:    status,
		SessionID: sessionID,
	}); err != nil {
		return types.SudokuGame{}, fmt.Errorf("save sudoku cell: %w", err)
	}
	if err := s.db.Queries.TouchChatSession(ctx, sessionID); err != nil {
		return types.SudokuGame{}, fmt.Errorf("touch sudoku session: %w", err)
	}

	_, game, err = s.loadGame(ctx, sessionID)
	if err != nil {
		return types.SudokuGame{}, err
	}
	return types.PublicSudokuGame(game), nil
}

// CheckSudokuGame returns indexes that do not match the solution.
func (s *SudokuService) CheckSudokuGame(sessionID string) (types.SudokuCheckResult, error) {
	ctx := context.Background()
	_, game, err := s.loadGame(ctx, sessionID)
	if err != nil {
		return types.SudokuCheckResult{}, err
	}

	conflicts := sudoku.Conflicts(game.Board, game.Solution)
	status := game.Status
	if sudoku.IsComplete(game.Board, game.Solution) {
		status = sudoku.StatusSolved
		if game.Status != sudoku.StatusSolved {
			if err := s.db.Queries.UpdateSudokuGameBoard(ctx, store.UpdateSudokuGameBoardParams{
				Board:     game.Board,
				Status:    status,
				SessionID: sessionID,
			}); err != nil {
				return types.SudokuCheckResult{}, fmt.Errorf("mark sudoku solved: %w", err)
			}
			if err := s.db.Queries.TouchChatSession(ctx, sessionID); err != nil {
				return types.SudokuCheckResult{}, fmt.Errorf("touch sudoku session: %w", err)
			}
			_, game, err = s.loadGame(ctx, sessionID)
			if err != nil {
				return types.SudokuCheckResult{}, err
			}
		}
	}

	if conflicts == nil {
		conflicts = []int{}
	}
	return types.SudokuCheckResult{
		Game:      types.PublicSudokuGame(game),
		Conflicts: conflicts,
		Solved:    status == sudoku.StatusSolved,
	}, nil
}

// NewSudokuPuzzle replaces the board on an existing Sudoku session.
func (s *SudokuService) NewSudokuPuzzle(sessionID string, difficulty string) (types.SudokuGame, error) {
	ctx := context.Background()
	session, _, err := s.loadGame(ctx, sessionID)
	if err != nil {
		return types.SudokuGame{}, err
	}

	puzzle, err := s.fetchPuzzle(difficulty)
	if err != nil {
		return types.SudokuGame{}, err
	}

	if err := s.db.Queries.ReplaceSudokuPuzzle(ctx, store.ReplaceSudokuPuzzleParams{
		Difficulty: puzzle.Difficulty,
		Puzzle:     puzzle.Puzzle,
		Solution:   puzzle.Solution,
		Board:      puzzle.Puzzle,
		SessionID:  sessionID,
	}); err != nil {
		return types.SudokuGame{}, fmt.Errorf("replace sudoku puzzle: %w", err)
	}

	title := sudoku.TitleForDifficulty(puzzle.Difficulty)
	if session.Title == "" || strings.HasPrefix(session.Title, "Sudoku ·") {
		if err := s.db.Queries.RenameChatSession(ctx, store.RenameChatSessionParams{
			Title: title,
			ID:    sessionID,
		}); err != nil {
			return types.SudokuGame{}, fmt.Errorf("rename sudoku session: %w", err)
		}
	} else if err := s.db.Queries.TouchChatSession(ctx, sessionID); err != nil {
		return types.SudokuGame{}, fmt.Errorf("touch sudoku session: %w", err)
	}

	_, game, err := s.loadGame(ctx, sessionID)
	if err != nil {
		return types.SudokuGame{}, err
	}
	return types.PublicSudokuGame(game), nil
}
