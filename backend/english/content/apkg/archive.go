package apkg

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const fieldSeparator = "\x1f"

// OpenArchive extracts or opens the collection database from an APKG zip.
func OpenArchive(apkgPath string) (*Collection, func(), error) {
	reader, err := zip.OpenReader(apkgPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open apkg zip: %w", err)
	}

	dbName := ""
	for _, name := range []string{"collection.anki21", "collection.anki2"} {
		for _, f := range reader.File {
			if f.Name == name {
				dbName = name
				break
			}
		}
		if dbName != "" {
			break
		}
	}
	if dbName == "" {
		reader.Close()
		return nil, nil, fmt.Errorf("apkg missing collection.anki21 or collection.anki2")
	}

	tmpDir, err := os.MkdirTemp("", "talus-apkg-*")
	if err != nil {
		reader.Close()
		return nil, nil, fmt.Errorf("create temp dir: %w", err)
	}

	cleanup := func() {
		reader.Close()
		os.RemoveAll(tmpDir)
	}

	if err := extractZipMember(reader, dbName, filepath.Join(tmpDir, dbName)); err != nil {
		cleanup()
		return nil, nil, err
	}

	mediaMap, err := readMediaMap(reader)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	for zipKey, originalName := range mediaMap {
		dest := filepath.Join(tmpDir, "medias", originalName)
		if err := extractZipMember(reader, zipKey, dest); err != nil {
			// Missing media entries are common; skip quietly.
			continue
		}
	}

	coll, err := OpenCollection(filepath.Join(tmpDir, dbName), filepath.Join(tmpDir, "medias"))
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return coll, cleanup, nil
}

func extractZipMember(reader *zip.ReadCloser, name, dest string) error {
	for _, f := range reader.File {
		if f.Name != name {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("mkdir media dir: %w", err)
		}
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open zip member %s: %w", name, err)
		}
		defer rc.Close()

		out, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("create %s: %w", dest, err)
		}
		_, copyErr := io.Copy(out, rc)
		closeErr := out.Close()
		if copyErr != nil {
			return fmt.Errorf("copy %s: %w", name, copyErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close %s: %w", dest, closeErr)
		}
		return nil
	}
	return fmt.Errorf("zip member %s not found", name)
}

func readMediaMap(reader *zip.ReadCloser) (map[string]string, error) {
	for _, f := range reader.File {
		if f.Name != "media" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open media map: %w", err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read media map: %w", err)
		}
		text := strings.TrimSpace(string(data))
		if text == "" {
			return map[string]string{}, nil
		}
		return parseMediaJSON(text)
	}
	return map[string]string{}, nil
}

func parseMediaJSON(text string) (map[string]string, error) {
	// Anki stores media as JSON object mapping zip keys to filenames.
	m := map[string]string{}
	if err := decodeJSON(text, &m); err != nil {
		return nil, fmt.Errorf("parse media json: %w", err)
	}
	return m, nil
}
