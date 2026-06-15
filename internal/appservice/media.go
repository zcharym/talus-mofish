package appservice

import (
	"path/filepath"

	"github.com/songwei.ma/talus-mofish/internal/content"
)

// MediaRoot returns the directory where imported media files are stored.
func (s *Service) MediaRoot() (string, error) {
	return content.DefaultMediaDir()
}

// MediaFilePath resolves a stored media relative path to an absolute file path.
func (s *Service) MediaFilePath(storedPath string) (string, error) {
	root, err := content.DefaultMediaDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, storedPath), nil
}
