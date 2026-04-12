package note

import (
	"context"
	"encoding/json"
	"path/filepath"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/utils"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type createNoteCommand struct {
	Title        string `json:"title"`
	TemplateName string `json:"templateName"`
	Directory    string `json:"directory"`
}

type CreateNoteHandler struct {
	uow *repositories.UnitOfWork
}

func NewCreateNoteHandler() *CreateNoteHandler {
	return &CreateNoteHandler{repositories.NewUnitOfWork()}
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

	tagRepo := repositories.NewTagRepository(dbCtx)
	noteRepo := repositories.NewNoteRepository()

	dbTags, err := tagRepo.GetByNames(ctx, template.FrontMatter.Tags)
	if err != nil {
		return nil, err
	}

	newTags := utils.Filter(template.FrontMatter.Tags, func(name string) bool {
		for _, t := range dbTags {
			if t.Name() == name {
				return false
			}
		}
		return true
	})

	tagModels, err := models.CreateTags(newTags)
	if err != nil {
		return nil, err
	}

	if err = tagRepo.Insert(ctx, tagModels); err != nil {
		return nil, err
	}

	tagModels = append(tagModels, dbTags...)

	note := models.CreateNote(notePath, cmd.Title, template.Slug)

	if err = noteRepo.Insert(ctx, dbCtx, note); err != nil {
		return nil, err
	}

	if err = tagRepo.InsertNoteTags(ctx, note.ID(), utils.Select(tagModels, func(t *models.Tag) string { return t.ID() })); err != nil {
		return nil, err
	}

	h.uow.FileStore.Stage(notePath, template.Content)

	if err = h.uow.Commit(); err != nil {
		return nil, err
	}

	return notePath, nil
}
