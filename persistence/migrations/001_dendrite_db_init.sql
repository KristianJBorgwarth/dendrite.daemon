CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY,
    path TEXT UNIQUE,
    title TEXT,
    slug TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS links (
    from_note_id TEXT NOT NULL,
    to_note_id TEXT NOT NULL,
    raw TEXT,
    FOREIGN KEY(from_note_id) REFERENCES notes(id) ON DELETE CASCADE,
    FOREIGN KEY(to_note_id) REFERENCES notes(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_notes_slug ON notes(slug);
CREATE INDEX IF NOT EXISTS idx_links_from ON links(from_note_id);
CREATE INDEX IF NOT EXISTS idx_links_to ON links(to_note_id);

