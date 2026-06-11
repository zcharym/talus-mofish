package appservice

import (
	"context"
	"fmt"

	"github.com/songwei.ma/talus-mofish/internal/store"
)

// ListDecks returns all SRS decks.
func (s *Service) ListDecks() ([]store.Deck, error) {
	ctx := context.Background()
	items, err := s.db.Queries.ListDecks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list decks: %w", err)
	}
	if items == nil {
		return []store.Deck{}, nil
	}
	return items, nil
}

// ListCardsByDeck returns cards in a deck.
func (s *Service) ListCardsByDeck(deckID string) ([]store.Card, error) {
	ctx := context.Background()
	items, err := s.db.Queries.ListCardsByDeck(ctx, deckID)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	if items == nil {
		return []store.Card{}, nil
	}
	return items, nil
}

// GetCard returns a single SRS card by ID.
func (s *Service) GetCard(id string) (store.Card, error) {
	ctx := context.Background()
	item, err := s.db.Queries.GetCard(ctx, id)
	if err != nil {
		return store.Card{}, fmt.Errorf("get card: %w", err)
	}
	return item, nil
}

// UpdateCardContent saves editable card fields.
func (s *Service) UpdateCardContent(input store.UpdateCardContentParams) error {
	ctx := context.Background()
	if err := s.db.Queries.UpdateCardContent(ctx, input); err != nil {
		return fmt.Errorf("update card content: %w", err)
	}
	return nil
}

// DeleteCard removes an SRS card.
func (s *Service) DeleteCard(id string) error {
	ctx := context.Background()
	if err := s.db.Queries.DeleteCard(ctx, id); err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	return nil
}
