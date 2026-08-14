package database

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/songwei.ma/talus-mofish/internal/store"
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

	if err := migrateUserAccountProviders(sqlDB); err != nil {
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

// migrateUserAccountProviders rebuilds user_account when an older CHECK constraint
// does not allow provider='debug'. CREATE TABLE IF NOT EXISTS never alters CHECKs.
func migrateUserAccountProviders(sqlDB *sql.DB) error {
	var createSQL string
	err := sqlDB.QueryRow(
		`SELECT sql FROM sqlite_master WHERE type = 'table' AND name = 'user_account'`,
	).Scan(&createSQL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect user_account schema: %w", err)
	}
	if strings.Contains(createSQL, "'debug'") {
		return nil
	}

	tx, err := sqlDB.Begin()
	if err != nil {
		return fmt.Errorf("begin user_account migration: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmts := []string{
		`CREATE TABLE user_account_new (
			id TEXT NOT NULL PRIMARY KEY,
			provider TEXT NOT NULL CHECK (provider IN ('github', 'google', 'debug')),
			provider_user_id TEXT NOT NULL,
			display_name TEXT NOT NULL,
			email TEXT NOT NULL DEFAULT '',
			avatar_url TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT (datetime ('now')),
			last_login_at TEXT NOT NULL DEFAULT (datetime ('now')),
			UNIQUE (provider, provider_user_id)
		)`,
		`INSERT INTO user_account_new (
			id, provider, provider_user_id, display_name, email, avatar_url, created_at, last_login_at
		)
		SELECT id, provider, provider_user_id, display_name, email, avatar_url, created_at, last_login_at
		FROM user_account`,
		`DROP TABLE user_account`,
		`ALTER TABLE user_account_new RENAME TO user_account`,
	}
	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("migrate user_account providers: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit user_account migration: %w", err)
	}
	return nil
}

// Ping verifies the database connection is alive.
func (db *DB) Ping(ctx context.Context) error {
	return db.SQL.PingContext(ctx)
}
