package types

import "github.com/songwei.ma/talus-mofish/backend/utils/aiclient"

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
