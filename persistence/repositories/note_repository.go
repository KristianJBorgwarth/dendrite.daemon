package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type NoteRepository interface {
	Insert(ctx context.Context, dbContext persistence.IDbContext, note *models.Note) error
	GetBySlug(ctx context.Context, dbContext persistence.IDbContext, slug string) (*models.Note, error)
}

type noteRepository struct{}

func NewNoteRepository() NoteRepository {
	return &noteRepository{}
}

func (r *noteRepository) Insert(ctx context.Context, dbContext persistence.IDbContext, note *models.Note) error {
	query := `
	INSERT INTO note (id, title, path, slug, created_at, updated_at)
	VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	ON CONFLICT (slug) DO UPDATE
	SET title = EXCLUDED.title,
	    path = EXCLUDED.path;
	`
	_, err := dbContext.ExecContext(ctx, query, note.ID(), note.Title(), note.Path(), note.Slug())
	return err
}

func (r *noteRepository) GetBySlug(ctx context.Context, dbContext persistence.IDbContext, slug string) (*models.Note, error) {
	query := `SELECT id, title, path, slug, created_at, updated_at FROM note WHERE slug = ?`
	row := dbContext.QueryRowContext(ctx, query, slug)

	var id, title, path, createdAt, updatedAt string
	err := row.Scan(&id, &title, &path, &slug, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return models.NewNote(id, path, title, slug, createdAt, updatedAt), nil
}
