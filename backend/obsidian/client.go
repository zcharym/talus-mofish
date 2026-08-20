package obsidian

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/songwei.ma/talus-mofish/backend/types"
)

const (
	noteJSONAccept     = "application/vnd.olrapi.note+json"
	defaultSearchCtx   = 100
	maxResponseBytes   = 8 << 20
	defaultHTTPTimeout = 30 * time.Second
)

// Client talks to Obsidian Local REST API.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient builds a client that never uses the system HTTP proxy.
// Loopback HTTPS skips TLS verification because the plugin uses a self-signed cert.
func NewClient(baseURL, apiKey string) *Client {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = types.DefaultObsidianBaseURL
	}
	return &Client{
		BaseURL: baseURL,
		APIKey:  strings.TrimSpace(apiKey),
		HTTPClient: &http.Client{
			Timeout:   defaultHTTPTimeout,
			Transport: loopbackTransport(baseURL),
		},
	}
}

func loopbackTransport(baseURL string) *http.Transport {
	return &http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 15 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: shouldSkipTLS(baseURL), //nolint:gosec // plugin ships a self-signed localhost cert
		},
	}
}

func shouldSkipTLS(baseURL string) bool {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	if !strings.EqualFold(parsed.Scheme, "https") {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	return host == "127.0.0.1" || host == "localhost" || host == "::1"
}

func (c *Client) httpClient() *http.Client {
	if c != nil && c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: defaultHTTPTimeout}
}

func (c *Client) baseURL() string {
	if c == nil {
		return types.DefaultObsidianBaseURL
	}
	base := strings.TrimRight(strings.TrimSpace(c.BaseURL), "/")
	if base == "" {
		return types.DefaultObsidianBaseURL
	}
	return base
}

func (c *Client) apiKey() string {
	if c == nil {
		return ""
	}
	return strings.TrimSpace(c.APIKey)
}

// Ping calls GET / and reports whether the bearer token is accepted.
func (c *Client) Ping(ctx context.Context) (Status, error) {
	if c == nil {
		return Status{}, fmt.Errorf("obsidian client is nil")
	}
	raw, err := c.do(ctx, http.MethodGet, "/", nil, nil)
	if err != nil {
		return Status{}, err
	}
	var status Status
	if err := json.Unmarshal(raw, &status); err != nil {
		return Status{}, fmt.Errorf("decode obsidian status: %w", err)
	}
	return status, nil
}

