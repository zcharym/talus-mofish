package aiclient

import "testing"

func TestConfigNormalizeAndEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		in       Config
		wantURL  string
		wantModel string
	}{
		{
			name:      "openai default",
			in:        Config{Provider: ProviderOpenAI},
			wantURL:   "https://api.openai.com/v1/chat/completions",
			wantModel: "gpt-4o-mini",
		},
		{
			name:      "deepseek default",
			in:        Config{Provider: ProviderDeepSeek},
			wantURL:   "https://api.deepseek.com/v1/chat/completions",
			wantModel: "deepseek-chat",
		},
		{
			name:      "ollama default",
			in:        Config{Provider: ProviderOllama, Model: "mistral"},
			wantURL:   "http://localhost:11434/v1/chat/completions",
			wantModel: "mistral",
		},
		{
			name:      "custom base url with v1",
			in:        Config{Provider: ProviderOpenAI, BaseURL: "https://proxy.example/v1"},
			wantURL:   "https://proxy.example/v1/chat/completions",
			wantModel: "gpt-4o-mini",
		},
		{
			name:      "moonshot default",
			in:        Config{Provider: ProviderMoonshot},
			wantURL:   "https://api.moonshot.cn/v1/chat/completions",
			wantModel: "kimi-k2.5",
		},
		{
			name: "moonshot wrong anthropic base url",
			in: Config{
				Provider: ProviderOpenAI,
				Model:    "kimi-k2.5",
				BaseURL:  "https://api.moonshot.cn/anthropic",
			},
			wantURL:   "https://api.moonshot.cn/v1/chat/completions",
			wantModel: "kimi-k2.5",
		},
		{
			name:      "custom host without v1 suffix",
			in:        Config{Provider: ProviderOpenAI, BaseURL: "https://proxy.example"},
			wantURL:   "https://proxy.example/v1/chat/completions",
			wantModel: "gpt-4o-mini",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Normalize()
			if got.Endpoint() != tt.wantURL {
				t.Fatalf("Endpoint() = %q, want %q", got.Endpoint(), tt.wantURL)
			}
			if got.Model != tt.wantModel {
				t.Fatalf("Model = %q, want %q", got.Model, tt.wantModel)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	if err := (Config{Provider: ProviderOpenAI, Model: "gpt-4o-mini"}).Validate(); err != ErrAPIKeyRequired {
		t.Fatalf("Validate() = %v, want ErrAPIKeyRequired", err)
	}
	if err := (Config{Provider: ProviderOllama, Model: "llama3.2"}).Validate(); err != nil {
		t.Fatalf("Validate() ollama = %v", err)
	}
}

func TestChatTemperature(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want float64
	}{
		{
			name: "openai default",
			cfg:  Config{Provider: ProviderOpenAI, Model: "gpt-4o-mini"},
			want: 0.7,
		},
		{
			name: "kimi model with openai provider",
			cfg:  Config{Provider: ProviderOpenAI, Model: "kimi-k2.5", BaseURL: "https://api.moonshot.cn/anthropic"},
			want: 1,
		},
		{
			name: "moonshot provider",
			cfg:  Config{Provider: ProviderMoonshot, Model: "kimi-k2.5"},
			want: 1,
		},
		{
			name: "kimi ignores requested temperature",
			cfg:  Config{Provider: ProviderOpenAI, Model: "kimi-k2.5", BaseURL: "https://api.moonshot.cn/v1"},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.ChatTemperature(0.7); got != tt.want {
				t.Fatalf("ChatTemperature() = %v, want %v", got, tt.want)
			}
		})
	}
}
