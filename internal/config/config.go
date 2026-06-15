package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/songwei.ma/talus-mofish/internal/aiclient"
)

// App holds user-facing application settings persisted in config.json.
type App struct {
	Theme            string           `json:"theme"`
	DailyGoalMinutes int              `json:"dailyGoalMinutes"`
	WordsPerSession  int              `json:"wordsPerSession"`
	AutoStart        bool             `json:"autoStart"`
	AI               aiclient.Config  `json:"ai"`
}

// Store loads and saves App settings at a fixed file path.
type Store struct {
	path string
	App  App
}

// DefaultApp returns factory defaults for a new installation.
func DefaultApp() App {
	return App{
		Theme:            "auto",
		DailyGoalMinutes: 30,
		WordsPerSession:  20,
		AI:               aiclient.DefaultConfig(),
	}
}

// Load reads config from path. Missing files are created with defaults.
func Load(path string) (*Store, error) {
	defaults := DefaultApp()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create config directory: %w", err)
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		store := &Store{path: path, App: defaults}
		if err := store.Save(); err != nil {
			return nil, err
		}
		return store, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var app App
	if err := json.Unmarshal(data, &app); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	app = mergeDefaults(app, defaults)
	return &Store{path: path, App: app}, nil
}

// Path returns the on-disk config.json file path.
func (s *Store) Path() string {
	return s.path
}

// Get returns the current in-memory settings.
func (s *Store) Get() App {
	return s.App
}

// Update replaces in-memory settings and persists them to disk.
func (s *Store) Update(app App) error {
	s.App = mergeDefaults(app, DefaultApp())
	return s.Save()
}

// Save writes the current settings to disk.
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.App, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func mergeDefaults(app, defaults App) App {
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
