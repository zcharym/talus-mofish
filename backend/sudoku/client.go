package sudoku

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mattn/go-ieproxy"
)

const DefaultAPIURL = "https://youdosudoku.com/api/"

// Puzzle is a generated board from YouDoSudoku.
type Puzzle struct {
	Difficulty string
	Puzzle     string
	Solution   string
}

type generateRequest struct {
	Difficulty string `json:"difficulty"`
	Solution   bool   `json:"solution"`
	Array      bool   `json:"array"`
}

type generateResponse struct {
	Difficulty string `json:"difficulty"`
	Puzzle     string `json:"puzzle"`
	Solution   string `json:"solution"`
}

// Client fetches puzzles from the YouDoSudoku HTTP API.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient builds a client that uses the system proxy, matching OAuth.
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: DefaultAPIURL,
		APIKey:  strings.TrimSpace(apiKey),
		HTTPClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: youDoTransport(),
		},
	}
}

func youDoTransport() *http.Transport {
	return &http.Transport{
		Proxy: ieproxy.GetProxyFunc(),
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 15 * time.Second,
		ForceAttemptHTTP2:   true,
	}
}

// Generate requests a unique puzzle and its solution.
func (c *Client) Generate(ctx context.Context, difficulty string) (Puzzle, error) {
	if c == nil {
		return Puzzle{}, fmt.Errorf("sudoku client is nil")
	}
	difficulty, err := NormalizeDifficulty(difficulty)
	if err != nil {
		return Puzzle{}, err
	}

	body, err := json.Marshal(generateRequest{
		Difficulty: difficulty,
		Solution:   true,
		Array:      false,
	})
	if err != nil {
		return Puzzle{}, fmt.Errorf("encode YouDoSudoku request: %w", err)
	}

	baseURL := strings.TrimSpace(c.BaseURL)
	if baseURL == "" {
		baseURL = DefaultAPIURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return Puzzle{}, fmt.Errorf("create YouDoSudoku request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if key := strings.TrimSpace(c.APIKey); key != "" {
		req.Header.Set("x-api-key", key)
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return Puzzle{}, fmt.Errorf("call YouDoSudoku: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Puzzle{}, fmt.Errorf("read YouDoSudoku response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return Puzzle{}, fmt.Errorf("YouDoSudoku API unauthorized (%d); add an API key from https://www.youdosudoku.com/ in Configuration", resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Puzzle{}, fmt.Errorf("YouDoSudoku API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var parsed generateResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return Puzzle{}, fmt.Errorf("decode YouDoSudoku response: %w", err)
	}

	puzzle, err := NormalizeGrid(parsed.Puzzle)
	if err != nil {
		return Puzzle{}, fmt.Errorf("puzzle: %w", err)
	}
	solution, err := NormalizeGrid(parsed.Solution)
	if err != nil {
		return Puzzle{}, fmt.Errorf("solution: %w", err)
	}

	outDifficulty := parsed.Difficulty
	if outDifficulty == "" {
		outDifficulty = difficulty
	} else if normalized, err := NormalizeDifficulty(outDifficulty); err == nil {
		outDifficulty = normalized
	}

	return Puzzle{
		Difficulty: outDifficulty,
		Puzzle:     puzzle,
		Solution:   solution,
	}, nil
}
