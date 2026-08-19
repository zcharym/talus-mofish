package apkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// NoteModel is an Anki note type from col.models JSON.
type NoteModel struct {
	ID     int64
	Name   string
	Type   int
	CSS    string
	Fields []ModelField
	Tmpls  []CardTemplate
}

type ModelField struct {
	Name string
	Ord  int
}

type CardTemplate struct {
	Name string
	Ord  int
	QFmt string
	AFmt string
}

// DeckInfo is an Anki deck from col.decks JSON.
type DeckInfo struct {
	ID   int64
	Name string
}

// NoteRow is a note with its primary deck (MIN did from cards).
type NoteRow struct {
	ID     int64
	GUID   string
	MID    int64
	Tags   string
	Fields []string
	DeckID int64
}

// CardRow links a note to a deck and template.
type CardRow struct {
	ID     int64
	NID    int64
	DID    int64
	Ord    int64
	Queue  int64
	Type   int64
}

// Collection holds parsed Anki collection data.
type Collection struct {
	db        *sql.DB
	mediaDir  string
	Models    map[int64]*NoteModel
	Decks     map[int64]*DeckInfo
	Notes     []NoteRow
	Cards     []CardRow
	notesByID map[int64]*NoteRow
}

func OpenCollection(dbPath, mediaDir string) (*Collection, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open collection db: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping collection db: %w", err)
	}

	c := &Collection{
		db:       db,
		mediaDir: mediaDir,
		Models:   map[int64]*NoteModel{},
		Decks:    map[int64]*DeckInfo{},
		notesByID: map[int64]*NoteRow{},
	}

	if err := c.loadModels(); err != nil {
		db.Close()
		return nil, err
	}
	if err := c.loadDecks(); err != nil {
		db.Close()
		return nil, err
	}
	if err := c.loadNotes(); err != nil {
		db.Close()
		return nil, err
	}
	if err := c.loadCards(); err != nil {
		db.Close()
		return nil, err
	}

	return c, nil
}

func (c *Collection) Close() error {
	if c.db == nil {
		return nil
	}
	return c.db.Close()
}

func (c *Collection) MediaDir() string {
	return c.mediaDir
}

func (c *Collection) NoteByID(id int64) *NoteRow {
	return c.notesByID[id]
}

func (c *Collection) ModelByID(id int64) *NoteModel {
	return c.Models[id]
}

func decodeJSON(text string, v interface{}) error {
	return json.Unmarshal([]byte(text), v)
}

func (c *Collection) loadModels() error {
	var modelsJSON string
	err := c.db.QueryRow("SELECT models FROM col").Scan(&modelsJSON)
	if err != nil {
		return fmt.Errorf("read col.models: %w", err)
	}

	raw := map[string]json.RawMessage{}
	if err := decodeJSON(modelsJSON, &raw); err != nil {
		return fmt.Errorf("parse models json: %w", err)
	}

	for key, blob := range raw {
		id, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			continue
		}
		var m struct {
			Name string `json:"name"`
			Type int    `json:"type"`
			CSS  string `json:"css"`
			Flds []struct {
				Name string `json:"name"`
				Ord  int    `json:"ord"`
			} `json:"flds"`
			Tmpls []struct {
				Name string `json:"name"`
				Ord  int    `json:"ord"`
				QFmt string `json:"qfmt"`
				AFmt string `json:"afmt"`
			} `json:"tmpls"`
		}
		if err := json.Unmarshal(blob, &m); err != nil {
			return fmt.Errorf("parse model %s: %w", key, err)
		}

		model := &NoteModel{
			ID:   id,
			Name: m.Name,
			Type: m.Type,
			CSS:  m.CSS,
		}
		for _, f := range m.Flds {
			model.Fields = append(model.Fields, ModelField{Name: f.Name, Ord: f.Ord})
		}
		for _, t := range m.Tmpls {
			model.Tmpls = append(model.Tmpls, CardTemplate{
				Name: t.Name,
				Ord:  t.Ord,
				QFmt: t.QFmt,
				AFmt: t.AFmt,
			})
		}
		c.Models[id] = model
	}
	return nil
}

func (c *Collection) loadDecks() error {
	var decksJSON string
	err := c.db.QueryRow("SELECT decks FROM col").Scan(&decksJSON)
	if err != nil {
		return fmt.Errorf("read col.decks: %w", err)
	}

	raw := map[string]json.RawMessage{}
	if err := decodeJSON(decksJSON, &raw); err != nil {
		return fmt.Errorf("parse decks json: %w", err)
	}

	for key, blob := range raw {
		id, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			continue
		}
		var d struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(blob, &d); err != nil {
			return fmt.Errorf("parse deck %s: %w", key, err)
		}
		c.Decks[id] = &DeckInfo{ID: id, Name: d.Name}
	}
	return nil
}

func (c *Collection) loadNotes() error {
	rows, err := c.db.Query(`
		SELECT notes.id, notes.guid, notes.mid, notes.tags, notes.flds, MIN(cards.did)
		FROM notes
		JOIN cards ON cards.nid = notes.id
		GROUP BY notes.id
	`)
	if err != nil {
		return fmt.Errorf("query notes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, mid, deckID int64
		var guid, tags, flds string
		if err := rows.Scan(&id, &guid, &mid, &tags, &flds, &deckID); err != nil {
			return fmt.Errorf("scan note: %w", err)
		}
		note := NoteRow{
			ID:     id,
			GUID:   guid,
			MID:    mid,
			Tags:   tags,
			Fields: splitFields(flds),
			DeckID: deckID,
		}
		c.Notes = append(c.Notes, note)
		c.notesByID[id] = &c.Notes[len(c.Notes)-1]
	}
	return rows.Err()
}

func (c *Collection) loadCards() error {
	rows, err := c.db.Query(`
		SELECT id, nid, did, ord, queue, type FROM cards
	`)
	if err != nil {
		return fmt.Errorf("query cards: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var card CardRow
		if err := rows.Scan(&card.ID, &card.NID, &card.DID, &card.Ord, &card.Queue, &card.Type); err != nil {
			return fmt.Errorf("scan card: %w", err)
		}
		c.Cards = append(c.Cards, card)
	}
	return rows.Err()
}

func splitFields(flds string) []string {
	if flds == "" {
		return []string{}
	}
	return strings.Split(flds, fieldSeparator)
}
