package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/songwei.ma/talus-mofish/backend/types"
	"github.com/songwei.ma/talus-mofish/backend/utils/aiclient"
)

func TestLoadConfigCreatesDefaultConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	store, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
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

	store, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	updated := types.App{
		Theme:            "dark",
		DailyGoalMinutes: 45,
		WordsPerSession:  15,
		AI:               aiclient.DefaultConfig(),
		Obsidian:         types.Obsidian{BaseURL: types.DefaultObsidianBaseURL},
	}
	if err := store.Update(updated); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	reloaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() second error = %v", err)
	}
	if reloaded.App != updated {
		t.Fatalf("got %+v, want %+v", reloaded.App, updated)
	}
}
