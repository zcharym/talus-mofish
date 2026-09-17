package feeds

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	htmlChannelID = regexp.MustCompile(`"channelId"\s*:\s*"(UC[\w-]{20,})"`)
	htmlRSSLink   = regexp.MustCompile(`feeds/videos\.xml\?channel_id=(UC[\w-]{20,})`)
)

func (c *Client) enrichYouTube(ctx context.Context, spec Spec, apiKey string) (Spec, error) {
	if spec.Kind != KindYouTube || spec.LinkOnly {
		return spec, nil
	}
	if spec.FeedURL != "" && spec.RemoteID != "" {
		return spec, nil
	}
	channelID, err := c.resolveYouTubeChannelID(ctx, spec, apiKey)
	if err != nil {
		return spec, err
	}
	filled := youtubeChannelSpec(channelID)
	if spec.Title != "" && spec.Title != "YouTube channel" {
		filled.Title = spec.Title
	}
	if spec.URL != "" {
		filled.URL = spec.URL
	}
	return filled, nil
}

func (c *Client) resolveYouTubeChannelID(ctx context.Context, spec Spec, apiKey string) (string, error) {
	if channelIDPattern.MatchString(spec.RemoteID) {
		return spec.RemoteID, nil
	}
	if apiKey != "" {
		id, err := c.youtubeChannelIDFromAPI(ctx, spec, apiKey)
		if err == nil && id != "" {
			return id, nil
		}
		if spec.URL == "" {
			return "", err
		}
	}
	if spec.URL == "" {
		return "", fmt.Errorf("YouTube channel could not be resolved; add a Data API key in Config → Feeds or paste a /channel/UC… URL")
	}
	return c.youtubeChannelIDFromPage(ctx, spec.URL)
}

func (c *Client) youtubeService(ctx context.Context, apiKey string) (*youtube.Service, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("YouTube Data API key is not set")
	}
	// WithHTTPClient disables googleapi's built-in API-key injection, so attach
	// the key on the RoundTripper before handing the client to youtube.NewService.
	httpClient := withAPIKey(c.http, apiKey)
	svc, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("youtube api: %w", err)
	}
	return svc, nil
}

func (c *Client) youtubeChannelIDFromAPI(ctx context.Context, spec Spec, apiKey string) (string, error) {
	svc, err := c.youtubeService(ctx, apiKey)
	if err != nil {
		return "", err
	}
	handle := strings.TrimPrefix(strings.TrimPrefix(spec.Title, "@"), "@")
	if strings.Contains(spec.URL, "/@") {
		if _, after, ok := strings.Cut(spec.URL, "/@"); ok {
			handle = strings.Split(after, "/")[0]
		}
	}
	if handle != "" {
		call := svc.Channels.List([]string{"id", "snippet"}).ForHandle(handle)
		resp, err := call.Context(ctx).Do()
		if err == nil && resp != nil && len(resp.Items) > 0 {
			return resp.Items[0].Id, nil
		}
	}
	query := strings.TrimSpace(spec.Title)
	if query == "" {
		query = spec.URL
	}
	search, err := svc.Search.List([]string{"snippet"}).Q(query).Type("channel").MaxResults(1).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("youtube search: %w", err)
	}
	if search == nil || len(search.Items) == 0 || search.Items[0].Snippet == nil {
		return "", fmt.Errorf("no YouTube channel matched %q", query)
	}
	id := search.Items[0].Snippet.ChannelId
	if id == "" && search.Items[0].Id != nil {
		id = search.Items[0].Id.ChannelId
	}
	if id == "" {
		return "", fmt.Errorf("YouTube search returned an empty channel id")
	}
	return id, nil
}

func (c *Client) youtubeChannelIDFromPage(ctx context.Context, pageURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, pageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept-Language", "en")
	resp, err := c.do(ctx, req)
	if err != nil {
		return "", fmt.Errorf("fetch YouTube channel page: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("YouTube channel page: HTTP %d", resp.StatusCode)
	}
	text := string(body)
	if match := htmlRSSLink.FindStringSubmatch(text); len(match) > 1 {
		return match[1], nil
	}
	if match := htmlChannelID.FindStringSubmatch(text); len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("could not find a channel id on %s; add a YouTube Data API key in Config → Feeds", pageURL)
}

func (c *Client) pingYouTube(ctx context.Context, apiKey string) (PingResult, error) {
	svc, err := c.youtubeService(ctx, apiKey)
	if err != nil {
		return PingResult{}, err
	}
	resp, err := svc.Videos.List([]string{"id"}).Id("jNQXAC9IVRw").Context(ctx).Do()
	if err != nil {
		return PingResult{}, fmt.Errorf("youtube ping: %w", err)
	}
	if resp == nil || len(resp.Items) == 0 {
		return PingResult{}, fmt.Errorf("YouTube API key was accepted but returned no video metadata")
	}
	return PingResult{OK: true, Service: "youtube", Detail: "Data API key works"}, nil
}
