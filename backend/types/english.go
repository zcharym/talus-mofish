package types

// Vocabulary is an English Learning word entry on the Wails wire.
type Vocabulary struct {
	ID           string `json:"id"`
	Word         string `json:"word"`
	Phonetic     string `json:"phonetic"`
	Pos          string `json:"pos"`
	Definition   string `json:"definition"`
	DefinitionEn string `json:"definition_en"`
	Examples     string `json:"examples"`
	Level        string `json:"level"`
	Source       string `json:"source"`
	AnkiNoteGUID string `json:"anki_note_guid"`
	CreatedAt    string `json:"created_at"`
}

// VocabularyUpdate is the editable vocabulary fields.
type VocabularyUpdate struct {
	ID           string `json:"id"`
	Word         string `json:"word"`
	Phonetic     string `json:"phonetic"`
	Pos          string `json:"pos"`
	Definition   string `json:"definition"`
	DefinitionEn string `json:"definition_en"`
	Examples     string `json:"examples"`
	Level        string `json:"level"`
	Source       string `json:"source"`
}

// Article is a full reading item. List views use ArticleSummary.
type Article struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Content     string `json:"content"`
	Translation string `json:"translation"`
	Level       string `json:"level"`
	Source      string `json:"source"`
	WordCount   int64  `json:"word_count"`
	ModelCSS    string `json:"model_css"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Card is an SRS card. Internal Anki IDs stay off the wire.
type Card struct {
	ID              string `json:"id"`
	DeckID          string `json:"deck_id"`
	Front           string `json:"front"`
	Back            string `json:"back"`
	ExampleSentence string `json:"example_sentence"`
	Hints           string `json:"hints"`
	CardType        string `json:"card_type"`
	ModelCSS        string `json:"model_css"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// CardContentUpdate is the editable HTML/text fields on an SRS card.
type CardContentUpdate struct {
	ID              string `json:"id"`
	Front           string `json:"front"`
	Back            string `json:"back"`
	ExampleSentence string `json:"example_sentence"`
	Hints           string `json:"hints"`
	CardType        string `json:"card_type"`
	ModelCSS        string `json:"model_css"`
}

// Deck is an SRS deck.
type Deck struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int64  `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// AnkiImportRecord is a past APKG import session.
type AnkiImportRecord struct {
	ID           string `json:"id"`
	Filename     string `json:"filename"`
	ImportedAt   string `json:"imported_at"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	StatsJSON    string `json:"stats_json"`
}
