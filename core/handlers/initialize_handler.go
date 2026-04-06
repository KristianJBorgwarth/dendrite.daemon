package handlers

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type initializeCommand struct {
	VaultPath         string `json:"vaultPath"`
	TemplateDirectory string `json:"templateDirectory"`
	ScratchDirectory  string `json:"scratchDirectory"`
	ScratchTemplate   string `json:"scratchTemplate"`
	DailyDirectory    string `json:"dailyDirectory"`
	DailyTemplate     string `json:"dailyTemplate"`
}

type InitializeHandler struct{}

func NewInitializeHandler() *InitializeHandler {
	return &InitializeHandler{}
}

func (h InitializeHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd initializeCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}



	err := persistence.InitializeDBContext(cmd.VaultPath)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
