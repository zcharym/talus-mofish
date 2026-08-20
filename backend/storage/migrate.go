package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const userAccountProviderCheck = "provider IN ('github', 'google', 'email', 'debug')"

func tableHasColumn(sqlDB *sql.DB, table, column string) (bool, error) {
	rows, err := sqlDB.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return false, fmt.Errorf("pragma table_info(%s): %w", table, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); err != nil {
			return false, fmt.Errorf("scan table_info(%s): %w", table, err)
		}
		if strings.EqualFold(name, column) {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("table_info(%s): %w", table, err)
	}
	return false, nil
}

// migrateChatSessionKind adds chat_sessions.kind for databases created before Sudoku sessions.
func migrateChatSessionKind(sqlDB *sql.DB) error {
	hasKind, err := tableHasColumn(sqlDB, "chat_sessions", "kind")
	if err != nil {
		return err
	}
	if hasKind {
		return nil
	}
	if _, err := sqlDB.Exec(`ALTER TABLE chat_sessions ADD COLUMN kind TEXT NOT NULL DEFAULT 'chat'`); err != nil {
		return fmt.Errorf("add chat_sessions.kind: %w", err)
	}
	return nil
}

// migrateUserAccountProvider rebuilds user_account when an older CHECK constraint
// is missing 'email' or 'debug'. CREATE TABLE IF NOT EXISTS never alters CHECKs.
func migrateUserAccountProvider(sqlDB *sql.DB) error {
	var tableSQL string
	err := sqlDB.QueryRow(
		`SELECT sql FROM sqlite_master WHERE type = 'table' AND name = 'user_account'`,
	).Scan(&tableSQL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read user_account schema: %w", err)
	}

	if tableSQL == "" ||
		(strings.Contains(tableSQL, "'email'") && strings.Contains(tableSQL, "'debug'")) {
		return nil
	}

	tx, err := sqlDB.Begin()
	if err != nil {
		return fmt.Errorf("begin user_account migration: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statements := []string{
		`CREATE TABLE user_account_new (
			id TEXT NOT NULL PRIMARY KEY,
			provider TEXT NOT NULL CHECK (` + userAccountProviderCheck + `),
			provider_user_id TEXT NOT NULL,
			display_name TEXT NOT NULL,
			email TEXT NOT NULL DEFAULT '',
			avatar_url TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			last_login_at TEXT NOT NULL DEFAULT (datetime('now')),
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

	for _, stmt := range statements {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("migrate user_account provider: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit user_account migration: %w", err)
	}
	return nil
}
