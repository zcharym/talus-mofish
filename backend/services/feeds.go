package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/songwei.ma/talus-mofish/backend/feeds"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/storage/store"
)

const feedsTimeout = 25 * time.Second

// FeedsService exposes the Agent Feeds inbox (RSS, YouTube, Bilibili, read later).
type FeedsService struct {
	db     *storage.DB
	config *storage.ConfigStore
	client *feeds.Client
}

// NewFeedsService creates the Feeds Wails service.
func NewFeedsService(db *storage.DB, cfg *storage.ConfigStore) *FeedsService {
	return &FeedsService{db: db, config: cfg}
}

func (s *FeedsService) apiClient() *feeds.Client {
	if s.client != nil {
		return s.client
	}
	return feeds.NewClient(feeds.NewHTTPClient(s.config.Get().Feeds.ProxyURL))
}

func (s *FeedsService) creds() feeds.Credentials {
	cfg := s.config.Get().Feeds
	return feeds.Credentials{
		YouTubeAPIKey:    cfg.YouTubeAPIKey,
		BilibiliSESSDATA: cfg.BilibiliSESSDATA,
	}
}

// ListSources returns subscriptions (the built-in Read later bucket is omitted).
func (s *FeedsService) ListSources() ([]feeds.Source, error) {
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	rows, err := s.db.Queries.ListFeedSources(ctx)
	if err != nil {
		return nil, fmt.Errorf("list feed sources: %w", err)
	}
	out := make([]feeds.Source, 0, len(rows))
	for _, row := range rows {
		if row.Kind == feeds.KindReadLater {
			continue
		}
		out = append(out, mapSource(row))
	}
	return out, nil
}

// AddSource subscribes to an RSS / YouTube / Bilibili URL, or saves a single video/link.
func (s *FeedsService) AddSource(input string) (feeds.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()

	spec, err := feeds.Resolve(input)
	if err != nil {
		return feeds.Item{}, err
	}
	spec, err = s.apiClient().Enrich(ctx, spec, s.creds())
	if err != nil {
		return feeds.Item{}, err
	}
	if spec.LinkOnly {
		return s.saveLink(ctx, spec.Title, spec.URL, spec.Kind)
	}

	existing, err := s.db.Queries.ListFeedSources(ctx)
	if err != nil {
		return feeds.Item{}, fmt.Errorf("list feed sources: %w", err)
	}
	for _, row := range existing {
		if row.Kind == spec.Kind && ((spec.RemoteID != "" && row.RemoteID == spec.RemoteID) || (spec.URL != "" && row.Url == spec.URL)) {
			if err := s.refreshSource(ctx, mapSource(row)); err != nil {
				return feeds.Item{}, err
			}
			return feeds.Item{SourceID: row.ID, SourceTitle: row.Title, Title: row.Title}, nil
		}
	}

	source := feeds.Source{
		ID:       uuid.NewString(),
		Kind:     spec.Kind,
		Title:    spec.Title,
		URL:      spec.URL,
		RemoteID: spec.RemoteID,
		Enabled:  true,
	}
	if spec.FeedURL != "" && spec.Kind == feeds.KindRSS {
		source.URL = spec.FeedURL
	}
	if err := s.db.Queries.CreateFeedSource(ctx, store.CreateFeedSourceParams{
		ID:       source.ID,
		Kind:     source.Kind,
		Title:    source.Title,
		Url:      source.URL,
		RemoteID: source.RemoteID,
	}); err != nil {
		return feeds.Item{}, fmt.Errorf("create feed source: %w", err)
	}
	if err := s.refreshSource(ctx, source); err != nil {
		return feeds.Item{}, err
	}
	return feeds.Item{SourceID: source.ID, SourceTitle: source.Title, Title: source.Title, SourceKind: source.Kind}, nil
}

// DeleteSource removes a subscription and its cached items.
func (s *FeedsService) DeleteSource(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("source id is required")
	}
	if id == feeds.ReadLaterSourceID {
		return fmt.Errorf("the Read later list cannot be deleted")
	}
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	if err := s.db.Queries.DeleteFeedSource(ctx, id); err != nil {
		return fmt.Errorf("delete feed source: %w", err)
	}
	return nil
}

