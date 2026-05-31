package note

import (
	"context"
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type saveNoteCommand struct {
	Path string `json:"path"`
}

type SaveNoteHandler struct {
	uow         *repositories.UnitOfWork
	noteRepo    repositories.INoteRepository
	tagService  services.ITagService
	noteService services.INoteService
	linkService services.ILinkService
	cfeSvc			 services.ICfService
}

func NewSaveNoteHandler(
	uow *repositories.UnitOfWork,
	nr repositories.INoteRepository,
	ts services.ITagService,
	ns services.INoteService,
	ls services.ILinkService,
	cfs services.ICfService,
) *SaveNoteHandler {
	return &SaveNoteHandler{uow, nr, ts, ns, ls, cfs}
}

func (h *SaveNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd saveNoteCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	vaultRoot := store.GetVaultStore().Config.VaultPath()

	file, err := filehandling.ReadFile(vaultRoot, cmd.Path)
	if err != nil {
		return nil, err
	}

	tx, err := h.uow.Begin()
	if err != nil {
		return nil, err
	}
	defer h.uow.Rollback()

	note, err := h.noteRepo.GetBySlug(ctx, file.Slug)
	if err != nil {
		return nil, err
	}

	if note == nil {
		note, err = h.noteService.CreateNote(ctx, tx, file.Path, file.Title, file.Slug)
		if err != nil {
			return nil, err
		}
	} else {
		if err = h.noteService.DeleteNoteMetaData(ctx, tx, note.ID()); err != nil {
			return nil, err
		}
		if err = h.noteService.UpdateNote(ctx, tx, note.ID(), file.Path, file.Title, file.Slug); err != nil {
			return nil, err
		}
	}

	tagModels, err := h.tagService.CreateTags(ctx, tx, file.FrontMatter.Tags)
	if err != nil {
		return nil, err
	}

	if err = h.tagService.CreateNoteTags(ctx, tx, note.ID(), tagModels); err != nil {
		return nil, err
	}

	if err = h.linkService.CreateLinks(ctx, tx, note.ID(), file.ExtractedLinks); err != nil {
		return nil, err
	}

	if err = h.cfeSvc.AddCfe(ctx, tx, note.ID(), file.FrontMatter.Custom); err != nil {
		return nil, err
	}

	if err := h.uow.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}
