package feeds

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	youtubeHost = map[string]struct{}{
		"youtube.com":       {},
		"www.youtube.com":   {},
		"m.youtube.com":     {},
		"music.youtube.com": {},
		"youtu.be":          {},
		"www.youtu.be":      {},
	}
	bilibiliHost = map[string]struct{}{
		"bilibili.com":       {},
		"www.bilibili.com":   {},
		"m.bilibili.com":     {},
		"space.bilibili.com": {},
		"b23.tv":             {},
	}
	channelIDPattern = regexp.MustCompile(`^UC[\w-]{20,}$`)
	bvPattern        = regexp.MustCompile(`(?i)BV[0-9A-Za-z]{10}`)
)

// Resolve classifies a pasted URL or shortcut into a subscribe spec or a single link.
func Resolve(raw string) (Spec, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Spec{}, fmt.Errorf("paste an RSS, YouTube, or Bilibili URL")
	}

	if spec, ok := resolveShortcut(raw); ok {
		return spec, nil
	}

	if strings.HasPrefix(raw, "@") && !strings.Contains(raw, "/") {
		handle := strings.TrimPrefix(raw, "@")
		if handle == "" {
			return Spec{}, fmt.Errorf("YouTube handle is empty")
		}
		return Spec{
			Kind:  KindYouTube,
			Title: "@" + handle,
			URL:   "https://www.youtube.com/@" + handle,
		}, nil
	}

	parsed, err := parseHTTPURL(raw)
	if err != nil {
		return Spec{}, err
	}

	host := strings.ToLower(parsed.Hostname())
	if _, ok := youtubeHost[host]; ok {
		return resolveYouTubeURL(parsed)
	}
	if _, ok := bilibiliHost[host]; ok {
		return resolveBilibiliURL(parsed)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return Spec{}, fmt.Errorf("unsupported URL scheme %q", parsed.Scheme)
	}

	return Spec{
		Kind:    KindRSS,
		Title:   host,
		URL:     parsed.String(),
		FeedURL: parsed.String(),
	}, nil
}

func resolveShortcut(raw string) (Spec, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "watch later", "稍后再看", "bilibili watch later", "bilibili:toview", "bilibili:watchlater":
		return Spec{
			Kind:     KindBilibili,
			Title:    "Bilibili watch later",
			URL:      "https://www.bilibili.com/watchlater/",
			RemoteID: RemoteBilibiliToView,
		}, true
	default:
		return Spec{}, false
	}
}

func parseHTTPURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}
	if parsed.Scheme == "" {
		parsed, err = url.Parse("https://" + raw)
		if err != nil {
			return nil, fmt.Errorf("parse URL: %w", err)
		}
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("URL is missing a host")
	}
	return parsed, nil
}

func resolveYouTubeURL(parsed *url.URL) (Spec, error) {
	host := strings.ToLower(parsed.Hostname())
	query := parsed.Query()
	path := strings.Trim(parsed.Path, "/")
	parts := splitPath(path)

	if host == "youtu.be" || host == "www.youtu.be" {
		if len(parts) == 0 {
			return Spec{}, fmt.Errorf("YouTube short URL is missing a video id")
		}
		videoID := parts[0]
		return Spec{
			Kind:     KindYouTube,
			Title:    videoID,
			URL:      "https://www.youtube.com/watch?v=" + videoID,
			RemoteID: videoID,
			LinkOnly: true,
		}, nil
	}

	if playlistID := strings.TrimSpace(query.Get("list")); playlistID != "" && playlistID != "WL" {
		return Spec{
			Kind:     KindYouTube,
			Title:    "YouTube playlist",
			URL:      "https://www.youtube.com/playlist?list=" + playlistID,
			RemoteID: playlistID,
			FeedURL:  youtubePlaylistRSS(playlistID),
		}, nil
	}
	if playlistID := strings.TrimSpace(query.Get("playlist_id")); playlistID != "" {
		return Spec{
			Kind:     KindYouTube,
			Title:    "YouTube playlist",
			URL:      "https://www.youtube.com/playlist?list=" + playlistID,
			RemoteID: playlistID,
			FeedURL:  youtubePlaylistRSS(playlistID),
		}, nil
	}
	if channelID := strings.TrimSpace(query.Get("channel_id")); channelIDPattern.MatchString(channelID) {
		return youtubeChannelSpec(channelID), nil
	}

	if len(parts) >= 2 && parts[0] == "feeds" && parts[1] == "videos.xml" {
		if channelID := query.Get("channel_id"); channelIDPattern.MatchString(channelID) {
			return youtubeChannelSpec(channelID), nil
		}
		if playlistID := query.Get("playlist_id"); playlistID != "" {
			return Spec{
				Kind:     KindYouTube,
				Title:    "YouTube playlist",
				URL:      "https://www.youtube.com/playlist?list=" + playlistID,
				RemoteID: playlistID,
				FeedURL:  youtubePlaylistRSS(playlistID),
			}, nil
		}
	}

	if len(parts) >= 2 && parts[0] == "channel" && channelIDPattern.MatchString(parts[1]) {
		return youtubeChannelSpec(parts[1]), nil
	}
	if len(parts) >= 1 && strings.HasPrefix(parts[0], "@") {
		handle := strings.TrimPrefix(parts[0], "@")
		return Spec{
			Kind:  KindYouTube,
			Title: "@" + handle,
			URL:   "https://www.youtube.com/@" + handle,
		}, nil
	}
	if len(parts) >= 2 && (parts[0] == "c" || parts[0] == "user") {
		name := parts[1]
		return Spec{
			Kind:  KindYouTube,
			Title: name,
			URL:   "https://www.youtube.com/" + parts[0] + "/" + name,
		}, nil
	}

	if videoID := strings.TrimSpace(query.Get("v")); videoID != "" {
		return Spec{
			Kind:     KindYouTube,
			Title:    videoID,
			URL:      "https://www.youtube.com/watch?v=" + videoID,
			RemoteID: videoID,
			LinkOnly: true,
		}, nil
	}
	if len(parts) >= 2 && parts[0] == "shorts" {
		return Spec{
			Kind:     KindYouTube,
			Title:    parts[1],
			URL:      "https://www.youtube.com/shorts/" + parts[1],
			RemoteID: parts[1],
			LinkOnly: true,
		}, nil
	}
	if len(parts) >= 2 && parts[0] == "playlist" {
		return Spec{}, fmt.Errorf("YouTube playlist URL is missing list=")
	}

	return Spec{}, fmt.Errorf("unrecognized YouTube URL")
}

