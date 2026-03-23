package handlers

import (
	"encoding/json"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type initializeCommand struct {
	VaultPath   string `json:"vaultPath"`
}

type InitializeHandler struct{}

func (h InitializeHandler) Handle(raw json.RawMessage) (any, *rpc.Error) {
	var params initializeCommand

	if err := json.Unmarshal(raw, &params); err != nil {
		return nil, &rpc.Error{
			Code:    -32602,
			Message: "invalid params: " + err.Error(),
		}
	}


	err := persistence.InitializeIndex(params.VaultPath)
	if err != nil {
		return nil, &rpc.Error{
			Code:    -1,
			Message: "failed to initialize index: " + err.Error(),
		}
	}

	return nil, nil
}

