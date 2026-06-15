package appservice

// WindowManager controls showing and focusing application windows.
type WindowManager interface {
	ShowAgentWindow()
	ShowManagementWindow()
	EmitAgentEvent(name string, data any)
}

// ConfigureWindows attaches the window manager for native window switching.
func ConfigureWindows(wm WindowManager) {
	globalWindowManager = wm
}

var globalWindowManager WindowManager

// ShowAgentWindow shows and focuses the agent chat window.
func (s *Service) ShowAgentWindow() {
	if globalWindowManager != nil {
		globalWindowManager.ShowAgentWindow()
	}
}

// ShowManagementWindow shows and focuses the management window.
func (s *Service) ShowManagementWindow() {
	if globalWindowManager != nil {
		globalWindowManager.ShowManagementWindow()
	}
}
