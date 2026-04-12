package note

import (
	"context"
	"encoding/json"
	"path/filepath"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type createNoteCommand struct {
	Title        string `json:"title"`
	TemplateName string `json:"templateName"`
	Directory    string `json:"directory"`
}

type CreateNoteHandler struct {
	uow        *repositories.UnitOfWork
	tagService services.ITagService
	noteRepo   repositories.NoteRepository
}

func NewCreateNoteHandler(
	uow *repositories.UnitOfWork,
	tagRepo services.ITagService,
	noteRepo repositories.NoteRepository,
) *CreateNoteHandler {
	return &CreateNoteHandler{uow, tagRepo, noteRepo}
}

func (h *CreateNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd createNoteCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	template, err := filehandling.NewTemplate(cmd.TemplateName, cmd.Title)
	if err != nil {
		return nil, err
	}

	notePath := filepath.Join(store.GetVaultStore().Config.VaultPath(), cmd.Directory, template.Slug+".md")

	dbCtx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer h.uow.Rollback()

	tagModels, err := h.tagService.CreateTags(ctx, dbCtx, template.FrontMatter.Tags)
	if err != nil {
		return nil, err
	}

	note := models.CreateNote(notePath, cmd.Title, template.Slug)

	if err = h.noteRepo.Insert(ctx, dbCtx, note); err != nil {
		return nil, err
	}

	if err = h.tagService.CreateNoteTags(ctx, dbCtx, note.ID(), tagModels); err != nil {
		return nil, err
	}

	h.uow.FileStore.Stage(notePath, template.Content)

	if err = h.uow.Commit(); err != nil {
		return nil, err
	}

	return notePath, nil
}
