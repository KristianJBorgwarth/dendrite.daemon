CREATE TABLE custom_frontmatter(
  note_id TEXT,
  key TEXT,
  value TEXT,
  PRIMARY KEY (note_id, key, value),
  FOREIGN KEY (note_id) REFERENCES note(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cfe_key_value on custom_frontmatter(key, value);
