package appservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/songwei.ma/talus-mofish/internal/store"
)

// ListSettings returns all stored key/value settings.
func (s *Service) ListSettings() ([]store.Setting, error) {
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
func (s *Service) GetSetting(key string) (store.Setting, error) {
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
func (s *Service) SetSetting(key, value string) error {
	ctx := context.Background()
	if err := s.db.Queries.UpsertSetting(ctx, store.UpsertSettingParams{
		Key:   key,
		Value: value,
	}); err != nil {
		return fmt.Errorf("set setting: %w", err)
	}
	return nil
}
