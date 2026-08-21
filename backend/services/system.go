package services

import (
	"fmt"
	"path/filepath"

	"github.com/songwei.ma/talus-mofish/backend/english/content"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/utils/autostart"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// SystemService exposes system-level operations including file dialogs,
// window management, autostart configuration, and media file paths.
// This service provides the bridge between the Go backend and OS-level features.
type SystemService struct {
	db        *storage.DB
	autostart *autostart.Manager
	wailsApp  *application.App
	windows   WindowManager
}

// NewSystemService creates a new system service with the given dependencies.
// The autostart manager handles OS-level login item configuration.
func NewSystemService(db *storage.DB, autostartManager *autostart.Manager) *SystemService {
	return &SystemService{
		db:        db,
		autostart: autostartManager,
	}
}

// DatabasePath returns the absolute path to the SQLite database file on disk.
// This is useful for debugging and manual database inspection.
func (s *SystemService) DatabasePath() string {
	return s.db.Path
}

// GetAutostartStatus retrieves the current OS-level autostart configuration.
// It checks whether the application is registered to launch on system startup.
func (s *SystemService) GetAutostartStatus() (autostart.Status, error) {
	status, err := s.autostart.Status()
	if err != nil {
		return autostart.Status{}, fmt.Errorf("failed to retrieve autostart status: %w", err)
	}
	return status, nil
}

// ShowAgentWindow brings the agent chat window to the foreground and focuses it.
// This is typically called from system tray menu actions or keyboard shortcuts.
func (s *SystemService) ShowAgentWindow() {
	if s.windows != nil {
		s.windows.ShowAgentWindow()
	}
}

// ShowManagementWindow brings the management window to the foreground and focuses it.
// The management window provides access to domain-specific tools and settings.
func (s *SystemService) ShowManagementWindow() {
	if s.windows != nil {
		s.windows.ShowManagementWindow()
	}
}

// PickAnkiAPKG opens a native file picker dialog for selecting an Anki deck file.
// Returns the selected file path, an empty string if cancelled, or an error on failure.
func (s *SystemService) PickAnkiAPKG() (string, error) {
	if s.wailsApp == nil {
		return "", fmt.Errorf("file dialog is not available (wails app not initialized)")
	}
	
	path, err := s.wailsApp.Dialog.OpenFile().
		SetTitle("Select Anki deck (.apkg)").
		AddFilter("Anki deck", "*.apkg").
		AddFilter("All files", "*.*").
		PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("failed to open file dialog: %w", err)
	}
	
	// Empty path means user cancelled
	if path == "" {
		return "", nil
	}
	return path, nil
}

// MediaRoot returns the root directory where imported media files are stored.
// This directory is typically within the application's user data directory.
func (s *SystemService) MediaRoot() (string, error) {
	root, err := content.DefaultMediaDir()
	if err != nil {
		return "", fmt.Errorf("failed to get media root directory: %w", err)
	}
	return root, nil
}

// MediaFilePath converts a relative media path (as stored in the database)
// to an absolute filesystem path. This is used by the frontend to load media files.
func (s *SystemService) MediaFilePath(storedPath string) (string, error) {
	root, err := content.DefaultMediaDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve media directory: %w", err)
	}
	return filepath.Join(root, storedPath), nil
}
