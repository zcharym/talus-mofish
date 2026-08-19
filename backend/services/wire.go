package services

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WireSystemRuntime attaches Wails and window dependencies after app startup.
func WireSystemRuntime(s *SystemService, app *application.App, wm WindowManager) {
	s.wailsApp = app
	s.windows = wm
}

// WireChatRuntime attaches Wails and window dependencies after app startup.
func WireChatRuntime(s *ChatService, app *application.App, wm WindowManager) {
	s.wailsApp = app
	s.windows = wm
}
