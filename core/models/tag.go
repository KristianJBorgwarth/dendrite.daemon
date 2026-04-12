package models

import "github.com/google/uuid"

type Tag struct {
	id   string
	name string
}

func NewTag(id, name string) *Tag {
	return &Tag{
		id:   id,
		name: name,
	}
}

func CreateTags(tags []string) ([]*Tag) {
	tagModels := make([]*Tag, len(tags))
	for i, tag := range tags {
		id, _ := uuid.NewV7()
		tagModels[i] = NewTag(id.String(), tag)
	}

	return tagModels
}

func (t *Tag) ID() string {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}
