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

func setupSystemTray(app *application.App, wm *WindowManager) *application.SystemTray {
	systemTray := app.SystemTray.New()

	switch runtime.GOOS {
	case "windows":
		systemTray.SetIcon(appIconICO)
	default:
		systemTray.SetIcon(appIconJPG)
	}
	systemTray.OnClick(wm.ShowAgentWindow)

	menu := app.NewMenu()
	menu.Add("Agent").OnClick(func(_ *application.Context) {
		wm.ShowAgentWindow()
	})
	menu.Add("Manage").OnClick(func(_ *application.Context) {
		wm.ShowManagementWindow()
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
