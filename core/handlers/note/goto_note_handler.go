package note

import "github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"

type gotoNoteCommand struct {
	Slug string `json:"slug"`
}

type gotoNoteResult struct {
	Path string `json:"path"`
}

type GotoNoteHandler struct{
	noteRepo repositories.INoteRepository
}

func(h *GotoNoteHandler) Handle(command gotoNoteCommand) (gotoNoteResult, error) {
	note, err := h.noteRepo.GetBySlug(command.Slug)
	if err != nil {
		return gotoNoteResult{}, err
	}
}
