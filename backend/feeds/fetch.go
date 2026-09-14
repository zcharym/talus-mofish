package feeds

import (
	"context"
	"fmt"
)

// Enrich fills YouTube channel RSS URLs when the paste was a handle or custom URL.
func (c *Client) Enrich(ctx context.Context, spec Spec, creds Credentials) (Spec, error) {
	if spec.Kind == KindYouTube && !spec.LinkOnly {
		return c.enrichYouTube(ctx, spec, creds.YouTubeAPIKey)
	}
	return spec, nil
}

// Fetch pulls items for a resolved source.
func (c *Client) Fetch(ctx context.Context, spec Spec, creds Credentials) (string, []FetchedItem, error) {
	switch spec.Kind {
	case KindRSS:
		return c.fetchRSS(ctx, spec.FeedURL, spec.Title)
	case KindYouTube:
		if spec.LinkOnly {
			return "", nil, fmt.Errorf("single YouTube videos are saved to Read later, not subscribed")
		}
		feedURL := spec.FeedURL
		if feedURL == "" && spec.RemoteID != "" {
			if channelIDPattern.MatchString(spec.RemoteID) {
				feedURL = youtubeChannelRSS(spec.RemoteID)
			} else {
				feedURL = youtubePlaylistRSS(spec.RemoteID)
			}
		}
		return c.fetchRSS(ctx, feedURL, spec.Title)
	case KindBilibili:
		return c.fetchBilibili(ctx, spec, creds.BilibiliSESSDATA)
	case KindReadLater:
		return spec.Title, nil, nil
	default:
		return "", nil, fmt.Errorf("unknown feed kind %q", spec.Kind)
	}
}

// PingYouTube verifies a Data API key.
func (c *Client) PingYouTube(ctx context.Context, apiKey string) (PingResult, error) {
	return c.pingYouTube(ctx, apiKey)
}

// PingBilibili verifies a SESSDATA cookie.
func (c *Client) PingBilibili(ctx context.Context, sessdata string) (PingResult, error) {
	return c.pingBilibili(ctx, sessdata)
}
