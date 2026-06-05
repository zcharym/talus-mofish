package main

import (
	"embed"
	"log"

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

	appService := NewAppService(db, cfg)

	app := application.New(application.Options{
		Name:        "talus-mofish",
		Description: "English learning companion",
		Services: []application.Service{
			application.NewService(appService),
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Talus Echo",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
