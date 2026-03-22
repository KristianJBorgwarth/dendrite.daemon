package handlers

import (
	"encoding/json"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
)

type createNoteCommand struct {
	Title        string            `json:"title"`
	TemplatePath string            `json:"templatePath"`
	Path         string            `json:"path"`
	Vars         map[string]string `json:"vars"`
}

type CreateNoteHandler struct{
	repository NoteRepository
}

func (h CreateNoteHandler) Handle(params []byte) (any, *rpc.Error) {
	var cmd createNoteCommand
	if err := json.Unmarshal(params, &cmd); err != nil {
		return nil, &rpc.Error{Code: -32602, Message: "invalid params"}
	}

}
