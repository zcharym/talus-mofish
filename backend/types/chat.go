package types

import "github.com/songwei.ma/talus-mofish/backend/storage/store"

// StartChatTurnResult contains the persisted user message and streaming assistant placeholder.
type StartChatTurnResult struct {
	UserMessage      store.ChatMessage `json:"user_message"`
	AssistantMessage store.ChatMessage `json:"assistant_message"`
}
