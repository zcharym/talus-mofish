package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

type Region struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	W int `yaml:"w"`
	H int `yaml:"h"`
}

type Rule struct {
	ID      string   `yaml:"id"`
	AnyText []string `yaml:"any_text"`
	Regex   string   `yaml:"regex"`
	Severity string  `yaml:"severity"`
}

type Target struct {
	Name              string `yaml:"name"`
	ProcessName       string `yaml:"process_name"`
	WindowTitleRegex  string `yaml:"window_title_regex"`
	Region            Region `yaml:"region"`
	RequireVisible    *bool  `yaml:"require_visible"`
	Rules             []Rule `yaml:"rules"`

	titlePattern *regexp.Regexp
}

func (t *Target) TitlePattern() (*regexp.Regexp, error) {
	if t.titlePattern != nil {
		return t.titlePattern, nil
	}
	if t.WindowTitleRegex == "" {
		return nil, fmt.Errorf("target %q: window_title_regex is required", t.Name)
	}
	re, err := regexp.Compile(t.WindowTitleRegex)
	if err != nil {
		return nil, fmt.Errorf("target %q: invalid window_title_regex: %w", t.Name, err)
	}
	t.titlePattern = re
	return re, nil
}

func (t *Target) RequireVisibleWindow() bool {
	if t.RequireVisible == nil {
		return true
	}
	return *t.RequireVisible
}

type Config struct {
	WorkerURL    string        `yaml:"worker_url"`
	Secret       string        `yaml:"secret"`
	PollInterval time.Duration `yaml:"poll_interval"`
	Cooldown     time.Duration `yaml:"cooldown"`
	Targets      []Target      `yaml:"targets"`
}

func DefaultConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "TalusEcho", "echo-watch.yaml"), nil
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		PollInterval: 3 * time.Second,
		Cooldown:     5 * time.Minute,
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.WorkerURL == "" {
		return nil, fmt.Errorf("worker_url is required")
	}
	if cfg.Secret == "" {
		return nil, fmt.Errorf("secret is required")
	}
	if len(cfg.Targets) == 0 {
		return nil, fmt.Errorf("at least one target is required")
	}
	if cfg.PollInterval < time.Second {
		cfg.PollInterval = time.Second
	}

	for i := range cfg.Targets {
		if _, err := cfg.Targets[i].TitlePattern(); err != nil {
			return nil, err
		}
		for j := range cfg.Targets[i].Rules {
			rule := &cfg.Targets[i].Rules[j]
			if rule.ID == "" {
				return nil, fmt.Errorf("target %q rule %d: id is required", cfg.Targets[i].Name, j)
			}
			if len(rule.AnyText) == 0 && rule.Regex == "" {
				return nil, fmt.Errorf("target %q rule %q: any_text or regex required", cfg.Targets[i].Name, rule.ID)
			}
			if rule.Regex != "" {
				if _, err := regexp.Compile(rule.Regex); err != nil {
					return nil, fmt.Errorf("target %q rule %q: invalid regex: %w", cfg.Targets[i].Name, rule.ID, err)
				}
			}
		}
	}

	return cfg, nil
}
