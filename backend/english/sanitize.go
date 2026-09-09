package english

import "github.com/microcosm-cc/bluemonday"

var htmlPolicy = bluemonday.UGCPolicy()

// SanitizeHTML strips unsafe markup from imported or user-edited Anki HTML.
func SanitizeHTML(raw string) string {
	if raw == "" {
		return ""
	}
	return htmlPolicy.Sanitize(raw)
}
