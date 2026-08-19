package types

import "github.com/songwei.ma/talus-mofish/backend/storage/store"

// ArticleSummary is a lightweight article row for list views.
type ArticleSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Source    string `json:"source"`
	WordCount int64  `json:"word_count"`
	CreatedAt string `json:"created_at"`
}

// ArticlePageResult is a paginated list of article summaries.
type ArticlePageResult struct {
	Items    []ArticleSummary `json:"items"`
	Total    int64            `json:"total"`
	Page     int64            `json:"page"`
	PageSize int64            `json:"page_size"`
}

// VocabularyPageResult is a paginated list of vocabulary entries.
type VocabularyPageResult struct {
	Items    []store.Vocabulary `json:"items"`
	Total    int64              `json:"total"`
	Page     int64              `json:"page"`
	PageSize int64              `json:"page_size"`
}
