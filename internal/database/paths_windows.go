//go:build windows

package database

import (
	"fmt"
	"os"
	"path/filepath"
)

// userDataDir returns the base directory for app data on Windows.
// Prefer LOCALAPPDATA (non-roaming, typical for SQLite); fall back to UserConfigDir then UserHomeDir.
func userDataDir() (string, error) {
	if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
		return filepath.Clean(dir), nil
	}

	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Clean(dir), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve windows data directory: %w", err)
	}
	return filepath.Join(home, "AppData", "Local"), nil
}
