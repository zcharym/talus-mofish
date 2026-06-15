package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/songwei.ma/talus-mofish/internal/aiclient"
)

const persistDebounce = 2 * time.Second

// MessageStore persists assistant message content during a turn.
type MessageStore interface {
	UpdateMessageContent(ctx context.Context, messageID, content string) error
}

// Orchestrator runs chat turns and emits streaming events.
type Orchestrator struct {
	emitter  EventEmitter
	registry *TurnRegistry
	store    MessageStore
}

// NewOrchestrator creates a chat turn orchestrator.
func NewOrchestrator(emitter EventEmitter, registry *TurnRegistry, store MessageStore) *Orchestrator {
	return &Orchestrator{
		emitter:  emitter,
		registry: registry,
		store:    store,
	}
}

// RunTurnParams identifies one assistant generation turn.
type RunTurnParams struct {
	SessionID string
	MessageID string
	History   []aiclient.Message
	AI        aiclient.Config
}

// RunTurn streams a model response and emits agent events until done or cancelled.
func (o *Orchestrator) RunTurn(parent context.Context, params RunTurnParams) {
	ctx, cancel := context.WithCancel(parent)
	o.registry.Register(params.MessageID, cancel)
	defer o.registry.Unregister(params.MessageID)

	client, err := aiclient.NewClient(params.AI)
	if err != nil {
		o.emitError(params.SessionID, params.MessageID, err.Error())
		o.persistFinal(ctx, params.MessageID, "")
		return
	}

	stream, err := client.ChatStream(ctx, aiclient.ChatRequest{Messages: params.History})
	if err != nil {
		o.emitError(params.SessionID, params.MessageID, err.Error())
		o.persistFinal(ctx, params.MessageID, "")
		return
	}

	var content strings.Builder
	var lastPersist time.Time
	var persistMu sync.Mutex

	persist := func(force bool) {
		persistMu.Lock()
		defer persistMu.Unlock()
		now := time.Now()
		if !force && now.Sub(lastPersist) < persistDebounce {
			return
		}
		lastPersist = now
		text := content.String()
		if err := o.store.UpdateMessageContent(ctx, params.MessageID, text); err != nil && ctx.Err() == nil {
			o.emitError(params.SessionID, params.MessageID, fmt.Sprintf("persist message: %v", err))
		}
	}

	for part := range stream {
		if ctx.Err() != nil {
			final := content.String()
			o.emitCancelled(params.SessionID, params.MessageID, final)
			o.persistFinal(context.Background(), params.MessageID, final)
			return
		}

		switch part.Type {
		case aiclient.StreamPartTextDelta:
			content.WriteString(part.Text)
			o.emitter.EmitAgentEvent(EventStreamChunk, StreamChunkEvent{
				SessionID: params.SessionID,
				MessageID: params.MessageID,
				Chunk: StreamChunk{
					Type: string(part.Type),
					Text: part.Text,
				},
			})
			persist(false)
		case aiclient.StreamPartFinish:
			final := content.String()
			o.persistFinal(ctx, params.MessageID, final)
			o.emitter.EmitAgentEvent(EventTurnDone, TurnDoneEvent{
				SessionID: params.SessionID,
				MessageID: params.MessageID,
				Content:   final,
			})
			return
		case aiclient.StreamPartError:
			msg := part.Error
			if msg == "" {
				msg = "stream error"
			}
			final := content.String()
			o.emitError(params.SessionID, params.MessageID, msg)
			o.persistFinal(context.Background(), params.MessageID, final)
			return
		}
	}

	final := content.String()
	o.persistFinal(ctx, params.MessageID, final)
	o.emitter.EmitAgentEvent(EventTurnDone, TurnDoneEvent{
		SessionID: params.SessionID,
		MessageID: params.MessageID,
		Content:   final,
	})
}

func (o *Orchestrator) emitError(sessionID, messageID, msg string) {
	o.emitter.EmitAgentEvent(EventTurnError, TurnErrorEvent{
		SessionID: sessionID,
		MessageID: messageID,
		Error:     msg,
	})
}

func (o *Orchestrator) emitCancelled(sessionID, messageID, content string) {
	o.emitter.EmitAgentEvent(EventTurnCancelled, TurnCancelledEvent{
		SessionID: sessionID,
		MessageID: messageID,
		Content:   content,
	})
}

func (o *Orchestrator) persistFinal(ctx context.Context, messageID, content string) {
	if err := o.store.UpdateMessageContent(ctx, messageID, content); err != nil && ctx.Err() == nil {
		// Best-effort on shutdown/cancel paths.
		_ = o.store.UpdateMessageContent(context.Background(), messageID, content)
	}
}
