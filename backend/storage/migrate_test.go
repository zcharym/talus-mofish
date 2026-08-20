package storage

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateUserAccountProviderAddsEmailAndDebug(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	db, err := Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if _, err := db.SQL.Exec(`DROP TABLE user_account`); err != nil {
		t.Fatalf("drop user_account: %v", err)
	}
	if _, err := db.SQL.Exec(`
		CREATE TABLE user_account (
			id TEXT NOT NULL PRIMARY KEY,
			provider TEXT NOT NULL CHECK (provider IN ('github', 'google')),
			provider_user_id TEXT NOT NULL,
			display_name TEXT NOT NULL,
			email TEXT NOT NULL DEFAULT '',
			avatar_url TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			last_login_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE (provider, provider_user_id)
		)`); err != nil {
		t.Fatalf("create legacy user_account: %v", err)
	}

	if err := migrateUserAccountProvider(db.SQL); err != nil {
		t.Fatalf("migrate user_account: %v", err)
	}

	var tableSQL string
	if err := db.SQL.QueryRow(
		`SELECT sql FROM sqlite_master WHERE type = 'table' AND name = 'user_account'`,
	).Scan(&tableSQL); err != nil {
		t.Fatalf("read user_account schema: %v", err)
	}
	if !strings.Contains(tableSQL, "'email'") || !strings.Contains(tableSQL, "'debug'") {
		t.Fatalf("expected email and debug in CHECK, got: %s", tableSQL)
	}

	for _, provider := range []string{"email", "debug"} {
		if _, err := db.SQL.Exec(
			`INSERT INTO user_account (id, provider, provider_user_id, display_name, email)
			 VALUES (?, ?, ?, 'Test User', 'test@example.com')`,
			"u-"+provider, provider, "uid-"+provider,
		); err != nil {
			t.Fatalf("insert %s provider user: %v", provider, err)
		}
	}
}

func TestMigrateUserAccountProviderPreservesEmailWhenAddingDebug(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	db, err := Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if _, err := db.SQL.Exec(`DROP TABLE user_account`); err != nil {
		t.Fatalf("drop user_account: %v", err)
	}
	if _, err := db.SQL.Exec(`
		CREATE TABLE user_account (
			id TEXT NOT NULL PRIMARY KEY,
			provider TEXT NOT NULL CHECK (provider IN ('github', 'google', 'email')),
			provider_user_id TEXT NOT NULL,
			display_name TEXT NOT NULL,
			email TEXT NOT NULL DEFAULT '',
			avatar_url TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			last_login_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE (provider, provider_user_id)
		)`); err != nil {
		t.Fatalf("create email-only user_account: %v", err)
	}
	if _, err := db.SQL.Exec(
		`INSERT INTO user_account (id, provider, provider_user_id, display_name, email)
		 VALUES ('u1', 'email', 'u1', 'Test User', 'test@example.com')`,
	); err != nil {
		t.Fatalf("seed email user: %v", err)
	}

	if err := migrateUserAccountProvider(db.SQL); err != nil {
		t.Fatalf("migrate user_account: %v", err)
	}

	var count int
	if err := db.SQL.QueryRow(
		`SELECT COUNT(*) FROM user_account WHERE provider = 'email'`,
	).Scan(&count); err != nil {
		t.Fatalf("count email users: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 email user preserved, got %d", count)
	}
}

func TestMigrateChatSessionKindAddsColumn(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "legacy.db")

	sqlDB, err := sql.Open(driverName, sqliteDSN(path))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if _, err := sqlDB.Exec(`
		CREATE TABLE chat_sessions (
			id TEXT NOT NULL PRIMARY KEY,
			title TEXT NOT NULL DEFAULT 'New chat',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now'))
		);
		INSERT INTO chat_sessions (id, title) VALUES ('s1', 'Old chat');
	`); err != nil {
		_ = sqlDB.Close()
		t.Fatalf("create legacy chat_sessions: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close legacy db: %v", err)
	}

	db, err := Open(path)
	if err != nil {
		t.Fatalf("open migrated database: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	hasKind, err := tableHasColumn(db.SQL, "chat_sessions", "kind")
	if err != nil {
		t.Fatalf("tableHasColumn: %v", err)
	}
	if !hasKind {
		t.Fatal("expected chat_sessions.kind after Open")
	}

	var kind string
	if err := db.SQL.QueryRow(`SELECT kind FROM chat_sessions WHERE id = 's1'`).Scan(&kind); err != nil {
		t.Fatalf("read migrated kind: %v", err)
	}
	if kind != "chat" {
		t.Fatalf("kind = %q, want chat", kind)
	}

	var table string
	if err := db.SQL.QueryRow(
		`SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'sudoku_games'`,
	).Scan(&table); err != nil {
		t.Fatalf("sudoku_games missing: %v", err)
	}
}
