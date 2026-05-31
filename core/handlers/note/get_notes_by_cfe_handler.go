package note

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/dtos"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type getNotesByCfeQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetNotesByCfeHandler struct {
	cfeRepo repositories.ICfeRepository
}

func NewGetNotesByCfeHandler(
	cfeRepo repositories.ICfeRepository,
) *GetNotesByCfeHandler {
	return &GetNotesByCfeHandler{cfeRepo: cfeRepo}
}

func (h *GetNotesByCfeHandler) Handle(
	ctx context.Context,
	raw json.RawMessage,
) (any, error) {
	var query getNotesByCfeQuery
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
