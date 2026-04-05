package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/frontmatter"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/template"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/utils"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type createNoteCommand struct {
	Title        string            `json:"title"`
	TemplatePath string            `json:"templatePath"`
	Path         string            `json:"path"`
}

type CreateNoteHandler struct{
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

	data, err := template.RenderTemplate(cmd.TemplatePath, cmd.Title, slug)
	if err != nil {
		return nil, err
	}

	slog.Debug("rendered template", "data", string(data))

	tags, err := frontmatter.ParseTags(data)
	if err != nil {
		return nil, err
	}

	slog.Debug("parsed tags", "tags", tags)

	tx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer h.uow.Rollback()

	tagRepo := repositories.NewTagRepository(tx)
	noteRepo := repositories.NewNoteRepository(tx)

	tagModels, err := models.CreateTags(tags)
	if err != nil {
		return nil, err
	}

	slog.Debug("Creating tags", "tags", tags)

	if err = tagRepo.Upsert(ctx, tagModels); err != nil {
		return nil, err
	}

	note := models.CreateNote(cmd.Path, cmd.Title, slug)

	if err = noteRepo.Upsert(ctx, note); err != nil {
		return nil, err
	}

	if err = tagRepo.UpsertNoteTags(ctx, note.ID(), utils.Select(tagModels, func(t *models.Tag) string { return t.ID() })); err != nil {
		return nil, err
	}

	h.uow.FileStore.Stage(cmd.Path, data)

	if err = h.uow.Commit(); err != nil {
		return nil, err
	}

	return cmd.Path, nil
}
