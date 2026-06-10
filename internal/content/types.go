package content

// ImportDeckConfig describes how one Anki deck should be imported.
type ImportDeckConfig struct {
	AnkiDeckID   int64          `json:"ankiDeckId"`
	AnkiDeckName string         `json:"ankiDeckName"`
	TargetType   string         `json:"targetType"`
	AnkiModelID  int64          `json:"ankiModelId"`
	FieldMapping map[string]int `json:"fieldMapping"`
}

// AnkiDeckPreview is returned before import so the UI can configure mappings.
type AnkiDeckPreview struct {
	AnkiDeckID  int64    `json:"ankiDeckId"`
	Name        string   `json:"name"`
	CardCount   int      `json:"cardCount"`
	NoteCount   int      `json:"noteCount"`
	AnkiModelID int64    `json:"ankiModelId"`
	ModelName   string   `json:"modelName"`
	Fields      []string `json:"fields"`
	IsCloze     bool     `json:"isCloze"`
}

// AnkiPreview is the full preview payload for an APKG file.
type AnkiPreview struct {
	Filename string            `json:"filename"`
	Decks    []AnkiDeckPreview `json:"decks"`
}

// ImportStats summarizes a completed import run.
type ImportStats struct {
	DecksProcessed int `json:"decksProcessed"`
	VocabCreated   int `json:"vocabCreated"`
	CardsCreated   int `json:"cardsCreated"`
	ArticlesCreated int `json:"articlesCreated"`
	SkippedNotes  int `json:"skippedNotes"`
	SkippedCards  int `json:"skippedCards"`
	MediaFiles     int `json:"mediaFiles"`
}

// ImportResult is returned after import completes.
type ImportResult struct {
	ImportID string      `json:"importId"`
	Status   string      `json:"status"`
	Stats    ImportStats `json:"stats"`
}
