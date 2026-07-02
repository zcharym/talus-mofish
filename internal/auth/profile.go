package auth

// UserProfile is the signed-in user exposed to the frontend.
type UserProfile struct {
	ID          string `json:"id"`
	Provider    string `json:"provider"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
}
