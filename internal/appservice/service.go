package appservice

import (
	"context"
	"fmt"

	"github.com/songwei.ma/talus-mofish/internal/agent"
	"github.com/songwei.ma/talus-mofish/internal/auth"
	"github.com/songwei.ma/talus-mofish/internal/autostart"
	"github.com/songwei.ma/talus-mofish/internal/config"
	"github.com/songwei.ma/talus-mofish/internal/database"
	"github.com/songwei.ma/talus-mofish/internal/store"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Service exposes application and persistence APIs to the frontend.
//
// Domain-specific APIs are grouped by file:
//   - English Learning (domain.English): import.go, articles.go, vocabulary.go, srs.go
//   - Agent (core): chat.go
//   - App shell: config.go, settings.go, auth.go, windows.go, dialogs.go, media.go
type Service struct {
	db           *database.DB
	config       *config.Store
	autostart    *autostart.Manager
	wailsApp     *application.App
	turnRegistry *agent.TurnRegistry
	orchestrator *agent.Orchestrator
	auth         *auth.Service
}

func New(db *database.DB, cfg *config.Store, autostartManager *autostart.Manager) *Service {
	registry := agent.NewTurnRegistry()
	s := &Service{
		db:           db,
		config:       cfg,
		autostart:    autostartManager,
		turnRegistry: registry,
		auth:         auth.New(db.Queries, cfg),
	}
	s.orchestrator = agent.NewOrchestrator(agentEventEmitter{}, registry, chatMessageStore{s})
	return s
}

type agentEventEmitter struct{}

func (agentEventEmitter) EmitAgentEvent(name string, data any) {
	if globalWindowManager != nil {
		globalWindowManager.EmitAgentEvent(name, data)
	}
}

type chatMessageStore struct {
	s *Service
}

func (c chatMessageStore) UpdateMessageContent(ctx context.Context, messageID, content string) error {
	if err := c.s.db.Queries.UpdateChatMessageContent(ctx, store.UpdateChatMessageContentParams{
		Content: content,
		ID:      messageID,
	}); err != nil {
		return fmt.Errorf("update chat message content: %w", err)
	}
	return nil
}

// SetWailsApp attaches the Wails application for native dialogs.
func (s *Service) SetWailsApp(app *application.App) {
	s.wailsApp = app
}

// DatabasePath returns the on-disk SQLite database file path.
func (s *Service) DatabasePath() string {
	return s.db.Path
}
