package note

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type gotoNoteCommand struct {
	Link string `json:"link"`
}

type gotoNoteResult struct {
	Target string `json:"target"`
	Type   string `json:"type,omitempty"`
}

type GotoNoteHandler struct {
	noteRepo repositories.INoteRepository
}

func NewGotoNoteHandler(noteRepo repositories.INoteRepository) *GotoNoteHandler {
	return &GotoNoteHandler{noteRepo}
}

func (h *GotoNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd gotoNoteCommand
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	parsedLink := filehandling.ParseLink(cmd.Link)
	switch parsedLink.Kind {

	case filehandling.Note:
		note, err := h.noteRepo.GetBySlug(ctx, parsedLink.Target)
		if err != nil {
			return nil, err
		}
		if note == nil {
			return nil, nil
		}
		return gotoNoteResult{Target: note.Path(), Type: "note"}, nil

	case filehandling.URL:
		return gotoNoteResult{Target: parsedLink.Target, Type: "url"}, nil

	default:
		return nil, errors.New("unsupported link type")
	}
}
