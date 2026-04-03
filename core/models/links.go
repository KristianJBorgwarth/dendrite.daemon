package models

type Link struct {
	fromNoteID string
	toNoteID   string
	raw        string
}

func NewLink(fromNoteID, toNoteID, raw string) *Link {
	return &Link{
		fromNoteID: fromNoteID,
		toNoteID:   toNoteID,
		raw:        raw,
	}
}

func (l *Link) FromNoteID() string {
	return l.fromNoteID
}

func (l *Link) ToNoteID() string {
	return l.toNoteID
}

func (l *Link) Raw() string {
	return l.raw
}
