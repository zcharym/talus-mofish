package main

import (
	_ "embed"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.jpg
var appIconJPG []byte

//go:embed build/windows/icon.ico
var appIconICO []byte

const navigateToConfigEvent = "frontend:navigate"

func setupSystemTray(app *application.App, window *application.WebviewWindow) *application.SystemTray {
	systemTray := app.SystemTray.New()

	switch runtime.GOOS {
	case "windows":
		systemTray.SetIcon(appIconICO)
	default:
		systemTray.SetIcon(appIconJPG)
	}
	systemTray.SetLabel("Talus MoFish")

	menu := app.NewMenu()
	menu.Add("Config").OnClick(func(_ *application.Context) {
		window.Show()
		window.Focus()
		app.Event.Emit(navigateToConfigEvent, "config")
	})
	menu.AddSeparator()
	menu.Add("Exit").OnClick(func(_ *application.Context) {
		app.Quit()
	})
	systemTray.SetMenu(menu)

	app.OnShutdown(func() {
		systemTray.Destroy()
	})

	return systemTray
}
