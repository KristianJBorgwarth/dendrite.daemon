package models

type NoteTag struct {
	noteID int
	tagID  int
}

func NewNoteTag(noteID, tagID int) *NoteTag {
	return &NoteTag{
		noteID: noteID,
		tagID:  tagID,
	}
}

func (nt *NoteTag) NoteID() int {
	return nt.noteID
}

func (nt *NoteTag) TagID() int {
	return nt.tagID
}
