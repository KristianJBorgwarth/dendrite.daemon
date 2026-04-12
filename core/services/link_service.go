package services

import (
	"context"
	"log/slog"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type ILinkService interface {
	CreateLinks(ctx context.Context, dbCtx persistence.IDbContext, noteID string, links []*filehandling.ExtractedLink) error
}

type linkService struct {
	linkRepo repositories.ILinkRepository
}

func NewLinkService(linkRepo repositories.ILinkRepository) ILinkService {
	return &linkService{linkRepo}
}

func (s *linkService) CreateLinks(ctx context.Context, dbCtx persistence.IDbContext, noteID string, links []*filehandling.ExtractedLink) error {
	if len(links) == 0 {
		return nil
	}

	linkModels := models.MapToLinkModel(noteID, links)
	slog.Debug("Creating links", "noteID", noteID, "links", linkModels)
	return s.linkRepo.Insert(ctx, dbCtx, linkModels)
}

