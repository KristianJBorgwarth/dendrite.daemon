package services

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type INoteService interface {
	CreateNote(ctx context.Context, dbCtx persistence.IDbContext, path, title, slug string) (*models.Note, error)
	DeleteNoteMetaData(ctx context.Context, dbCtx persistence.IDbContext, noteID string) error
}

type noteService struct {
	tagRepo  repositories.ITagRepository
	linkRepo repositories.ILinkRepository
	noteRepo repositories.INoteRepository
}

func NewNoteService(
	tagRepo repositories.ITagRepository,
	linkRepo repositories.ILinkRepository,
	noteRepo repositories.INoteRepository,
) INoteService {
	return &noteService{tagRepo, linkRepo, noteRepo}
}

func (s *noteService) CreateNote(ctx context.Context, dbCtx persistence.IDbContext, path, title, slug string) (*models.Note, error) {
	note:= models.CreateNote(path, title, slug)
	if err := s.noteRepo.Insert(ctx, dbCtx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *noteService) DeleteNoteMetaData(ctx context.Context, dbCtx persistence.IDbContext, noteID string) error {
	if err := s.linkRepo.Delete(ctx, dbCtx, noteID); err != nil {
		return err
	}

	if err := s.tagRepo.DeleteNoteTags(ctx, dbCtx, noteID); err != nil {
		return err
	}

	return nil
}
