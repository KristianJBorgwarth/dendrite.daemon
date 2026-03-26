package repositories

import (
	"context"
)

type NoteRepository interface {
	Upsert(ctx context.Context, title string, path string, slug string) error
}

type noteRepository struct {
	dbCtx DBContext
}

func NewNoteRepository(dbContext DBContext) NoteRepository {
	return &noteRepository{dbCtx: dbContext}
}

func (r *noteRepository) Upsert(ctx context.Context, title string, path string, slug string) error {
	query := `
	INSERT INTO notes (title, path, slug)
	VALUES ($1, $2, $3)
	ON CONFLICT (slug) DO UPDATE
	SET title = EXCLUDED.title,
	    path = EXCLUDED.path;
	`
	_, err := r.dbCtx.ExecContext(ctx, query, title, path, slug)
	return err
}