// Refresh pulls every enabled source and upserts items.
func (s *FeedsService) Refresh() (feeds.RefreshResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	rows, err := s.db.Queries.ListFeedSources(ctx)
	if err != nil {
		return feeds.RefreshResult{}, fmt.Errorf("list feed sources: %w", err)
	}
	result := feeds.RefreshResult{
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		Sources:   []feeds.SourceRefresh{},
	}
	for _, row := range rows {
		if row.Kind == feeds.KindReadLater || row.Enabled == 0 {
			continue
		}
		source := mapSource(row)
		entry := feeds.SourceRefresh{ID: source.ID, Title: source.Title, Kind: source.Kind}
		if err := s.refreshSource(ctx, source); err != nil {
			entry.Error = err.Error()
		} else {
			entry.ItemCount = 1
		}
		result.Sources = append(result.Sources, entry)
	}
	items, err := s.db.Queries.ListFeedItems(ctx)
	if err != nil {
		return result, fmt.Errorf("list feed items: %w", err)
	}
	result.ItemCount = len(items)
	return result, nil
}

// ListItems returns the inbox. filter is all, saved, youtube, bilibili, or rss.
func (s *FeedsService) ListItems(filter string) ([]feeds.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	filter = strings.ToLower(strings.TrimSpace(filter))
	if filter == "" {
		filter = "all"
	}

	var (
		rows []store.ListFeedItemsRow
		err  error
	)
	switch filter {
	case "all":
		rows, err = s.db.Queries.ListFeedItems(ctx)
	case "saved":
		saved, savedErr := s.db.Queries.ListSavedFeedItems(ctx)
		if savedErr != nil {
			err = savedErr
			break
		}
		rows = savedToList(saved)
	case feeds.KindYouTube, feeds.KindBilibili, feeds.KindRSS:
		byKind, kindErr := s.db.Queries.ListFeedItemsByKind(ctx, filter)
		if kindErr != nil {
			err = kindErr
			break
		}
		rows = kindToList(byKind)
	default:
		return nil, fmt.Errorf("unknown filter %q", filter)
	}
	if err != nil {
		return nil, fmt.Errorf("list feed items: %w", err)
	}
	out := make([]feeds.Item, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapItem(row))
	}
	return out, nil
}

// GetInbox returns sources plus items for the Agent page.
func (s *FeedsService) GetInbox(filter string) (feeds.Inbox, error) {
	sources, err := s.ListSources()
	if err != nil {
		return feeds.Inbox{}, err
	}
	items, err := s.ListItems(filter)
	if err != nil {
		return feeds.Inbox{}, err
	}
	if filter == "" {
		filter = "all"
	}
	return feeds.Inbox{
		Sources:   sources,
		Items:     items,
		Filter:    filter,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// SetItemSaved marks an item as read later / unsaved.
func (s *FeedsService) SetItemSaved(id string, saved bool) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("item id is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	flag := int64(0)
	if saved {
		flag = 1
	}
	if err := s.db.Queries.SetFeedItemSaved(ctx, store.SetFeedItemSavedParams{Saved: flag, ID: id}); err != nil {
		return fmt.Errorf("save feed item: %w", err)
	}
	return nil
}

// SetItemRead marks an item as read / unread.
func (s *FeedsService) SetItemRead(id string, read bool) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("item id is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	flag := int64(0)
	if read {
		flag = 1
	}
	if err := s.db.Queries.SetFeedItemRead(ctx, store.SetFeedItemReadParams{Read: flag, ID: id}); err != nil {
		return fmt.Errorf("read feed item: %w", err)
	}
	return nil
}

// SaveLink adds a URL to the local Read later list.
func (s *FeedsService) SaveLink(title, rawURL string) (feeds.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	return s.saveLink(ctx, title, rawURL, feeds.KindReadLater)
}

// DeleteItem removes a cached or saved item.
func (s *FeedsService) DeleteItem(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("item id is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	if err := s.db.Queries.DeleteFeedItem(ctx, id); err != nil {
		return fmt.Errorf("delete feed item: %w", err)
	}
	return nil
}

// PingYouTube tests the saved Data API key.
func (s *FeedsService) PingYouTube() (feeds.PingResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	return s.apiClient().PingYouTube(ctx, s.creds().YouTubeAPIKey)
}

// PingBilibili tests the saved SESSDATA cookie.
func (s *FeedsService) PingBilibili() (feeds.PingResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), feedsTimeout)
	defer cancel()
	return s.apiClient().PingBilibili(ctx, s.creds().BilibiliSESSDATA)
}

