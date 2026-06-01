//go:build darwin || linux

package database

import (
	"fmt"
	"os"
	"path/filepath"
)

// userDataDir returns the base directory for app data on macOS and Linux.
// macOS: ~/Library/Application Support
// Linux: ~/.config (XDG_CONFIG_HOME when set is respected by UserConfigDir)
func userDataDir() (string, error) {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Clean(dir), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve unix data directory: %w", err)
	}

	// UserConfigDir failed; use conventional fallbacks.
	switch {
	case fileExists(filepath.Join(home, "Library")):
		return filepath.Join(home, "Library", "Application Support"), nil
	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Clean(xdg), nil
		}
		return filepath.Join(home, ".config"), nil
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
