package services

import (
	"context"

	"github.com/songwei.ma/talus-mofish/backend/auth"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/types"
)

// AuthService exposes sign-in and sign-out APIs.
type AuthService struct {
	auth *auth.Service
}

// NewAuthService creates the auth Wails service.
func NewAuthService(db *storage.DB, cfg *storage.ConfigStore) *AuthService {
	return &AuthService{auth: auth.New(db.Queries, cfg)}
}

func toUserProfile(profile *auth.UserProfile) *types.UserProfile {
	if profile == nil {
		return nil
	}
	return &types.UserProfile{
		ID:          profile.ID,
		Provider:    profile.Provider,
		DisplayName: profile.DisplayName,
		Email:       profile.Email,
		AvatarURL:   profile.AvatarURL,
	}
}

// GetCurrentUser returns the active signed-in user, or nil when signed out.
func (s *AuthService) GetCurrentUser() (*types.UserProfile, error) {
	profile, err := s.auth.GetCurrentUser(context.Background())
	if err != nil {
		return nil, err
	}
	return toUserProfile(profile), nil
}

// SignInWithGitHub runs the GitHub OAuth flow and persists the signed-in user.
func (s *AuthService) SignInWithGitHub() (*types.UserProfile, error) {
	profile, err := s.auth.SignInWithGitHub(context.Background())
	if err != nil {
		return nil, err
	}
	return toUserProfile(profile), nil
}

// SignInWithGoogle runs the Google OAuth flow and persists the signed-in user.
func (s *AuthService) SignInWithGoogle() (*types.UserProfile, error) {
	profile, err := s.auth.SignInWithGoogle(context.Background())
	if err != nil {
		return nil, err
	}
	return toUserProfile(profile), nil
}

// SignInWithEmail starts the email magic-link flow and persists the signed-in user.
func (s *AuthService) SignInWithEmail(email string) (*types.UserProfile, error) {
	profile, err := s.auth.SignInWithEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return toUserProfile(profile), nil
}

// SignInAsDebug creates a local admin user when debugMode is enabled.
func (s *AuthService) SignInAsDebug() (*types.UserProfile, error) {
	profile, err := s.auth.EnsureDebugUser(context.Background())
	if err != nil {
		return nil, err
	}
	return toUserProfile(profile), nil
}

// SignOut clears the signed-in user and stored OAuth tokens.
func (s *AuthService) SignOut() error {
	return s.auth.SignOut(context.Background())
}
