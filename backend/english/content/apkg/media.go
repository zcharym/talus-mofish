package apkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var reImgSrc = regexp.MustCompile(`(?i)(src\s*=\s*["'])([^"']+)(["'])`)

// MediaStore copies Anki media into the app media directory and rewrites HTML references.
type MediaStore struct {
	root      string
	urlPrefix string
	byName    map[string]string // original filename -> stored relative path
}

func NewMediaStore(root string) (*MediaStore, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create media root: %w", err)
	}
	return &MediaStore{
		root:      root,
		urlPrefix: fileURLPrefix(root),
		byName:    map[string]string{},
	}, nil
}

func (m *MediaStore) ImportFromAnki(mediaDir string) error {
	entries, err := os.ReadDir(mediaDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		src := filepath.Join(mediaDir, e.Name())
		if err := m.importFile(e.Name(), src); err != nil {
			return err
		}
	}
	return nil
}

func (m *MediaStore) importFile(originalName, srcPath string) error {
	if _, ok := m.byName[originalName]; ok {
		return nil
	}
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read media %s: %w", originalName, err)
	}
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:8])
	safeName := sanitizeFilename(originalName)
	rel := filepath.Join(hashStr, safeName)
	dest := filepath.Join(m.root, rel)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("write media %s: %w", rel, err)
	}
	m.byName[originalName] = rel
	return nil
}

func (m *MediaStore) RewriteHTML(html string) string {
	return reImgSrc.ReplaceAllStringFunc(html, func(match string) string {
		parts := reImgSrc.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}
		name := parts[2]
		if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") || strings.HasPrefix(name, "data:") {
			return match
		}
		rel, ok := m.byName[name]
		if !ok {
			return match
		}
		url := m.urlPrefix + filepath.ToSlash(rel)
		return parts[1] + url + parts[3]
	})
}

func (m *MediaStore) Count() int {
	return len(m.byName)
}

func (m *MediaStore) Entries() map[string]string {
	out := make(map[string]string, len(m.byName))
	for k, v := range m.byName {
		out[k] = v
	}
	return out
}

func (m *MediaStore) CopyExisting(existing map[string]string) {
	for k, v := range existing {
		m.byName[k] = v
	}
}

func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, ":", "_")
	return name
}

func fileURLPrefix(root string) string {
	abs, err := filepath.Abs(root)
	if err != nil {
		abs = root
	}
	abs = filepath.ToSlash(abs)
	if !strings.HasSuffix(abs, "/") {
		abs += "/"
	}
	return "file:///" + strings.TrimPrefix(abs, "/")
}

// CopyFile is a helper for importing a single media file by path.
func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
