package note

import (
	"context"
	"encoding/json"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type saveNoteCommand struct {
	Path      string   `json:"path"`
}

type SaveNoteHandler struct {
	uow *repositories.UnitOfWork
}

func NewSaveNoteHandler() *SaveNoteHandler {
	return &SaveNoteHandler{repositories.NewUnitOfWork()}
}

func (h *SaveNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd saveNoteCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	file, err := filehandling.ReadFile(cmd.Path)
	if err != nil {
		return nil, err
	}

	tx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer h.uow.Rollback()
	tagRepo := repositories.NewTagRepository(tx)
	noteRepo := repositories.NewNoteRepository(tx)

	
	note, err := noteRepo.GetBySlug(ctx, file.Slug)
	if err != nil {
		return nil, err
	}



	return nil, nil
}
