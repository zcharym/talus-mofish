package services

import (
	"database/sql"

	"github.com/songwei.ma/talus-mofish/backend/storage/store"
	"github.com/songwei.ma/talus-mofish/backend/types"
)

func chatSessionDTO(row store.ChatSession) types.ChatSession {
	return types.ChatSession{
		ID:        row.ID,
		Title:     row.Title,
		Kind:      row.Kind,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func chatMessageDTO(row store.ChatMessage) types.ChatMessage {
	return types.ChatMessage{
		ID:        row.ID,
		SessionID: row.SessionID,
		Role:      row.Role,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
	}
}

func chatSessionsDTO(rows []store.ChatSession) []types.ChatSession {
	out := make([]types.ChatSession, 0, len(rows))
	for _, row := range rows {
		out = append(out, chatSessionDTO(row))
	}
	return out
}

func chatMessagesDTO(rows []store.ChatMessage) []types.ChatMessage {
	out := make([]types.ChatMessage, 0, len(rows))
	for _, row := range rows {
		out = append(out, chatMessageDTO(row))
	}
	return out
}

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
