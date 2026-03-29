package handlers

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
)

type Handler interface {
	Handle(ctx context.Context, raw json.RawMessage) (any, *rpc.Error)
}
