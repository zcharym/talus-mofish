package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/songwei.ma/talus-mofish/backend/consts"
	"github.com/songwei.ma/talus-mofish/backend/english/content"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/storage/store"
	"github.com/songwei.ma/talus-mofish/backend/types"
)

const defaultSearchLimit = 50

// EnglishService exposes English Learning domain APIs.
type EnglishService struct {
	db *storage.DB
}

// NewEnglishService creates the English Learning Wails service.
func NewEnglishService(db *storage.DB) *EnglishService {
	return &EnglishService{db: db}
}

func normalizePageParams(page, pageSize int64) (int64, int64) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = consts.DefaultPageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// PreviewAnkiAPKG returns deck and model metadata for import configuration.
func (s *EnglishService) PreviewAnkiAPKG(apkgPath string) (content.AnkiPreview, error) {
	importer, err := content.NewImporter(s.db)
	if err != nil {
		return content.AnkiPreview{}, fmt.Errorf("init importer: %w", err)
	}
	preview, err := importer.Preview(apkgPath)
	if err != nil {
		return content.AnkiPreview{}, fmt.Errorf("preview apkg: %w", err)
	}
	return preview, nil
}

// ImportAnkiAPKG imports an APKG using per-deck configuration from the UI.
func (s *EnglishService) ImportAnkiAPKG(apkgPath string, configs []content.ImportDeckConfig) (content.ImportResult, error) {
	importer, err := content.NewImporter(s.db)
	if err != nil {
		return content.ImportResult{}, fmt.Errorf("init importer: %w", err)
	}
	result, err := importer.Import(apkgPath, configs)
	if err != nil {
		return content.ImportResult{}, fmt.Errorf("import apkg: %w", err)
	}
	return result, nil
}

// ListAnkiImports returns past import sessions.
func (s *EnglishService) ListAnkiImports() ([]store.AnkiImport, error) {
	ctx := context.Background()
	items, err := s.db.Queries.ListAnkiImports(ctx)
	if err != nil {
		return nil, fmt.Errorf("list imports: %w", err)
	}
	if items == nil {
		return []store.AnkiImport{}, nil
	}
	return items, nil
}

// ListArticlesPage returns a paginated list of article summaries.
func (s *EnglishService) ListArticlesPage(page, pageSize int64) (types.ArticlePageResult, error) {
	ctx := context.Background()
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := s.db.Queries.CountArticles(ctx)
	if err != nil {
		return types.ArticlePageResult{}, fmt.Errorf("count articles: %w", err)
	}

	rows, err := s.db.Queries.ListArticlesPage(ctx, store.ListArticlesPageParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return types.ArticlePageResult{}, fmt.Errorf("list articles page: %w", err)
	}

	items := make([]types.ArticleSummary, 0, len(rows))
	for _, row := range rows {
		items = append(items, types.ArticleSummary{
			ID:        row.ID,
			Title:     row.Title,
			Source:    row.Source,
			WordCount: row.WordCount,
			CreatedAt: row.CreatedAt,
		})
	}

	return types.ArticlePageResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetArticle returns a single article by ID.
func (s *EnglishService) GetArticle(id string) (store.Article, error) {
	ctx := context.Background()
	article, err := s.db.Queries.GetArticle(ctx, id)
	if err != nil {
		return store.Article{}, fmt.Errorf("get article: %w", err)
	}
	return article, nil
}

// ListVocabularyPage returns a paginated list of vocabulary entries.
func (s *EnglishService) ListVocabularyPage(page, pageSize int64) (types.VocabularyPageResult, error) {
	ctx := context.Background()
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := s.db.Queries.CountVocabulary(ctx)
	if err != nil {
		return types.VocabularyPageResult{}, fmt.Errorf("count vocabulary: %w", err)
	}

	items, err := s.db.Queries.ListVocabularyPage(ctx, store.ListVocabularyPageParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return types.VocabularyPageResult{}, fmt.Errorf("list vocabulary page: %w", err)
	}
	if items == nil {
		items = []store.Vocabulary{}
	}

	return types.VocabularyPageResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetVocabulary returns a single vocabulary entry by ID.
func (s *EnglishService) GetVocabulary(id string) (store.Vocabulary, error) {
	ctx := context.Background()
	item, err := s.db.Queries.GetVocabulary(ctx, id)
	if err != nil {
		return store.Vocabulary{}, fmt.Errorf("get vocabulary: %w", err)
	}
	return item, nil
}

// UpdateVocabulary saves vocabulary field changes.
func (s *EnglishService) UpdateVocabulary(input store.UpdateVocabularyParams) error {
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
func (s *EnglishService) DeleteVocabulary(id string) error {
	ctx := context.Background()
	if err := s.db.Queries.DeleteVocabulary(ctx, id); err != nil {
		return fmt.Errorf("delete vocabulary: %w", err)
	}
	return nil
}

// SearchVocabulary finds vocabulary entries matching a query string.
func (s *EnglishService) SearchVocabulary(query string, limit int64) ([]store.Vocabulary, error) {
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
func (s *EnglishService) ListCardsForVocab(vocabID string) ([]store.Card, error) {
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

// ListDecks returns all SRS decks.
func (s *EnglishService) ListDecks() ([]store.Deck, error) {
	ctx := context.Background()
	items, err := s.db.Queries.ListDecks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list decks: %w", err)
	}
	if items == nil {
		return []store.Deck{}, nil
	}
	return items, nil
}

// ListCardsByDeck returns cards in a deck.
func (s *EnglishService) ListCardsByDeck(deckID string) ([]store.Card, error) {
	ctx := context.Background()
	items, err := s.db.Queries.ListCardsByDeck(ctx, deckID)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	if items == nil {
		return []store.Card{}, nil
	}
	return items, nil
}

// GetCard returns a single SRS card by ID.
func (s *EnglishService) GetCard(id string) (store.Card, error) {
	ctx := context.Background()
	item, err := s.db.Queries.GetCard(ctx, id)
	if err != nil {
		return store.Card{}, fmt.Errorf("get card: %w", err)
	}
	return item, nil
}

// UpdateCardContent saves editable card fields.
func (s *EnglishService) UpdateCardContent(input store.UpdateCardContentParams) error {
	ctx := context.Background()
	if err := s.db.Queries.UpdateCardContent(ctx, input); err != nil {
		return fmt.Errorf("update card content: %w", err)
	}
	return nil
}

// DeleteCard removes an SRS card.
func (s *EnglishService) DeleteCard(id string) error {
	ctx := context.Background()
	if err := s.db.Queries.DeleteCard(ctx, id); err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	return nil
}
