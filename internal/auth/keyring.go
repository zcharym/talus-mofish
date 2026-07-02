package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
)

const keyringService = "talus-mofish"

type storedTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Provider     string `json:"provider"`
}

func saveTokens(provider, accessToken, refreshToken string) error {
	payload, err := json.Marshal(storedTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Provider:     provider,
	})
	if err != nil {
		return fmt.Errorf("marshal oauth tokens: %w", err)
	}
	if err := keyring.Set(keyringService, "oauth_tokens", string(payload)); err != nil {
		return fmt.Errorf("store oauth tokens: %w", err)
	}
	return nil
}

func clearTokens() error {
	if err := keyring.Delete(keyringService, "oauth_tokens"); err != nil {
		if err == keyring.ErrNotFound {
			return nil
		}
		return fmt.Errorf("clear oauth tokens: %w", err)
	}
	return nil
}
