package dtos

type LinkDiagnosticDto struct {
	NoteID   string `json:"noteId"`
	NotePath string `json:"notePath"`
	Line     int    `json:"line"`
	Col      int    `json:"col"`
	Raw      string `json:"raw"`
	Target   string `json:"target"`
	Message  string `json:"message"` 
}
