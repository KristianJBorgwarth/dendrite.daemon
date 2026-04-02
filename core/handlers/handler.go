package handlers

import (
	"context"
	"encoding/json"
)

type Handler interface {
	Handle(ctx context.Context, raw json.RawMessage) (any, error)
}
