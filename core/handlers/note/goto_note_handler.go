package note

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type gotoNoteCommand struct {
	Slug string `json:"slug"`
}

type gotoNoteResult struct {
	Path string `json:"path"`
}

type GotoNoteHandler struct {
	noteRepo repositories.INoteRepository
}

func (h *GotoNoteHandler) Handle(ctx context.Context, command gotoNoteCommand) (gotoNoteResult, error) {
	note, err := h.noteRepo.GetBySlug(ctx, command.Slug)
	if err != nil {
		return gotoNoteResult{}, err
	}
	return gotoNoteResult{Path: note.Path()}, nil
}
