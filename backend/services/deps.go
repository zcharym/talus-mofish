package services

// WindowManager controls showing and focusing application windows.
type WindowManager interface {
	ShowAgentWindow()
	ShowManagementWindow()
	EmitAgentEvent(name string, data any)
}
