package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	id        string
	path      string
	title     string
	slug      string
	createdAt string
	updatedAt string
}

func NewNote(id, path, title, slug, createdAt, updatedAt string) *Note {
	return &Note{
		id:        id,
		path:      path,
		title:     title,
		slug:      slug,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func CreateNote(path, title, slug string) *Note {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return &Note{
		id:        id.String(),
		path:      path,
		title:     title,
		slug:      slug,
		createdAt: time.Now().Format("2006-01-02 15:04:05"),
		updatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (n *Note) ID() string {
	return n.id
}

func (n *Note) Path() string {
	return n.path
}

func (n *Note) Title() string {
	return n.title
}

func (n *Note) Slug() string {
	return n.slug
}

func (n *Note) CreatedAt() string {
	return n.createdAt
}

func (n *Note) UpdatedAt() string {
	return n.updatedAt
}
