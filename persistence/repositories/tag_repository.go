package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
)

type ITagRepository interface {
	Upsert(ctx context.Context, tags []*models.Tag) error
	UpsertNoteTags(noteID string, tagIDs []string) error
}

type tagRepository struct {
	Transaction *sql.Tx
}

func NewTagRepository(tx *sql.Tx) ITagRepository {
	return &tagRepository{Transaction: tx}
}

func (r *tagRepository) Upsert(ctx context.Context, tags []*models.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(tags))
	args := make([]any, 0, len(tags))

	for _, tag := range tags {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, tag.ID(), tag.Name())
	}

	query := "INSERT OR IGNORE INTO tags(id, name) VALUES " + strings.Join(placeholders, ",")

	_, err := r.Transaction.ExecContext(ctx, query, args...)
	return err
}

func (r *tagRepository) UpsertNoteTags(noteID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(tagIDs))
	args := make([]any, 0, len(tagIDs))

	for _, tagID := range tagIDs {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, noteID, tagID)
	}

	query := "INSERT OR IGNORE INTO note_tags(note_id, tag_id) VALUES " + strings.Join(placeholders, ",")

	_, err := r.Transaction.ExecContext(context.Background(), query, args...)
	if err != nil {
		return err
	}

	return nil
}
