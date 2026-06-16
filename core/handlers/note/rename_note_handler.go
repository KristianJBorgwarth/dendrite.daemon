package note

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type renameNoteCommand struct {
	Path    string `json:"path"`
	NewPath string `json:"newPath"`
}

type RenameNoteHandler struct {
	uow         *repositories.UnitOfWork
	tagService  services.ITagService
	noteService services.INoteService
}

func (h *RenameNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd renameNoteCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	dbCtx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer h.uow.Rollback()

	if err := h.noteService.RenameNote(ctx, dbCtx, cmd.Path, cmd.NewPath); err != nil {
		return nil, err
	}

	if err := h.uow.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
