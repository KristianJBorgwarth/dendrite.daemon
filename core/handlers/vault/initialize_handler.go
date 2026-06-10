package vault

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type initializeCommand struct {
	Config Config `json:"config"`
}

type Config struct {
	VaultName              string     `json:"vault_name"`
	VaultPath              string     `json:"vault_path"`
	TemplatesDir           string     `json:"templates_dir"`
	ExcludeIndexFiles      []string   `json:"exclude_index_files"`
	OverrideDefaultIgnores bool       `json:"override_default_ignores"`
	DailyNotes             DailyNotes `json:"daily_notes"`
}

type DailyNotes struct {
	Dir            string `json:"dir"`
	FilenameFormat string `json:"filename_format"`
	TemplateName   string `json:"template_name"`
}

type InitializeHandler struct {
	idxRebuilder services.IIndexRebuilder
	noteService  services.INoteService
}

func NewInitializeHandler(idxR services.IIndexRebuilder, nsv services.INoteService) *InitializeHandler {
	return &InitializeHandler{idxRebuilder: idxR, noteService: nsv}
}

func (h InitializeHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
	var cmd initializeCommand

	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, err
	}

	if cmd.Config.VaultPath == "" {
		return nil, errors.New("vault_path is empty — check your dendrite.nvim config")
	}

	store := store.NewVaultStore(
		cmd.Config.VaultName,
		cmd.Config.VaultPath,
		cmd.Config.TemplatesDir,
		cmd.Config.ExcludeIndexFiles,
		cmd.Config.OverrideDefaultIgnores,
		cmd.Config.DailyNotes.Dir,
		cmd.Config.DailyNotes.FilenameFormat,
		cmd.Config.DailyNotes.TemplateName,
	)

	err := persistence.InitializeDBContext(store.Config.VaultName())
	if err != nil {
		return nil, err
	}

	noteCount, err := h.noteService.GetNoteCount(ctx)
	if err != nil {
		return nil, err
	}

	if noteCount > 0 {
		slog.Debug("Vault already initialized, skipping index rebuild")
		return nil, nil
	}

	err = h.idxRebuilder.RebuildIndex(ctx, cmd.Config.VaultPath)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
