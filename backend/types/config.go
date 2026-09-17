package types

import (
	"strings"

	"github.com/songwei.ma/talus-mofish/backend/utils/aiclient"
)

// App holds user-facing application settings persisted in config.json.
type App struct {
	Theme            string          `json:"theme"`
	DailyGoalMinutes int             `json:"dailyGoalMinutes"` // English Learning domain
	WordsPerSession  int             `json:"wordsPerSession"`  // English Learning domain
	AutoStart        bool            `json:"autoStart"`
	DebugMode        bool            `json:"debugMode"`
	AI               aiclient.Config `json:"ai"`
	Auth             Auth            `json:"auth"`
	OAuth            OAuth           `json:"oauth"`
	Sudoku           Sudoku          `json:"sudoku"`
	Obsidian         Obsidian        `json:"obsidian"`
	Cloudflare       Cloudflare      `json:"cloudflare"`
	Feeds            Feeds           `json:"feeds"`
}

// Auth holds settings for email magic-link authentication.
type Auth struct {
	AuthServerURL string `json:"authServerUrl"`
}

// OAuth holds OAuth client credentials for third-party sign-in providers.
type OAuth struct {
	GitHubClientID     string `json:"githubClientId"`
	GitHubClientSecret string `json:"githubClientSecret"`
	GoogleClientID     string `json:"googleClientId"`
	GoogleClientSecret string `json:"googleClientSecret"`
}

// Sudoku holds optional YouDoSudoku API credentials.
type Sudoku struct {
	APIKey string `json:"apiKey"`
}

// DefaultObsidianBaseURL is the Local REST API HTTPS endpoint (self-signed cert).
const DefaultObsidianBaseURL = "https://127.0.0.1:27124"

// Obsidian holds Local REST API connection settings.
type Obsidian struct {
	BaseURL string `json:"baseUrl"`
	APIKey  string `json:"apiKey"`
}

// NormalizeObsidian fills the default plugin URL when unset.
func NormalizeObsidian(o Obsidian) Obsidian {
	o.BaseURL = strings.TrimSpace(o.BaseURL)
	if o.BaseURL == "" {
		o.BaseURL = DefaultObsidianBaseURL
	}
	o.APIKey = strings.TrimSpace(o.APIKey)
	return o
}

// Cloudflare holds Account API credentials for the Agent dashboard.
type Cloudflare struct {
	AccountID string `json:"accountId"`
	APIToken  string `json:"apiToken"`
}

// NormalizeCloudflare trims stored Cloudflare credentials.
func NormalizeCloudflare(c Cloudflare) Cloudflare {
	c.AccountID = strings.TrimSpace(c.AccountID)
	c.APIToken = strings.TrimSpace(c.APIToken)
	return c
}

// CloudflareConfigured reports whether both an account ID and API token are set.
func CloudflareConfigured(c Cloudflare) bool {
	c = NormalizeCloudflare(c)
	return c.AccountID != "" && c.APIToken != ""
}

// Feeds holds optional credentials and HTTP proxy for YouTube / Bilibili / RSS.
type Feeds struct {
	YouTubeAPIKey    string `json:"youtubeApiKey"`
	BilibiliSESSDATA string `json:"bilibiliSessdata"`
	// ProxyURL is an optional HTTP(S) proxy for all Feeds fetches (YouTube, RSS, Bilibili).
	// Empty means use the OS/system proxy (Windows Internet Options / env). Example: http://127.0.0.1:7890
	ProxyURL string `json:"proxyUrl"`
}

// NormalizeFeeds trims stored feed credentials and proxy URL.
func NormalizeFeeds(f Feeds) Feeds {
	f.YouTubeAPIKey = strings.TrimSpace(f.YouTubeAPIKey)
	f.BilibiliSESSDATA = strings.TrimSpace(f.BilibiliSESSDATA)
	f.ProxyURL = strings.TrimSpace(f.ProxyURL)
	return f
}
