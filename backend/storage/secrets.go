package storage

import (
	"log"
	"sync"

	"github.com/zalando/go-keyring"
)

const configSecretService = "talus-mofish-config"

const (
	secretAIAPIKey           = "ai.apiKey"
	secretGitHubClientSecret = "oauth.githubClientSecret"
	secretGoogleClientSecret = "oauth.googleClientSecret"
	secretSudokuAPIKey       = "sudoku.apiKey"
	secretObsidianAPIKey     = "obsidian.apiKey"
	secretCloudflareAPIToken = "cloudflare.apiToken"
	secretYouTubeAPIKey      = "feeds.youtubeApiKey"
	secretBilibiliSESSDATA   = "feeds.bilibiliSessdata"
)

// SecretBackend stores config secrets outside config.json.
type SecretBackend interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type keyringSecrets struct{}

func (keyringSecrets) Get(key string) (string, error) {
	value, err := keyring.Get(configSecretService, key)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (keyringSecrets) Set(key, value string) error {
	if value == "" {
		return keyringSecrets{}.Delete(key)
	}
	return keyring.Set(configSecretService, key, value)
}

func (keyringSecrets) Delete(key string) error {
	err := keyring.Delete(configSecretService, key)
	if err == keyring.ErrNotFound {
		return nil
	}
	return err
}

type memorySecrets struct {
	mu   sync.Mutex
	data map[string]string
}

func newMemorySecrets() *memorySecrets {
	return &memorySecrets{data: map[string]string{}}
}

func (m *memorySecrets) Get(key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[key], nil
}

func (m *memorySecrets) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if value == "" {
		delete(m.data, key)
		return nil
	}
	m.data[key] = value
	return nil
}

func (m *memorySecrets) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

var (
	secretsMu      sync.Mutex
	secretBackend  SecretBackend = keyringSecrets{}
	secretsEnabled               = true
)

// UseMemorySecrets swaps in an in-process store (tests).
func UseMemorySecrets() {
	secretsMu.Lock()
	defer secretsMu.Unlock()
	secretBackend = newMemorySecrets()
	secretsEnabled = true
}

func persistConfigSecrets(aiKey, githubSecret, googleSecret, sudokuKey, obsidianKey, cloudflareToken, youtubeKey, biliSess string) {
	secretsMu.Lock()
	backend := secretBackend
	enabled := secretsEnabled
	secretsMu.Unlock()
	if !enabled || backend == nil {
		return
	}

	pairs := []struct {
		key   string
		value string
	}{
		{secretAIAPIKey, aiKey},
		{secretGitHubClientSecret, githubSecret},
		{secretGoogleClientSecret, googleSecret},
		{secretSudokuAPIKey, sudokuKey},
		{secretObsidianAPIKey, obsidianKey},
		{secretCloudflareAPIToken, cloudflareToken},
		{secretYouTubeAPIKey, youtubeKey},
		{secretBilibiliSESSDATA, biliSess},
	}
	for _, pair := range pairs {
		if err := backend.Set(pair.key, pair.value); err != nil {
			log.Printf("store config secret %s: %v (integration tokens will not persist)", pair.key, err)
			secretsMu.Lock()
			secretsEnabled = false
			secretsMu.Unlock()
			return
		}
	}
}

func loadConfigSecret(key string) string {
	secretsMu.Lock()
	backend := secretBackend
	enabled := secretsEnabled
	secretsMu.Unlock()
	if !enabled || backend == nil {
		return ""
	}
	value, err := backend.Get(key)
	if err != nil {
		log.Printf("read config secret %s: %v", key, err)
		secretsMu.Lock()
		secretsEnabled = false
		secretsMu.Unlock()
		return ""
	}
	return value
}

func overlayConfigSecrets(aiKey, githubSecret, googleSecret, sudokuKey, obsidianKey, cloudflareToken, youtubeKey, biliSess *string) {
	if overlay := loadConfigSecret(secretAIAPIKey); overlay != "" {
		*aiKey = overlay
	}
	if overlay := loadConfigSecret(secretGitHubClientSecret); overlay != "" {
		*githubSecret = overlay
	}
	if overlay := loadConfigSecret(secretGoogleClientSecret); overlay != "" {
		*googleSecret = overlay
	}
	if overlay := loadConfigSecret(secretSudokuAPIKey); overlay != "" {
		*sudokuKey = overlay
	}
	if overlay := loadConfigSecret(secretObsidianAPIKey); overlay != "" {
		*obsidianKey = overlay
	}
	if overlay := loadConfigSecret(secretCloudflareAPIToken); overlay != "" {
		*cloudflareToken = overlay
	}
	if overlay := loadConfigSecret(secretYouTubeAPIKey); overlay != "" {
		*youtubeKey = overlay
	}
	if overlay := loadConfigSecret(secretBilibiliSESSDATA); overlay != "" {
		*biliSess = overlay
	}
}

func stripConfigSecretsForDisk(aiKey, githubSecret, googleSecret, sudokuKey, obsidianKey, cloudflareToken, youtubeKey, biliSess *string) {
	*aiKey = ""
	*githubSecret = ""
	*googleSecret = ""
	*sudokuKey = ""
	*obsidianKey = ""
	*cloudflareToken = ""
	*youtubeKey = ""
	*biliSess = ""
}
