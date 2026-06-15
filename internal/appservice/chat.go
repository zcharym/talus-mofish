package appservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/songwei.ma/talus-mofish/internal/agent"
	"github.com/songwei.ma/talus-mofish/internal/store"
)

const (
	defaultChatSessionTitle = "New chat"
	maxSessionTitleLength   = 48
)

// StartChatTurnResult contains the persisted user message and streaming assistant placeholder.
type StartChatTurnResult struct {
	UserMessage      store.ChatMessage `json:"user_message"`
	AssistantMessage store.ChatMessage `json:"assistant_message"`
}

// ListChatSessions returns chat sessions ordered by most recently updated.
func (s *Service) ListChatSessions() ([]store.ChatSession, error) {
	ctx := context.Background()
	sessions, err := s.db.Queries.ListChatSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("list chat sessions: %w", err)
	}
	if sessions == nil {
		return []store.ChatSession{}, nil
	}
	return sessions, nil
}

// CreateChatSession inserts a new chat session.
func (s *Service) CreateChatSession(title string) (store.ChatSession, error) {
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
	}); err != nil {
		return store.ChatSession{}, fmt.Errorf("create chat session: %w", err)
	}

	created, err := s.db.Queries.GetChatSession(ctx, session.ID)
	if err != nil {
		return store.ChatSession{}, fmt.Errorf("get created chat session: %w", err)
	}
	return created, nil
}

// RenameChatSession updates a session title.
func (s *Service) RenameChatSession(id, title string) error {
	ctx := context.Background()
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("title is required")
	}
	if err := s.db.Queries.RenameChatSession(ctx, store.RenameChatSessionParams{
		Title: title,
		ID:    id,
	}); err != nil {
		return fmt.Errorf("rename chat session: %w", err)
	}
	return nil
}

// DeleteChatSession removes a session and its messages.
func (s *Service) DeleteChatSession(id string) error {
	ctx := context.Background()
	if err := s.db.Queries.DeleteChatSession(ctx, id); err != nil {
		return fmt.Errorf("delete chat session: %w", err)
	}
	return nil
}

// ListChatMessages returns messages for a session in chronological order.
func (s *Service) ListChatMessages(sessionID string) ([]store.ChatMessage, error) {
	ctx := context.Background()
	messages, err := s.db.Queries.ListChatMessages(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list chat messages: %w", err)
	}
	if messages == nil {
		return []store.ChatMessage{}, nil
	}
	return messages, nil
}

// StartChatTurn persists the user message, creates an assistant placeholder, and begins streaming.
func (s *Service) StartChatTurn(sessionID, content string) (StartChatTurnResult, error) {
	ctx := context.Background()
	content = strings.TrimSpace(content)
	if content == "" {
		return StartChatTurnResult{}, fmt.Errorf("message content is required")
	}

	session, err := s.db.Queries.GetChatSession(ctx, sessionID)
	if err != nil {
		return StartChatTurnResult{}, fmt.Errorf("get chat session: %w", err)
	}

	priorMessages, err := s.db.Queries.ListChatMessages(ctx, sessionID)
	if err != nil {
		return StartChatTurnResult{}, fmt.Errorf("list chat messages: %w", err)
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
		return StartChatTurnResult{}, fmt.Errorf("create user message: %w", err)
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
		return StartChatTurnResult{}, fmt.Errorf("create assistant message: %w", err)
	}

	if err := s.db.Queries.TouchChatSession(ctx, sessionID); err != nil {
		return StartChatTurnResult{}, fmt.Errorf("touch chat session: %w", err)
	}

	if session.Title == defaultChatSessionTitle {
		title := sessionTitleFromMessage(content)
		if err := s.db.Queries.RenameChatSession(ctx, store.RenameChatSessionParams{
			Title: title,
			ID:    sessionID,
		}); err != nil {
			return StartChatTurnResult{}, fmt.Errorf("auto-title chat session: %w", err)
		}
	}

	userMessage, err = s.db.Queries.GetChatMessage(ctx, userMessage.ID)
	if err != nil {
		return StartChatTurnResult{}, fmt.Errorf("reload user message: %w", err)
	}
	assistantMessage, err = s.db.Queries.GetChatMessage(ctx, assistantMessage.ID)
	if err != nil {
		return StartChatTurnResult{}, fmt.Errorf("reload assistant message: %w", err)
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

	return StartChatTurnResult{
		UserMessage:      userMessage,
		AssistantMessage: assistantMessage,
	}, nil
}

// CancelChatTurn aborts an in-flight assistant response.
func (s *Service) CancelChatTurn(sessionID, messageID string) error {
	sessionID = strings.TrimSpace(sessionID)
	messageID = strings.TrimSpace(messageID)
	if sessionID == "" || messageID == "" {
		return fmt.Errorf("session id and message id are required")
	}
	if !s.turnRegistry.Cancel(messageID) {
		return fmt.Errorf("no active turn for message %s", messageID)
	}
	return nil
}

func sessionTitleFromMessage(content string) string {
	content = strings.Join(strings.Fields(content), " ")
	if len(content) <= maxSessionTitleLength {
		return content
	}
	return content[:maxSessionTitleLength-1] + "…"
}
