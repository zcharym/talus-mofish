package appservice

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/songwei.ma/talus-mofish/internal/store"
)

const defaultSearchLimit = 50

// ListVocabularyPage returns a paginated list of vocabulary entries.
func (s *Service) ListVocabularyPage(page, pageSize int64) (store.VocabularyPageResult, error) {
	ctx := context.Background()
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := s.db.Queries.CountVocabulary(ctx)
	if err != nil {
		return store.VocabularyPageResult{}, fmt.Errorf("count vocabulary: %w", err)
	}

	items, err := s.db.Queries.ListVocabularyPage(ctx, store.ListVocabularyPageParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return store.VocabularyPageResult{}, fmt.Errorf("list vocabulary page: %w", err)
	}
	if items == nil {
		items = []store.Vocabulary{}
	}

	return store.VocabularyPageResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetVocabulary returns a single vocabulary entry by ID.
func (s *Service) GetVocabulary(id string) (store.Vocabulary, error) {
	ctx := context.Background()
	item, err := s.db.Queries.GetVocabulary(ctx, id)
	if err != nil {
		return store.Vocabulary{}, fmt.Errorf("get vocabulary: %w", err)
	}
	return item, nil
}

// UpdateVocabulary saves vocabulary field changes.
func (s *Service) UpdateVocabulary(input store.UpdateVocabularyParams) error {
	ctx := context.Background()
	input.Word = strings.TrimSpace(input.Word)
	if input.Word == "" {
		return fmt.Errorf("word is required")
	}
	if err := s.db.Queries.UpdateVocabulary(ctx, input); err != nil {
		return fmt.Errorf("update vocabulary: %w", err)
	}
	return nil
}

// DeleteVocabulary removes a vocabulary entry.
func (s *Service) DeleteVocabulary(id string) error {
	ctx := context.Background()
	if err := s.db.Queries.DeleteVocabulary(ctx, id); err != nil {
		return fmt.Errorf("delete vocabulary: %w", err)
	}
	return nil
}

// SearchVocabulary finds vocabulary entries matching a query string.
func (s *Service) SearchVocabulary(query string, limit int64) ([]store.Vocabulary, error) {
	ctx := context.Background()
	query = strings.TrimSpace(query)
	if query == "" {
		return []store.Vocabulary{}, nil
	}
	if limit <= 0 {
		limit = defaultSearchLimit
	}

	items, err := s.db.Queries.SearchVocabulary(ctx, store.SearchVocabularyParams{
		Column1: sql.NullString{String: query, Valid: true},
		Column2: sql.NullString{String: query, Valid: true},
		Limit:   limit,
	})
	if err != nil {
		return nil, fmt.Errorf("search vocabulary: %w", err)
	}
	if items == nil {
		return []store.Vocabulary{}, nil
	}
	return items, nil
}

// ListCardsForVocab returns SRS cards linked to a vocabulary entry.
func (s *Service) ListCardsForVocab(vocabID string) ([]store.Card, error) {
	ctx := context.Background()
	items, err := s.db.Queries.ListCardsForVocab(ctx, vocabID)
	if err != nil {
		return nil, fmt.Errorf("list cards for vocab: %w", err)
	}
	if items == nil {
		return []store.Card{}, nil
	}
	return items, nil
}
