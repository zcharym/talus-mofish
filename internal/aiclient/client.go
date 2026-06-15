package aiclient

import "context"

// Client streams chat completions from an LLM provider.
type Client interface {
	Config() Config
	ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamPart, error)
}

// NewClient builds a provider-backed client from config.
func NewClient(cfg Config) (Client, error) {
	cfg = cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newOpenAICompatibleClient(cfg), nil
}
