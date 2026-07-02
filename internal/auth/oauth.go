package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-ieproxy"
	"github.com/pkg/browser"
	"github.com/songwei.ma/talus-mofish/internal/config"
	"github.com/songwei.ma/talus-mofish/internal/store"
	"golang.org/x/oauth2"
)

const (
	ProviderGitHub = "github"
	ProviderGoogle = "google"
	callbackPath   = "/callback"
)

// GitHubLoopbackCallbackURL is the callback URL to register on the OAuth app.
// GitHub allows any port on 127.0.0.1 when this path is configured.
const GitHubLoopbackCallbackURL = "http://127.0.0.1/callback"

type Service struct {
	queries *store.Queries
	config  *config.Store
}

func New(queries *store.Queries, cfg *config.Store) *Service {
	return &Service{
		queries: queries,
		config:  cfg,
	}
}

func (s *Service) oauthHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   60 * time.Second,
		Transport: oauthTransport(s.config.Get().OAuth.HTTPProxy),
	}
}

func oauthTransport(configuredProxy string) *http.Transport {
	proxyFunc := ieproxy.GetProxyFunc()
	if configuredProxy = strings.TrimSpace(configuredProxy); configuredProxy != "" {
		proxyURL, err := url.Parse(configuredProxy)
		if err == nil {
			proxyFunc = http.ProxyURL(proxyURL)
		}
	}

	return &http.Transport{
		Proxy: proxyFunc,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 15 * time.Second,
		ForceAttemptHTTP2:   true,
	}
}

func (s *Service) oauthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, s.oauthHTTPClient())
}

type oauthResult struct {
	code  string
	state string
	err   error
}

func (s *Service) GetCurrentUser(ctx context.Context) (*UserProfile, error) {
	row, err := s.queries.GetCurrentUser(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get current user: %w", err)
	}
	return rowToProfile(row), nil
}

func (s *Service) SignInWithGitHub(ctx context.Context) (*UserProfile, error) {
	oauthCfg := s.config.Get().OAuth
	if strings.TrimSpace(oauthCfg.GitHubClientID) == "" {
		return nil, fmt.Errorf("GitHub client ID is not configured")
	}
	if strings.TrimSpace(oauthCfg.GitHubClientSecret) == "" {
		return nil, fmt.Errorf("GitHub client secret is not configured")
	}

	conf := &oauth2.Config{
		ClientID:     oauthCfg.GitHubClientID,
		ClientSecret: oauthCfg.GitHubClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		Scopes: []string{"user:email"},
	}

	token, err := s.runOAuthLoop(ctx, conf, ProviderGitHub)
	if err != nil {
		return nil, err
	}

	client := s.oauthHTTPClient()
	profile, err := fetchGitHubProfile(ctx, client, token.AccessToken)
	if err != nil {
		return nil, err
	}

	return s.persistUser(ctx, profile, token.AccessToken, token.RefreshToken)
}

func (s *Service) SignInWithGoogle(ctx context.Context) (*UserProfile, error) {
	oauthCfg := s.config.Get().OAuth
	if strings.TrimSpace(oauthCfg.GoogleClientID) == "" {
		return nil, fmt.Errorf("Google client ID is not configured")
	}
	if strings.TrimSpace(oauthCfg.GoogleClientSecret) == "" {
		return nil, fmt.Errorf("Google client secret is not configured")
	}

	conf := &oauth2.Config{
		ClientID:     oauthCfg.GoogleClientID,
		ClientSecret: oauthCfg.GoogleClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		Scopes: []string{"openid", "email", "profile"},
	}

	token, err := s.runOAuthLoop(ctx, conf, ProviderGoogle)
	if err != nil {
		return nil, err
	}

	client := s.oauthHTTPClient()
	profile, err := fetchGoogleProfile(ctx, client, token.AccessToken)
	if err != nil {
		return nil, err
	}

	return s.persistUser(ctx, profile, token.AccessToken, token.RefreshToken)
}

