package handlers

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/frontmatter"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/template"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/google/uuid"
)

type createNoteCommand struct {
	Title        string            `json:"title"`
	TemplatePath string            `json:"templatePath"`
	Path         string            `json:"path"`
	Vars         map[string]string `json:"vars"`
}

type CreateNoteHandler struct {
	uow *repositories.UnitOfWork
}

func NewCreateNoteHandler(uow *repositories.UnitOfWork) *CreateNoteHandler {
	return &CreateNoteHandler{uow: uow}
}

func (h *CreateNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd createNoteCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	slug := frontmatter.Slugify(cmd.Title)

	data, err := template.RenderTemplate(cmd.TemplatePath, cmd.Title, slug)
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

	tagModels := make([]models.Tag, len(tags))
	for i, tag := range tags {
		tagModels[i] = models.NewTag(uuidv7.New(), tag)
	}

	if err = tagRepo.Upsert(ctx, tags); err != nil {
		return nil, err
	}

	if err = noteRepo.Upsert(ctx, cmd.Title, cmd.Path, slug); err != nil {
		return nil, err
	}


	h.uow.FileStore.Stage(cmd.Path, data)

	if err = h.uow.Commit(); err != nil {
		return nil, err
	}

	return cmd.Path, nil
}
