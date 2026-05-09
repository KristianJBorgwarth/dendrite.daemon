package diagnostics

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/dtos"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type LinkDiagnosticCommand struct {
	Path string `json:"path"`
}

type LinkDiagnosticHandler struct {
	noteRepository repositories.INoteRepository
	linkRepository repositories.ILinkRepository
}

func NewLinkDiagnosticHandler(
	noteRepository repositories.INoteRepository,
	linkRepository repositories.ILinkRepository,
) *LinkDiagnosticHandler {
	return &LinkDiagnosticHandler{noteRepository, linkRepository}
}

func (h *LinkDiagnosticHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd LinkDiagnosticCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	note, err := h.noteRepository.GetByPath(ctx, cmd.Path)
	if err != nil {
		return nil, err
	}

	brokenLinks, err := h.linkRepository.GetBrokenLinks(ctx, note.ID())
	if err != nil {
		return nil, err
	}

	return mapToDiagnosticResult(note, brokenLinks), nil
}

func mapToDiagnosticResult(note *models.Note, links []*models.Link) []*dtos.LinkDiagnosticDto {
	diagnostics := make([]*dtos.LinkDiagnosticDto, 0, len(links))
	for _, link := range links {
		diagnostics = append(diagnostics, &dtos.LinkDiagnosticDto{
			NoteID:   link.ID(),
			NotePath: note.Path(),
			Line:     link.Line(),
			Col:      link.Col(),
			Raw:      link.Raw(),
			Target:   link.TargetSlug(),
			Message:  "Broken link: target note does not exist",
		})
	}

	return diagnostics
}
