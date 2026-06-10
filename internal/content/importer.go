package content

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/songwei.ma/talus-mofish/internal/content/apkg"
	"github.com/songwei.ma/talus-mofish/internal/database"
	"github.com/songwei.ma/talus-mofish/internal/store"
)

// Importer imports Anki APKG files into the application database.
type Importer struct {
	db        *database.DB
	mediaRoot string
}

func NewImporter(db *database.DB) (*Importer, error) {
	root, err := DefaultMediaDir()
	if err != nil {
		return nil, err
	}
	return &Importer{db: db, mediaRoot: root}, nil
}

// DefaultMediaDir returns the per-user media storage directory.
func DefaultMediaDir() (string, error) {
	base, err := database.UserDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "talus-mofish", "media"), nil
}

// Preview reads an APKG and returns deck/model metadata for the import UI.
func (im *Importer) Preview(apkgPath string) (AnkiPreview, error) {
	coll, cleanup, err := apkg.OpenArchive(apkgPath)
	if err != nil {
		return AnkiPreview{}, err
	}
	defer cleanup()

	filename := filepath.Base(apkgPath)
	deckStats := map[int64]*AnkiDeckPreview{}

	for _, card := range coll.Cards {
		p, ok := deckStats[card.DID]
		if !ok {
			name := "Unknown"
			if d, ok := coll.Decks[card.DID]; ok {
				name = d.Name
			}
			note := coll.NoteByID(card.NID)
			modelID := int64(0)
			modelName := ""
			fields := []string{}
			isCloze := false
			if note != nil {
				modelID = note.MID
				if m := coll.ModelByID(note.MID); m != nil {
					modelName = m.Name
					isCloze = m.Type == 1
					for _, f := range m.Fields {
						fields = append(fields, f.Name)
					}
				}
			}
			p = &AnkiDeckPreview{
				AnkiDeckID:  card.DID,
				Name:        name,
				AnkiModelID: modelID,
				ModelName:   modelName,
				Fields:      fields,
				IsCloze:     isCloze,
			}
			deckStats[card.DID] = p
		}
		p.CardCount++
	}

	noteInDeck := map[int64]map[int64]bool{}
	for _, card := range coll.Cards {
		if noteInDeck[card.DID] == nil {
			noteInDeck[card.DID] = map[int64]bool{}
		}
		noteInDeck[card.DID][card.NID] = true
	}
	for did, notes := range noteInDeck {
		if p, ok := deckStats[did]; ok {
			p.NoteCount = len(notes)
		}
	}

	decks := make([]AnkiDeckPreview, 0, len(deckStats))
	for _, p := range deckStats {
		decks = append(decks, *p)
	}
	sortDecksByName(decks)

	return AnkiPreview{Filename: filename, Decks: decks}, nil
}

