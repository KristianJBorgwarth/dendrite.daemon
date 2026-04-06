package handlers

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/frontmatter"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/template"
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

	slug := frontmatter.Slugify(cmd.Title)

	templatePath := store.GetVaultStore().GetTemplatePath(cmd.TemplateName)

	data, err := template.RenderTemplate(templatePath, cmd.Title, slug)
	if err != nil {
		return nil, err
	}

	tags, err := frontmatter.ParseTags(data)
	if err != nil {
		return nil, err
	}

	tx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer h.uow.Rollback()

	tagRepo := repositories.NewTagRepository(tx)
	noteRepo := repositories.NewNoteRepository(tx)

	dbTags, err := tagRepo.GetByNames(ctx, tags)
	if err != nil {
		return nil, err
	}

	newTags := utils.Filter(tags, func(name string) bool {
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

	if err = tagRepo.Upsert(ctx, tagModels); err != nil {
		return nil, err
	}

	tagModels = append(tagModels, dbTags...)

	note := models.CreateNote(cmd.Directory, cmd.Title, slug)

	if err = noteRepo.Upsert(ctx, note); err != nil {
		return nil, err
	}

	if err = tagRepo.UpsertNoteTags(ctx, note.ID(), utils.Select(tagModels, func(t *models.Tag) string { return t.ID() })); err != nil {
		return nil, err
	}

	h.uow.FileStore.Stage(cmd.Directory, data)

	if err = h.uow.Commit(); err != nil {
		return nil, err
	}

	return cmd.Directory, nil
}
