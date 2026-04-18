package repositories

import (
	"context"
	"strings"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type ITagRepository interface {
	Insert(ctx context.Context, dbCtx persistence.IDbContext, tags []*models.Tag) error
	InsertRage(ctx context.Context, dbCtx persistence.IDbContext, tags []*models.Tag) error
	InsertNoteTags(ctx context.Context, dbCtx persistence.IDbContext, noteTags []*models.NoteTag) error
	GetByNames(ctx context.Context, dbCtx persistence.IDbContext, names []string) ([]*models.Tag, error)
	DeleteNoteTags(ctx context.Context, dbCtx persistence.IDbContext, noteID string) error
}

type tagRepository struct{}

func NewTagRepository() ITagRepository {
	return &tagRepository{}
}

func (r *tagRepository) Insert(ctx context.Context, dbCtx persistence.IDbContext, tags []*models.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(tags))
	args := make([]any, 0, len(tags))

	for _, tag := range tags {
		placeholders = append(placeholders, "(?)")
		args = append(args, tag.Name())
	}

	query := "INSERT OR IGNORE INTO tag(name) VALUES " + strings.Join(placeholders, ",")

	_, err := dbCtx.ExecContext(ctx, query, args...)
	return err
}

func (r *tagRepository) InsertRage(ctx context.Context, dbCtx persistence.IDbContext, tags []*models.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	tagStatement, err := dbCtx.Prepare(`INSERT OR IGNORE INTO tag (name) VALUES (?)`)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if _, err = tagStatement.ExecContext(ctx, tag.Name()); err != nil {
			return err
		}
	}
	
	return nil
}	

func (r *tagRepository) InsertNoteTags(ctx context.Context, dbCtx persistence.IDbContext, noteTags []*models.NoteTag) error {
	noteTagStmt, err := dbCtx.Prepare(`INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	for _, noteTag := range noteTags {
		if _, err = noteTagStmt.ExecContext(ctx, noteTag.NoteID(), noteTag.TagID()); err != nil {
			return err
		}
	}

	return nil
}

func (r *tagRepository) GetByNames(ctx context.Context, dbCtx persistence.IDbContext, names []string) ([]*models.Tag, error) {
	if len(names) == 0 {
		return []*models.Tag{}, nil
	}

	placeholders := make([]string, 0, len(names))
	args := make([]any, 0, len(names))

	for _, name := range names {
		placeholders = append(placeholders, "?")
		args = append(args, name)
	}

	query := "SELECT name FROM tag WHERE name IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := dbCtx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tags = append(tags, models.NewTag(name))
	}

	return tags, nil
}

func (r *tagRepository) DeleteNoteTags(ctx context.Context, dbCtx persistence.IDbContext, noteID string) error {
	query := "DELETE FROM note_tag WHERE note_id = ?"
	_, err := dbCtx.ExecContext(ctx, query, noteID)
	return err
}