// Import applies user deck mappings and writes vocabulary/cards/articles.
func (im *Importer) Import(apkgPath string, configs []ImportDeckConfig) (ImportResult, error) {
	coll, cleanup, err := apkg.OpenArchive(apkgPath)
	if err != nil {
		return ImportResult{}, err
	}
	defer cleanup()

	importID := uuid.NewString()
	stats := ImportStats{}
	ctx := context.Background()

	if err := im.db.Queries.CreateAnkiImport(ctx, store.CreateAnkiImportParams{
		ID:           importID,
		Filename:     filepath.Base(apkgPath),
		Status:       "pending",
		ErrorMessage: "",
		StatsJson:    "",
	}); err != nil {
		return ImportResult{}, fmt.Errorf("create import record: %w", err)
	}

	mediaStore, err := apkg.NewMediaStore(im.mediaRoot)
	if err != nil {
		return ImportResult{}, err
	}
	if err := mediaStore.ImportFromAnki(coll.MediaDir()); err != nil {
		return ImportResult{}, fmt.Errorf("import media: %w", err)
	}
	stats.MediaFiles = mediaStore.Count()

	tx, err := im.db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return ImportResult{}, err
	}
	defer tx.Rollback()

	qtx := im.db.Queries.WithTx(tx)

	for _, cfg := range configs {
		if cfg.TargetType == "skip" {
			continue
		}
		stats.DecksProcessed++

		var appDeckID sql.NullString
		switch cfg.TargetType {
		case "vocabulary":
			deckID, v, c, skippedNotes, skippedCards, err := im.importVocabularyDeck(ctx, qtx, coll, mediaStore, cfg)
			if err != nil {
				return ImportResult{}, err
			}
			appDeckID = sql.NullString{String: deckID, Valid: deckID != ""}
			stats.VocabCreated += v
			stats.CardsCreated += c
			stats.SkippedNotes += skippedNotes
			stats.SkippedCards += skippedCards
		case "reading":
			a, skipped, err := im.importReadingDeck(ctx, qtx, coll, mediaStore, cfg)
			if err != nil {
				return ImportResult{}, err
			}
			stats.ArticlesCreated += a
			stats.SkippedNotes += skipped
		}

		modelID := sql.NullInt64{}
		if cfg.AnkiModelID > 0 {
			modelID = sql.NullInt64{Int64: cfg.AnkiModelID, Valid: true}
		}
		fieldMappingJSON, _ := json.Marshal(cfg.FieldMapping)
		if err := qtx.CreateAnkiImportDeck(ctx, store.CreateAnkiImportDeckParams{
			ImportID:     importID,
			AnkiDeckID:   cfg.AnkiDeckID,
			AnkiDeckName: cfg.AnkiDeckName,
			TargetType:   cfg.TargetType,
			AnkiModelID:  modelID,
			FieldMapping: string(fieldMappingJSON),
			AppDeckID:    appDeckID,
		}); err != nil {
			return ImportResult{}, fmt.Errorf("record deck mapping: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return ImportResult{}, err
	}

	if err := im.recordMediaAssets(ctx, mediaStore); err != nil {
		return ImportResult{}, err
	}

	statsJSON, _ := json.Marshal(stats)
	status := "completed"
	if err := im.db.Queries.UpdateAnkiImportStatus(ctx, store.UpdateAnkiImportStatusParams{
		Status:       status,
		ErrorMessage: "",
		StatsJson:    string(statsJSON),
		ID:           importID,
	}); err != nil {
		return ImportResult{}, err
	}

	return ImportResult{
		ImportID: importID,
		Status:   status,
		Stats:    stats,
	}, nil
}

func (im *Importer) importVocabularyDeck(
	ctx context.Context,
	q *store.Queries,
	coll *apkg.Collection,
	media *apkg.MediaStore,
	cfg ImportDeckConfig,
) (appDeckID string, vocabCreated, cardsCreated, skippedNotes, skippedCards int, err error) {
	model := coll.ModelByID(cfg.AnkiModelID)
	if model == nil {
		return "", 0, 0, 0, 0, fmt.Errorf("model %d not found for deck %s", cfg.AnkiModelID, cfg.AnkiDeckName)
	}

	appDeckID = uuid.NewString()
	if err := q.CreateDeck(ctx, store.CreateDeckParams{
		ID:          appDeckID,
		Name:        cfg.AnkiDeckName,
		Description: "Imported from Anki",
		ParentID:    sql.NullString{},
		SortOrder:   0,
	}); err != nil {
		return "", 0, 0, 0, 0, err
	}

	notesInDeck := map[int64]bool{}
	for _, card := range coll.Cards {
		if card.DID == cfg.AnkiDeckID {
			notesInDeck[card.NID] = true
		}
	}

	for nid := range notesInDeck {
		note := coll.NoteByID(nid)
		if note == nil || note.MID != cfg.AnkiModelID {
			continue
		}

		_, err := q.GetVocabularyByAnkiNoteGuid(ctx, sql.NullString{String: note.GUID, Valid: true})
		if err == nil {
			skippedNotes++
			continue
		}
		if err != nil && !errorsIsNoRows(err) {
			return "", 0, 0, 0, 0, err
		}

		vocabID := uuid.NewString()
		word := apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "word")
		if word == "" && len(note.Fields) > 0 {
			word = apkg.FieldValue(note.Fields, 0)
		}
		definition := apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "definition")
		if definition == "" && len(note.Fields) > 1 {
			definition = apkg.FieldValue(note.Fields, 1)
		}
		if definition == "" {
			definition = word
		}

		examples := apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "examples")
		if err := q.CreateVocabulary(ctx, store.CreateVocabularyParams{
			ID:           vocabID,
			Word:         word,
			Phonetic:     apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "phonetic"),
			Pos:          apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "pos"),
			Definition:   definition,
			DefinitionEn: apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "definition_en"),
			Examples:     examplesJSON(examples),
			Level:        "custom",
			Source:       "import:anki",
			AnkiNoteGuid: sql.NullString{String: note.GUID, Valid: true},
		}); err != nil {
			return "", 0, 0, 0, 0, err
		}
		vocabCreated++

		for _, cardRow := range coll.Cards {
			if cardRow.DID != cfg.AnkiDeckID || cardRow.NID != nid {
				continue
			}
			if _, err := q.GetCardByAnkiCardID(ctx, sql.NullInt64{Int64: cardRow.ID, Valid: true}); err == nil {
				skippedCards++
				continue
			}
			if err != nil && !errorsIsNoRows(err) {
				return "", 0, 0, 0, 0, err
			}

			front, back := apkg.RenderCard(model, note.Fields, int(cardRow.Ord))
			front = media.RewriteHTML(front)
			back = media.RewriteHTML(back)

			cardType := "vocabulary"
			if model.Type == 1 {
				cardType = "cloze"
			}

			cardID := uuid.NewString()
			if err := q.CreateCard(ctx, store.CreateCardParams{
				ID:              cardID,
				DeckID:          appDeckID,
				Front:           front,
				Back:            back,
				ExampleSentence: examples,
				Hints:           "",
				CardType:        cardType,
				Ease:            2.5,
				Interval:        0,
				Repetitions:     0,
				DueAt:           sql.NullString{},
				Lapsed:          0,
				Suspended:       0,
				Source:          "import:anki",
				AnkiCardID:      sql.NullInt64{Int64: cardRow.ID, Valid: true},
				AnkiNoteGuid:    sql.NullString{String: note.GUID, Valid: true},
				TemplateOrd:     cardRow.Ord,
				ModelCss:        model.CSS,
			}); err != nil {
				return "", 0, 0, 0, 0, err
			}
			if err := q.LinkCardVocab(ctx, store.LinkCardVocabParams{
				CardID:  cardID,
				VocabID: vocabID,
			}); err != nil {
				return "", 0, 0, 0, 0, err
			}
			cardsCreated++
		}
	}

	return appDeckID, vocabCreated, cardsCreated, skippedNotes, skippedCards, nil
}

