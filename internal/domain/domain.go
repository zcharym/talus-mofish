// Package domain defines logical product domains for Talus Echo.
//
// Talus Echo is a chat-oriented desktop agent. Domains are feature areas with
// their own Management routes, data, and agent tools. The first domain is
// English Learning (route prefix "english.*"); additional domains will follow
// the same pattern.
package domain

const (
	// English is the English Learning domain: vocabulary, reading, SRS, Anki import.
	English = "english"
)
