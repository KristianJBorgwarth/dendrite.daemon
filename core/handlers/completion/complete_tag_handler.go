package completion

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type completeTagCommand struct {
	Query string `json:"query"`
}

type CompleteTagHandler struct {
	tagRepo repositories.ITagRepository
}

func NewCompleteTagHandler(tr repositories.ITagRepository) *CompleteTagHandler {
	return &CompleteTagHandler{tagRepo: tr}
}

func (h *CompleteTagHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd completeTagCommand
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	tags, err := h.tagRepo.GetByName(ctx, cmd.Query)
	slog.Debug("Got tags", "query", cmd.Query, "count", len(tags))
	if err != nil {
		return nil, err
	}

	results := make([]string, len(tags))
	for i, tag := range tags {
		results[i] = tag.Name()
	}

	return results, nil
}
