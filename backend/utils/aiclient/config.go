package aiclient

import (
	"net/url"
	"strings"
)

// Provider identifies an LLM backend.
type Provider string

const (
	ProviderOpenAI   Provider = "openai"
	ProviderDeepSeek Provider = "deepseek"
	ProviderMoonshot Provider = "moonshot"
	ProviderOllama   Provider = "ollama"
)

// Config holds connection settings for the chat API client.
type Config struct {
	Provider Provider `json:"provider"`
	Model    string   `json:"model"`
	APIKey   string   `json:"apiKey"`
	BaseURL  string   `json:"baseURL"`
}

// DefaultConfig returns sensible defaults for a new installation.
func DefaultConfig() Config {
	return Config{
		Provider: ProviderOpenAI,
		Model:    "gpt-4o-mini",
	}
}

// Normalize fills empty fields with provider defaults.
func (c Config) Normalize() Config {
	if c.Provider == "" {
		c.Provider = ProviderOpenAI
	}
	if c.Model == "" {
		switch c.Provider {
		case ProviderDeepSeek:
			c.Model = "deepseek-chat"
		case ProviderMoonshot:
			c.Model = "kimi-k2.5"
		case ProviderOllama:
			c.Model = "llama3.2"
		default:
			c.Model = "gpt-4o-mini"
		}
	}
	c.BaseURL = strings.TrimRight(strings.TrimSpace(c.BaseURL), "/")
	return c
}

// Endpoint returns the chat-completions URL for this provider and base URL.
func (c Config) Endpoint() string {
	return joinChatCompletionsURL(c.resolveAPIRoot())
}

func (c Config) resolveAPIRoot() string {
	c = c.Normalize()
	if c.BaseURL != "" {
		if root, ok := moonshotAPIRoot(c.BaseURL); ok {
			return root
		}
		return c.BaseURL
	}
	switch c.Provider {
	case ProviderDeepSeek:
		return "https://api.deepseek.com/v1"
	case ProviderMoonshot:
		return "https://api.moonshot.cn/v1"
	case ProviderOllama:
		return "http://localhost:11434/v1"
	default:
		return "https://api.openai.com/v1"
	}
}

// moonshotAPIRoot normalizes Moonshot/Kimi hosts to the OpenAI-compatible /v1 root.
// Users sometimes paste /anthropic paths from other provider docs.
func moonshotAPIRoot(baseURL string) (string, bool) {
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Host == "" {
		return "", false
	}
	host := strings.ToLower(parsed.Hostname())
	if host != "api.moonshot.cn" && !strings.HasSuffix(host, ".moonshot.cn") {
		return "", false
	}
	return "https://api.moonshot.cn/v1", true
}

func joinChatCompletionsURL(apiRoot string) string {
	apiRoot = strings.TrimRight(apiRoot, "/")
	if strings.HasSuffix(apiRoot, "/chat/completions") {
		return apiRoot
	}
	if strings.HasSuffix(apiRoot, "/v1") {
		return apiRoot + "/chat/completions"
	}
	return apiRoot + "/v1/chat/completions"
}

// ChatTemperature returns the temperature to send for chat completions.
// Moonshot/Kimi models only accept 1; other providers default to 0.7 when unset.
func (c Config) ChatTemperature(requested float64) float64 {
	c = c.Normalize()
	if c.requiresFixedTemperature() {
		return 1
	}
	if requested <= 0 {
		return 0.7
	}
	return requested
}

func (c Config) requiresFixedTemperature() bool {
	if c.Provider == ProviderMoonshot {
		return true
	}
	model := strings.ToLower(c.Model)
	if strings.HasPrefix(model, "kimi-") || strings.HasPrefix(model, "moonshot-") {
		return true
	}
	if _, ok := moonshotAPIRoot(c.BaseURL); ok {
		return true
	}
	return false
}

// RequiresAPIKey reports whether the provider expects an API key.
func (c Config) RequiresAPIKey() bool {
	return c.Normalize().Provider != ProviderOllama
}

// Validate checks that the config is usable for a chat request.
func (c Config) Validate() error {
	c = c.Normalize()
	if c.Model == "" {
		return ErrModelRequired
	}
	if c.RequiresAPIKey() && strings.TrimSpace(c.APIKey) == "" {
		return ErrAPIKeyRequired
	}
	return nil
}
