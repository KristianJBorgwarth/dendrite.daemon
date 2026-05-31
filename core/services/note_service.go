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
	DeleteNote(ctx context.Context, dbCtx persistence.IDbContext, path string) error
	GetByPath(ctx context.Context, dbCtx persistence.IDbContext, path string) (*models.Note, error)
	UpdateNote(ctx context.Context, dbCtx persistence.IDbContext, noteID, path, title, slug string) error
	GetNoteCount(ctx context.Context) (int, error)
}

type noteService struct {
	tagRepo  repositories.ITagRepository
	linkRepo repositories.ILinkRepository
	noteRepo repositories.INoteRepository
	cfeRepo  repositories.ICfRepository
}

func NewNoteService(
	tagRepo repositories.ITagRepository,
	linkRepo repositories.ILinkRepository,
	noteRepo repositories.INoteRepository,
	cfeRepo repositories.ICfRepository,
) INoteService {
	return &noteService{tagRepo, linkRepo, noteRepo, cfeRepo}
}

func (s *noteService) CreateNote(ctx context.Context, dbCtx persistence.IDbContext, path, title, slug string) (*models.Note, error) {
	note := models.CreateNote(path, title, slug)
	if err := s.noteRepo.Insert(ctx, dbCtx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *noteService) UpdateNote(ctx context.Context, dbCtx persistence.IDbContext, noteID, path, title, slug string) error {
	if err := s.noteRepo.Update(ctx, dbCtx, noteID, path, title, slug); err != nil {
		return err
	}
	return nil
}

func (s *noteService) DeleteNoteMetaData(ctx context.Context, dbCtx persistence.IDbContext, noteID string) error {
	if err := s.linkRepo.Delete(ctx, dbCtx, noteID); err != nil {
		return err
	}

	if err := s.tagRepo.DeleteNoteTags(ctx, dbCtx, noteID); err != nil {
		return err
	}

	if err := s.cfeRepo.Delete(ctx, dbCtx, noteID); err != nil {
		return err
	}

	return nil
}

func (s *noteService) DeleteNote(ctx context.Context, dbCtx persistence.IDbContext, noteID string) error {
	if err := s.DeleteNoteMetaData(ctx, dbCtx, noteID); err != nil {
		return err
	}
	
	if err := s.noteRepo.Delete(ctx, dbCtx, noteID); err != nil {
		return err
	}

	return nil
}

func (s *noteService) GetNoteCount(ctx context.Context) (int, error) {
	return s.noteRepo.GetNoteCount(ctx)
}

func (s *noteService) GetByPath(ctx context.Context, dbCtx persistence.IDbContext, path string) (*models.Note, error) {
	return s.noteRepo.GetByPath(ctx,path)
}
