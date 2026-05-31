package services

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type ICfService interface {
	AddCfe(ctx context.Context, dbCtx persistence.IDbContext, noteID string, cfe map[string]any) error
}

type cfService struct {
	cfRepo repositories.ICfRepository
}

func NewCfeService(cfRepo repositories.ICfRepository) ICfService {
	return &cfService{cfRepo}
}

func (s *cfService) AddCfe(
	ctx context.Context,
	dbCtx persistence.IDbContext,
	noteID string,
	cfe map[string]any,
) error {
	if len(cfe) == 0 {
		return nil
	}
	cfModels, err := models.MapToCustomFrontmatter(noteID, cfe)
	if err != nil {
		return err
	}
	return s.cfRepo.InsertRange(ctx, dbCtx, cfModels)
}
