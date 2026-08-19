package vdiupload

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Region is a client-relative capture rectangle for detectors.
type Region struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	W int `yaml:"w"`
	H int `yaml:"h"`
}

// WindowMatch locates the VDI client top-level window.
type WindowMatch struct {
	ProcessNames       []string `yaml:"process_names"`
	TitleRegex         string   `yaml:"title_regex"`
	ClassName          string   `yaml:"class_name"`
	RequireVisible     *bool    `yaml:"require_visible"`
	RestoreIfMinimized bool     `yaml:"restore_if_minimized"`
}

func (m WindowMatch) RequireVisibleWindow() bool {
	if m.RequireVisible == nil {
		return true
	}
	return *m.RequireVisible
}

// ClipboardConfig controls CF_HDROP staging.
type ClipboardConfig struct {
	Format           string `yaml:"format"` // "hdrop"
	AlsoSetTextPaths bool   `yaml:"also_set_text_paths"`
	RestorePrevious  bool   `yaml:"restore_previous"`
	MaxFiles         int    `yaml:"max_files"`
	MaxTotalBytes    int64  `yaml:"max_total_bytes"`
}

// OCRWait is a completion detector rule.
type OCRWait struct {
	Region   Region        `yaml:"region"`
	AnyText  []string      `yaml:"any_text"`
	DenyText []string      `yaml:"deny_text"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Step is one profile automation action.
type Step struct {
	Focus   bool     `yaml:"focus"`
	Keys    []string `yaml:"keys"`
	WaitMs  int      `yaml:"wait_ms"`
	WaitOCR *OCRWait `yaml:"wait_ocr"`
}

// RetryConfig governs orchestrator retries.
type RetryConfig struct {
	Attempts  int           `yaml:"attempts"`
	BaseDelay time.Duration `yaml:"base_delay"`
	MaxDelay  time.Duration `yaml:"max_delay"`
}

// Profile is a named VDI client automation recipe.
type Profile struct {
	Window    WindowMatch     `yaml:"window"`
	Clipboard ClipboardConfig `yaml:"clipboard"`
	Steps     []Step          `yaml:"steps"`
	Retry     RetryConfig     `yaml:"retry"`
}

// Config is the top-level YAML document.
type Config struct {
	DefaultProfile string             `yaml:"default_profile"`
	LogLevel       string             `yaml:"log_level"`
	MutexName      string             `yaml:"mutex_name"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

// DefaultConfigPath returns %AppData%/TalusEcho/vdi-upload.yaml (or OS equivalent).
func DefaultConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "TalusEcho", "vdi-upload.yaml"), nil
}

// LoadConfig reads and validates YAML config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		LogLevel:  "info",
		MutexName: `Local\TalusEcho.VDIUpload`,
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks required fields and applies defaults.
func (c *Config) Validate() error {
	if len(c.Profiles) == 0 {
		return fmt.Errorf("at least one profile is required")
	}
	if c.DefaultProfile == "" {
		for name := range c.Profiles {
			c.DefaultProfile = name
			break
		}
	}
	if _, ok := c.Profiles[c.DefaultProfile]; !ok {
		return fmt.Errorf("default_profile %q not found", c.DefaultProfile)
	}
	if c.MutexName == "" {
		c.MutexName = `Local\TalusEcho.VDIUpload`
	}

	for name, p := range c.Profiles {
		if len(p.Window.ProcessNames) == 0 && p.Window.TitleRegex == "" {
			return fmt.Errorf("profile %q: process_names or title_regex required", name)
		}
		if p.Clipboard.Format == "" {
			p.Clipboard.Format = "hdrop"
		}
		if p.Clipboard.Format != "hdrop" {
			return fmt.Errorf("profile %q: unsupported clipboard format %q", name, p.Clipboard.Format)
		}
		if p.Clipboard.MaxFiles == 0 {
			p.Clipboard.MaxFiles = 20
		}
		if p.Clipboard.MaxTotalBytes == 0 {
			p.Clipboard.MaxTotalBytes = 500 * 1024 * 1024
		}
		if p.Retry.Attempts == 0 {
			p.Retry.Attempts = 3
		}
		if p.Retry.BaseDelay == 0 {
			p.Retry.BaseDelay = time.Second
		}
		if p.Retry.MaxDelay == 0 {
			p.Retry.MaxDelay = 30 * time.Second
		}
		if len(p.Steps) == 0 {
			return fmt.Errorf("profile %q: at least one step is required", name)
		}
		c.Profiles[name] = p
	}
	return nil
}

// ProfileOrDefault returns a named profile or the default.
func (c *Config) ProfileOrDefault(name string) (Profile, error) {
	if name == "" {
		name = c.DefaultProfile
	}
	p, ok := c.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found", name)
	}
	return p, nil
}
