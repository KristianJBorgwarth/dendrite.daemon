CREATE TABLE IF NOT EXISTS link (
    id TEXT PRIMARY KEY,
    from_note_id TEXT NOT NULL,
    target_slug TEXT NOT NULL,
    display TEXT,
    raw TEXT NOT NULL,
    line INTEGER NOT NULL,
    col INTEGER NOT NULL,
    FOREIGN KEY(from_note_id) REFERENCES note(id) ON DELETE CASCADE
);

CREATE INDEX idx_link_from ON link(from_note_id);
CREATE INDEX idx_link_target ON link(target_slug);

