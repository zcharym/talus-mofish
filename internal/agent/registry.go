package agent

import "sync"

// TurnRegistry tracks in-flight chat turns for cancellation.
type TurnRegistry struct {
	mu    sync.Mutex
	turns map[string]contextCancel
}

type contextCancel struct {
	cancel func()
}

// NewTurnRegistry creates an empty turn registry.
func NewTurnRegistry() *TurnRegistry {
	return &TurnRegistry{turns: make(map[string]contextCancel)}
}

// Register stores a cancel function for a message ID.
func (r *TurnRegistry) Register(messageID string, cancel func()) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.turns[messageID] = contextCancel{cancel: cancel}
}

// Cancel invokes the cancel function for a message ID, if registered.
func (r *TurnRegistry) Cancel(messageID string) bool {
	r.mu.Lock()
	entry, ok := r.turns[messageID]
	if ok {
		delete(r.turns, messageID)
	}
	r.mu.Unlock()
	if !ok {
		return false
	}
	entry.cancel()
	return true
}

// Unregister removes a message ID without invoking cancel.
func (r *TurnRegistry) Unregister(messageID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.turns, messageID)
}
