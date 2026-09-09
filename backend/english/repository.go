package english

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/songwei.ma/talus-mofish/backend/consts"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/storage/store"
	"github.com/songwei.ma/talus-mofish/backend/types"
)

const defaultSearchLimit = 50

// Repository owns English Learning persistence. The Wails service stays a thin RPC façade.
type Repository struct {
	db *storage.DB
}

// NewRepository wraps the application SQLite handle.
func NewRepository(db *storage.DB) *Repository {
	return &Repository{db: db}
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

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func vocabularyDTO(row store.Vocabulary) types.Vocabulary {
	return types.Vocabulary{
		ID:           row.ID,
		Word:         row.Word,
		Phonetic:     row.Phonetic,
		Pos:          row.Pos,
		Definition:   row.Definition,
		DefinitionEn: row.DefinitionEn,
		Examples:     row.Examples,
		Level:        row.Level,
		Source:       row.Source,
		AnkiNoteGUID: nullString(row.AnkiNoteGuid),
		CreatedAt:    row.CreatedAt,
	}
}

func vocabulariesDTO(rows []store.Vocabulary) []types.Vocabulary {
	out := make([]types.Vocabulary, 0, len(rows))
	for _, row := range rows {
		out = append(out, vocabularyDTO(row))
	}
	return out
}

func articleDTO(row store.Article) types.Article {
	return types.Article{
		ID:          row.ID,
		Title:       row.Title,
		URL:         row.Url,
		Content:     row.Content,
		Translation: row.Translation,
		Level:       row.Level,
		Source:      row.Source,
		WordCount:   row.WordCount,
		ModelCSS:    row.ModelCss,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func cardDTO(row store.Card) types.Card {
	return types.Card{
		ID:              row.ID,
		DeckID:          row.DeckID,
		Front:           row.Front,
		Back:            row.Back,
		ExampleSentence: row.ExampleSentence,
		Hints:           row.Hints,
		CardType:        row.CardType,
		ModelCSS:        row.ModelCss,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func cardsDTO(rows []store.Card) []types.Card {
	out := make([]types.Card, 0, len(rows))
	for _, row := range rows {
		out = append(out, cardDTO(row))
	}
	return out
}

func deckDTO(row store.Deck) types.Deck {
	return types.Deck{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		SortOrder:   row.SortOrder,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func importDTO(row store.AnkiImport) types.AnkiImportRecord {
	return types.AnkiImportRecord{
		ID:           row.ID,
		Filename:     row.Filename,
		ImportedAt:   row.ImportedAt,
		Status:       row.Status,
		ErrorMessage: row.ErrorMessage,
		StatsJSON:    row.StatsJson,
	}
}

// ListAnkiImports returns past import sessions.
func (r *Repository) ListAnkiImports(ctx context.Context) ([]types.AnkiImportRecord, error) {
	items, err := r.db.Queries.ListAnkiImports(ctx)
	if err != nil {
		return nil, fmt.Errorf("list imports: %w", err)
	}
	out := make([]types.AnkiImportRecord, 0, len(items))
	for _, row := range items {
		out = append(out, importDTO(row))
	}
	return out, nil
}

// ListArticlesPage returns a paginated list of article summaries.
func (r *Repository) ListArticlesPage(ctx context.Context, page, pageSize int64) (types.ArticlePageResult, error) {
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := r.db.Queries.CountArticles(ctx)
	if err != nil {
		return types.ArticlePageResult{}, fmt.Errorf("count articles: %w", err)
	}

	rows, err := r.db.Queries.ListArticlesPage(ctx, store.ListArticlesPageParams{
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
func (r *Repository) GetArticle(ctx context.Context, id string) (types.Article, error) {
	article, err := r.db.Queries.GetArticle(ctx, id)
	if err != nil {
		return types.Article{}, fmt.Errorf("get article: %w", err)
	}
	return articleDTO(article), nil
}

// ListVocabularyPage returns a paginated list of vocabulary entries.
func (r *Repository) ListVocabularyPage(ctx context.Context, page, pageSize int64) (types.VocabularyPageResult, error) {
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := r.db.Queries.CountVocabulary(ctx)
	if err != nil {
		return types.VocabularyPageResult{}, fmt.Errorf("count vocabulary: %w", err)
	}

	items, err := r.db.Queries.ListVocabularyPage(ctx, store.ListVocabularyPageParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return types.VocabularyPageResult{}, fmt.Errorf("list vocabulary page: %w", err)
	}

	return types.VocabularyPageResult{
		Items:    vocabulariesDTO(items),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetVocabulary returns a single vocabulary entry by ID.
func (r *Repository) GetVocabulary(ctx context.Context, id string) (types.Vocabulary, error) {
	item, err := r.db.Queries.GetVocabulary(ctx, id)
	if err != nil {
		return types.Vocabulary{}, fmt.Errorf("get vocabulary: %w", err)
	}
	return vocabularyDTO(item), nil
}

// UpdateVocabulary saves vocabulary field changes.
func (r *Repository) UpdateVocabulary(ctx context.Context, input types.VocabularyUpdate) error {
	input.Word = strings.TrimSpace(input.Word)
	if input.Word == "" {
		return fmt.Errorf("word is required")
	}
	if err := r.db.Queries.UpdateVocabulary(ctx, store.UpdateVocabularyParams{
		Word:         input.Word,
		Phonetic:     input.Phonetic,
		Pos:          input.Pos,
		Definition:   input.Definition,
		DefinitionEn: input.DefinitionEn,
		Examples:     input.Examples,
		Level:        input.Level,
		Source:       input.Source,
		ID:           input.ID,
	}); err != nil {
		return fmt.Errorf("update vocabulary: %w", err)
	}
	return nil
}

// DeleteVocabulary removes a vocabulary entry.
func (r *Repository) DeleteVocabulary(ctx context.Context, id string) error {
	if err := r.db.Queries.DeleteVocabulary(ctx, id); err != nil {
		return fmt.Errorf("delete vocabulary: %w", err)
	}
	return nil
}

// SearchVocabulary finds vocabulary entries matching a query string.
func (r *Repository) SearchVocabulary(ctx context.Context, query string, limit int64) ([]types.Vocabulary, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []types.Vocabulary{}, nil
	}
	if limit <= 0 {
		limit = defaultSearchLimit
	}

	items, err := r.db.Queries.SearchVocabulary(ctx, store.SearchVocabularyParams{
		Column1: sql.NullString{String: query, Valid: true},
		Column2: sql.NullString{String: query, Valid: true},
		Limit:   limit,
	})
	if err != nil {
		return nil, fmt.Errorf("search vocabulary: %w", err)
	}
	return vocabulariesDTO(items), nil
}

// ListCardsForVocab returns SRS cards linked to a vocabulary entry.
func (r *Repository) ListCardsForVocab(ctx context.Context, vocabID string) ([]types.Card, error) {
	items, err := r.db.Queries.ListCardsForVocab(ctx, vocabID)
	if err != nil {
		return nil, fmt.Errorf("list cards for vocab: %w", err)
	}
	return cardsDTO(items), nil
}

// ListDecks returns all SRS decks.
func (r *Repository) ListDecks(ctx context.Context) ([]types.Deck, error) {
	items, err := r.db.Queries.ListDecks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list decks: %w", err)
	}
	out := make([]types.Deck, 0, len(items))
	for _, row := range items {
		out = append(out, deckDTO(row))
	}
	return out, nil
}

// ListCardsByDeck returns cards in a deck.
func (r *Repository) ListCardsByDeck(ctx context.Context, deckID string) ([]types.Card, error) {
	items, err := r.db.Queries.ListCardsByDeck(ctx, deckID)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	return cardsDTO(items), nil
}

// GetCard returns a single SRS card by ID.
func (r *Repository) GetCard(ctx context.Context, id string) (types.Card, error) {
	item, err := r.db.Queries.GetCard(ctx, id)
	if err != nil {
		return types.Card{}, fmt.Errorf("get card: %w", err)
	}
	return cardDTO(item), nil
}

// UpdateCardContent saves editable card fields after HTML sanitization.
func (r *Repository) UpdateCardContent(ctx context.Context, input types.CardContentUpdate) error {
	if err := r.db.Queries.UpdateCardContent(ctx, store.UpdateCardContentParams{
		Front:           SanitizeHTML(input.Front),
		Back:            SanitizeHTML(input.Back),
		ExampleSentence: input.ExampleSentence,
		Hints:           input.Hints,
		CardType:        input.CardType,
		ModelCss:        input.ModelCSS,
		ID:              input.ID,
	}); err != nil {
		return fmt.Errorf("update card content: %w", err)
	}
	return nil
}

// DeleteCard removes an SRS card.
func (r *Repository) DeleteCard(ctx context.Context, id string) error {
	if err := r.db.Queries.DeleteCard(ctx, id); err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	return nil
}
