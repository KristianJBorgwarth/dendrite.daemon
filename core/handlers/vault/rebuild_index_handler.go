package vault

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type rebuildIndexCommand struct {
}

type RebuildIndexHandler struct {
	idxRebuilder services.IIndexRebuilder
}

func NewRebuildIndexHandler(idxR services.IIndexRebuilder) *RebuildIndexHandler {
	return &RebuildIndexHandler{idxRebuilder: idxR}
}

func (h RebuildIndexHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd rebuildIndexCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	err := h.idxRebuilder.RebuildIndex(ctx, store.GetVaultStore().Config.VaultPath())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

