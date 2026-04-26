package note

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/dtos"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)
type getNotesByTagCommand struct {
	Tag string `json:"tag"`
}

type GetNotesByTagHandler struct {
	noteRepo repositories.INoteRepository
}

func NewGetNotesByTagHandler(noteRepo repositories.INoteRepository) *GetNotesByTagHandler {
	return &GetNotesByTagHandler{noteRepo: noteRepo}
}

func (h *GetNotesByTagHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd getNotesByTagCommand
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	notes, err := h.noteRepo.GetByTag(ctx, cmd.Tag)
	if err != nil {
		return nil, err
	}

	noteDtos := make([]*dtos.NoteDto, len(notes))
	for i, note := range notes {
		noteDtos[i] = dtos.NewNoteDto(note)
	}

	return noteDtos, nil
}