func (s *Service) SignOut(ctx context.Context) error {
	if err := s.queries.DeleteAllUserAccounts(ctx); err != nil {
		return fmt.Errorf("delete user accounts: %w", err)
	}
	if err := clearTokens(); err != nil {
		return err
	}
	return nil
}

func (s *Service) persistUser(
	ctx context.Context,
	profile providerProfile,
	accessToken, refreshToken string,
) (*UserProfile, error) {
	if err := s.queries.DeleteAllUserAccounts(ctx); err != nil {
		return nil, fmt.Errorf("clear existing user: %w", err)
	}

	userID := uuid.NewString()
	if err := s.queries.InsertUserAccount(ctx, store.InsertUserAccountParams{
		ID:             userID,
		Provider:       profile.Provider,
		ProviderUserID: profile.ProviderUserID,
		DisplayName:    profile.DisplayName,
		Email:          profile.Email,
		AvatarUrl:      profile.AvatarURL,
	}); err != nil {
		return nil, fmt.Errorf("insert user account: %w", err)
	}

	if err := saveTokens(profile.Provider, accessToken, refreshToken); err != nil {
		return nil, err
	}

	return &UserProfile{
		ID:          userID,
		Provider:    profile.Provider,
		DisplayName: profile.DisplayName,
		Email:       profile.Email,
		AvatarURL:   profile.AvatarURL,
	}, nil
}

func (s *Service) runOAuthLoop(ctx context.Context, conf *oauth2.Config, provider string) (*oauth2.Token, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("start oauth callback listener: %w", err)
	}
	defer listener.Close()

	redirectURL := fmt.Sprintf("http://%s%s", listener.Addr().String(), callbackPath)
	oauthConf := *conf
	oauthConf.RedirectURL = redirectURL

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, err
	}

	state, err := generateState()
	if err != nil {
		return nil, err
	}

	authOpts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", codeChallenge(codeVerifier)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	}
	if provider == ProviderGoogle {
		authOpts = append(authOpts, oauth2.AccessTypeOffline)
	}

	authURL := oauthConf.AuthCodeURL(state, authOpts...)

	resultCh := make(chan oauthResult, 1)
	server := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != callbackPath {
				http.NotFound(w, r)
				return
			}

			if errMsg := r.URL.Query().Get("error"); errMsg != "" {
				resultCh <- oauthResult{err: fmt.Errorf("oauth denied: %s", errMsg)}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte(oauthFailureHTML(errMsg)))
				return
			}

			if r.URL.Query().Get("state") != state {
				resultCh <- oauthResult{err: fmt.Errorf("oauth state mismatch")}
				http.Error(w, "Invalid OAuth state.", http.StatusBadRequest)
				return
			}

			code := r.URL.Query().Get("code")
			if code == "" {
				resultCh <- oauthResult{err: fmt.Errorf("oauth callback missing code")}
				http.Error(w, "Missing authorization code.", http.StatusBadRequest)
				return
			}

			resultCh <- oauthResult{code: code, state: state}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(oauthSuccessHTML(provider)))
		}),
	}

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			resultCh <- oauthResult{err: fmt.Errorf("oauth callback server: %w", err)}
		}
	}()

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	if err := browser.OpenURL(authURL); err != nil {
		return nil, fmt.Errorf("open browser for %s sign-in: %w", provider, err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultCh:
		if result.err != nil {
			return nil, result.err
		}

		token, err := oauthConf.Exchange(
			s.oauthContext(ctx),
			result.code,
			oauth2.SetAuthURLParam("code_verifier", codeVerifier),
		)
		if err != nil {
			return nil, fmt.Errorf("exchange oauth code: %w", wrapOAuthExchangeError(provider, err))
		}
		return token, nil
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("oauth sign-in timed out")
	}
}

