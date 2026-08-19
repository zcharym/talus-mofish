package auth

import (
	"context"
	"fmt"
)

const (
	// ProviderDebug is the local identity used when debugMode is enabled.
	ProviderDebug = "debug"

	// DebugToken is the sentinel token stored for the debug admin user.
	// It is not an OAuth credential — it only marks the session as debug.
	DebugToken = "debug-token"

	debugProviderUserID = "admin"
	debugDisplayName    = "Debug Admin"
	debugEmail          = "debug@localhost"
)

// EnsureDebugUser creates (or replaces with) a local admin user for debug mode.
// It is a no-op error when debugMode is disabled.
func (s *Service) EnsureDebugUser(ctx context.Context) (*UserProfile, error) {
	if !s.config.Get().DebugMode {
		return nil, fmt.Errorf("debug sign-in requires debugMode in config")
	}

	return s.persistUser(ctx, providerProfile{
		Provider:       ProviderDebug,
		ProviderUserID: debugProviderUserID,
		DisplayName:    debugDisplayName,
		Email:          debugEmail,
	}, DebugToken, "")
}
