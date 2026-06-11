package appservice

import (
	"github.com/songwei.ma/talus-mofish/internal/autostart"
	"github.com/songwei.ma/talus-mofish/internal/config"
	"github.com/songwei.ma/talus-mofish/internal/database"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Service exposes application and persistence APIs to the frontend.
type Service struct {
	db        *database.DB
	config    *config.Store
	autostart *autostart.Manager
	wailsApp  *application.App
}

func New(db *database.DB, cfg *config.Store, autostartManager *autostart.Manager) *Service {
	return &Service{db: db, config: cfg, autostart: autostartManager}
}

// SetWailsApp attaches the Wails application for native dialogs.
func (s *Service) SetWailsApp(app *application.App) {
	s.wailsApp = app
}

// DatabasePath returns the on-disk SQLite database file path.
func (s *Service) DatabasePath() string {
	return s.db.Path
}
