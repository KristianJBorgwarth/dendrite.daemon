package handlers

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type initializeCommand struct {
	VaultPath         string `json:"vaultPath"`
	TemplateDirectory string `json:"templateDirectory"`
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

	store := store.NewVaultStore(cmd.VaultPath, cmd.TemplateDirectory)

	err := persistence.InitializeDBContext(store.Config.VaultPath())
	if err != nil {
		return nil, err
	}

	return nil, nil
}
