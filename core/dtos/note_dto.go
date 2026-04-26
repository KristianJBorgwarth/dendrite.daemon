package dtos

import "github.com/KristianJBorgwarth/dendrite.daemon/core/models"

type NoteDto struct {
	ID    string `json:"id"`
	Path  string `json:"path"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

func NewNoteDto(note *models.Note) *NoteDto {
	return &NoteDto{
		ID:    note.ID(),
		Path:  note.Path(),
		Title: note.Title(),
		Slug:  note.Slug(),
	}
}