func (s *FeedsService) refreshSource(ctx context.Context, source feeds.Source) error {
	spec := feeds.Spec{
		Kind:     source.Kind,
		Title:    source.Title,
		URL:      source.URL,
		RemoteID: source.RemoteID,
		FeedURL:  source.URL,
	}
	if source.Kind == feeds.KindYouTube {
		if channelIDPatternOK(source.RemoteID) {
			spec.FeedURL = "https://www.youtube.com/feeds/videos.xml?channel_id=" + source.RemoteID
		} else if source.RemoteID != "" {
			spec.FeedURL = "https://www.youtube.com/feeds/videos.xml?playlist_id=" + source.RemoteID
		}
	}
	title, items, err := s.apiClient().Fetch(ctx, spec, s.creds())
	errText := ""
	if err != nil {
		errText = err.Error()
	}
	if title == "" {
		title = source.Title
	}
	if updErr := s.db.Queries.UpdateFeedSourceFetch(ctx, store.UpdateFeedSourceFetchParams{
		Title:     title,
		LastError: errText,
		ID:        source.ID,
	}); updErr != nil {
		return fmt.Errorf("update feed source: %w", updErr)
	}
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := s.upsertItem(ctx, source.ID, item, false); err != nil {
			return err
		}
	}
	return nil
}

func (s *FeedsService) upsertItem(ctx context.Context, sourceID string, item feeds.FetchedItem, saved bool) error {
	externalID := strings.TrimSpace(item.ExternalID)
	if externalID == "" {
		externalID = strings.TrimSpace(item.URL)
	}
	if externalID == "" {
		externalID = uuid.NewString()
	}
	published := ""
	if !item.PublishedAt.IsZero() {
		published = item.PublishedAt.UTC().Format(time.RFC3339)
	}
	existing, err := s.db.Queries.GetFeedItemBySourceExternal(ctx, store.GetFeedItemBySourceExternalParams{
		SourceID:   sourceID,
		ExternalID: externalID,
	})
	if err == nil {
		return s.db.Queries.UpdateFeedItemMetadata(ctx, store.UpdateFeedItemMetadataParams{
			Title:        item.Title,
			Url:          item.URL,
			Author:       item.Author,
			Summary:      item.Summary,
			ThumbnailUrl: item.ThumbnailURL,
			PublishedAt:  published,
			ID:           existing.ID,
		})
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("lookup feed item: %w", err)
	}
	savedInt := int64(0)
	if saved {
		savedInt = 1
	}
	return s.db.Queries.CreateFeedItem(ctx, store.CreateFeedItemParams{
		ID:           uuid.NewString(),
		SourceID:     sourceID,
		ExternalID:   externalID,
		Title:        item.Title,
		Url:          item.URL,
		Author:       item.Author,
		Summary:      item.Summary,
		ThumbnailUrl: item.ThumbnailURL,
		PublishedAt:  published,
		Saved:        savedInt,
	})
}

