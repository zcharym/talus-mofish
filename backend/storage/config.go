package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/songwei.ma/talus-mofish/backend/types"
	"github.com/songwei.ma/talus-mofish/backend/utils/aiclient"
)

// ConfigStore loads and saves App settings at a fixed file path.
type ConfigStore struct {
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
	return &ConfigStore{path: path, App: app}, nil
}

// Path returns the on-disk config.json file path.
func (s *ConfigStore) Path() string {
	return s.path
}

// Get returns the current in-memory settings.
func (s *ConfigStore) Get() types.App {
	return s.App
}

// Update replaces in-memory settings and persists them to disk.
func (s *ConfigStore) Update(app types.App) error {
	s.App = mergeDefaults(app, DefaultApp())
	return s.Save()
}

// Save writes the current settings to disk.
func (s *ConfigStore) Save() error {
	data, err := json.MarshalIndent(s.App, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
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
	return app
}
