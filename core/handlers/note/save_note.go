package note

import (
	"context"
	"encoding/json"
	"log/slog"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type saveNoteCommand struct {
	Path string `json:"path"`
}

type SaveNoteHandler struct {
	uow *repositories.UnitOfWork
}

func NewSaveNoteHandler() *SaveNoteHandler {
	return &SaveNoteHandler{repositories.NewUnitOfWork()}
}

func (h *SaveNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd saveNoteCommand
	slog.Debug("Handling SaveNoteCommand", "raw", string(raw))

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	file, err := filehandling.ReadFile(cmd.Path)
	if err != nil {
		slog.Debug("Failed to read file", "path", cmd.Path, "error", err)
		return nil, err
	}

	tx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer h.uow.Rollback()
	noteRepo := repositories.NewNoteRepository(tx)
	linkRepo := repositories.NewLinkRepository(tx)
	tagRepo := repositories.NewTagRepository(tx)

	note, err := noteRepo.GetBySlug(ctx, file.Slug)
	if err != nil {
		slog.Debug("Failed to get note by slug", "slug", file.Slug, "error", err)
		return nil, err
	}

	if note == nil {
		if err := h.handleNewNote(ctx, noteRepo, linkRepo, tagRepo, file); err != nil {
			return nil, err
		}
	} else {
		if err := h.handleExistingNote(ctx, linkRepo, tagRepo, note, file); err != nil {
			return nil, err
		}
	}

	if err := h.uow.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *SaveNoteHandler) handleNewNote(
	ctx context.Context,
	noteRepo repositories.NoteRepository,
	linkRepo repositories.ILinkRepository,
	tagRepo repositories.ITagRepository,
	file *filehandling.File,
) error {
	note := models.CreateNote(file.Path, file.Title, file.Slug)
	slog.Debug("EXTRACTED FILE", "path", file.Path, "title", file.Title, "slug", file.Slug, "links", file.ExtractedLinks, "tags", file.FrontMatter.Tags)

	if err := noteRepo.Insert(ctx, note); err != nil {
		slog.Debug("Failed to insert new note", "noteID", note.ID(), "error", err)
		return err
	}

	links := models.MapToLinkModel(note.ID(), file.ExtractedLinks)

	if err := linkRepo.Insert(ctx, links); err != nil {
		slog.Debug("Failed to insert links for new note", "noteID", note.ID(), "error", err)
		return err
	}

	if err := tagRepo.InsertNoteTags(ctx, note.ID(), file.FrontMatter.Tags); err != nil {
		slog.Debug("Failed to insert tags for new note", "noteID", note.ID(), "error", err)
		return err
	}

	return nil
}

func (h *SaveNoteHandler) handleExistingNote(ctx context.Context,
	linkRepo repositories.ILinkRepository,
	tagRepo repositories.ITagRepository,
	note *models.Note,
	file *filehandling.File,
) error {
	err := h.deleteExistingNoteRelations(ctx, linkRepo, tagRepo, note.ID())
	if err != nil {
		slog.Debug("Failed to delete existing note relations", "noteID", note.ID(), "error", err)
		return err
	}

	links := models.MapToLinkModel(note.ID(), file.ExtractedLinks)

	if err = linkRepo.Insert(ctx, links); err != nil {
		slog.Debug("Failed to insert links for existing note", "noteID", note.ID(), "error", err)
		return err
	}

	if err = tagRepo.InsertNoteTags(ctx, note.ID(), file.FrontMatter.Tags); err != nil {
		slog.Debug("Failed to insert tags for existing note", "noteID", note.ID(), "error", err)
		return err
	}

	return nil
}

func (h *SaveNoteHandler) deleteExistingNoteRelations(ctx context.Context, linkRepo repositories.ILinkRepository, tagRepo repositories.ITagRepository, noteID string) error {
	if err := linkRepo.Delete(ctx, noteID); err != nil {
		return err
	}
	if err := tagRepo.DeleteNoteTags(ctx, noteID); err != nil {
		return err
	}
	return nil
}

func (h *SaveNoteHandler) handleTags(ctx context.Context, tagRepo repositories.ITagRepository, noteID string, tags []string) error {
	panic("not implemented")
}
