package services

import (
	"fmt"
	"path/filepath"

	"github.com/songwei.ma/talus-mofish/backend/english/content"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/utils/autostart"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// SystemService exposes shell paths, native dialogs, and window management.
type SystemService struct {
	db        *storage.DB
	autostart *autostart.Manager
	wailsApp  *application.App
	windows   WindowManager
}

// NewSystemService creates the system Wails service.
func NewSystemService(db *storage.DB, autostartManager *autostart.Manager) *SystemService {
	return &SystemService{
		db:        db,
		autostart: autostartManager,
	}
}

// DatabasePath returns the on-disk SQLite database file path.
func (s *SystemService) DatabasePath() string {
	return s.db.Path
}

// GetAutostartStatus returns the OS-level login autostart registration.
func (s *SystemService) GetAutostartStatus() (autostart.Status, error) {
	status, err := s.autostart.Status()
	if err != nil {
		return autostart.Status{}, fmt.Errorf("autostart status: %w", err)
	}
	return status, nil
}

// ShowAgentWindow shows and focuses the agent chat window.
func (s *SystemService) ShowAgentWindow() {
	if s.windows != nil {
		s.windows.ShowAgentWindow()
	}
}

// ShowManagementWindow shows and focuses the management window.
func (s *SystemService) ShowManagementWindow() {
	if s.windows != nil {
		s.windows.ShowManagementWindow()
	}
}

// PickAnkiAPKG opens a file dialog to select an Anki APKG file.
func (s *SystemService) PickAnkiAPKG() (string, error) {
	if s.wailsApp == nil {
		return "", fmt.Errorf("file dialog unavailable")
	}
	path, err := s.wailsApp.Dialog.OpenFile().
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

// MediaRoot returns the directory where imported media files are stored.
func (s *SystemService) MediaRoot() (string, error) {
	return content.DefaultMediaDir()
}

// MediaFilePath resolves a stored media relative path to an absolute file path.
func (s *SystemService) MediaFilePath(storedPath string) (string, error) {
	root, err := content.DefaultMediaDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, storedPath), nil
}
