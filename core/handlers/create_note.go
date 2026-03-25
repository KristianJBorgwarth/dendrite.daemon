package handlers

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/frontmatter"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type createNoteCommand struct {
	Title        string            `json:"title"`
	TemplatePath string            `json:"templatePath"`
	Path         string            `json:"path"`
	Vars         map[string]string `json:"vars"`
}

type CreateNoteHandler struct {
	noteRepo repositories.NoteRepository
	tagRepo  repositories.TagRepository
}

func (h CreateNoteHandler) Handle(params []byte) (*rpc.Response, *rpc.Error) {
	var cmd createNoteCommand

	if err := json.Unmarshal(params, &cmd); err != nil {
		return nil, h.ReturnError(err)
	}

	data, err := os.ReadFile(cmd.TemplatePath)
	if err != nil {
		return nil, h.ReturnError(err)
	}

	feMatter, err := frontmatter.ParseFrontMatter(bytes.NewReader(data))
	if err != nil {
		return nil, h.ReturnError(err)
	}

	tags, err := frontmatter.ExtractTags(feMatter)
	if err != nil {
		return nil, h.ReturnError(err)
	}

	err = h.tagRepo.Upsert(tags)
	if err != nil {
		return nil, h.ReturnError(err)
	}

	err = h.noteRepo.Upsert(cmd.Title, cmd.Path, frontmatter.Slugify(cmd.Title))
	if err != nil {
		return nil, h.ReturnError(err)
	}

	return &rpc.Response{Jsonrpc: "2.0", Result: cmd.Path}, nil
}

func (h *CreateNoteHandler) ReturnError(err error) *rpc.Error {
	return &rpc.Error{Code: -1, Message: err.Error()}
}
