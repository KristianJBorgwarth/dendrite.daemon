package note

import (
	"context"
	"encoding/json"
	"log/slog"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type saveNoteCommand struct {
	Path string `json:"path"`
}

type SaveNoteHandler struct {
	uow         *repositories.UnitOfWork
	noteRepo    repositories.INoteRepository
	tagService  services.ITagService
	noteService services.INoteService
	linkService services.ILinkService
}

func NewSaveNoteHandler(
	uow *repositories.UnitOfWork,
	nr repositories.INoteRepository,
	ts services.ITagService,
	ns services.INoteService,
	ls services.ILinkService,
) *SaveNoteHandler {
	return &SaveNoteHandler{uow, nr, ts, ns, ls}
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

	slog.Debug("Parsed file", "path", cmd.Path, "slug", file.Slug, "title", file.Title, "tags", file.FrontMatter.Tags, "links", file.ExtractedLinks)

	tx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}
	defer h.uow.Rollback()

	note, err := h.noteRepo.GetBySlug(ctx, tx, file.Slug)
	if err != nil {
		slog.Debug("Failed to get note by slug", "slug", file.Slug, "error", err)
		return nil, err
	}

	if note == nil {
		note, err = h.noteService.CreateNote(ctx, tx, file.Path, file.Title, file.Slug)
		if err != nil {
			slog.Debug("Failed to save new note", "slug", file.Slug, "error", err)
			return nil, err
		}
	} else {
		if err = h.noteService.DeleteNoteMetaData(ctx, tx, note.ID()); err != nil {
			slog.Debug("Failed to delete existing note metadata", "noteID", note.ID(), "error", err)
			return nil, err
		}
		if err = h.noteService.UpdateNote(ctx, tx, note.ID(), file.Path, file.Title, file.Slug); err != nil {
			slog.Debug("Failed to update existing note", "noteID", note.ID(), "error", err)
			return nil, err
		}
	}

	tagModels, err := h.tagService.CreateTags(ctx, tx, file.FrontMatter.Tags)
	if err != nil {
		slog.Debug("Failed to create tags", "tags", file.FrontMatter.Tags, "error", err)
		return nil, err
	}

	if err = h.tagService.CreateNoteTags(ctx, tx, note.ID(), tagModels); err != nil {
		slog.Debug("Failed to create note tags", "noteID", note.ID(), "tags", file.FrontMatter.Tags, "error", err)
		return nil, err
	}

	if err = h.linkService.CreateLinks(ctx, tx, note.ID(), file.ExtractedLinks); err != nil {
		slog.Debug("Failed to create links", "noteID", note.ID(), "links", file.ExtractedLinks, "error", err)
		return nil, err
	}

	if err := h.uow.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}
