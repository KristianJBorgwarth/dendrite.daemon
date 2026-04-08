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

func (r *noteRepository) GetBySlug(ctx context.Context, slug string) (*models.Note, error) {
	query := `SELECT id, title, path, slug, created_at, updated_at FROM notes WHERE slug = ?`
	row := r.Transaction.QueryRowContext(ctx, query, slug)

	var id, title, path, createdAt, updatedAt string
	err := row.Scan(&id, &title, &path, &slug, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return models.NewNote(id, path, title, slug, createdAt, updatedAt), nil
}
