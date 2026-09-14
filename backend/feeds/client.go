package feeds

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

const defaultUserAgent = "TalusMofish/1.0 (+https://github.com/songwei.ma/talus-mofish)"

var (
	strictHTML = bluemonday.StrictPolicy()
	spaceRE    = strings.NewReplacer("\n", " ", "\t", " ", "\r", " ")
)

// Client talks to RSS, YouTube, and Bilibili over HTTP.
type Client struct {
	http   *http.Client
	now    func() time.Time
	parser *gofeed.Parser
}

// NewClient builds a feed client. httpClient may be nil.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 20 * time.Second}
	}
	parser := gofeed.NewParser()
	parser.UserAgent = defaultUserAgent
	parser.Client = httpClient
	return &Client{
		http:   httpClient,
		now:    time.Now,
		parser: parser,
	}
}

func (c *Client) clock() time.Time {
	if c != nil && c.now != nil {
		return c.now()
	}
	return time.Now()
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if c == nil || c.http == nil {
		return nil, fmt.Errorf("feeds HTTP client is not configured")
	}
	return c.http.Do(req.WithContext(ctx))
}

func (c *Client) getJSON(ctx context.Context, rawURL, cookie string, extraHeaders map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	for key, value := range extraHeaders {
		req.Header.Set(key, value)
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET %s: HTTP %d", rawURL, resp.StatusCode)
	}
	return body, nil
}

func sanitizeSummary(raw string) string {
	text := strings.TrimSpace(spaceRE.Replace(strictHTML.Sanitize(raw)))
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}
	return truncateRunes(text, 400)
}

func truncateRunes(s string, max int) string {
	if max <= 0 || s == "" {
		return s
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max]) + "…"
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
