package main

import (
	"github.com/songwei.ma/talus-mofish/internal/appservice"
	"github.com/songwei.ma/talus-mofish/internal/autostart"
	"github.com/songwei.ma/talus-mofish/internal/config"
	"github.com/songwei.ma/talus-mofish/internal/database"
)

// AppService is the Wails-bound facade over internal/appservice.
// Methods are promoted from the embedded Service so frontend bindings stay stable.
type AppService struct {
	*appservice.Service
}

func NewAppService(db *database.DB, cfg *config.Store, autostartManager *autostart.Manager) *AppService {
	return &AppService{Service: appservice.New(db, cfg, autostartManager)}
}
