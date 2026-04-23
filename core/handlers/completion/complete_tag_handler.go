package completion

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type completeTagCommand struct {
	Query string `json:"query"`
}

type completeTagResult struct {
	Name string `json:"name"`
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

	panic("not implemented")
}
