package sudoku

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientGeneratePostsExpectedBody(t *testing.T) {
	t.Parallel()

	puzzle := strings.Repeat("0", 80) + "1"
	solution := strings.Repeat("123456789", 9)

	var gotKey string
	var gotBody generateRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("x-api-key")
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Errorf("decode body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(generateResponse{
			Difficulty: "medium",
			Puzzle:     puzzle,
			Solution:   solution,
		})
	}))
	t.Cleanup(server.Close)

	client := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}
	got, err := client.Generate(context.Background(), "medium")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if gotKey != "test-key" {
		t.Fatalf("x-api-key = %q", gotKey)
	}
	if !gotBody.Solution || gotBody.Array || gotBody.Difficulty != "medium" {
		t.Fatalf("request body = %+v", gotBody)
	}
	if got.Puzzle != puzzle || got.Solution != solution || got.Difficulty != "medium" {
		t.Fatalf("puzzle = %+v", got)
	}
}

func TestClientGenerateUnauthorized(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := client.Generate(context.Background(), "easy")
	if err == nil || !strings.Contains(err.Error(), "youdosudoku.com") {
		t.Fatalf("error = %v", err)
	}
}
