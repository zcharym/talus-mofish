package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/songwei.ma/talus_echo_loop/internal/store"
)

//go:embed schema.sql
var schemaSQL string

// driverName is the database/sql driver registered by modernc.org/sqlite (pure Go, no CGO).
const driverName = "sqlite"

// DB wraps the SQLite connection and sqlc-generated queries.
type DB struct {
	SQL     *sql.DB
	Queries *store.Queries
	Path    string
}

// Open opens (or creates) the database at path, applies schema, and returns a DB handle.
func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	sqlDB, err := sql.Open(driverName, sqliteDSN(path))
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	sqlDB.SetMaxOpenConns(1)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	if err := applySchema(sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return &DB{
		SQL:     sqlDB,
		Queries: store.New(sqlDB),
		Path:    path,
	}, nil
}

// OpenDefault opens the database at the default user data path.
func OpenDefault() (*DB, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, err
	}
	return Open(path)
}

// Close closes the underlying database connection.
func (db *DB) Close() error {
	if db == nil || db.SQL == nil {
		return nil
	}
	return db.SQL.Close()
}

func applySchema(sqlDB *sql.DB) error {
	if _, err := sqlDB.Exec(schemaSQL); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}
	return nil
}

// Ping verifies the database connection is alive.
func (db *DB) Ping(ctx context.Context) error {
	return db.SQL.PingContext(ctx)
}
