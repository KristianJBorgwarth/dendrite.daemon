package repositories

import (
	"context"
	"database/sql"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
)

type NoteRepository interface {
	Upsert(ctx context.Context, note *models.Note) error
	GetBySlug(ctx context.Context, slug string) (*models.Note, error)
}

type noteRepository struct {
	Transaction *sql.Tx
}

func NewNoteRepository(tx *sql.Tx) NoteRepository {
	return &noteRepository{Transaction: tx}
}

func (r *noteRepository) Upsert(ctx context.Context, note *models.Note) error {
	query := `
	INSERT INTO notes (id, title, path, slug, created_at, updated_at)
	VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	ON CONFLICT (slug) DO UPDATE
	SET title = EXCLUDED.title,
	    path = EXCLUDED.path;
	`
	_, err := r.Transaction.ExecContext(ctx, query, note.ID(), note.Title(), note.Path(), note.Slug())
	return err
}
