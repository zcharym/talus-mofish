package types

import "testing"

func TestCloudflareConfigured(t *testing.T) {
	if CloudflareConfigured(Cloudflare{}) {
		t.Fatal("empty credentials should not be configured")
	}
	got := NormalizeCloudflare(Cloudflare{AccountID: "  acc  ", APIToken: "  tok  "})
	if got.AccountID != "acc" || got.APIToken != "tok" {
		t.Fatalf("NormalizeCloudflare = %+v", got)
	}
	if !CloudflareConfigured(got) {
		t.Fatal("expected configured after trim")
	}
}

func TestNormalizeObsidianDefaultURL(t *testing.T) {
	got := NormalizeObsidian(Obsidian{})
	if got.BaseURL != DefaultObsidianBaseURL {
		t.Fatalf("BaseURL = %q", got.BaseURL)
	}
}

func TestNormalizeFeeds(t *testing.T) {
	got := NormalizeFeeds(Feeds{YouTubeAPIKey: "  yt  ", BilibiliSESSDATA: "  sess  "})
	if got.YouTubeAPIKey != "yt" || got.BilibiliSESSDATA != "sess" {
		t.Fatalf("NormalizeFeeds = %+v", got)
	}
}
