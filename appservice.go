package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/songwei.ma/talus-mofish/internal/autostart"
	"github.com/songwei.ma/talus-mofish/internal/config"
	"github.com/songwei.ma/talus-mofish/internal/database"
	"github.com/songwei.ma/talus-mofish/internal/store"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// AppService exposes application and persistence APIs to the frontend.
type AppService struct {
	db        *database.DB
	config    *config.Store
	autostart *autostart.Manager
	wailsApp  *application.App
}

func NewAppService(db *database.DB, cfg *config.Store, autostartManager *autostart.Manager) *AppService {
	return &AppService{db: db, config: cfg, autostart: autostartManager}
}

// DatabasePath returns the on-disk SQLite database file path.
func (a *AppService) DatabasePath() string {
	return a.db.Path
}

// ConfigPath returns the on-disk config.json file path.
func (a *AppService) ConfigPath() string {
	return a.config.Path()
}

// GetConfig returns the current application settings.
func (a *AppService) GetConfig() config.App {
	return a.config.Get()
}

// SaveConfig persists updated application settings to config.json.
func (a *AppService) SaveConfig(cfg config.App) error {
	if err := a.config.Update(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	if err := a.autostart.Sync(cfg.AutoStart); err != nil {
		return fmt.Errorf("apply autostart: %w", err)
	}
	return nil
}

// GetAutostartStatus returns the OS-level login autostart registration.
func (a *AppService) GetAutostartStatus() (autostart.Status, error) {
	status, err := a.autostart.Status()
	if err != nil {
		return autostart.Status{}, fmt.Errorf("autostart status: %w", err)
	}
	return status, nil
}

// ListSettings returns all stored key/value settings.
func (a *AppService) ListSettings() ([]store.Setting, error) {
	ctx := context.Background()
	items, err := a.db.Queries.ListSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("list settings: %w", err)
	}
	if items == nil {
		return []store.Setting{}, nil
	}
	return items, nil
}

// GetSetting returns a single setting by key.
func (a *AppService) GetSetting(key string) (store.Setting, error) {
	ctx := context.Background()
	setting, err := a.db.Queries.GetSetting(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.Setting{}, fmt.Errorf("setting %q not found", key)
		}
		return store.Setting{}, fmt.Errorf("get setting: %w", err)
	}
	return setting, nil
}

// SetSetting creates or updates a key/value setting.
func (a *AppService) SetSetting(key, value string) error {
	ctx := context.Background()
	if err := a.db.Queries.UpsertSetting(ctx, store.UpsertSettingParams{
		Key:   key,
		Value: value,
	}); err != nil {
		return fmt.Errorf("set setting: %w", err)
	}
	return nil
}
