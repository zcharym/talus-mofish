package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
)

const authSessionKey = "auth_session"

type storedAuthSession struct {
	AccessToken string `json:"access_token"`
	Provider    string `json:"provider"`
}

func saveAuthSession(provider, accessToken string) error {
	payload, err := json.Marshal(storedAuthSession{
		AccessToken: accessToken,
		Provider:    provider,
	})
	if err != nil {
		return fmt.Errorf("marshal auth session: %w", err)
	}
	if err := keyring.Set(keyringService, authSessionKey, string(payload)); err != nil {
		return fmt.Errorf("store auth session: %w", err)
	}
	return nil
}

func loadAuthSession() (string, string, error) {
	value, err := keyring.Get(keyringService, authSessionKey)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", "", nil
		}
		return "", "", fmt.Errorf("load auth session: %w", err)
	}
	var session storedAuthSession
	if err := json.Unmarshal([]byte(value), &session); err != nil {
		return "", "", fmt.Errorf("parse auth session: %w", err)
	}
	return session.Provider, session.AccessToken, nil
}

func clearAuthSession() error {
	if err := keyring.Delete(keyringService, authSessionKey); err != nil {
		if err == keyring.ErrNotFound {
			return nil
		}
		return fmt.Errorf("clear auth session: %w", err)
	}
	return nil
}
