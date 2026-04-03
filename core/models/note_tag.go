package models

type NoteTag struct {
	noteID string
	tagID  string
}

func NewNoteTag(noteID, tagID string) *NoteTag {
	return &NoteTag{
		noteID: noteID,
		tagID:  tagID,
	}
}

func (nt *NoteTag) NoteID() string {
	return nt.noteID
}

func (nt *NoteTag) TagID() string {
	return nt.tagID
}
