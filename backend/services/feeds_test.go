package services

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/songwei.ma/talus-mofish/backend/feeds"
	"github.com/songwei.ma/talus-mofish/backend/storage"
)

func TestFeedsServiceAddRSSAndSaveLink(t *testing.T) {
	dir := t.TempDir()
	db, err := storage.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	storage.UseMemorySecrets()

	cfg, err := storage.LoadConfig(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}

	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"><channel><title>Demo</title>
<item><title>Hello</title><link>https://example.com/hello</link><guid>g1</guid></item>
</channel></rss>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(xml))
	}))
	t.Cleanup(server.Close)

	svc := NewFeedsService(db, cfg)
	svc.client = feeds.NewClient(server.Client())

	if _, err := svc.AddSource(server.URL); err != nil {
		t.Fatalf("AddSource: %v", err)
	}
	items, err := svc.ListItems("rss")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Title != "Hello" {
		t.Fatalf("items = %+v", items)
	}

	saved, err := svc.SaveLink("Later", "https://example.com/later")
	if err != nil {
		t.Fatal(err)
	}
	if !saved.Saved || saved.Title != "Later" {
		t.Fatalf("saved = %+v", saved)
	}
	if err := svc.SetItemSaved(items[0].ID, true); err != nil {
		t.Fatal(err)
	}
}
