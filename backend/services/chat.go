package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/songwei.ma/talus-mofish/backend/agent"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/storage/store"
	"github.com/songwei.ma/talus-mofish/backend/types"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	defaultChatSessionTitle = "New chat"
	maxSessionTitleLength   = 48
)

// ChatService exposes chat session and streaming turn APIs.
// It manages chat sessions, messages, and orchestrates streaming responses from LLM providers.
// The service integrates with the agent orchestrator for turn-based conversations.
type ChatService struct {
	db           *storage.DB
	config       *storage.ConfigStore
	wailsApp     *application.App
	windows      WindowManager
	turnRegistry *agent.TurnRegistry
	orchestrator *agent.Orchestrator
}

// NewChatService creates the chat Wails service.
func NewChatService(db *storage.DB, cfg *storage.ConfigStore) *ChatService {
	registry := agent.NewTurnRegistry()
	s := &ChatService{
		db:           db,
		config:       cfg,
		turnRegistry: registry,
	}
	s.orchestrator = agent.NewOrchestrator(chatEventEmitter{s}, registry, chatMessageStore{s})
	return s
}

type chatEventEmitter struct {
	s *ChatService
}

func (e chatEventEmitter) EmitAgentEvent(name string, data any) {
	if e.s.windows != nil {
		e.s.windows.EmitAgentEvent(name, data)
	}
}

type chatMessageStore struct {
	s *ChatService
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

// ListChatSessions returns all chat sessions ordered by most recently updated.
// Sessions are returned in descending order of their last update time.
// Returns an empty slice if no sessions exist.
func (s *ChatService) ListChatSessions() ([]store.ChatSession, error) {
	ctx := context.Background()
	sessions, err := s.db.Queries.ListChatSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat sessions: %w", err)
	}
	if sessions == nil {
		return []store.ChatSession{}, nil
	}
	return sessions, nil
}

// CreateChatSession creates a new chat session with the given title.
// If the title is empty or whitespace-only, a default title is used.
// The session is persisted to the database and returned with timestamps populated.
func (s *ChatService) CreateChatSession(title string) (store.ChatSession, error) {
	ctx := context.Background()
	title = strings.TrimSpace(title)
	if title == "" {
		title = defaultChatSessionTitle
	}

	session := store.ChatSession{
		ID:    uuid.NewString(),
		Title: title,
	}
	if err := s.db.Queries.CreateChatSession(ctx, store.CreateChatSessionParams{
		ID:    session.ID,
		Title: session.Title,
		Kind:  "chat",
	}); err != nil {
		return store.ChatSession{}, fmt.Errorf("failed to create chat session: %w", err)
	}

	created, err := s.db.Queries.GetChatSession(ctx, session.ID)
	if err != nil {
		return store.ChatSession{}, fmt.Errorf("failed to retrieve created session: %w", err)
	}
	return created, nil
}

// RenameChatSession updates the title of an existing chat session.
// The title must not be empty after trimming whitespace.
// Returns an error if the session does not exist or the title is invalid.
func (s *ChatService) RenameChatSession(id, title string) error {
	ctx := context.Background()
	
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("session id is required")
	}
	
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	
	if err := s.db.Queries.RenameChatSession(ctx, store.RenameChatSessionParams{
		Title: title,
		ID:    id,
	}); err != nil {
		return fmt.Errorf("failed to rename chat session %q: %w", id, err)
	}
	return nil
}

// DeleteChatSession permanently removes a chat session and all its messages.
// This operation cannot be undone. Returns an error if the session does not exist.
func (s *ChatService) DeleteChatSession(id string) error {
	ctx := context.Background()
	
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("session id is required")
	}
	
	if err := s.db.Queries.DeleteChatSession(ctx, id); err != nil {
		return fmt.Errorf("failed to delete chat session %q: %w", id, err)
	}
	return nil
}

// ListChatMessages returns all messages for a session in chronological order.
// Messages are ordered from oldest to newest to maintain conversation flow.
// Returns an empty slice if the session has no messages or does not exist.
func (s *ChatService) ListChatMessages(sessionID string) ([]store.ChatMessage, error) {
	ctx := context.Background()
	
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, fmt.Errorf("session id is required")
	}
	
	messages, err := s.db.Queries.ListChatMessages(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages for session %q: %w", sessionID, err)
	}
	if messages == nil {
		return []store.ChatMessage{}, nil
	}
	return messages, nil
}

