package feeds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func (c *Client) fetchRSS(ctx context.Context, feedURL, fallbackTitle string) (string, []FetchedItem, error) {
	if strings.TrimSpace(feedURL) == "" {
		return "", nil, fmt.Errorf("RSS URL is required")
	}
	feed, err := c.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return "", nil, fmt.Errorf("parse RSS: %w", err)
	}
	title := strings.TrimSpace(feed.Title)
	if title == "" {
		title = fallbackTitle
	}
	items := make([]FetchedItem, 0, len(feed.Items))
	for _, entry := range feed.Items {
		if entry == nil {
			continue
		}
		item := fetchedFromRSS(entry)
		if item.Title == "" && item.URL == "" {
			continue
		}
		items = append(items, item)
	}
	return title, items, nil
}

func fetchedFromRSS(entry *gofeed.Item) FetchedItem {
	published := time.Time{}
	if entry.PublishedParsed != nil {
		published = *entry.PublishedParsed
	} else if entry.UpdatedParsed != nil {
		published = *entry.UpdatedParsed
	}
	author := ""
	if entry.Author != nil {
		author = strings.TrimSpace(entry.Author.Name)
	}
	if author == "" && len(entry.Authors) > 0 && entry.Authors[0] != nil {
		author = strings.TrimSpace(entry.Authors[0].Name)
	}
	id := strings.TrimSpace(entry.GUID)
	if id == "" {
		id = strings.TrimSpace(entry.Link)
	}
	if id == "" {
		id = strings.TrimSpace(entry.Title)
	}
	return FetchedItem{
		ExternalID:   id,
		Title:        strings.TrimSpace(entry.Title),
		URL:          strings.TrimSpace(entry.Link),
		Author:       author,
		Summary:      sanitizeSummary(firstNonEmpty(entry.Description, entry.Content)),
		ThumbnailURL: rssThumbnail(entry),
		PublishedAt:  published,
	}
}

func rssThumbnail(entry *gofeed.Item) string {
	if entry.Image != nil {
		if url := strings.TrimSpace(entry.Image.URL); url != "" {
			return url
		}
	}
	if media, ok := entry.Extensions["media"]; ok {
		if thumbs := media["thumbnail"]; len(thumbs) > 0 {
			if url := strings.TrimSpace(thumbs[0].Attrs["url"]); url != "" {
				return url
			}
		}
		if groups := media["group"]; len(groups) > 0 {
			if thumbs := groups[0].Children["thumbnail"]; len(thumbs) > 0 {
				if url := strings.TrimSpace(thumbs[0].Attrs["url"]); url != "" {
					return url
				}
			}
		}
	}
	for _, enc := range entry.Enclosures {
		if enc == nil {
			continue
		}
		if strings.HasPrefix(enc.Type, "image/") && enc.URL != "" {
			return enc.URL
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
