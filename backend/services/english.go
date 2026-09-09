package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/songwei.ma/talus-mofish/backend/english"
	"github.com/songwei.ma/talus-mofish/backend/english/content"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/types"
)

// EnglishService exposes English Learning domain APIs.
type EnglishService struct {
	db   *storage.DB
	repo *english.Repository
}

// NewEnglishService creates the English Learning Wails service.
func NewEnglishService(db *storage.DB) *EnglishService {
	return &EnglishService{db: db, repo: english.NewRepository(db)}
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
func (s *EnglishService) ListAnkiImports() ([]types.AnkiImportRecord, error) {
	return s.repo.ListAnkiImports(context.Background())
}

// ListArticlesPage returns a paginated list of article summaries.
func (s *EnglishService) ListArticlesPage(page, pageSize int64) (types.ArticlePageResult, error) {
	return s.repo.ListArticlesPage(context.Background(), page, pageSize)
}

// GetArticle returns a single article by ID.
func (s *EnglishService) GetArticle(id string) (types.Article, error) {
	return s.repo.GetArticle(context.Background(), id)
}

// ListVocabularyPage returns a paginated list of vocabulary entries.
func (s *EnglishService) ListVocabularyPage(page, pageSize int64) (types.VocabularyPageResult, error) {
	return s.repo.ListVocabularyPage(context.Background(), page, pageSize)
}

// GetVocabulary returns a single vocabulary entry by ID.
func (s *EnglishService) GetVocabulary(id string) (types.Vocabulary, error) {
	return s.repo.GetVocabulary(context.Background(), id)
}

// UpdateVocabulary saves vocabulary field changes.
func (s *EnglishService) UpdateVocabulary(input types.VocabularyUpdate) error {
	return s.repo.UpdateVocabulary(context.Background(), input)
}

// DeleteVocabulary removes a vocabulary entry.
func (s *EnglishService) DeleteVocabulary(id string) error {
	return s.repo.DeleteVocabulary(context.Background(), id)
}

// SearchVocabulary finds vocabulary entries matching a query string.
func (s *EnglishService) SearchVocabulary(query string, limit int64) ([]types.Vocabulary, error) {
	return s.repo.SearchVocabulary(context.Background(), query, limit)
}

// ListCardsForVocab returns SRS cards linked to a vocabulary entry.
func (s *EnglishService) ListCardsForVocab(vocabID string) ([]types.Card, error) {
	return s.repo.ListCardsForVocab(context.Background(), vocabID)
}

// ListDecks returns all SRS decks.
func (s *EnglishService) ListDecks() ([]types.Deck, error) {
	return s.repo.ListDecks(context.Background())
}

// ListCardsByDeck returns cards in a deck.
func (s *EnglishService) ListCardsByDeck(deckID string) ([]types.Card, error) {
	return s.repo.ListCardsByDeck(context.Background(), deckID)
}

// GetCard returns a single SRS card by ID.
func (s *EnglishService) GetCard(id string) (types.Card, error) {
	return s.repo.GetCard(context.Background(), id)
}

// UpdateCardContent saves editable card fields.
func (s *EnglishService) UpdateCardContent(input types.CardContentUpdate) error {
	return s.repo.UpdateCardContent(context.Background(), input)
}

// DeleteCard removes an SRS card.
func (s *EnglishService) DeleteCard(id string) error {
	return s.repo.DeleteCard(context.Background(), id)
}

// MediaRoot returns the directory where imported media files are stored.
func (s *EnglishService) MediaRoot() (string, error) {
	return content.DefaultMediaDir()
}

// MediaFilePath resolves a stored media relative path to an absolute file path.
func (s *EnglishService) MediaFilePath(storedPath string) (string, error) {
	root, err := content.DefaultMediaDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, storedPath), nil
}
