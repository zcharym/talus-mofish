package feeds

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchRSSParsesAtom(t *testing.T) {
	t.Parallel()

	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Demo Feed</title>
  <entry>
    <id>item-1</id>
    <title>Hello</title>
    <link href="https://example.com/hello"/>
    <updated>2024-01-02T03:04:05Z</updated>
    <author><name>Ada</name></author>
    <summary>A short note.</summary>
  </entry>
</feed>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		_, _ = w.Write([]byte(xml))
	}))
	t.Cleanup(server.Close)

	client := NewClient(server.Client())
	title, items, err := client.Fetch(context.Background(), Spec{
		Kind:    KindRSS,
		Title:   "fallback",
		FeedURL: server.URL,
	}, Credentials{})
	if err != nil {
		t.Fatal(err)
	}
	if title != "Demo Feed" {
		t.Fatalf("title = %q", title)
	}
	if len(items) != 1 || items[0].Title != "Hello" || items[0].Author != "Ada" {
		t.Fatalf("items = %+v", items)
	}
	if items[0].PublishedAt.UTC() != time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC) {
		t.Fatalf("published = %v", items[0].PublishedAt)
	}
}

func TestYouTubeChannelIDFromPage(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html><link rel="alternate" type="application/rss+xml" href="https://www.youtube.com/feeds/videos.xml?channel_id=UCXuqSBlHAE6Xw-yeJA0Tunw"></html>`))
	}))
	t.Cleanup(server.Close)

	client := NewClient(server.Client())
	id, err := client.youtubeChannelIDFromPage(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if id != "UCXuqSBlHAE6Xw-yeJA0Tunw" {
		t.Fatalf("id = %q", id)
	}
}

func TestBilibiliToView(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Cookie"), "SESSDATA=abc") {
			t.Errorf("cookie = %q", r.Header.Get("Cookie"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"0","data":{"count":1,"list":[{"bvid":"BV1xx411c7mD","title":"Clip","pic":"https://i0.hdslb.com/x.jpg","desc":"d","pubdate":1700000000,"owner":{"name":"UP"}}]}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient(server.Client())
	body, err := client.getJSON(context.Background(), server.URL, biliCookie("abc"), nil)
	if err != nil {
		t.Fatal(err)
	}
	title, items, err := parseBiliToView(body)
	if err != nil {
		t.Fatal(err)
	}
	if title != "Bilibili watch later" || len(items) != 1 || items[0].ExternalID != "BV1xx411c7mD" {
		t.Fatalf("%s %#v", title, items)
	}
}
