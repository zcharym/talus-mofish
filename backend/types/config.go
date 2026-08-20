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

// Normalize fills the default plugin URL when unset.
func (o Obsidian) Normalize() Obsidian {
	o.BaseURL = strings.TrimSpace(o.BaseURL)
	if o.BaseURL == "" {
		o.BaseURL = DefaultObsidianBaseURL
	}
	o.APIKey = strings.TrimSpace(o.APIKey)
	return o
}
