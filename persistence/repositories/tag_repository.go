
package repositories

import "database/sql"

type TagRepository interface {
	Upsert(title string, path string) error
}

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Upsert(title string, path string) error {
	query := `
	INSERT INTO tags (title, path)
	VALUES ($1, $2)
	ON CONFLICT (slug) DO UPDATE
	SET title = EXCLUDED.title,
	    path = EXCLUDED.path;
	`
	_, err := r.db.Exec(query, title, path)
	return err
}
