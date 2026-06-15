package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/songwei.ma/talus-mofish/internal/aiclient"
)

func TestLoadCreatesDefaultConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	store, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	defaults := DefaultApp()
	if store.App != defaults {
		t.Fatalf("got %+v, want %+v", store.App, defaults)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected config file to be written")
	}
}

func TestUpdatePersistsChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	store, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	updated := App{
		Theme:            "dark",
		DailyGoalMinutes: 45,
		WordsPerSession:  15,
		AI:               aiclient.DefaultConfig(),
	}
	if err := store.Update(updated); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	reloaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() second error = %v", err)
	}
	if reloaded.App != updated {
		t.Fatalf("got %+v, want %+v", reloaded.App, updated)
	}
}
