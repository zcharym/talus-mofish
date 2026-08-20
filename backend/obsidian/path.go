package obsidian

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// NormalizeVaultPath rejects absolute paths and ".." and returns a vault-relative path.
// The empty string means the vault root. Trailing slashes are stripped.
func NormalizeVaultPath(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	raw = strings.ReplaceAll(raw, "\\", "/")
	if strings.HasPrefix(raw, "/") {
		return "", fmt.Errorf("absolute paths are not allowed")
	}

	parts := strings.Split(raw, "/")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			return "", fmt.Errorf("path must not contain '..'")
		}
		if strings.Contains(part, ":") {
			return "", fmt.Errorf("absolute paths are not allowed")
		}
		out = append(out, part)
	}
	return strings.Join(out, "/"), nil
}

// EncodeVaultPath percent-encodes each path segment, keeping "/" as separators.
func EncodeVaultPath(vaultPath string) string {
	if vaultPath == "" {
		return ""
	}
	parts := strings.Split(vaultPath, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

// JoinVaultPath joins a directory and a listing name (which may end with "/").
func JoinVaultPath(dir, name string) string {
	name = strings.TrimSuffix(strings.TrimSpace(name), "/")
	dir = strings.Trim(dir, "/")
	if dir == "" {
		return name
	}
	if name == "" {
		return dir
	}
	return dir + "/" + name
}

// IsMarkdown reports whether path looks like a markdown note.
func IsMarkdown(vaultPath string) bool {
	return strings.EqualFold(path.Ext(vaultPath), ".md")
}
