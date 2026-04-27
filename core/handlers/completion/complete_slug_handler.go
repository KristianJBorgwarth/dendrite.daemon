package completion

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type completeSlugCommand struct {
	Query string `json:"query"`
}

type CompleteSlugHandler struct {
	noteRepo repositories.INoteRepository
}

func NewCompleteSlugHandler(nr repositories.INoteRepository) *CompleteSlugHandler {
	return &CompleteSlugHandler{noteRepo: nr}
}

func (h *CompleteSlugHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd completeSlugCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	notes, err := h.noteRepo.Search(ctx, cmd.Query)
	if err != nil {
		return nil, err
	}

	slugs := make([]string, len(notes))
	for i, note := range notes {
		slugs[i] = note.Slug()
	}

	return slugs, nil
}


