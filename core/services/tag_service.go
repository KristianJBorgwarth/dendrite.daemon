package services

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/utils"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type ITagService interface {
	CreateTags(ctx context.Context, dbCtx persistence.IDbContext, names []string) ([]*models.Tag, error)
	CreateNoteTags(ctx context.Context, dbCtx persistence.IDbContext, noteID string, tags []*models.Tag) error
}

type tagService struct {
	tagRepo repositories.ITagRepository
}

func NewTagService(tagRepo repositories.ITagRepository) ITagService {
	return &tagService{tagRepo}
}

func (s *tagService) CreateTags(ctx context.Context, dbCtx persistence.IDbContext, names []string) ([]*models.Tag, error) {
	existingTags, err := s.tagRepo.GetByNames(ctx, dbCtx, names)
	if err != nil {
		return nil, err
	}

	newTags := utils.Filter(names, func(name string) bool {
		for _, tag := range existingTags {
			if tag.Name() == name {
				return false
			}
		}
		return true
	})

	newTagModels := models.CreateTags(newTags)

	if err = s.tagRepo.Insert(ctx, dbCtx,  newTagModels); err != nil {
		return nil, err
	}

	tags := append(existingTags, newTagModels...)

	return tags, nil
}

func(s *tagService) CreateNoteTags(ctx context.Context, dbCtx persistence.IDbContext, noteID string, tags []*models.Tag) error {
	tagIds := utils.Select(tags, func(t *models.Tag) string {
		return t.ID()
	});

	return s.tagRepo.InsertNoteTags(ctx, dbCtx, noteID, tagIds)
}
