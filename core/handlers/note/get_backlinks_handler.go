package note

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type getBacklinksCommand struct {
	Path string `json:"path"`
}

type GetBackLinksHandler struct {
	linkRepo repositories.ILinkRepository
	noteRepo repositories.INoteRepository
}

func NewGetBackLinksHandler(lr repositories.ILinkRepository, nr repositories.INoteRepository) *GetBackLinksHandler {
	return &GetBackLinksHandler{linkRepo: lr, noteRepo: nr}
}

func (h *GetBackLinksHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd getBacklinksCommand
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	note, err := h.noteRepo.GetByPath(ctx, cmd.Path)
	if err != nil {
		return nil, err
	}

	links, err := h.linkRepo.GetBacklinks(ctx, note.Slug())
	if err != nil {
		return nil, err
	}

	return links, nil
}
