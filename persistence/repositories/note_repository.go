package repositories

import "database/sql"

type NoteRepository interface {
	Upsert(title string, path string, slug string) error
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) NoteRepository {
	return &noteRepository{db: db}
}


func (r *noteRepository) Upsert(title string, path string, slug string) error {
	query := `
	INSERT INTO notes (title, path, slug)
	VALUES ($1, $2, $3)
	ON CONFLICT (slug) DO UPDATE
	SET title = EXCLUDED.title,
	    path = EXCLUDED.path;
	`
	_, err := r.db.Exec(query, title, path, slug)
	return err
}
