package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

const DefaultIdentifier = "com.songwei.ma.talus-mofish"

// Options configures how the app registers for login autostart.
type Options struct {
	Identifier string
	Arguments  []string
}

// Status reports the current OS-level autostart registration.
type Status struct {
	Enabled  bool   `json:"enabled"`
	Path     string `json:"path"`
	Strategy string `json:"strategy"`
}

// Manager registers the application to launch at user login.
type Manager struct {
	identifier string
}

// New returns a Manager using the default identifier when id is empty.
func New(identifier string) *Manager {
	if identifier == "" {
		identifier = DefaultIdentifier
	}
	return &Manager{identifier: identifier}
}

// Enable registers the app to launch at login.
func (m *Manager) Enable() error {
	return m.EnableWithOptions(Options{Identifier: m.identifier})
}

// EnableWithOptions registers with a custom identifier and launch arguments.
func (m *Manager) EnableWithOptions(opts Options) error {
	id := opts.Identifier
	if id == "" {
		id = m.identifier
	}
	exe, err := executablePath()
	if err != nil {
		return fmt.Errorf("resolve executable: %w", err)
	}
	return m.enable(id, exe, opts.Arguments)
}

// Disable removes the login autostart registration.
func (m *Manager) Disable() error {
	return m.disable(m.identifier)
}

// IsEnabled reports whether login autostart is active.
func (m *Manager) IsEnabled() (bool, error) {
	status, err := m.Status()
	if err != nil {
		return false, err
	}
	return status.Enabled, nil
}

// Status returns the current registration details.
func (m *Manager) Status() (Status, error) {
	return m.status(m.identifier)
}

// Sync enables or disables autostart to match the desired setting.
func (m *Manager) Sync(enabled bool) error {
	if enabled {
		return m.Enable()
	}
	return m.Disable()
}

func executablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return filepath.Clean(exe), nil
	}
	return filepath.Clean(resolved), nil
}
