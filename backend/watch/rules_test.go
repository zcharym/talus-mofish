package watch

import (
	"os"
	"testing"
	"time"
)

func TestRuleEngineEvaluate(t *testing.T) {
	engine := NewRuleEngine(time.Minute)
	target := Target{
		Name: "build",
		Rules: []Rule{{
			ID:      "success",
			AnyText: []string{"successful"},
		}},
	}

	match, ok := engine.Evaluate(target, "Build finished: Successful")
	if !ok || match == nil {
		t.Fatalf("expected match")
	}
	if match.Rule.ID != "success" {
		t.Fatalf("got rule %q", match.Rule.ID)
	}

	_, ok = engine.Evaluate(target, "Build finished: Successful")
	if ok {
		t.Fatalf("expected cooldown to suppress duplicate")
	}
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/echo-watch.yaml"
	content := `
worker_url: https://example.workers.dev
secret: test-token
poll_interval: 2s
cooldown: 1m
targets:
  - name: vdi-build-result
    process_name: mstsc.exe
    window_title_regex: ".*Remote Desktop.*"
    region: { x: 0, y: 0, w: 100, h: 50 }
    rules:
      - id: success
        any_text: ["successful"]
`
	if err := writeFile(path, content); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PollInterval != 2*time.Second {
		t.Fatalf("poll interval: %s", cfg.PollInterval)
	}
	if !cfg.Targets[0].RequireVisibleWindow() {
		t.Fatal("expected require_visible default true")
	}
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600)
}
