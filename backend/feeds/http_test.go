package feeds

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestProxyFuncExplicitURL(t *testing.T) {
	t.Parallel()

	fn := proxyFunc("http://127.0.0.1:7890")
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/youtube/v3/videos", nil)
	if err != nil {
		t.Fatal(err)
	}
	proxy, err := fn(req)
	if err != nil {
		t.Fatal(err)
	}
	if proxy == nil || proxy.String() != "http://127.0.0.1:7890" {
		t.Fatalf("proxy = %v", proxy)
	}
}

func TestProxyFuncRejectsBadURL(t *testing.T) {
	t.Parallel()

	fn := proxyFunc("not-a-url")
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fn(req); err == nil {
		t.Fatal("expected error for invalid proxy URL")
	}
}

func TestNewHTTPClientUsesExplicitProxy(t *testing.T) {
	t.Parallel()

	client := NewHTTPClient("http://127.0.0.1:7890")
	transport, ok := client.Transport.(*http.Transport)
	if !ok || transport.Proxy == nil {
		t.Fatal("expected transport with Proxy")
	}
	req, _ := http.NewRequest(http.MethodGet, "https://www.youtube.com/", nil)
	proxy, err := transport.Proxy(req)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := url.Parse("http://127.0.0.1:7890")
	if proxy.String() != want.String() {
		t.Fatalf("got %s want %s", proxy, want)
	}
}

func TestWithAPIKeyAddsQueryParam(t *testing.T) {
	t.Parallel()

	var gotURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.String()
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	base := server.Client()
	client := withAPIKey(base, "test-key-123")
	resp, err := client.Get(server.URL + "/youtube/v3/videos?part=id")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if !strings.Contains(gotURL, "key=test-key-123") {
		t.Fatalf("URL missing api key: %s", gotURL)
	}
}
