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

func NewNoteTags(noteID string, tagIDs []string) []*NoteTag {
	noteTags := make([]*NoteTag, len(tagIDs))
	for i, tagID := range tagIDs {
		noteTags[i] = NewNoteTag(noteID, tagID)
	}
	return noteTags
}

func (nt *NoteTag) NoteID() string {
	return nt.noteID
}

func (nt *NoteTag) TagID() string {
	return nt.tagID
}
