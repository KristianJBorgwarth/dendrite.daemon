package repositories

import (
	"context"
	"database/sql"
)

type NoteRepository interface {
	Upsert(ctx context.Context, title string, path string, slug string) error
}

type noteRepository struct {
	Transaction *sql.Tx
}

func NewNoteRepository(tx *sql.Tx) NoteRepository {
	return &noteRepository{Transaction: tx}
}

func (r *noteRepository) Upsert(ctx context.Context, title string, path string, slug string) error {
	query := `
	INSERT INTO notes (title, path, slug, created_at, updated_at)
	VALUES (?, ?, ?, datetime('now'), datetime('now'))
	ON CONFLICT (slug) DO UPDATE
	SET title = EXCLUDED.title,
	    path = EXCLUDED.path;
	`
	_, err := r.Transaction.ExecContext(ctx, query, title, path, slug)
	return err
}
