package services

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type ICfeService interface {
	AddCfe(ctx context.Context, dbCtx persistence.IDbContext, noteID string, cfe map[string]any) error
}

type cfeService struct {
	cfeRepo repositories.ICfeRepository
}

func NewCfeService(cfeRepo repositories.ICfeRepository) ICfeService {
	return &cfeService{cfeRepo}
}

func (s *cfeService) AddCfe(
	ctx context.Context,
	dbCtx persistence.IDbContext,
	noteID string,
	cfe map[string]any,
) error {
	if len(cfe) == 0 {
		return nil
	}
	cfeModels := models.MapToCfe(noteID, cfe)
	return s.cfeRepo.InsertRange(ctx, dbCtx, cfeModels)
}
