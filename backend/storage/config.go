package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/songwei.ma/talus-mofish/backend/types"
	"github.com/songwei.ma/talus-mofish/backend/utils/aiclient"
)

// ConfigStore loads and saves App settings at a fixed file path.
// Product settings live in config.json; integration secrets are stored via SecretBackend
// (OS keyring) and overlaid in memory. SQLite settings remain ephemeral/debug KV.
type ConfigStore struct {
	mu   sync.RWMutex
	path string
	App  types.App
}

// DefaultApp returns factory defaults for a new installation.
func DefaultApp() types.App {
	return types.App{
		Theme:            "auto",
		DailyGoalMinutes: 30,
		WordsPerSession:  20,
		AI:               aiclient.DefaultConfig(),
		Obsidian:         types.Obsidian{BaseURL: types.DefaultObsidianBaseURL},
	}
}

// LoadConfig reads config from path. Missing files are created with defaults.
func LoadConfig(path string) (*ConfigStore, error) {
	defaults := DefaultApp()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create config directory: %w", err)
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		store := &ConfigStore{path: path, App: defaults}
		if err := store.Save(); err != nil {
			return nil, err
		}
		return store, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var app types.App
	if err := json.Unmarshal(data, &app); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	app = mergeDefaults(app, defaults)
	hadFileSecrets := configHasSecrets(app)
	overlayConfigSecrets(
		&app.AI.APIKey,
		&app.OAuth.GitHubClientSecret,
		&app.OAuth.GoogleClientSecret,
		&app.Sudoku.APIKey,
		&app.Obsidian.APIKey,
		&app.Cloudflare.APIToken,
	)
	store := &ConfigStore{path: path, App: app}
	if hadFileSecrets {
		if err := store.Save(); err != nil {
			return nil, err
		}
	}
	return store, nil
}

// Path returns the on-disk config.json file path.
func (s *ConfigStore) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.path
}

// Get returns the current in-memory settings (secrets included).
func (s *ConfigStore) Get() types.App {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.App
}

// Update replaces in-memory settings and persists them to disk.
func (s *ConfigStore) Update(app types.App) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.App = mergeDefaults(app, DefaultApp())
	return s.saveLocked()
}

// Save writes the current settings to disk with secrets stripped.
func (s *ConfigStore) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked()
}

func (s *ConfigStore) saveLocked() error {
	app := s.App
	path := s.path

	persistConfigSecrets(
		app.AI.APIKey,
		app.OAuth.GitHubClientSecret,
		app.OAuth.GoogleClientSecret,
		app.Sudoku.APIKey,
		app.Obsidian.APIKey,
		app.Cloudflare.APIToken,
	)

	disk := app
	stripConfigSecretsForDisk(
		&disk.AI.APIKey,
		&disk.OAuth.GitHubClientSecret,
		&disk.OAuth.GoogleClientSecret,
		&disk.Sudoku.APIKey,
		&disk.Obsidian.APIKey,
		&disk.Cloudflare.APIToken,
	)

	data, err := json.MarshalIndent(disk, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func mergeDefaults(app, defaults types.App) types.App {
	if app.Theme == "" {
		app.Theme = defaults.Theme
	}
	if app.DailyGoalMinutes <= 0 {
		app.DailyGoalMinutes = defaults.DailyGoalMinutes
	}
	if app.WordsPerSession <= 0 {
		app.WordsPerSession = defaults.WordsPerSession
	}
	app.AI = app.AI.Normalize()
	app.Obsidian = types.NormalizeObsidian(app.Obsidian)
	app.Cloudflare = types.NormalizeCloudflare(app.Cloudflare)
	return app
}

func configHasSecrets(app types.App) bool {
	return app.AI.APIKey != "" ||
		app.OAuth.GitHubClientSecret != "" ||
		app.OAuth.GoogleClientSecret != "" ||
		app.Sudoku.APIKey != "" ||
		app.Obsidian.APIKey != "" ||
		app.Cloudflare.APIToken != ""
}
