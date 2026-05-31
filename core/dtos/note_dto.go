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

func NewNoteDtos(notes []*models.Note) []*NoteDto {
	dtos := make([]*NoteDto, len(notes))
	for i, note := range notes {
		dtos[i] = NewNoteDto(note)
	}
	return dtos
}
