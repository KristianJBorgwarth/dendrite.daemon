package handlers

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type saveNoteCommand struct {
	Title     string   `json:"title"`
	Path      string   `json:"path"`
	Directory string   `json:"directory"`
	Tags      []string `json:"tags"`
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

	return nil, nil
}
