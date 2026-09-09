package services

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"

	"github.com/songwei.ma/talus-mofish/backend/consts"
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

// Platform returns "mac" | "windows" | GOOS for native chrome decisions.
func (s *SystemService) Platform() string {
	if runtime.GOOS == "darwin" {
		return "mac"
	}
	return runtime.GOOS
}

// Version returns the application version string.
func (s *SystemService) Version() string {
	return consts.AppVersion
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

// ServiceShutdown closes SQLite after other services have stopped.
func (s *SystemService) ServiceShutdown() error {
	if s.db == nil {
		return nil
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}

// OpenURL opens an external URL in the system browser.
func (s *SystemService) OpenURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("url is required")
	}
	if s.wailsApp == nil {
		return fmt.Errorf("browser unavailable")
	}
	parsed, err := url.Parse(raw)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("unsupported url")
	}
	return s.wailsApp.Browser.OpenURL(parsed.String())
}
