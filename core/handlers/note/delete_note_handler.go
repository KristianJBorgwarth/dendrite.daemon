package note

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type deleteNoteCommand struct {
	Path string `json:"path"`
}

type DeleteNoteHandler struct {
	uow *repositories.UnitOfWork
	tagService  services.ITagService
	noteService services.INoteService
}

func NewDeleteNoteHandler(
	uow *repositories.UnitOfWork,
	tagService services.ITagService,
	noteService services.INoteService,
) *DeleteNoteHandler {
	return &DeleteNoteHandler{uow, tagService, noteService}
}

func (h *DeleteNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd deleteNoteCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	tx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}
	defer h.uow.Rollback()

	note, err := h.noteService.GetByPath(ctx, tx, cmd.Path)
	if err != nil {
		return nil, err
	}

	err = h.noteService.DeleteNote(ctx, tx, note.ID())
	if err != nil {
		return nil, err
	}

	if h.uow.Commit() != nil {
		return nil, err
	}

	return nil, nil
}