type providerProfile struct {
	Provider       string
	ProviderUserID string
	DisplayName    string
	Email          string
	AvatarURL      string
}

func rowToProfile(row store.UserAccount) *UserProfile {
	return &UserProfile{
		ID:          row.ID,
		Provider:    row.Provider,
		DisplayName: row.DisplayName,
		Email:       row.Email,
		AvatarURL:   row.AvatarUrl,
	}
}

func fetchGitHubProfile(ctx context.Context, client *http.Client, accessToken string) (providerProfile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return providerProfile{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return providerProfile{}, fmt.Errorf("fetch github profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return providerProfile{}, fmt.Errorf("fetch github profile: %s", strings.TrimSpace(string(body)))
	}

	var payload struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return providerProfile{}, fmt.Errorf("decode github profile: %w", err)
	}

	displayName := strings.TrimSpace(payload.Name)
	if displayName == "" {
		displayName = payload.Login
	}

	email := strings.TrimSpace(payload.Email)
	if email == "" {
		email, _ = fetchGitHubPrimaryEmail(ctx, client, accessToken)
	}

	return providerProfile{
		Provider:       ProviderGitHub,
		ProviderUserID: fmt.Sprintf("%d", payload.ID),
		DisplayName:    displayName,
		Email:          email,
		AvatarURL:      payload.AvatarURL,
	}, nil
}

func fetchGitHubPrimaryEmail(ctx context.Context, client *http.Client, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, item := range emails {
		if item.Primary && item.Verified {
			return item.Email, nil
		}
	}
	for _, item := range emails {
		if item.Verified {
			return item.Email, nil
		}
	}
	return "", nil
}

func fetchGoogleProfile(ctx context.Context, client *http.Client, accessToken string) (providerProfile, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo",
		nil,
	)
	if err != nil {
		return providerProfile{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return providerProfile{}, fmt.Errorf("fetch google profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return providerProfile{}, fmt.Errorf("fetch google profile: %s", strings.TrimSpace(string(body)))
	}

	var payload struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return providerProfile{}, fmt.Errorf("decode google profile: %w", err)
	}

	displayName := strings.TrimSpace(payload.Name)
	if displayName == "" {
		displayName = strings.TrimSpace(payload.Email)
	}

	return providerProfile{
		Provider:       ProviderGoogle,
		ProviderUserID: payload.ID,
		DisplayName:    displayName,
		Email:          payload.Email,
		AvatarURL:      payload.Picture,
	}, nil
}

func oauthSuccessHTML(provider string) string {
	label := provider
	if label != "" {
		label = strings.ToUpper(label[:1]) + label[1:]
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><title>Signed in</title></head>
<body style="font-family: system-ui, sans-serif; padding: 2rem; text-align: center;">
  <h1>Signed in with %s</h1>
  <p>You can close this tab and return to Talus Agent.</p>
</body>
</html>`, label)
}

func oauthFailureHTML(message string) string {
	safe := html.EscapeString(message)
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><title>Sign-in failed</title></head>
<body style="font-family: system-ui, sans-serif; padding: 2rem; text-align: center;">
  <h1>Sign-in failed</h1>
  <p>%s</p>
</body>
</html>`, safe)
}

func wrapOAuthExchangeError(provider string, err error) error {
	message := err.Error()
	if !strings.Contains(message, "connect") &&
		!strings.Contains(message, "timeout") &&
		!strings.Contains(message, "connection refused") {
		return err
	}

	target := "the provider"
	if provider == ProviderGoogle {
		target = "Google (oauth2.googleapis.com)"
	}
	if provider == ProviderGitHub {
		target = "GitHub"
	}

	return fmt.Errorf(
		"%w: could not reach %s from the app. Your browser may use a system proxy, but the app needs the same route — enable Windows proxy settings or set OAuth HTTP proxy in Config (e.g. http://127.0.0.1:7890)",
		err,
		target,
	)
}
