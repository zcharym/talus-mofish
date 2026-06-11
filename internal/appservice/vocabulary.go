package appservice

import (
	"context"
	"fmt"

	"github.com/songwei.ma/talus-mofish/internal/store"
)

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
