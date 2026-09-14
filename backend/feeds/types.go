// Package feeds fetches RSS, YouTube, and Bilibili watch-later items for the Agent inbox.
package feeds

import "time"

const (
	KindRSS       = "rss"
	KindYouTube   = "youtube"
	KindBilibili  = "bilibili"
	KindReadLater = "readlater"

	RemoteBilibiliToView = "toview"
	ReadLaterSourceID    = "feed-readlater"
)

// Credentials are optional API keys / cookies from config.json (keyring overlay).
type Credentials struct {
	YouTubeAPIKey    string
	BilibiliSESSDATA string
}

// Spec is a parsed subscribe or save-link target.
type Spec struct {
	Kind     string `json:"kind"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	RemoteID string `json:"remoteId"`
	FeedURL  string `json:"feedUrl"`
	LinkOnly bool   `json:"linkOnly"`
}

// Source is a persisted subscription shown in the Agent Feeds tab.
type Source struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	RemoteID      string `json:"remoteId"`
	Enabled       bool   `json:"enabled"`
	LastError     string `json:"lastError"`
	LastFetchedAt string `json:"lastFetchedAt"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// Item is one timeline / watch-later / saved entry.
type Item struct {
	ID           string `json:"id"`
	SourceID     string `json:"sourceId"`
	SourceKind   string `json:"sourceKind"`
	SourceTitle  string `json:"sourceTitle"`
	ExternalID   string `json:"externalId"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Author       string `json:"author"`
	Summary      string `json:"summary"`
	ThumbnailURL string `json:"thumbnailUrl"`
	PublishedAt  string `json:"publishedAt"`
	Saved        bool   `json:"saved"`
	Read         bool   `json:"read"`
	CreatedAt    string `json:"createdAt"`
}

// RefreshResult summarizes a pull of every enabled source.
type RefreshResult struct {
	FetchedAt string          `json:"fetchedAt"`
	ItemCount int             `json:"itemCount"`
	Sources   []SourceRefresh `json:"sources"`
}

// SourceRefresh is per-source fetch outcome.
type SourceRefresh struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Kind      string `json:"kind"`
	ItemCount int    `json:"itemCount"`
	Error     string `json:"error"`
}

// PingResult is the Config tab test-connection payload.
type PingResult struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
	Detail  string `json:"detail"`
}

// FetchedItem is a provider-normalized entry before SQLite upsert.
type FetchedItem struct {
	ExternalID   string
	Title        string
	URL          string
	Author       string
	Summary      string
	ThumbnailURL string
	PublishedAt  time.Time
}

// Inbox is the Agent page payload.
type Inbox struct {
	Sources   []Source `json:"sources"`
	Items     []Item   `json:"items"`
	Filter    string   `json:"filter"`
	FetchedAt string   `json:"fetchedAt"`
}
