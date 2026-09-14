package feeds

import (
	"regexp"
	"strings"
	"testing"
)

func TestResolveRSS(t *testing.T) {
	spec, err := Resolve("https://example.com/feed.xml")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Kind != KindRSS || spec.FeedURL != "https://example.com/feed.xml" || spec.LinkOnly {
		t.Fatalf("%+v", spec)
	}
}

func TestResolveYouTubeChannelPlaylistAndVideo(t *testing.T) {
	channel, err := Resolve("https://www.youtube.com/channel/UCXuqSBlHAE6Xw-yeJA0Tunw")
	if err != nil {
		t.Fatal(err)
	}
	if channel.Kind != KindYouTube || channel.RemoteID != "UCXuqSBlHAE6Xw-yeJA0Tunw" || channel.FeedURL == "" {
		t.Fatalf("channel %+v", channel)
	}

	playlist, err := Resolve("https://www.youtube.com/playlist?list=PLtest123")
	if err != nil {
		t.Fatal(err)
	}
	if playlist.RemoteID != "PLtest123" || !strings.Contains(playlist.FeedURL, "playlist_id=PLtest123") {
		t.Fatalf("playlist %+v", playlist)
	}

	handle, err := Resolve("https://www.youtube.com/@veritasium")
	if err != nil {
		t.Fatal(err)
	}
	if handle.Title != "@veritasium" || handle.FeedURL != "" {
		t.Fatalf("handle %+v", handle)
	}

	video, err := Resolve("https://youtu.be/dQw4w9WgXcQ")
	if err != nil {
		t.Fatal(err)
	}
	if !video.LinkOnly || video.RemoteID != "dQw4w9WgXcQ" {
		t.Fatalf("video %+v", video)
	}
}

func TestResolveBilibili(t *testing.T) {
	watch, err := Resolve("稍后再看")
	if err != nil {
		t.Fatal(err)
	}
	if watch.Kind != KindBilibili || watch.RemoteID != RemoteBilibiliToView {
		t.Fatalf("watch %+v", watch)
	}

	space, err := Resolve("https://space.bilibili.com/2")
	if err != nil {
		t.Fatal(err)
	}
	if space.RemoteID != "mid:2" {
		t.Fatalf("space %+v", space)
	}

	fav, err := Resolve("https://space.bilibili.com/2/favlist?fid=123")
	if err != nil {
		t.Fatal(err)
	}
	if fav.RemoteID != "fav:123" {
		t.Fatalf("fav %+v", fav)
	}

	video, err := Resolve("https://www.bilibili.com/video/BV1xx411c7mD")
	if err != nil {
		t.Fatal(err)
	}
	if !video.LinkOnly || video.RemoteID != "BV1xx411c7mD" {
		t.Fatalf("video %+v", video)
	}
}

func TestResolveRejectsEmpty(t *testing.T) {
	if _, err := Resolve("  "); err == nil {
		t.Fatal("expected error")
	}
}

func TestYouTubeChannelIDPattern(t *testing.T) {
	if !regexp.MustCompile(`^UC[\w-]{20,}$`).MatchString("UCXuqSBlHAE6Xw-yeJA0Tunw") {
		t.Fatal("expected channel id to match")
	}
}
