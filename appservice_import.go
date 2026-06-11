package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/songwei.ma/talus-mofish/internal/content"
	"github.com/songwei.ma/talus-mofish/internal/store"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// SetWailsApp attaches the Wails application for native dialogs.
func (a *AppService) SetWailsApp(app *application.App) {
	a.wailsApp = app
}

// PickAnkiAPKG opens a file dialog to select an Anki APKG file.
func (a *AppService) PickAnkiAPKG() (string, error) {
	if a.wailsApp == nil {
		return "", fmt.Errorf("file dialog unavailable")
	}
	path, err := a.wailsApp.Dialog.OpenFile().
		SetTitle("Select Anki deck (.apkg)").
		AddFilter("Anki deck", "*.apkg").
		AddFilter("All files", "*.*").
		PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("open file dialog: %w", err)
	}
	if path == "" {
		return "", nil
	}
	return path, nil
}

// PreviewAnkiAPKG returns deck and model metadata for import configuration.
func (a *AppService) PreviewAnkiAPKG(apkgPath string) (content.AnkiPreview, error) {
	importer, err := content.NewImporter(a.db)
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
func (a *AppService) ImportAnkiAPKG(apkgPath string, configs []content.ImportDeckConfig) (content.ImportResult, error) {
	importer, err := content.NewImporter(a.db)
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
func (a *AppService) ListAnkiImports() ([]store.AnkiImport, error) {
	ctx := context.Background()
	items, err := a.db.Queries.ListAnkiImports(ctx)
	if err != nil {
		return nil, fmt.Errorf("list imports: %w", err)
	}
	if items == nil {
		return []store.AnkiImport{}, nil
	}
	return items, nil
}

// ListVocabularyPage returns a paginated list of vocabulary entries.
func (a *AppService) ListVocabularyPage(page, pageSize int64) (store.VocabularyPageResult, error) {
	ctx := context.Background()
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := a.db.Queries.CountVocabulary(ctx)
	if err != nil {
		return store.VocabularyPageResult{}, fmt.Errorf("count vocabulary: %w", err)
	}

	items, err := a.db.Queries.ListVocabularyPage(ctx, store.ListVocabularyPageParams{
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

// ListArticlesPage returns a paginated list of article summaries.
func (a *AppService) ListArticlesPage(page, pageSize int64) (store.ArticlePageResult, error) {
	ctx := context.Background()
	page, pageSize = normalizePageParams(page, pageSize)

	total, err := a.db.Queries.CountArticles(ctx)
	if err != nil {
		return store.ArticlePageResult{}, fmt.Errorf("count articles: %w", err)
	}

	rows, err := a.db.Queries.ListArticlesPage(ctx, store.ListArticlesPageParams{
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
func (a *AppService) GetArticle(id string) (store.Article, error) {
	ctx := context.Background()
	article, err := a.db.Queries.GetArticle(ctx, id)
	if err != nil {
		return store.Article{}, fmt.Errorf("get article: %w", err)
	}
	return article, nil
}

func normalizePageParams(page, pageSize int64) (int64, int64) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = store.DefaultPageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// ListDecks returns all SRS decks.
func (a *AppService) ListDecks() ([]store.Deck, error) {
	ctx := context.Background()
	items, err := a.db.Queries.ListDecks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list decks: %w", err)
	}
	if items == nil {
		return []store.Deck{}, nil
	}
	return items, nil
}

// ListCardsByDeck returns cards in a deck.
func (a *AppService) ListCardsByDeck(deckID string) ([]store.Card, error) {
	ctx := context.Background()
	items, err := a.db.Queries.ListCardsByDeck(ctx, deckID)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	if items == nil {
		return []store.Card{}, nil
	}
	return items, nil
}

// MediaRoot returns the directory where imported media files are stored.
func (a *AppService) MediaRoot() (string, error) {
	return content.DefaultMediaDir()
}

// MediaFilePath resolves a stored media relative path to an absolute file path.
func (a *AppService) MediaFilePath(storedPath string) (string, error) {
	root, err := content.DefaultMediaDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, storedPath), nil
}
