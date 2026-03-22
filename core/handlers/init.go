package handlers

import (
	"encoding/json"
	"github.com/KristianJBorgwarth/dendrite.daemon/config"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type initializeCommand struct {
	VaultPath   string `json:"vaultPath"`
	TemplateDir string `json:"templateDir"`

	ScratchNote struct {
		Dir          string `json:"dir"`
		TemplateName string `json:"templateName"`
	} `json:"scratchNote"`

	DailyNote struct {
		Dir            string `json:"dir"`
		TemplateName   string `json:"templateName"`
		FilenameFormat string `json:"filenameFormat"`
	} `json:"dailyNote"`
}

type InitializeHandler struct{}

func (h InitializeHandler) Handle(raw json.RawMessage) (any, *rpc.Error) {
	var params initializeCommand

	if err := json.Unmarshal(raw, &params); err != nil {
		return nil, &rpc.Error{
			Code:    -32602,
			Message: "invalid params: " + err.Error(),
		}
	}

	cfg := &config.Config{
		VaultPath:   params.VaultPath,
		TemplateDir: params.TemplateDir,
		ScratchNote: config.ScratchConfig{
			Dir:          params.ScratchNote.Dir,
			TemplateName: params.ScratchNote.TemplateName,
		},
		DailyNote: config.DailyConfig{
			Dir:            params.DailyNote.Dir,
			TemplateName:   params.DailyNote.TemplateName,
			FilenameFormat: params.DailyNote.FilenameFormat,
		},
	}

	cfg.SetDefaults()

	err := persistence.InitializeIndex(cfg.VaultPath)
	if err != nil {
		return nil, &rpc.Error{
			Code:    -1,
			Message: "failed to initialize index: " + err.Error(),
		}
	}

	return cfg, nil
}

