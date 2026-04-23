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
	noteService  services.INoteService
}

func NewInitializeHandler(idxR services.IIndexRebuilder, nsv services.INoteService) *InitializeHandler {
	return &InitializeHandler{idxRebuilder: idxR, noteService: nsv}
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

	noteCount, err := h.noteService.GetNoteCount(ctx)
	if err != nil {
		return nil, err
	}

	if noteCount > 0 {
		return nil, nil
	}

	err = h.idxRebuilder.RebuildIndex(ctx, cmd.VaultPath)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
