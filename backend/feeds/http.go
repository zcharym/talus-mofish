package feeds

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mattn/go-ieproxy"
)

// NewHTTPClient builds an HTTP client for feed fetches.
// When proxyURL is empty, the OS/system proxy is used (Windows IE settings via ieproxy).
// When set (e.g. http://127.0.0.1:7890), that proxy is used for all Feeds requests.
func NewHTTPClient(proxyURL string) *http.Client {
	return &http.Client{
		Timeout:   20 * time.Second,
		Transport: feedsTransport(proxyURL),
	}
}

func feedsTransport(proxyURL string) *http.Transport {
	return &http.Transport{
		Proxy: proxyFunc(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 15 * time.Second,
		ForceAttemptHTTP2:   true,
	}
}

func proxyFunc(proxyURL string) func(*http.Request) (*url.URL, error) {
	proxyURL = strings.TrimSpace(proxyURL)
	if proxyURL == "" {
		return ieproxy.GetProxyFunc()
	}
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return func(*http.Request) (*url.URL, error) {
			return nil, fmt.Errorf("feeds proxy URL: %w", err)
		}
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return func(*http.Request) (*url.URL, error) {
			return nil, fmt.Errorf("feeds proxy URL must include scheme and host (e.g. http://127.0.0.1:7890)")
		}
	}
	return http.ProxyURL(parsed)
}

// withAPIKey returns a shallow copy of client whose transport adds ?key= for Google APIs.
// Required when using a custom HTTP client (proxy): option.WithAPIKey is ignored then.
func withAPIKey(client *http.Client, apiKey string) *http.Client {
	if client == nil {
		client = NewHTTPClient("")
	}
	base := client.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	return &http.Client{
		Transport:     &apiKeyRoundTripper{base: base, apiKey: strings.TrimSpace(apiKey)},
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
		Timeout:       client.Timeout,
	}
}

type apiKeyRoundTripper struct {
	base   http.RoundTripper
	apiKey string
}

func (t *apiKeyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	q := clone.URL.Query()
	q.Set("key", t.apiKey)
	clone.URL.RawQuery = q.Encode()
	return t.base.RoundTrip(clone)
}
