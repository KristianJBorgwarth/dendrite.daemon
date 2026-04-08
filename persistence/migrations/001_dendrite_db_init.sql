CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY,
    path TEXT UNIQUE,
    title TEXT,
    slug TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS links (
    id TEXT PRIMARY KEY,
    from_note_id TEXT NOT NULL,
    target_slug TEXT NOT NULL,
    display TEXT,
    raw TEXT NOT NULL,
    line INTEGER NOT NULL,
    col INTEGER NOT NULL,
    FOREIGN KEY(from_note_id) REFERENCES notes(id) ON DELETE CASCADE,
);

CREATE INDEX IF NOT EXISTS idx_notes_slug ON notes(slug);
CREATE INDEX idx_links_from ON links(from_note_id);
CREATE INDEX idx_links_target ON links(target_slug);
