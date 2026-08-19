//go:build linux

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (m *Manager) enable(identifier, exe string, args []string) error {
	desktopPath, err := desktopEntryPath(identifier)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(desktopPath), 0o755); err != nil {
		return fmt.Errorf("create autostart directory: %w", err)
	}

	execLine := shellJoin(exe, args...)
	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=Talus MoFish
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
`, execLine)

	if err := os.WriteFile(desktopPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write desktop entry: %w", err)
	}
	return nil
}

func (m *Manager) disable(identifier string) error {
	desktopPath, err := desktopEntryPath(identifier)
	if err != nil {
		return err
	}
	if err := os.Remove(desktopPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove desktop entry: %w", err)
	}
	return nil
}

func (m *Manager) status(identifier string) (Status, error) {
	desktopPath, err := desktopEntryPath(identifier)
	if err != nil {
		return Status{}, err
	}
	if _, err := os.Stat(desktopPath); err != nil {
		if os.IsNotExist(err) {
			return Status{Enabled: false, Strategy: "xdg-autostart"}, nil
		}
		return Status{}, fmt.Errorf("stat desktop entry: %w", err)
	}
	return Status{
		Enabled:  true,
		Path:     desktopPath,
		Strategy: "xdg-autostart",
	}, nil
}

func desktopEntryPath(identifier string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve config directory: %w", err)
	}
	return filepath.Join(configDir, "autostart", identifier+".desktop"), nil
}

func shellJoin(parts ...string) string {
	quoted := make([]string, len(parts))
	for i, part := range parts {
		quoted[i] = shellQuote(part)
	}
	result := ""
	for i, part := range quoted {
		if i > 0 {
			result += " "
		}
		result += part
	}
	return result
}

func shellQuote(value string) string {
	if value == "" {
		return `''`
	}
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}
