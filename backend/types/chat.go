package types

// ChatSession is a sidebar row. Full message bodies are loaded separately.
type ChatSession struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Kind      string `json:"kind"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ChatMessage is one turn in a chat session.
type ChatMessage struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// StartChatTurnResult contains the persisted user message and streaming assistant placeholder.
type StartChatTurnResult struct {
	UserMessage      ChatMessage `json:"user_message"`
	AssistantMessage ChatMessage `json:"assistant_message"`
}

// OverlayChatTurnRequest starts an in-memory overlay turn (no sidebar session).
type OverlayChatTurnRequest struct {
	ConversationID string        `json:"conversation_id"`
	Content        string        `json:"content"`
	DomainContext  string        `json:"domain_context"`
	History        []ChatMessage `json:"history"`
}
