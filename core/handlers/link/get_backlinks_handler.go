package handlers

import (
	"context"
	"encoding/json"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type getBacklinksCommand struct {
	Path string `json:"path"`
}

type GetBackLinksHandler struct {
	linkRepo repositories.ILinkRepository
}

func NewGetBackLinksHandler(lr repositories.ILinkRepository) *GetBackLinksHandler {
	return &GetBackLinksHandler{linkRepo: lr}
}

func (h *GetBackLinksHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd getBacklinksCommand
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	slug := filehandling.Slugify(cmd.Path)

	links, err := h.linkRepo.GetBacklinks(ctx, slug)
	if err != nil {
		return nil, err
	}

	return links, nil
}