func (im *Importer) importReadingDeck(
	ctx context.Context,
	q *store.Queries,
	coll *apkg.Collection,
	media *apkg.MediaStore,
	cfg ImportDeckConfig,
) (articlesCreated, skipped int, err error) {
	model := coll.ModelByID(cfg.AnkiModelID)
	if model == nil {
		return 0, 0, fmt.Errorf("model %d not found for deck %s", cfg.AnkiModelID, cfg.AnkiDeckName)
	}

	notesInDeck := map[int64]bool{}
	for _, card := range coll.Cards {
		if card.DID == cfg.AnkiDeckID {
			notesInDeck[card.NID] = true
		}
	}

	for nid := range notesInDeck {
		note := coll.NoteByID(nid)
		if note == nil || note.MID != cfg.AnkiModelID {
			continue
		}

		if _, err := q.GetArticleByAnkiNoteGuid(ctx, sql.NullString{String: note.GUID, Valid: true}); err == nil {
			skipped++
			continue
		}
		if err != nil && !errorsIsNoRows(err) {
			return 0, 0, err
		}

		title := apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "title")
		if title == "" && len(note.Fields) > 0 {
			title = apkg.FieldValue(note.Fields, 0)
		}
		content := apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "content")
		if content == "" && len(note.Fields) > 1 {
			content = apkg.FieldValue(note.Fields, 1)
		}
		translation := apkg.FieldValueByName(note.Fields, cfg.FieldMapping, "translation")

		content = media.RewriteHTML(content)
		translation = media.RewriteHTML(translation)

		rawFields, _ := json.Marshal(map[string]interface{}{
			"fieldNames": modelFieldNames(model),
			"fields":     note.Fields,
			"tags":       note.Tags,
		})

		articleID := uuid.NewString()
		if err := q.CreateArticle(ctx, store.CreateArticleParams{
			ID:           articleID,
			Title:        title,
			Url:          "",
			Content:      content,
			Translation:  translation,
			Level:        "intermediate",
			Source:       "import:anki",
			WordCount:    int64(apkg.WordCountHTML(content)),
			AnkiNoteGuid: sql.NullString{String: note.GUID, Valid: true},
			ModelCss:     model.CSS,
			RawFields:    string(rawFields),
		}); err != nil {
			return 0, 0, err
		}
		articlesCreated++
	}

	return articlesCreated, skipped, nil
}

func (im *Importer) recordMediaAssets(ctx context.Context, media *apkg.MediaStore) error {
	for original, rel := range media.Entries() {
		_, err := im.db.Queries.GetMediaAssetByFilename(ctx, store.GetMediaAssetByFilenameParams{
			OriginalFilename: original,
			Source:           "import:anki",
		})
		if err == nil {
			continue
		}
		if err != nil && !errorsIsNoRows(err) {
			return err
		}
		fullPath := filepath.Join(im.mediaRoot, rel)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		hash := sha256Hex(data)
		if err := im.db.Queries.CreateMediaAsset(ctx, store.CreateMediaAssetParams{
			ID:               uuid.NewString(),
			OriginalFilename: original,
			StoredPath:       rel,
			MimeType:         mimeFromName(original),
			Sha256:           hash,
			Source:           "import:anki",
		}); err != nil {
			return err
		}
	}
	return nil
}

func modelFieldNames(m *apkg.NoteModel) []string {
	out := make([]string, len(m.Fields))
	for i, f := range m.Fields {
		out[i] = f.Name
	}
	return out
}

func examplesJSON(examples string) string {
	if examples == "" {
		return "[]"
	}
	if strings.HasPrefix(strings.TrimSpace(examples), "[") {
		return examples
	}
	b, _ := json.Marshal([]string{examples})
	return string(b)
}

func sortDecksByName(decks []AnkiDeckPreview) {
	for i := 0; i < len(decks); i++ {
		for j := i + 1; j < len(decks); j++ {
			if decks[j].Name < decks[i].Name {
				decks[i], decks[j] = decks[j], decks[i]
			}
		}
	}
}

func errorsIsNoRows(err error) bool {
	return err == sql.ErrNoRows
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func mimeFromName(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lower, ".png"):
		return "image/png"
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lower, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lower, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(lower, ".mp3"):
		return "audio/mpeg"
	case strings.HasSuffix(lower, ".wav"):
		return "audio/wav"
	default:
		return "application/octet-stream"
	}
}
