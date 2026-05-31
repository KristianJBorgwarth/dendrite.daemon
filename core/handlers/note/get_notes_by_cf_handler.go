package note

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/dtos"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type getNotesByCfQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetNotesByCfHandler struct {
	cfeRepo repositories.ICfRepository
}

func NewGetNotesByCfeHandler(
	cfeRepo repositories.ICfRepository,
) *GetNotesByCfHandler {
	return &GetNotesByCfHandler{cfeRepo: cfeRepo}
}

func (h *GetNotesByCfHandler) Handle(
	ctx context.Context,
	raw json.RawMessage,
) (any, error) {
	var query getNotesByCfQuery
	if err := json.Unmarshal(raw, &query); err != nil {
		return nil, err
	}

	notes, err := h.cfeRepo.Search(ctx, query.Key, query.Value)
	if err != nil {
		return nil, err
	}

	if len(notes) == 0 {
		return make([]*dtos.NoteDto,0), nil
	}

	noteDtos := dtos.NewNoteDtos(notes)

	return noteDtos, nil
}
