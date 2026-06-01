package database

import (
	"fmt"
	"path/filepath"
)

// sqliteDSN builds a connection string for modernc.org/sqlite (file URI + pragmas).
func sqliteDSN(path string) string {
	p := filepath.ToSlash(path)
	return fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", p)
}
