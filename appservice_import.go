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

// ListVocabulary returns all vocabulary entries.
func (a *AppService) ListVocabulary() ([]store.Vocabulary, error) {
	ctx := context.Background()
	items, err := a.db.Queries.ListVocabulary(ctx)
	if err != nil {
		return nil, fmt.Errorf("list vocabulary: %w", err)
	}
	if items == nil {
		return []store.Vocabulary{}, nil
	}
	return items, nil
}

// ListArticles returns all reading articles.
func (a *AppService) ListArticles() ([]store.Article, error) {
	ctx := context.Background()
	items, err := a.db.Queries.ListArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("list articles: %w", err)
	}
	if items == nil {
		return []store.Article{}, nil
	}
	return items, nil
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
