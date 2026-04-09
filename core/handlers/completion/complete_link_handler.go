package completion

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/utils"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type completeLinkCommand struct {
	LinkQuery string `json:"linkQuery"`
}

type CompleteLinkHandler struct{
	uow *repositories.UnitOfWork
}

func NewCompleteLinkHandler() *CompleteLinkHandler {
	return &CompleteLinkHandler{repositories.NewUnitOfWork()}
}

func (h *CompleteLinkHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd completeLinkCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	tx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer h.uow.Rollback()
	linkRepo := repositories.NewLinkRepository(tx)

	links, err := linkRepo.Search(ctx, cmd.LinkQuery)
	if err != nil {
		return nil, err
	}

	slugs := utils.Select(links, func(l *models.Link) string {return l.TargetSlug()});

	return slugs, nil
}
