package agent

const (
	EventStreamChunk   = "agent:stream-chunk"
	EventTurnDone      = "agent:turn-done"
	EventTurnError     = "agent:turn-error"
	EventTurnCancelled = "agent:turn-cancelled"
)

// StreamChunkEvent is emitted for each streaming delta.
type StreamChunkEvent struct {
	SessionID string      `json:"sessionId"`
	MessageID string      `json:"messageId"`
	Chunk     StreamChunk `json:"chunk"`
}

// StreamChunk mirrors aiclient stream parts for the frontend.
type StreamChunk struct {
	Type         string `json:"type"`
	Text         string `json:"text,omitempty"`
	FinishReason string `json:"finishReason,omitempty"`
	Error        string `json:"error,omitempty"`
}

// TurnDoneEvent is emitted when a chat turn completes successfully.
type TurnDoneEvent struct {
	SessionID string `json:"sessionId"`
	MessageID string `json:"messageId"`
	Content   string `json:"content"`
}

// TurnErrorEvent is emitted when a chat turn fails.
type TurnErrorEvent struct {
	SessionID string `json:"sessionId"`
	MessageID string `json:"messageId"`
	Error     string `json:"error"`
}

// TurnCancelledEvent is emitted when a chat turn is cancelled.
type TurnCancelledEvent struct {
	SessionID string `json:"sessionId"`
	MessageID string `json:"messageId"`
	Content   string `json:"content"`
}

// EventEmitter pushes agent events to the UI layer.
type EventEmitter interface {
	EmitAgentEvent(name string, data any)
}
