package note

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type gotoNoteCommand struct {
	Link string `json:"link"`
}

type gotoNoteResult struct {
	Path string `json:"path"`
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

	cmd.Link = h.resolveLink(cmd.Link)

	note, err := h.noteRepo.GetBySlug(ctx, cmd.Link)
	if err != nil {
		return nil, err
	}
	if note == nil {
		slog.Debug("Note not found for link", "link", cmd.Link)
		return nil, nil
	}
	return gotoNoteResult{Path: note.Path()}, nil
}

func (h *GotoNoteHandler) resolveLink(link string) string {
	if len(link) < 5 || link[:2] != "[[" || link[len(link)-2:] != "]]" {
		return link
	}
	content := link[2 : len(link)-2]
	if before, _, ok := strings.Cut(content, "|"); ok {
		return before
	}
	return content
}
