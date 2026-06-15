package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const (
	windowNameManagement = "management"
	windowNameAgent      = "agent"
)

// WindowManager tracks the management and agent windows for show/focus operations.
type WindowManager struct {
	app              *application.App
	managementWindow *application.WebviewWindow
	agentWindow      *application.WebviewWindow
}

func NewWindowManager(app *application.App) *WindowManager {
	return &WindowManager{app: app}
}

func (wm *WindowManager) CreateWindows() {
	wm.managementWindow = wm.createWindow(windowNameManagement, "Talus Echo — Manage", "/")
	wm.agentWindow = wm.createWindow(windowNameAgent, "Talus Agent", "/agent")

	wm.managementWindow.Hide()
	wm.agentWindow.Show()
}

func (wm *WindowManager) createWindow(name, title, url string) *application.WebviewWindow {
	window := wm.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  name,
		Title: title,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              url,
	})

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	return window
}

func (wm *WindowManager) ShowAgentWindow() {
	if wm.agentWindow != nil {
		wm.agentWindow.Show()
		wm.agentWindow.Focus()
	}
}

func (wm *WindowManager) ShowManagementWindow() {
	if wm.managementWindow != nil {
		wm.managementWindow.Show()
		wm.managementWindow.Focus()
	}
}

// EmitAgentEvent sends an event to the agent chat window only.
func (wm *WindowManager) EmitAgentEvent(name string, data any) {
	if wm.agentWindow != nil {
		wm.agentWindow.EmitEvent(name, data)
	}
}
