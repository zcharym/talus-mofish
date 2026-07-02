package appservice

import (
	"context"

	"github.com/songwei.ma/talus-mofish/internal/auth"
)

// UserProfile is exposed to the frontend for the signed-in user.
type UserProfile struct {
	ID          string `json:"id"`
	Provider    string `json:"provider"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
}

func toUserProfile(profile *auth.UserProfile) *UserProfile {
	if profile == nil {
		return nil
	}
	return &UserProfile{
		ID:          profile.ID,
		Provider:    profile.Provider,
		DisplayName: profile.DisplayName,
		Email:       profile.Email,
		AvatarURL:   profile.AvatarURL,
	}
}

// GetCurrentUser returns the active signed-in user, or nil when signed out.
func (s *Service) GetCurrentUser() (*UserProfile, error) {
	profile, err := s.auth.GetCurrentUser(context.Background())
	if err != nil {
		return nil, err
	}
	return toUserProfile(profile), nil
}

// SignInWithGitHub runs the GitHub OAuth flow and persists the signed-in user.
func (s *Service) SignInWithGitHub() (*UserProfile, error) {
	profile, err := s.auth.SignInWithGitHub(context.Background())
	if err != nil {
		return nil, err
	}
	return toUserProfile(profile), nil
}

// SignInWithGoogle runs the Google OAuth flow and persists the signed-in user.
func (s *Service) SignInWithGoogle() (*UserProfile, error) {
	profile, err := s.auth.SignInWithGoogle(context.Background())
	if err != nil {
		return nil, err
	}
	return toUserProfile(profile), nil
}

// SignOut clears the signed-in user and stored OAuth tokens.
func (s *Service) SignOut() error {
	return s.auth.SignOut(context.Background())
}
