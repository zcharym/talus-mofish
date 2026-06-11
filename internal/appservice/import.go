package appservice

import (
	"context"
	"fmt"

	"github.com/songwei.ma/talus-mofish/internal/content"
	"github.com/songwei.ma/talus-mofish/internal/store"
)

// PreviewAnkiAPKG returns deck and model metadata for import configuration.
func (s *Service) PreviewAnkiAPKG(apkgPath string) (content.AnkiPreview, error) {
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
func (s *Service) ImportAnkiAPKG(apkgPath string, configs []content.ImportDeckConfig) (content.ImportResult, error) {
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
func (s *Service) ListAnkiImports() ([]store.AnkiImport, error) {
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
