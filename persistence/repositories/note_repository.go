package repositories

import "database/sql"

type NoteRepository interface {
	Upsert(title string, path string, slug string) error
}

type noteRepository struct {
	db *sql.DB
}

