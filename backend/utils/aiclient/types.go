package aiclient

// Role is a chat message role.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message is one turn in a chat conversation.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// StreamPartType identifies a chunk from ChatStream.
type StreamPartType string

const (
	StreamPartTextDelta StreamPartType = "text-delta"
	StreamPartFinish    StreamPartType = "finish"
	StreamPartError     StreamPartType = "error"
)

// StreamPart is one event from a streaming chat completion.
type StreamPart struct {
	Type         StreamPartType `json:"type"`
	Text         string         `json:"text,omitempty"`
	FinishReason string         `json:"finishReason,omitempty"`
	Error        string         `json:"error,omitempty"`
}

// ChatRequest is input to ChatStream.
type ChatRequest struct {
	Messages    []Message
	Temperature float64
}
