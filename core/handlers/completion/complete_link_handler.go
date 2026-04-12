package completion

import (
	"context"
	"encoding/json"
	"log/slog"
)

type completeLinkCommand struct {
	Query string `json:"query"`
}

type completionItem struct {
	Slug string `json:"slug"`
}

type CompleteLinkHandler struct{}

func NewCompleteLinkHandler() *CompleteLinkHandler {
	return &CompleteLinkHandler{}
}

func (h *CompleteLinkHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd completeLinkCommand
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}
	slog.Debug("handling complete link command", "query", cmd.Query)

	// this is just a stub and the client is expecting the title/display too
	return []completionItem{
		{Slug: "standard-streams"},
		{Slug: "unit-of-work"},
		{Slug: "treesitter-basics"},
	}, nil
}
