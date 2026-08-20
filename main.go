package main

import (
	"embed"
	"log"

	"github.com/songwei.ma/talus-mofish/backend/services"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/utils"
	"github.com/songwei.ma/talus-mofish/backend/utils/autostart"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := utils.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	db, err := storage.OpenDefault()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	configPath, err := storage.DefaultConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := storage.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	autostartManager := autostart.New(autostart.DefaultIdentifier)
	if err := autostartManager.Sync(cfg.Get().AutoStart); err != nil {
		log.Printf("apply autostart: %v", err)
	}

	systemSvc := services.NewSystemService(db, autostartManager)
	configSvc := services.NewConfigService(db, cfg, autostartManager)
	authSvc := services.NewAuthService(db, cfg)
	chatSvc := services.NewChatService(db, cfg)
	englishSvc := services.NewEnglishService(db)
	sudokuSvc := services.NewSudokuService(db, cfg)

	app := application.New(application.Options{
		Name:        "talus-mofish",
		Description: "Chat-oriented desktop agent",
		Services: []application.Service{
			application.NewService(systemSvc),
			application.NewService(configSvc),
			application.NewService(authSvc),
			application.NewService(chatSvc),
			application.NewService(englishSvc),
			application.NewService(sudokuSvc),
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

	windowManager := NewWindowManager(app)
	windowManager.CreateWindows()
	services.WireSystemRuntime(systemSvc, app, windowManager)
	services.WireChatRuntime(chatSvc, app, windowManager)

	setupSystemTray(app, windowManager)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