func youtubeChannelSpec(channelID string) Spec {
	return Spec{
		Kind:     KindYouTube,
		Title:    "YouTube channel",
		URL:      "https://www.youtube.com/channel/" + channelID,
		RemoteID: channelID,
		FeedURL:  youtubeChannelRSS(channelID),
	}
}

func youtubeChannelRSS(channelID string) string {
	return "https://www.youtube.com/feeds/videos.xml?channel_id=" + url.QueryEscape(channelID)
}

func youtubePlaylistRSS(playlistID string) string {
	return "https://www.youtube.com/feeds/videos.xml?playlist_id=" + url.QueryEscape(playlistID)
}

func resolveBilibiliURL(parsed *url.URL) (Spec, error) {
	host := strings.ToLower(parsed.Hostname())
	path := strings.Trim(parsed.Path, "/")
	parts := splitPath(path)
	query := parsed.Query()

	if host == "space.bilibili.com" || (len(parts) >= 1 && parts[0] == "space") {
		uidParts := parts
		if parts[0] == "space" {
			uidParts = parts[1:]
		}
		if len(uidParts) == 0 {
			return Spec{}, fmt.Errorf("Bilibili space URL is missing a uid")
		}
		mid := uidParts[0]
		if _, err := strconv.ParseInt(mid, 10, 64); err != nil {
			return Spec{}, fmt.Errorf("Bilibili uid is invalid")
		}
		if len(uidParts) >= 2 && uidParts[1] == "favlist" {
			fid := strings.TrimSpace(query.Get("fid"))
			if fid == "" {
				return Spec{}, fmt.Errorf("Bilibili favorites URL is missing fid=")
			}
			return Spec{
				Kind:     KindBilibili,
				Title:    "Bilibili favorites",
				URL:      "https://space.bilibili.com/" + mid + "/favlist?fid=" + fid,
				RemoteID: "fav:" + fid,
			}, nil
		}
		return Spec{
			Kind:     KindBilibili,
			Title:    "Bilibili " + mid,
			URL:      "https://space.bilibili.com/" + mid,
			RemoteID: "mid:" + mid,
		}, nil
	}

	if strings.Contains(path, "watchlater") {
		return Spec{
			Kind:     KindBilibili,
			Title:    "Bilibili watch later",
			URL:      "https://www.bilibili.com/watchlater/",
			RemoteID: RemoteBilibiliToView,
		}, nil
	}

	if bv := bvPattern.FindString(path); bv != "" {
		return Spec{
			Kind:     KindBilibili,
			Title:    bv,
			URL:      "https://www.bilibili.com/video/" + bv,
			RemoteID: bv,
			LinkOnly: true,
		}, nil
	}

	return Spec{}, fmt.Errorf("unrecognized Bilibili URL")
}

func splitPath(path string) []string {
	if path == "" {
		return nil
	}
	var parts []string
	for _, part := range strings.Split(path, "/") {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}
