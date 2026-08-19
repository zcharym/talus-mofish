package auth_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/songwei.ma/talus-mofish/backend/auth"
	"github.com/songwei.ma/talus-mofish/backend/storage"
)

func TestGetCurrentUserCreatesDebugAdmin(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	cfgPath := filepath.Join(t.TempDir(), "config.json")
	cfg, err := storage.LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	app := cfg.Get()
	app.DebugMode = true
	if err := cfg.Update(app); err != nil {
		t.Fatalf("enable debugMode: %v", err)
	}

	svc := auth.New(db.Queries, cfg)
	profile, err := svc.GetCurrentUser(context.Background())
	if err != nil {
		t.Fatalf("GetCurrentUser: %v", err)
	}
	if profile == nil {
		t.Fatal("expected debug admin user, got nil")
	}
	if profile.Provider != auth.ProviderDebug {
		t.Fatalf("provider = %q, want %q", profile.Provider, auth.ProviderDebug)
	}
	if profile.DisplayName != "Debug Admin" {
		t.Fatalf("display name = %q, want Debug Admin", profile.DisplayName)
	}
}

func TestGetCurrentUserClearsDebugUserWhenDebugOff(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	cfgPath := filepath.Join(t.TempDir(), "config.json")
	cfg, err := storage.LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	app := cfg.Get()
	app.DebugMode = true
	if err := cfg.Update(app); err != nil {
		t.Fatalf("enable debugMode: %v", err)
	}

	svc := auth.New(db.Queries, cfg)
	if _, err := svc.EnsureDebugUser(context.Background()); err != nil {
		t.Fatalf("EnsureDebugUser: %v", err)
	}

	app.DebugMode = false
	if err := cfg.Update(app); err != nil {
		t.Fatalf("disable debugMode: %v", err)
	}

	profile, err := svc.GetCurrentUser(context.Background())
	if err != nil {
		t.Fatalf("GetCurrentUser: %v", err)
	}
	if profile != nil {
		t.Fatalf("expected nil after debugMode off, got %+v", profile)
	}
}

func TestEnsureDebugUserRequiresDebugMode(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	cfgPath := filepath.Join(t.TempDir(), "config.json")
	cfg, err := storage.LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	svc := auth.New(db.Queries, cfg)
	if _, err := svc.EnsureDebugUser(context.Background()); err == nil {
		t.Fatal("expected error when debugMode is false")
	}
}
