package agent

import (
	"strings"

	"github.com/songwei.ma/talus-mofish/internal/aiclient"
	"github.com/songwei.ma/talus-mofish/internal/store"
)

const (
	defaultSystemPrompt = `You are Talus Agent, a helpful desktop assistant.
Answer clearly and concisely. Use examples when helpful.
You can help with many topics; when the user asks about English learning, vocabulary, reading, or IELTS practice, draw on their local library (decks, articles, vocabulary) when relevant.
If the user writes in Chinese, you may reply in Chinese for explanations but include English examples when discussing English.`
	maxContextMessages = 40
)

// BuildMessages assembles provider messages from history plus the new user turn.
func BuildMessages(history []store.ChatMessage, userContent string) []aiclient.Message {
	msgs := make([]aiclient.Message, 0, len(history)+2)
	msgs = append(msgs, aiclient.Message{
		Role:    aiclient.RoleSystem,
		Content: defaultSystemPrompt,
	})

	start := 0
	if len(history) > maxContextMessages {
		start = len(history) - maxContextMessages
	}
	for _, msg := range history[start:] {
		content := strings.TrimSpace(msg.Content)
		if content == "" {
			continue
		}
		role := aiclient.RoleUser
		if msg.Role == "assistant" {
			role = aiclient.RoleAssistant
		}
		msgs = append(msgs, aiclient.Message{Role: role, Content: content})
	}

	userContent = strings.TrimSpace(userContent)
	if userContent != "" {
		msgs = append(msgs, aiclient.Message{Role: aiclient.RoleUser, Content: userContent})
	}

	return msgs
}
