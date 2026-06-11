package appservice

import (
	"context"
	"fmt"

	"github.com/songwei.ma/talus-mofish/internal/store"
)

// ListArticlesPage returns a paginated list of article summaries.
func (s *Service) ListArticlesPage(page, pageSize int64) (store.ArticlePageResult, error) {
	ctx := context.Background()
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := s.db.Queries.CountArticles(ctx)
	if err != nil {
		return store.ArticlePageResult{}, fmt.Errorf("count articles: %w", err)
	}

	rows, err := s.db.Queries.ListArticlesPage(ctx, store.ListArticlesPageParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return store.ArticlePageResult{}, fmt.Errorf("list articles page: %w", err)
	}

	items := make([]store.ArticleSummary, 0, len(rows))
	for _, row := range rows {
		items = append(items, store.ArticleSummary{
			ID:        row.ID,
			Title:     row.Title,
			Source:    row.Source,
			WordCount: row.WordCount,
			CreatedAt: row.CreatedAt,
		})
	}

	return store.ArticlePageResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetArticle returns a single article by ID.
func (s *Service) GetArticle(id string) (store.Article, error) {
	ctx := context.Background()
	article, err := s.db.Queries.GetArticle(ctx, id)
	if err != nil {
		return store.Article{}, fmt.Errorf("get article: %w", err)
	}
	return article, nil
}
