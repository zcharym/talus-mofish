package database

import (
	"path/filepath"
	"testing"
)

func TestMigrateUserAccountProviderAddsEmail(t *testing.T) {
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

	if _, err := db.SQL.Exec(
		`INSERT INTO user_account (id, provider, provider_user_id, display_name, email)
		 VALUES ('u1', 'email', 'u1', 'Test User', 'test@example.com')`,
	); err != nil {
		t.Fatalf("insert email provider user: %v", err)
	}
}
