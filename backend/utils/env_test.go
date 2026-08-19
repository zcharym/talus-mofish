package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvDevelopment(t *testing.T) {
	dir := chdirTemp(t)
	t.Setenv("TALUS_ENV", "development")
	unsetEnv(t, "TALUS_AUTH_SERVER_URL")

	content := "TALUS_AUTH_SERVER_URL=https://dev.example.com\n"
	if err := os.WriteFile(filepath.Join(dir, ".env.development"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := LoadEnv(); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("TALUS_AUTH_SERVER_URL"); got != "https://dev.example.com" {
		t.Fatalf("TALUS_AUTH_SERVER_URL = %q, want https://dev.example.com", got)
	}
	if got := EnvName(); got != EnvDevelopment {
		t.Fatalf("EnvName() = %q, want %q", got, EnvDevelopment)
	}
}

func TestLoadEnvProduction(t *testing.T) {
	dir := chdirTemp(t)
	t.Setenv("TALUS_ENV", "production")
	unsetEnv(t, "TALUS_AUTH_SERVER_URL")

	content := "TALUS_AUTH_SERVER_URL=https://prod.example.com\n"
	if err := os.WriteFile(filepath.Join(dir, ".env.production"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := LoadEnv(); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("TALUS_AUTH_SERVER_URL"); got != "https://prod.example.com" {
		t.Fatalf("TALUS_AUTH_SERVER_URL = %q, want https://prod.example.com", got)
	}
	if got := EnvName(); got != EnvProduction {
		t.Fatalf("EnvName() = %q, want %q", got, EnvProduction)
	}
}

func TestLoadEnvDoesNotOverrideExisting(t *testing.T) {
	dir := chdirTemp(t)
	t.Setenv("TALUS_ENV", "development")
	t.Setenv("TALUS_AUTH_SERVER_URL", "https://already.set")

	content := "TALUS_AUTH_SERVER_URL=https://from.file\n"
	if err := os.WriteFile(filepath.Join(dir, ".env.development"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := LoadEnv(); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("TALUS_AUTH_SERVER_URL"); got != "https://already.set" {
		t.Fatalf("TALUS_AUTH_SERVER_URL = %q, want https://already.set", got)
	}
}

func chdirTemp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prevWD)
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return dir
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	prev, had := os.LookupEnv(key)
	_ = os.Unsetenv(key)
	t.Cleanup(func() {
		if had {
			_ = os.Setenv(key, prev)
			return
		}
		_ = os.Unsetenv(key)
	})
}
