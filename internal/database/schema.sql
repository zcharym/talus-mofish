PRAGMA foreign_keys = ON;

-- Internal app settings (key/value store for engine/runtime state).
CREATE TABLE IF NOT EXISTS settings (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
