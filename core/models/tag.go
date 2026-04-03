package models

import "github.com/google/uuid"

type Tag struct {
	id   string
	name string
}

func NewTag(id uuid.UUID, name string) *Tag {
	return &Tag{
		id:   id.String(),
		name: name,
	}
}

func (t *Tag) ID() string {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}
