package database

import "path/filepath"

const (
	appDirName = "talus_echo_loop"
	dbFileName = "talus_echo_loop.db"
)

// DefaultPath returns the per-user SQLite database file path for the current OS.
func DefaultPath() (string, error) {
	base, err := userDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appDirName, dbFileName), nil
}