// ListDirectory lists files in a vault folder. Empty path is the vault root.
func (c *Client) ListDirectory(ctx context.Context, dir string) ([]FileEntry, error) {
	if c == nil {
		return nil, fmt.Errorf("obsidian client is nil")
	}
	dir, err := NormalizeVaultPath(dir)
	if err != nil {
		return nil, err
	}

	apiPath := "/vault/"
	if dir != "" {
		apiPath = "/vault/" + EncodeVaultPath(dir) + "/"
	}

	raw, err := c.do(ctx, http.MethodGet, apiPath, nil, nil)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Files []string `json:"files"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("decode vault listing: %w", err)
	}

	entries := make([]FileEntry, 0, len(parsed.Files))
	for _, name := range parsed.Files {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		isDir := strings.HasSuffix(name, "/")
		entries = append(entries, FileEntry{
			Name:  strings.TrimSuffix(name, "/"),
			Path:  JoinVaultPath(dir, name),
			IsDir: isDir,
		})
	}
	return entries, nil
}

// ReadNote fetches a markdown note as structured JSON.
func (c *Client) ReadNote(ctx context.Context, vaultPath string) (Note, error) {
	if c == nil {
		return Note{}, fmt.Errorf("obsidian client is nil")
	}
	vaultPath, err := NormalizeVaultPath(vaultPath)
	if err != nil {
		return Note{}, err
	}
	if vaultPath == "" {
		return Note{}, fmt.Errorf("note path is required")
	}
	if !IsMarkdown(vaultPath) {
		return Note{}, fmt.Errorf("not a markdown note")
	}

	headers := map[string]string{"Accept": noteJSONAccept}
	raw, err := c.do(ctx, http.MethodGet, "/vault/"+EncodeVaultPath(vaultPath), headers, nil)
	if err != nil {
		return Note{}, err
	}

	var note Note
	if err := json.Unmarshal(raw, &note); err != nil {
		return Note{}, fmt.Errorf("decode note: %w", err)
	}
	if note.Path == "" {
		note.Path = vaultPath
	}
	if note.Tags == nil {
		note.Tags = []string{}
	}
	return note, nil
}

// WriteNote creates or overwrites a markdown file.
func (c *Client) WriteNote(ctx context.Context, vaultPath, content string) error {
	if c == nil {
		return fmt.Errorf("obsidian client is nil")
	}
	vaultPath, err := NormalizeVaultPath(vaultPath)
	if err != nil {
		return err
	}
	if vaultPath == "" {
		return fmt.Errorf("note path is required")
	}
	if !IsMarkdown(vaultPath) {
		return fmt.Errorf("not a markdown note")
	}

	headers := map[string]string{"Content-Type": "text/markdown"}
	_, err = c.do(ctx, http.MethodPut, "/vault/"+EncodeVaultPath(vaultPath), headers, []byte(content))
	return err
}

// SearchSimple runs the plugin's built-in full-text search.
func (c *Client) SearchSimple(ctx context.Context, query string, contextLength int) ([]SearchHit, error) {
	if c == nil {
		return nil, fmt.Errorf("obsidian client is nil")
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}
	if contextLength <= 0 {
		contextLength = defaultSearchCtx
	}

	params := url.Values{}
	params.Set("query", query)
	params.Set("contextLength", fmt.Sprintf("%d", contextLength))
	apiPath := "/search/simple/?" + params.Encode()

	raw, err := c.do(ctx, http.MethodPost, apiPath, nil, nil)
	if err != nil {
		return nil, err
	}

	var wire []searchHitWire
	if err := json.Unmarshal(raw, &wire); err != nil {
		return nil, fmt.Errorf("decode search results: %w", err)
	}

	hits := make([]SearchHit, 0, len(wire))
	for _, item := range wire {
		matches := make([]SearchMatch, 0, len(item.Matches))
		for _, match := range item.Matches {
			matches = append(matches, SearchMatch{
				Context: match.Context,
				Start:   int(match.Match.Start),
				End:     int(match.Match.End),
			})
		}
		hits = append(hits, SearchHit{
			Filename: item.Filename,
			Score:    item.Score,
			Matches:  matches,
		})
	}
	return hits, nil
}

type searchHitWire struct {
	Filename string  `json:"filename"`
	Score    float64 `json:"score"`
	Matches  []struct {
		Context string `json:"context"`
		Match   struct {
			Start float64 `json:"start"`
			End   float64 `json:"end"`
		} `json:"match"`
	} `json:"matches"`
}

type apiError struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

func (c *Client) do(ctx context.Context, method, apiPath string, headers map[string]string, body []byte) ([]byte, error) {
	fullURL := c.baseURL() + apiPath
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reader)
	if err != nil {
		return nil, fmt.Errorf("create obsidian request: %w", err)
	}
	if key := c.apiKey(); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("call obsidian API: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("read obsidian response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, mapResponseError(resp.StatusCode, raw)
	}
	return raw, nil
}

func mapResponseError(status int, body []byte) error {
	var ae apiError
	if json.Unmarshal(body, &ae) == nil && ae.Message != "" {
		if ae.ErrorCode != 0 {
			return fmt.Errorf("obsidian API %d (%d): %s", status, ae.ErrorCode, ae.Message)
		}
		return fmt.Errorf("obsidian API %d: %s", status, ae.Message)
	}
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		return fmt.Errorf("obsidian API returned %d", status)
	}
	return fmt.Errorf("obsidian API %d: %s", status, msg)
}
