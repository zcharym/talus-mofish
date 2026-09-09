package services

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WindowManager controls showing and focusing application windows.
type WindowManager interface {
	ShowAgentWindow()
	ShowManagementWindow()
	EmitAgentEvent(name string, data any)
}

// AppEmitter broadcasts application-wide events to every window.
type AppEmitter interface {
	Emit(name string, data any)
}

type wailsAppEmitter struct {
	app *application.App
}

func (e wailsAppEmitter) Emit(name string, data any) {
	if e.app != nil {
		e.app.Event.Emit(name, data)
	}
}
