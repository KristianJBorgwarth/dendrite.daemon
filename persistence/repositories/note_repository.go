package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type INoteRepository interface {
	Insert(ctx context.Context, dbContext persistence.IDbContext, note *models.Note) error
	InsertRange(ctx context.Context, dbContext persistence.IDbContext, note []*models.Note) error
	Update(ctx context.Context, dbContext persistence.IDbContext, noteID, path, title, slug string) error
	Search(ctx context.Context, query string) ([]*models.Note, error)
	GetBySlug(ctx context.Context, slug string) (*models.Note, error)
	GetByPath(ctx context.Context, path string) (*models.Note, error)
	GetByTag(ctx context.Context, tag string) ([]*models.Note, error)
	GetAll(ctx context.Context, skip, top int) ([]*models.Note, error)
	GetNoteCount(ctx context.Context) (int, error)
}

type noteRepository struct {
	readDBContext persistence.ReadContext
}

func NewNoteRepository(rdb persistence.ReadContext) INoteRepository {
	return &noteRepository{readDBContext: rdb}
}

func (r *noteRepository) Insert(ctx context.Context, dbContext persistence.IDbContext, note *models.Note) error {
	query := `
	INSERT INTO note (id, title, path, slug, created_at, updated_at)
	VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	ON CONFLICT DO NOTHING;
	`
	_, err := dbContext.ExecContext(ctx, query, note.ID(), note.Title(), note.Path(), note.Slug())
	return err
}

func (r *noteRepository) InsertRange(ctx context.Context, dbContext persistence.IDbContext, notes []*models.Note) error {
	if len(notes) == 0 {
		return nil
	}
	noteStatement, err := dbContext.Prepare(`
	INSERT INTO note (id, title, path, slug, created_at, updated_at)
	VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))`)
	if err != nil {
		return err
	}
	for _, note := range notes {
		if _, err := noteStatement.ExecContext(ctx, note.ID(), note.Title(), note.Path(), note.Slug()); err != nil {
			return err
		}
	}
	return nil
}

func (r *noteRepository) Update(ctx context.Context, dbContext persistence.IDbContext, noteID, path, title, slug string) error {
	query := `
	UPDATE note
	SET path = ?, title = ?, slug = ?, updated_at = datetime('now')
	WHERE id = ?;
	`
	_, err := dbContext.ExecContext(ctx, query, path, title, slug, noteID)
	if err != nil {
		return err
	}
	return nil
}

func (r *noteRepository) GetBySlug(ctx context.Context, slug string) (*models.Note, error) {
	query := `SELECT id, title, path, slug, created_at, updated_at FROM note WHERE slug = ?`
	row := r.readDBContext.QueryRowContext(ctx, query, slug)

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

func (r *noteRepository) GetByPath(ctx context.Context, path string) (*models.Note, error) {
	query := `SELECT id, title, path, slug, created_at, updated_at FROM note WHERE path = ?`
	row := r.readDBContext.QueryRowContext(ctx, query, path)

	var id, title, slug, createdAt, updatedAt string
	err := row.Scan(&id, &title, &path, &slug, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return models.NewNote(id, path, title, slug, createdAt, updatedAt), nil
}

func (r *noteRepository) GetNoteCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM note`
	row := r.readDBContext.QueryRowContext(ctx, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *noteRepository) GetByTag(ctx context.Context, tag string) ([]*models.Note, error) {
	rows, err := r.readDBContext.QueryContext(ctx, `
	SELECT DISTINCT n.id, n.title, n.path, n.slug, n.created_at, n.updated_at
	FROM note n
	JOIN note_tag nt ON n.id = nt.note_id
	JOIN tag t ON t.name = nt.tag_id
	WHERE t.name LIKE ?`, "%"+tag+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]*models.Note, 0)
	for rows.Next() {
		var id, title, path, slug, createdAt, updatedAt string
		if err := rows.Scan(&id, &title, &path, &slug, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, models.NewNote(id, path, title, slug, createdAt, updatedAt))
	}

	return notes, nil
}

func (r *noteRepository) GetAll(ctx context.Context, skip, top int) ([]*models.Note, error) {
	rows, err := r.readDBContext.QueryContext(ctx, `
	SELECT id, title, path, slug, created_at, updated_at
	FROM note
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?`, top, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]*models.Note, 0)
	for rows.Next() {
		var id, title, path, slug, createdAt, updatedAt string
		if err := rows.Scan(&id, &title, &path, &slug, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, models.NewNote(id, path, title, slug, createdAt, updatedAt))
	}

	return notes, nil
}

func (r *noteRepository) Search(ctx context.Context, query string) ([]*models.Note, error) {
	rows, err := r.readDBContext.QueryContext(ctx, `
	SELECT id, title, path, slug, created_at, updated_at
	FROM note
	WHERE title LIKE ? OR path LIKE ? OR slug LIKE ?
	ORDER BY created_at DESC`, "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]*models.Note, 0)
	for rows.Next() {
		var id, title, path, slug, createdAt, updatedAt string
		if err := rows.Scan(&id, &title, &path, &slug, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, models.NewNote(id, path, title, slug, createdAt, updatedAt))
	}

	return notes, nil
}
