CREATE TABLE IF NOT EXISTS note (
    id TEXT PRIMARY KEY,
    path TEXT UNIQUE,
    title TEXT,
    slug TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);


CREATE INDEX IF NOT EXISTS idx_note_slug ON note(slug);

