package database

import "path/filepath"

const (
	appDirName     = "talus-mofish"
	dbFileName     = "talus-mofish.db"
	configFileName = "config.json"
)

// UserDataDir returns the OS user data directory used for app storage.
func UserDataDir() (string, error) {
	return userDataDir()
}

// DefaultPath returns the per-user SQLite database file path for the current OS.
func DefaultPath() (string, error) {
	base, err := userDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appDirName, dbFileName), nil
}

// DefaultConfigPath returns the per-user config.json file path for the current OS.
func DefaultConfigPath() (string, error) {
	base, err := userDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appDirName, configFileName), nil
}
