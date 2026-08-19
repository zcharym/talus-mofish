package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/storage/store"
	"github.com/songwei.ma/talus-mofish/backend/types"
	"github.com/songwei.ma/talus-mofish/backend/utils/autostart"
)

// ConfigService exposes config.json and key/value settings APIs.
type ConfigService struct {
	db        *storage.DB
	config    *storage.ConfigStore
	autostart *autostart.Manager
}

// NewConfigService creates the config Wails service.
func NewConfigService(db *storage.DB, cfg *storage.ConfigStore, autostartManager *autostart.Manager) *ConfigService {
	return &ConfigService{
		db:        db,
		config:    cfg,
		autostart: autostartManager,
	}
}

// ConfigPath returns the on-disk config.json file path.
func (s *ConfigService) ConfigPath() string {
	return s.config.Path()
}

// GetConfig returns the current application settings.
func (s *ConfigService) GetConfig() types.App {
	return s.config.Get()
}

// SaveConfig persists updated application settings to config.json.
// Auth.AuthServerURL is application-level (.env.* / TALUS_AUTH_SERVER_URL / config.json)
// and is not managed by the Config UI, so it is preserved from the current store.
func (s *ConfigService) SaveConfig(cfg types.App) error {
	cfg.Auth = s.config.Get().Auth
	if err := s.config.Update(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	if err := s.autostart.Sync(cfg.AutoStart); err != nil {
		return fmt.Errorf("apply autostart: %w", err)
	}
	return nil
}

// ListSettings returns all stored key/value settings.
func (s *ConfigService) ListSettings() ([]store.Setting, error) {
	ctx := context.Background()
	items, err := s.db.Queries.ListSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("list settings: %w", err)
	}
	if items == nil {
		return []store.Setting{}, nil
	}
	return items, nil
}

// GetSetting returns a single setting by key.
func (s *ConfigService) GetSetting(key string) (store.Setting, error) {
	ctx := context.Background()
	setting, err := s.db.Queries.GetSetting(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.Setting{}, fmt.Errorf("setting %q not found", key)
		}
		return store.Setting{}, fmt.Errorf("get setting: %w", err)
	}
	return setting, nil
}

// SetSetting creates or updates a key/value setting.
func (s *ConfigService) SetSetting(key, value string) error {
	ctx := context.Background()
	if err := s.db.Queries.UpsertSetting(ctx, store.UpsertSettingParams{
		Key:   key,
		Value: value,
	}); err != nil {
		return fmt.Errorf("set setting: %w", err)
	}
	return nil
}
