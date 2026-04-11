package models

import "github.com/google/uuid"

type Link struct {
	id         string
	fromNoteID string
	targetSlug string
	raw        string
	display    string
	line       int
	col        int
}

func NewLink(id, fromNoteID, targetSlug, raw, display string, line, col int) *Link {
	return &Link{
		id:         id,
		fromNoteID: fromNoteID,
		targetSlug: targetSlug,
		raw:        raw,
		display:    display,
		line:       line,
		col:        col,
	}
}

func CreateLink(fromNoteID, targetSlug, raw, display string, line, col int) *Link {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return &Link{
		id:         id.String(),
		fromNoteID: fromNoteID,
		targetSlug: targetSlug,
		raw:        raw,
		display:    display,
		line:       line,
		col:        col,
	}
}

func (l *Link) ID() string {
	return l.id
}

func (l *Link) FromNoteID() string {
	return l.fromNoteID
}

func (l *Link) TargetSlug() string {
	return l.targetSlug
}

func (l *Link) Raw() string {
	return l.raw
}

func (l *Link) Display() string {
	return l.display
}

func (l *Link) Line() int {
	return l.line
}

func (l *Link) Col() int {
	return l.col
}
