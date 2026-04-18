package vault

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type initializeCommand struct {
	VaultName         string `json:"vaultName"`
	VaultPath         string `json:"vaultPath"`
	TemplateDirectory string `json:"templateDirectory"`
}

type InitializeHandler struct {
	idxRebuilder services.IIndexRebuilder
}

func NewInitializeHandler(idxR services.IIndexRebuilder) *InitializeHandler {
	return &InitializeHandler{idxRebuilder: idxR}
}

func (h InitializeHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd initializeCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	store := store.NewVaultStore(cmd.VaultName, cmd.VaultPath, cmd.TemplateDirectory)

	err := persistence.InitializeDBContext(store.Config.VaultName())
	if err != nil {
		return nil, err
	}

	if err = h.idxRebuilder.RebuildIndex(ctx, cmd.VaultPath); err != nil {
		return nil, err
	}

	return nil, nil
}
