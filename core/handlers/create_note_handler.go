package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/frontmatter"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
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

func (h CreateNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd createNoteCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	data, err := h.getTemplate(cmd.TemplatePath)
	if err != nil {
		return nil, err
	}

	tags, err := frontmatter.ParseTags(data)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(tx *sql.Tx) error {
		tagRepo := repositories.NewTagRepository(tx)
		noteRepo := repositories.NewNoteRepository(tx)

		if err = tagRepo.Upsert(ctx, tags); err != nil {
			return err
		}

		if err = noteRepo.Upsert(ctx, cmd.Title, cmd.Path, frontmatter.Slugify(cmd.Title)); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return cmd.Path, nil
}

func (h *CreateNoteHandler) getTemplate(templatePath string) ([]byte, error) {
	if templatePath == "" {
		return nil, nil
	}
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
