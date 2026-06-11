PRAGMA foreign_keys = ON;

-- Internal app settings (key/value store for engine/runtime state).
CREATE TABLE IF NOT EXISTS settings (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TEXT NOT NULL DEFAULT (datetime ('now'))
);

-- === SRS Engine ===
CREATE TABLE IF NOT EXISTS decks (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    parent_id TEXT REFERENCES decks (id),
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime ('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime ('now'))
);

CREATE TABLE IF NOT EXISTS cards (
    id TEXT NOT NULL PRIMARY KEY,
    deck_id TEXT NOT NULL REFERENCES decks (id),
    front TEXT NOT NULL,
    back TEXT NOT NULL,
    example_sentence TEXT NOT NULL DEFAULT '',
    hints TEXT NOT NULL DEFAULT '',
    card_type TEXT NOT NULL DEFAULT 'vocabulary' CHECK (
        card_type IN ('vocabulary', 'sentence', 'cloze', 'listening')
    ),
    ease REAL NOT NULL DEFAULT 2.5,
    interval INTEGER NOT NULL DEFAULT 0,
    repetitions INTEGER NOT NULL DEFAULT 0,
    due_at TEXT,
    lapsed INTEGER NOT NULL DEFAULT 0,
    suspended INTEGER NOT NULL DEFAULT 0,
    source TEXT NOT NULL DEFAULT '',
    anki_card_id INTEGER,
    anki_note_guid TEXT,
    template_ord INTEGER NOT NULL DEFAULT 0,
    model_css TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime ('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime ('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cards_anki_card_id ON cards (anki_card_id)
WHERE
    anki_card_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS review_log (
    id TEXT NOT NULL PRIMARY KEY,
    card_id TEXT NOT NULL REFERENCES cards (id),
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 4),
    elapsed_ms INTEGER NOT NULL,
    review_at TEXT NOT NULL DEFAULT (datetime ('now')),
    day_seq INTEGER NOT NULL,
    ease_before REAL,
    interval_before INTEGER
);

CREATE INDEX IF NOT EXISTS idx_cards_due ON cards (deck_id, due_at)
WHERE
    suspended = 0;

CREATE INDEX IF NOT EXISTS idx_review_log_card ON review_log (card_id);

CREATE INDEX IF NOT EXISTS idx_review_log_date ON review_log (review_at);

-- === Vocabulary ===
CREATE TABLE IF NOT EXISTS vocabulary (
    id TEXT NOT NULL PRIMARY KEY,
    word TEXT NOT NULL,
    phonetic TEXT NOT NULL DEFAULT '',
    pos TEXT NOT NULL DEFAULT '',
    definition TEXT NOT NULL,
    definition_en TEXT NOT NULL DEFAULT '',
    examples TEXT NOT NULL DEFAULT '',
    level TEXT NOT NULL DEFAULT 'custom' CHECK (
        level IN ('cet4', 'cet6', 'ielts', 'toefl', 'gre', 'custom')
    ),
    source TEXT NOT NULL DEFAULT 'manual',
    anki_note_guid TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime ('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_vocabulary_anki_note_guid ON vocabulary (anki_note_guid)
WHERE
    anki_note_guid IS NOT NULL;

CREATE TABLE IF NOT EXISTS card_vocab (
    card_id TEXT NOT NULL REFERENCES cards (id) ON DELETE CASCADE,
    vocab_id TEXT NOT NULL REFERENCES vocabulary (id) ON DELETE CASCADE,
    PRIMARY KEY (card_id, vocab_id)
);

-- === Reading / Articles ===
CREATE TABLE IF NOT EXISTS articles (
    id TEXT NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL,
    translation TEXT NOT NULL DEFAULT '',
    level TEXT NOT NULL DEFAULT 'intermediate' CHECK (
        level IN (
            'beginner',
            'elementary',
            'intermediate',
            'advanced'
        )
    ),
    source TEXT NOT NULL DEFAULT 'manual' CHECK (
        source IN (
            'manual',
            'clipboard',
            'web',
            'youtube',
            'import:anki'
        )
    ),
    word_count INTEGER NOT NULL DEFAULT 0,
    anki_note_guid TEXT,
    model_css TEXT NOT NULL DEFAULT '',
    raw_fields TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime ('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime ('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_articles_anki_note_guid ON articles (anki_note_guid)
WHERE
    anki_note_guid IS NOT NULL;

CREATE TABLE IF NOT EXISTS clippings (
    id TEXT NOT NULL PRIMARY KEY,
    article_id TEXT REFERENCES articles (id),
    text TEXT NOT NULL,
    translation TEXT NOT NULL DEFAULT '',
    explanation TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    saved_to_card INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime ('now'))
);

-- === Media Assets ===
CREATE TABLE IF NOT EXISTS media_assets (
    id TEXT NOT NULL PRIMARY KEY,
    original_filename TEXT NOT NULL,
    stored_path TEXT NOT NULL,
    mime_type TEXT NOT NULL DEFAULT '',
    sha256 TEXT NOT NULL DEFAULT '',
    source TEXT NOT NULL DEFAULT 'import:anki',
    created_at TEXT NOT NULL DEFAULT (datetime ('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_media_assets_filename ON media_assets (original_filename, source);

-- === Anki Import ===
CREATE TABLE IF NOT EXISTS anki_imports (
    id TEXT NOT NULL PRIMARY KEY,
    filename TEXT NOT NULL,
    imported_at TEXT NOT NULL DEFAULT (datetime ('now')),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (
        status IN ('pending', 'completed', 'failed', 'partial')
    ),
    error_message TEXT NOT NULL DEFAULT '',
    stats_json TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS anki_import_decks (
    import_id TEXT NOT NULL REFERENCES anki_imports (id) ON DELETE CASCADE,
    anki_deck_id INTEGER NOT NULL,
    anki_deck_name TEXT NOT NULL,
    target_type TEXT NOT NULL CHECK (target_type IN ('vocabulary', 'reading', 'skip')),
    anki_model_id INTEGER,
    field_mapping TEXT NOT NULL DEFAULT '',
    app_deck_id TEXT REFERENCES decks (id),
    PRIMARY KEY (import_id, anki_deck_id)
);
