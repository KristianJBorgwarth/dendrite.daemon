
package repositories

import "database/sql"

type TagRepository interface {
	Upsert(title string, path string, slug string) error
}

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepository{db: db}
}

