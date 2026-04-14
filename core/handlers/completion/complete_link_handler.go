package completion

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type completeLinkCommand struct {
	Query string `json:"query"`
}

type completeLinkResult struct {
	Slug    string `json:"slug"`
	Display string `json:"display"`
}

type CompleteLinkHandler struct {
	linkRepo repositories.ILinkRepository
}

func NewCompleteLinkHandler(lr repositories.ILinkRepository) *CompleteLinkHandler {
	return &CompleteLinkHandler{linkRepo: lr}
}

func (h *CompleteLinkHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd completeLinkCommand
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}
	links, err := h.linkRepo.Search(ctx, cmd.Query)
	if err != nil {
		slog.Error("Failed to search links", "error", err)
		return nil, err
	}

	items := make([]completeLinkResult, len(links))
	for i, link := range links {
		items[i] = completeLinkResult{Slug: link.TargetSlug(), Display: link.Display()}
	}

	return items, nil
}
