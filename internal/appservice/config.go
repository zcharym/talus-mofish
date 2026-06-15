package appservice

import (
	"fmt"

	"github.com/songwei.ma/talus-mofish/internal/autostart"
	"github.com/songwei.ma/talus-mofish/internal/config"
)

// ConfigPath returns the on-disk config.json file path.
func (s *Service) ConfigPath() string {
	return s.config.Path()
}

// GetConfig returns the current application settings.
func (s *Service) GetConfig() config.App {
	return s.config.Get()
}

// SaveConfig persists updated application settings to config.json.
func (s *Service) SaveConfig(cfg config.App) error {
	if err := s.config.Update(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	if err := s.autostart.Sync(cfg.AutoStart); err != nil {
		return fmt.Errorf("apply autostart: %w", err)
	}
	return nil
}

// GetAutostartStatus returns the OS-level login autostart registration.
func (s *Service) GetAutostartStatus() (autostart.Status, error) {
	status, err := s.autostart.Status()
	if err != nil {
		return autostart.Status{}, fmt.Errorf("autostart status: %w", err)
	}
	return status, nil
}
