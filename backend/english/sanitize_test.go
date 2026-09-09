package english

import (
	"strings"
	"testing"
)

func TestSanitizeHTMLStripsScript(t *testing.T) {
	got := SanitizeHTML(`<p>ok</p><script>alert(1)</script><img src="x" onerror="alert(1)">`)
	if got == "" {
		t.Fatal("expected sanitized markup")
	}
	lower := strings.ToLower(got)
	if strings.Contains(lower, "<script") || strings.Contains(lower, "onerror") {
		t.Fatalf("unsafe markup remained: %s", got)
	}
}