func (s *FeedsService) saveLink(ctx context.Context, title, rawURL, kind string) (feeds.Item, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return feeds.Item{}, fmt.Errorf("URL is required")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		title = rawURL
	}
	if err := s.ensureReadLater(ctx); err != nil {
		return feeds.Item{}, err
	}
	fetched := feeds.FetchedItem{
		ExternalID: rawURL,
		Title:      title,
		URL:        rawURL,
		Author:     kind,
	}
	if err := s.upsertItem(ctx, feeds.ReadLaterSourceID, fetched, true); err != nil {
		return feeds.Item{}, err
	}
	row, err := s.db.Queries.GetFeedItemBySourceExternal(ctx, store.GetFeedItemBySourceExternalParams{
		SourceID:   feeds.ReadLaterSourceID,
		ExternalID: rawURL,
	})
	if err != nil {
		return feeds.Item{}, err
	}
	return feeds.Item{
		ID:          row.ID,
		SourceID:    row.SourceID,
		SourceKind:  feeds.KindReadLater,
		SourceTitle: "Read later",
		ExternalID:  row.ExternalID,
		Title:       row.Title,
		URL:         row.Url,
		Author:      row.Author,
		Saved:       true,
		CreatedAt:   row.CreatedAt,
	}, nil
}

func (s *FeedsService) ensureReadLater(ctx context.Context) error {
	_, err := s.db.Queries.GetFeedSource(ctx, feeds.ReadLaterSourceID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return s.db.Queries.CreateFeedSource(ctx, store.CreateFeedSourceParams{
		ID:    feeds.ReadLaterSourceID,
		Kind:  feeds.KindReadLater,
		Title: "Read later",
	})
}

func mapSource(row store.FeedSource) feeds.Source {
	return feeds.Source{
		ID:            row.ID,
		Kind:          row.Kind,
		Title:         row.Title,
		URL:           row.Url,
		RemoteID:      row.RemoteID,
		Enabled:       row.Enabled != 0,
		LastError:     row.LastError,
		LastFetchedAt: row.LastFetchedAt,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func mapItem(row store.ListFeedItemsRow) feeds.Item {
	return feeds.Item{
		ID:           row.ID,
		SourceID:     row.SourceID,
		SourceKind:   row.SourceKind,
		SourceTitle:  row.SourceTitle,
		ExternalID:   row.ExternalID,
		Title:        row.Title,
		URL:          row.Url,
		Author:       row.Author,
		Summary:      row.Summary,
		ThumbnailURL: row.ThumbnailUrl,
		PublishedAt:  row.PublishedAt,
		Saved:        row.Saved != 0,
		Read:         row.Read != 0,
		CreatedAt:    row.CreatedAt,
	}
}

func savedToList(rows []store.ListSavedFeedItemsRow) []store.ListFeedItemsRow {
	out := make([]store.ListFeedItemsRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, store.ListFeedItemsRow{
			ID:           row.ID,
			SourceID:     row.SourceID,
			ExternalID:   row.ExternalID,
			Title:        row.Title,
			Url:          row.Url,
			Author:       row.Author,
			Summary:      row.Summary,
			ThumbnailUrl: row.ThumbnailUrl,
			PublishedAt:  row.PublishedAt,
			Saved:        row.Saved,
			Read:         row.Read,
			CreatedAt:    row.CreatedAt,
			SourceKind:   row.SourceKind,
			SourceTitle:  row.SourceTitle,
		})
	}
	return out
}

func kindToList(rows []store.ListFeedItemsByKindRow) []store.ListFeedItemsRow {
	out := make([]store.ListFeedItemsRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, store.ListFeedItemsRow{
			ID:           row.ID,
			SourceID:     row.SourceID,
			ExternalID:   row.ExternalID,
			Title:        row.Title,
			Url:          row.Url,
			Author:       row.Author,
			Summary:      row.Summary,
			ThumbnailUrl: row.ThumbnailUrl,
			PublishedAt:  row.PublishedAt,
			Saved:        row.Saved,
			Read:         row.Read,
			CreatedAt:    row.CreatedAt,
			SourceKind:   row.SourceKind,
			SourceTitle:  row.SourceTitle,
		})
	}
	return out
}

func channelIDPatternOK(id string) bool {
	return len(id) >= 22 && strings.HasPrefix(id, "UC")
}
