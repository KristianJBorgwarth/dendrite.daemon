package models

type Link struct {
	fromNoteID int
	toNoteID   int
	raw        string
}

func NewLink(fromNoteID, toNoteID int, raw string) *Link {
	return &Link{
		fromNoteID: fromNoteID,
		toNoteID:   toNoteID,
		raw:        raw,
	}
}

func (l *Link) FromNoteID() int {
	return l.fromNoteID
}

func (l *Link) ToNoteID() int {
	return l.toNoteID
}

func (l *Link) Raw() string {
	return l.raw
}
