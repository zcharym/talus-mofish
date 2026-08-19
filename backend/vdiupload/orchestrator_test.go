package vdiupload

import (
	"testing"
	"time"
)

func TestConfigValidateDefaults(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"h3c": {
				Window: WindowMatch{TitleRegex: "(?i)workspace"},
				Steps:  []Step{{Focus: true}},
			},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.DefaultProfile != "h3c" {
		t.Fatalf("default profile: %q", cfg.DefaultProfile)
	}
	p := cfg.Profiles["h3c"]
	if p.Clipboard.Format != "hdrop" {
		t.Fatalf("format: %q", p.Clipboard.Format)
	}
	if p.Retry.Attempts != 3 {
		t.Fatalf("attempts: %d", p.Retry.Attempts)
	}
	if p.Retry.BaseDelay != time.Second {
		t.Fatalf("base delay: %v", p.Retry.BaseDelay)
	}
}

func TestErrorCodeExit(t *testing.T) {
	if CodeOK.ExitCode() != 0 {
		t.Fatal()
	}
	if CodeConfigInvalid.ExitCode() != 3 {
		t.Fatal()
	}
	if !CodeWindowNotFound.Retryable() {
		t.Fatal()
	}
	if CodeConfigInvalid.Retryable() {
		t.Fatal()
	}
}

func TestBackoff(t *testing.T) {
	r := RetryConfig{BaseDelay: time.Second, MaxDelay: 5 * time.Second}
	if got := backoff(r, 1); got != time.Second {
		t.Fatalf("got %v", got)
	}
	if got := backoff(r, 4); got != 5*time.Second {
		t.Fatalf("got %v", got)
	}
}
