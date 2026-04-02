package handlers

import (
	"context"
	"encoding/json"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type initializeCommand struct {
	VaultPath string `json:"vaultPath"`
}

type InitializeHandler struct{}

func (h InitializeHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var params initializeCommand

	if err := json.Unmarshal(raw, &params); err != nil {
		return nil, err
	}

	err := persistence.InitializeIndex(params.VaultPath)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
