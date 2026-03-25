package repositories

import (
	"database/sql"
	"strings"
)

type TagRepository interface {
	Upsert(names []string) error
}

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Upsert(names []string) error {
	if len(names) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	placeholders := make([]string, 0, len(names))
	args := make([]any, 0, len(names))

	for _, name := range names {
		placeholders = append(placeholders, "(?)")
		args = append(args, name)
	}

	query := "WITH input(name) AS (VALUES " +
		strings.Join(placeholders, ",") +
		") INSERT INTO tags(name) " +
		"SELECT name FROM input " +
		"ON CONFLICT(name) DO NOTHING;"

	if _, err := tx.Exec(query, args...); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *tagRepository) UpsertNoteTags(noteID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	placeholders := make([]string, 0, len(tagIDs))
	args := make([]any, 0, len(tagIDs))

	for _, tagID := range tagIDs {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, noteID, tagID)
	}

	query := "WITH input(note_id, tag_id) AS (VALUES " +
		strings.Join(placeholders, ",") +
		") INSERT INTO note_tags(note_id, tag_id) " +
		"SELECT note_id, tag_id FROM input " +
		"ON CONFLICT(note_id, tag_id) DO NOTHING" +
		"SELECT note_id, tag_id FROM input;"

	if _, err := tx.Exec(query, args...); err != nil {
		return err
	}

	return tx.Commit()
}
