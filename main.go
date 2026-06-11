package main

import (
	"embed"
	"log"

	"github.com/songwei.ma/talus-mofish/internal/appservice"
	"github.com/songwei.ma/talus-mofish/internal/autostart"
	"github.com/songwei.ma/talus-mofish/internal/config"
	"github.com/songwei.ma/talus-mofish/internal/database"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	db, err := database.OpenDefault()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	configPath, err := database.DefaultConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	autostartManager := autostart.New(autostart.DefaultIdentifier)
	if err := autostartManager.Sync(cfg.Get().AutoStart); err != nil {
		log.Printf("apply autostart: %v", err)
	}

	appService := NewAppService(db, cfg, autostartManager)

	app := application.New(application.Options{
		Name:        "talus-mofish",
		Description: "English learning companion",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler: newAssetHandler(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		Linux: application.LinuxOptions{
			DisableQuitOnLastWindowClosed: true,
		},
	})

	appService.SetWailsApp(app)

	windowManager := NewWindowManager(app)
	windowManager.CreateWindows()
	appservice.ConfigureWindows(windowManager)

	setupSystemTray(app, windowManager)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
