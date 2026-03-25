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
