package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/songwei.ma/talus-mofish/internal/store"
)

const (
	magicLinkPollInterval = 2 * time.Second
	magicLinkTimeout      = 10 * time.Minute
)

var (
	errAuthOffline          = errors.New("auth server unreachable")
	errAuthSessionInvalid   = errors.New("auth session invalid")
	errAuthServerNotConfigured = errors.New("auth server URL is not configured")
)

type magicLinkRequest struct {
	Email string `json:"email"`
}

type magicLinkResponse struct {
	RequestID string `json:"requestId"`
	ExpiresAt string `json:"expiresAt"`
}

type magicLinkStatusResponse struct {
	Status      string          `json:"status"`
	AccessToken string          `json:"accessToken"`
	User        *remoteAuthUser `json:"user"`
	Error       string          `json:"error"`
}

type remoteAuthUser struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type meResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

func (s *Service) SignInWithEmail(ctx context.Context, email string) (*UserProfile, error) {
	baseURL, err := s.authServerURL()
	if err != nil {
		return nil, err
	}

	normalized, err := normalizeEmailAddress(email)
	if err != nil {
		return nil, err
	}

	client := s.authHTTPClient()
	requestID, err := s.requestMagicLink(ctx, client, baseURL, normalized)
	if err != nil {
		return nil, err
	}

	remoteUser, accessToken, err := s.pollMagicLinkStatus(ctx, client, baseURL, requestID)
	if err != nil {
		return nil, err
	}

	return s.persistEmailUser(ctx, remoteUser, accessToken)
}

func (s *Service) authServerURL() (string, error) {
	baseURL := strings.TrimSpace(os.Getenv("TALUS_AUTH_SERVER_URL"))
	if baseURL == "" {
		baseURL = strings.TrimSpace(s.config.Get().Auth.AuthServerURL)
	}
	if baseURL == "" {
		return "", errAuthServerNotConfigured
	}
	return strings.TrimRight(baseURL, "/"), nil
}

func (s *Service) authHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: oauthTransport(),
	}
}

func normalizeEmailAddress(email string) (string, error) {
	address, err := mail.ParseAddress(strings.TrimSpace(email))
	if err != nil {
		return "", fmt.Errorf("invalid email address")
	}
	return strings.ToLower(address.Address), nil
}

func (s *Service) requestMagicLink(
	ctx context.Context,
	client *http.Client,
	baseURL, email string,
) (string, error) {
	body, err := json.Marshal(magicLinkRequest{Email: email})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/v1/auth/magic-link",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request magic link: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("request magic link: %s", strings.TrimSpace(string(respBody)))
	}

	var payload magicLinkResponse
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return "", fmt.Errorf("decode magic link response: %w", err)
	}
	if payload.RequestID == "" {
		return "", fmt.Errorf("magic link response missing request id")
	}
	return payload.RequestID, nil
}

func (s *Service) pollMagicLinkStatus(
	ctx context.Context,
	client *http.Client,
	baseURL, requestID string,
) (*remoteAuthUser, string, error) {
	deadline := time.Now().Add(magicLinkTimeout)
	ticker := time.NewTicker(magicLinkPollInterval)
	defer ticker.Stop()

	for {
		user, token, done, err := s.fetchMagicLinkStatus(ctx, client, baseURL, requestID)
		if err != nil {
			return nil, "", err
		}
		if done {
			if user == nil || token == "" {
				return nil, "", fmt.Errorf("magic link verified without user session")
			}
			return user, token, nil
		}

		if time.Now().After(deadline) {
			return nil, "", fmt.Errorf("magic link sign-in timed out after %s", magicLinkTimeout)
		}

		select {
		case <-ctx.Done():
			return nil, "", ctx.Err()
		case <-ticker.C:
		}
	}
}

func (s *Service) fetchMagicLinkStatus(
	ctx context.Context,
	client *http.Client,
	baseURL, requestID string,
) (*remoteAuthUser, string, bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v1/auth/magic-link/status/%s", baseURL, requestID),
		nil,
	)
	if err != nil {
		return nil, "", false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", false, fmt.Errorf("poll magic link status: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 {
		return nil, "", false, fmt.Errorf("poll magic link status: %s", strings.TrimSpace(string(body)))
	}

	var payload magicLinkStatusResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, "", false, fmt.Errorf("decode magic link status: %w", err)
	}

	switch payload.Status {
	case "verified":
		return payload.User, payload.AccessToken, true, nil
	case "expired":
		return nil, "", false, fmt.Errorf("magic link expired; request a new sign-in email")
	case "pending", "":
		return nil, "", false, nil
	default:
		return nil, "", false, fmt.Errorf("unexpected magic link status: %s", payload.Status)
	}
}

func (s *Service) persistEmailUser(
	ctx context.Context,
	remote *remoteAuthUser,
	accessToken string,
) (*UserProfile, error) {
	if remote == nil || remote.ID == "" {
		return nil, fmt.Errorf("remote user is missing")
	}

	displayName := strings.TrimSpace(remote.DisplayName)
	if displayName == "" {
		displayName = remote.Email
	}

	if err := s.queries.DeleteAllUserAccounts(ctx); err != nil {
		return nil, fmt.Errorf("clear existing user: %w", err)
	}

	if err := s.queries.InsertUserAccount(ctx, store.InsertUserAccountParams{
		ID:             remote.ID,
		Provider:       ProviderEmail,
		ProviderUserID: remote.ID,
		DisplayName:    displayName,
		Email:          remote.Email,
		AvatarUrl:      "",
	}); err != nil {
		return nil, fmt.Errorf("insert user account: %w", err)
	}

	if err := clearTokens(); err != nil {
		return nil, err
	}
	if err := saveAuthSession(ProviderEmail, accessToken); err != nil {
		return nil, err
	}

	return &UserProfile{
		ID:          remote.ID,
		Provider:    ProviderEmail,
		DisplayName: displayName,
		Email:       remote.Email,
		AvatarURL:   "",
	}, nil
}

func (s *Service) refreshEmailUser(ctx context.Context, cached *UserProfile) (*UserProfile, error) {
	_, token, err := loadAuthSession()
	if err != nil || token == "" {
		return cached, errAuthOffline
	}

	baseURL, err := s.authServerURL()
	if err != nil {
		return cached, errAuthOffline
	}

	client := s.authHTTPClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/v1/auth/me", nil)
	if err != nil {
		return cached, errAuthOffline
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return cached, errAuthOffline
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errAuthSessionInvalid
	}
	if resp.StatusCode >= 400 {
		return cached, errAuthOffline
	}

	var remote meResponse
	if err := json.NewDecoder(resp.Body).Decode(&remote); err != nil {
		return cached, errAuthOffline
	}

	displayName := strings.TrimSpace(remote.DisplayName)
	if displayName == "" {
		displayName = remote.Email
	}

	if err := s.queries.DeleteAllUserAccounts(ctx); err != nil {
		return nil, fmt.Errorf("clear existing user: %w", err)
	}
	if err := s.queries.InsertUserAccount(ctx, store.InsertUserAccountParams{
		ID:             remote.ID,
		Provider:       ProviderEmail,
		ProviderUserID: remote.ID,
		DisplayName:    displayName,
		Email:          remote.Email,
		AvatarUrl:      "",
	}); err != nil {
		return nil, fmt.Errorf("refresh user account: %w", err)
	}

	return &UserProfile{
		ID:          remote.ID,
		Provider:    ProviderEmail,
		DisplayName: displayName,
		Email:       remote.Email,
		AvatarURL:   "",
	}, nil
}
