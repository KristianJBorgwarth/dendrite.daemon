package handlers

import (
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/config"
)

type initializeParams struct {
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

func Initialize(raw json.RawMessage) (*config.Config, error) {
	var params initializeParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return nil, err
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

	return cfg, nil
}
