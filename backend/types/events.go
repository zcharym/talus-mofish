package types

const EventConfigChanged = "config:changed"

// ConfigChanged is broadcast to every window after SaveConfig.
type ConfigChanged struct {
	Theme                string `json:"theme"`
	DebugMode            bool   `json:"debugMode"`
	CloudflareConfigured bool   `json:"cloudflareConfigured"`
}
