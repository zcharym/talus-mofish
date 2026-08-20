package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/songwei.ma/talus-mofish/backend/obsidian"
	"github.com/songwei.ma/talus-mofish/backend/storage"
)

const defaultObsidianSearchContext = 100

// ObsidianService exposes Local REST API vault operations to the management window.
type ObsidianService struct {
	config *storage.ConfigStore
	client *obsidian.Client
}

// NewObsidianService creates the Obsidian Wails service.
func NewObsidianService(cfg *storage.ConfigStore) *ObsidianService {
	return &ObsidianService{config: cfg}
}

func (s *ObsidianService) apiClient() *obsidian.Client {
	if s.client != nil {
		return s.client
	}
	cfg := s.config.Get().Obsidian
	return obsidian.NewClient(cfg.BaseURL, cfg.APIKey)
}

// Ping checks that Obsidian Local REST API is reachable and the saved API key is accepted.
func (s *ObsidianService) Ping() (obsidian.Status, error) {
	ctx := context.Background()
	status, err := s.apiClient().Ping(ctx)
	if err != nil {
		return obsidian.Status{}, err
	}
	return status, nil
}

// ListDirectory lists files in a vault folder. Pass an empty path for the vault root.
func (s *ObsidianService) ListDirectory(path string) ([]obsidian.FileEntry, error) {
	ctx := context.Background()
	entries, err := s.apiClient().ListDirectory(ctx, path)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		return []obsidian.FileEntry{}, nil
	}
	return entries, nil
}

// ReadNote loads a markdown note by vault-relative path.
func (s *ObsidianService) ReadNote(path string) (obsidian.Note, error) {
	ctx := context.Background()
	note, err := s.apiClient().ReadNote(ctx, path)
	if err != nil {
		return obsidian.Note{}, err
	}
	return note, nil
}

// WriteNote saves markdown content to an existing or new vault path.
func (s *ObsidianService) WriteNote(path, content string) error {
	ctx := context.Background()
	if err := s.apiClient().WriteNote(ctx, path, content); err != nil {
		return err
	}
	return nil
}

// SearchSimple runs Obsidian's built-in full-text search.
func (s *ObsidianService) SearchSimple(query string, contextLength int) ([]obsidian.SearchHit, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}
	if contextLength <= 0 {
		contextLength = defaultObsidianSearchContext
	}
	ctx := context.Background()
	hits, err := s.apiClient().SearchSimple(ctx, query, contextLength)
	if err != nil {
		return nil, err
	}
	if hits == nil {
		return []obsidian.SearchHit{}, nil
	}
	return hits, nil
}
