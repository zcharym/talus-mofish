package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	// EnvDevelopment is the default local/dev environment name.
	EnvDevelopment = "development"
	// EnvProduction is the production environment name.
	EnvProduction = "production"
)

// EnvName returns the active app environment.
// TALUS_ENV wins when set; otherwise defaults to development.
func EnvName() string {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("TALUS_ENV")))
	switch env {
	case EnvProduction, "prod":
		return EnvProduction
	case EnvDevelopment, "dev", "":
		return EnvDevelopment
	default:
		return env
	}
}

// LoadEnv loads environment files for the current environment.
// Order (first set wins; existing OS env vars are never overwritten):
//  1. .env.<TALUS_ENV>  — e.g. .env.development or .env.production
//  2. .env              — shared fallback
//
// Missing files are ignored. Call this once at process start, before config.Load.
func LoadEnv() error {
	env := EnvName()
	files := []string{
		".env." + env,
		".env",
	}

	for _, file := range files {
		if err := godotenv.Load(file); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("load %s: %w", file, err)
		}
	}
	return nil
}
