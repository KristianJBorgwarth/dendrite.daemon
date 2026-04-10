package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
)

type ITagRepository interface {
	Insert(ctx context.Context, tags []*models.Tag) error
	InsertNoteTags(ctx context.Context, noteID string, tagIDs []string) error
	GetByNames(ctx context.Context, names []string) ([]*models.Tag, error)
}

type tagRepository struct {
	Transaction *sql.Tx
}

func NewTagRepository(tx *sql.Tx) ITagRepository {
	return &tagRepository{Transaction: tx}
}

func (r *tagRepository) Insert(ctx context.Context, tags []*models.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(tags))
	args := make([]any, 0, len(tags))

	for _, tag := range tags {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, tag.ID(), tag.Name())
	}

	query := "INSERT OR IGNORE INTO tag(id, name) VALUES " + strings.Join(placeholders, ",")

	_, err := r.Transaction.ExecContext(ctx, query, args...)
	return err
}

func (r *tagRepository) InsertNoteTags(ctx context.Context ,noteID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(tagIDs))
	args := make([]any, 0, len(tagIDs))

	for _, tagID := range tagIDs {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, noteID, tagID)
	}

	query := "INSERT OR IGNORE INTO note_tag(note_id, tag_id) VALUES " + strings.Join(placeholders, ",")

	_, err := r.Transaction.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *tagRepository) GetByNames(ctx context.Context, names []string) ([]*models.Tag, error) {
	if len(names) == 0 {
		return []*models.Tag{}, nil
	}

	placeholders := make([]string, 0, len(names))
	args := make([]any, 0, len(names))

	for _, name := range names {
		placeholders = append(placeholders, "?")
		args = append(args, name)
	}

	query := "SELECT id, name FROM tag WHERE name IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := r.Transaction.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		var id string
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		tags = append(tags, models.NewTag(id, name))
	}

	return tags, nil
}
