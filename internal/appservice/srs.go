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
