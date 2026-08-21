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

const (
	defaultSearchLimit = 50
	maxPageSize        = 100
)

// EnglishService exposes English Learning domain APIs.
// This service handles vocabulary management, article storage, SRS card operations,
// and Anki import functionality for the English Learning domain.
type EnglishService struct {
	db *storage.DB
}

// NewEnglishService creates a new English Learning Wails service.
// The service requires a database connection for all operations.
func NewEnglishService(db *storage.DB) *EnglishService {
	return &EnglishService{db: db}
}

// normalizePageParams ensures page and pageSize are within valid ranges.
// Page numbers start at 1, and page sizes are capped at 100 for performance.
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

// PreviewAnkiAPKG analyzes an APKG file and returns metadata for import configuration.
// This includes deck information, note types, and field mappings without importing any data.
// The preview helps users configure the import before committing to it.
func (s *EnglishService) PreviewAnkiAPKG(apkgPath string) (content.AnkiPreview, error) {
	apkgPath = strings.TrimSpace(apkgPath)
	if apkgPath == "" {
		return content.AnkiPreview{}, fmt.Errorf("apkg file path is required")
	}
	
	importer, err := content.NewImporter(s.db)
	if err != nil {
		return content.AnkiPreview{}, fmt.Errorf("failed to initialize importer: %w", err)
	}
	
	preview, err := importer.Preview(apkgPath)
	if err != nil {
		return content.AnkiPreview{}, fmt.Errorf("failed to preview apkg file %q: %w", apkgPath, err)
	}
	return preview, nil
}

// ImportAnkiAPKG imports an APKG file using per-deck configuration from the UI.
// Each deck can have custom field mappings configured by the user.
// Returns a summary of imported items including counts and any errors encountered.
func (s *EnglishService) ImportAnkiAPKG(apkgPath string, configs []content.ImportDeckConfig) (content.ImportResult, error) {
	apkgPath = strings.TrimSpace(apkgPath)
	if apkgPath == "" {
		return content.ImportResult{}, fmt.Errorf("apkg file path is required")
	}
	if len(configs) == 0 {
		return content.ImportResult{}, fmt.Errorf("at least one deck configuration is required")
	}
	
	importer, err := content.NewImporter(s.db)
	if err != nil {
		return content.ImportResult{}, fmt.Errorf("failed to initialize importer: %w", err)
	}
	
	result, err := importer.Import(apkgPath, configs)
	if err != nil {
		return content.ImportResult{}, fmt.Errorf("failed to import apkg file %q: %w", apkgPath, err)
	}
	return result, nil
}

// ListAnkiImports returns a history of past APKG import sessions.
// Each entry includes metadata about what was imported and when.
// Returns an empty slice if no imports have been performed.
func (s *EnglishService) ListAnkiImports() ([]store.AnkiImport, error) {
	ctx := context.Background()
	items, err := s.db.Queries.ListAnkiImports(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list import history: %w", err)
	}
	if items == nil {
		return []store.AnkiImport{}, nil
	}
	return items, nil
}

// ListArticlesPage returns a paginated list of article summaries.
// Articles are returned in reverse chronological order (newest first).
// The result includes pagination metadata for navigation.
func (s *EnglishService) ListArticlesPage(page, pageSize int64) (types.ArticlePageResult, error) {
	ctx := context.Background()
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := s.db.Queries.CountArticles(ctx)
	if err != nil {
		return types.ArticlePageResult{}, fmt.Errorf("failed to count articles: %w", err)
	}

	rows, err := s.db.Queries.ListArticlesPage(ctx, store.ListArticlesPageParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return types.ArticlePageResult{}, fmt.Errorf("failed to fetch articles page: %w", err)
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

// GetArticle retrieves the full content of an article by its ID.
// Returns an error if the article does not exist.
func (s *EnglishService) GetArticle(id string) (store.Article, error) {
	ctx := context.Background()
	
	id = strings.TrimSpace(id)
	if id == "" {
		return store.Article{}, fmt.Errorf("article id is required")
	}
	
	article, err := s.db.Queries.GetArticle(ctx, id)
	if err != nil {
		return store.Article{}, fmt.Errorf("failed to get article %q: %w", id, err)
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

// UpdateVocabulary persists changes to vocabulary entry fields.
// The word field is required and must not be empty after trimming.
func (s *EnglishService) UpdateVocabulary(input store.UpdateVocabularyParams) error {
	ctx := context.Background()
	
	input.Word = strings.TrimSpace(input.Word)
	if input.Word == "" {
		return fmt.Errorf("word field cannot be empty")
	}
	
	if err := s.db.Queries.UpdateVocabulary(ctx, input); err != nil {
		return fmt.Errorf("failed to update vocabulary entry: %w", err)
	}
	return nil
}

// DeleteVocabulary permanently removes a vocabulary entry from the database.
// Associated SRS cards are not automatically deleted; handle them separately if needed.
func (s *EnglishService) DeleteVocabulary(id string) error {
	ctx := context.Background()
	
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("vocabulary id is required")
	}
	
	if err := s.db.Queries.DeleteVocabulary(ctx, id); err != nil {
		return fmt.Errorf("failed to delete vocabulary entry %q: %w", id, err)
	}
	return nil
}

// SearchVocabulary finds vocabulary entries matching a search query.
// The search looks in both the word and definition fields.
// Returns an empty slice if no matches are found or if the query is empty.
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
		return nil, fmt.Errorf("failed to search vocabulary with query %q: %w", query, err)
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
