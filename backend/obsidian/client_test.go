package obsidian

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShouldSkipTLS(t *testing.T) {
	t.Parallel()

	if !shouldSkipTLS("https://127.0.0.1:27124") {
		t.Fatal("expected loopback HTTPS to skip TLS verify")
	}
	if shouldSkipTLS("http://127.0.0.1:27123") {
		t.Fatal("HTTP should not skip TLS")
	}
	if shouldSkipTLS("https://example.com") {
		t.Fatal("remote HTTPS should not skip TLS")
	}
}

func TestClientPingSendsBearerAndParsesStatus(t *testing.T) {
	t.Parallel()

	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Errorf("path = %q", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")
		_ = json.NewEncoder(w).Encode(Status{
			OK:            "OK",
			Authenticated: true,
			Service:       "Obsidian Local REST API",
			Versions:      StatusVersions{Self: "3.0.0", Obsidian: "1.0"},
		})
	}))
	t.Cleanup(server.Close)

	client := &Client{BaseURL: server.URL, APIKey: "test-key", HTTPClient: server.Client()}
	got, err := client.Ping(context.Background())
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if gotAuth != "Bearer test-key" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if !got.Authenticated || got.OK != "OK" || got.Versions.Self != "3.0.0" {
		t.Fatalf("status = %+v", got)
	}
}

func TestClientListDirectoryParsesFilesAndDirs(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/vault/" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"files": []string{"daily.md", "projects/"},
		})
	}))
	t.Cleanup(server.Close)

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	entries, err := client.ListDirectory(context.Background(), "")
	if err != nil {
		t.Fatalf("ListDirectory: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len = %d", len(entries))
	}
	if entries[0].Name != "daily.md" || entries[0].IsDir || entries[0].Path != "daily.md" {
		t.Fatalf("file = %+v", entries[0])
	}
	if entries[1].Name != "projects" || !entries[1].IsDir || entries[1].Path != "projects" {
		t.Fatalf("dir = %+v", entries[1])
	}
}

func TestClientListDirectoryEncodesNestedPath(t *testing.T) {
	t.Parallel()

	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(map[string]any{"files": []string{"note.md"}})
	}))
	t.Cleanup(server.Close)

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	entries, err := client.ListDirectory(context.Background(), "hello world")
	if err != nil {
		t.Fatalf("ListDirectory: %v", err)
	}
	if gotPath != "/vault/hello world/" {
		t.Fatalf("path = %q", gotPath)
	}
	if len(entries) != 1 || entries[0].Path != "hello world/note.md" {
		t.Fatalf("entries = %+v", entries)
	}
}

func TestClientReadNoteSetsAcceptAndPath(t *testing.T) {
	t.Parallel()

	var gotAccept, gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAccept = r.Header.Get("Accept")
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(Note{
			Path:    "notes/daily.md",
			Content: "# Hello",
			Tags:    []string{"inbox"},
			Stat:    NoteStat{Size: 7},
		})
	}))
	t.Cleanup(server.Close)

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	note, err := client.ReadNote(context.Background(), "notes/daily.md")
	if err != nil {
		t.Fatalf("ReadNote: %v", err)
	}
	if gotAccept != noteJSONAccept {
		t.Fatalf("Accept = %q", gotAccept)
	}
	if gotPath != "/vault/notes/daily.md" {
		t.Fatalf("path = %q", gotPath)
	}
	if note.Content != "# Hello" || len(note.Tags) != 1 {
		t.Fatalf("note = %+v", note)
	}
}

func TestClientWriteNotePutsMarkdown(t *testing.T) {
	t.Parallel()

	var gotPath, gotType, gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotType = r.Header.Get("Content-Type")
		raw, _ := io.ReadAll(r.Body)
		gotBody = string(raw)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	if err := client.WriteNote(context.Background(), "notes/daily.md", "updated"); err != nil {
		t.Fatalf("WriteNote: %v", err)
	}
	if gotPath != "/vault/notes/daily.md" {
		t.Fatalf("path = %q", gotPath)
	}
	if gotType != "text/markdown" {
		t.Fatalf("Content-Type = %q", gotType)
	}
	if gotBody != "updated" {
		t.Fatalf("body = %q", gotBody)
	}
}

func TestClientSearchSimple(t *testing.T) {
	t.Parallel()

	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		gotQuery = r.URL.Query().Get("query")
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"filename": "notes/daily.md",
				"score":    1.5,
				"matches": []map[string]any{
					{
						"context": "hello world",
						"match":   map[string]any{"start": 0, "end": 5},
					},
				},
			},
		})
	}))
	t.Cleanup(server.Close)

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	hits, err := client.SearchSimple(context.Background(), "hello", 80)
	if err != nil {
		t.Fatalf("SearchSimple: %v", err)
	}
	if gotQuery != "hello" {
		t.Fatalf("query = %q", gotQuery)
	}
	if len(hits) != 1 || hits[0].Filename != "notes/daily.md" || hits[0].Matches[0].End != 5 {
		t.Fatalf("hits = %+v", hits)
	}
}

func TestClientMapsJSONError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(apiError{ErrorCode: 40149, Message: "Invalid API key"})
	}))
	t.Cleanup(server.Close)

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := client.Ping(context.Background())
	if err == nil || !strings.Contains(err.Error(), "40149") || !strings.Contains(err.Error(), "Invalid API key") {
		t.Fatalf("error = %v", err)
	}
}

func TestClientRejectsParentPath(t *testing.T) {
	t.Parallel()

	client := &Client{BaseURL: "http://127.0.0.1"}
	_, err := client.ReadNote(context.Background(), "../secret.md")
	if err == nil || !strings.Contains(err.Error(), "..") {
		t.Fatalf("error = %v", err)
	}
}