// StartChatTurn initiates a new chat turn by persisting the user message,
// creating a placeholder for the assistant response, and starting the streaming
// generation process. The function returns immediately after starting the async
// turn; the orchestrator will emit streaming events as the response is generated.
//
// The session title will be auto-updated from the default if this is the first message.
// Returns both the user and assistant message records with their IDs for tracking.
func (s *ChatService) StartChatTurn(sessionID, content string) (types.StartChatTurnResult, error) {
	ctx := context.Background()
	
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return types.StartChatTurnResult{}, fmt.Errorf("session id is required")
	}
	
	content = strings.TrimSpace(content)
	if content == "" {
		return types.StartChatTurnResult{}, fmt.Errorf("message content cannot be empty")
	}

	session, err := s.db.Queries.GetChatSession(ctx, sessionID)
	if err != nil {
		return types.StartChatTurnResult{}, fmt.Errorf("session %q not found: %w", sessionID, err)
	}

	priorMessages, err := s.db.Queries.ListChatMessages(ctx, sessionID)
	if err != nil {
		return types.StartChatTurnResult{}, fmt.Errorf("failed to retrieve message history: %w", err)
	}
	if priorMessages == nil {
		priorMessages = []store.ChatMessage{}
	}

	userMessage := store.ChatMessage{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		Role:      "user",
		Content:   content,
	}
	if err := s.db.Queries.CreateChatMessage(ctx, store.CreateChatMessageParams{
		ID:        userMessage.ID,
		SessionID: userMessage.SessionID,
		Role:      userMessage.Role,
		Content:   userMessage.Content,
	}); err != nil {
		return types.StartChatTurnResult{}, fmt.Errorf("create user message: %w", err)
	}

	assistantMessage := store.ChatMessage{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		Role:      "assistant",
		Content:   "",
	}
	if err := s.db.Queries.CreateChatMessage(ctx, store.CreateChatMessageParams{
		ID:        assistantMessage.ID,
		SessionID: assistantMessage.SessionID,
		Role:      assistantMessage.Role,
		Content:   assistantMessage.Content,
	}); err != nil {
		return types.StartChatTurnResult{}, fmt.Errorf("create assistant message: %w", err)
	}

	if err := s.db.Queries.TouchChatSession(ctx, sessionID); err != nil {
		return types.StartChatTurnResult{}, fmt.Errorf("touch chat session: %w", err)
	}

	if session.Title == defaultChatSessionTitle {
		title := sessionTitleFromMessage(content)
		if err := s.db.Queries.RenameChatSession(ctx, store.RenameChatSessionParams{
			Title: title,
			ID:    sessionID,
		}); err != nil {
			return types.StartChatTurnResult{}, fmt.Errorf("auto-title chat session: %w", err)
		}
	}

	userMessage, err = s.db.Queries.GetChatMessage(ctx, userMessage.ID)
	if err != nil {
		return types.StartChatTurnResult{}, fmt.Errorf("reload user message: %w", err)
	}
	assistantMessage, err = s.db.Queries.GetChatMessage(ctx, assistantMessage.ID)
	if err != nil {
		return types.StartChatTurnResult{}, fmt.Errorf("reload assistant message: %w", err)
	}

	aiCfg := s.config.Get().AI
	history := agent.BuildMessages(priorMessages, content)
	params := agent.RunTurnParams{
		SessionID: sessionID,
		MessageID: assistantMessage.ID,
		History:   history,
		AI:        aiCfg,
	}

	parent := context.Background()
	if s.wailsApp != nil {
		parent = s.wailsApp.Context()
	}
	go s.orchestrator.RunTurn(parent, params)

	return types.StartChatTurnResult{
		UserMessage:      userMessage,
		AssistantMessage: assistantMessage,
	}, nil
}

// CancelChatTurn attempts to abort an in-flight assistant response.
// This will stop the streaming generation and persist the partial response.
// Returns an error if no active turn exists for the given message ID.
func (s *ChatService) CancelChatTurn(sessionID, messageID string) error {
	sessionID = strings.TrimSpace(sessionID)
	messageID = strings.TrimSpace(messageID)
	
	if sessionID == "" {
		return fmt.Errorf("session id is required")
	}
	if messageID == "" {
		return fmt.Errorf("message id is required")
	}
	
	if !s.turnRegistry.Cancel(messageID) {
		return fmt.Errorf("no active turn found for message %q in session %q", messageID, sessionID)
	}
	return nil
}

// sessionTitleFromMessage generates a session title from the first user message.
// It normalizes whitespace and truncates to maxSessionTitleLength if needed.
func sessionTitleFromMessage(content string) string {
	content = strings.Join(strings.Fields(content), " ")
	if len(content) <= maxSessionTitleLength {
		return content
	}
	return content[:maxSessionTitleLength-1] + "…"
}
