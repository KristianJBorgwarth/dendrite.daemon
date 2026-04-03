package models

import "github.com/google/uuid"

type Tag struct {
	id   string
	name string
}

func NewTag(name string) *Tag {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return &Tag{
		id:   id.String(),
		name: name,
	}
}

func CreateTags(tags []string) ([]*Tag, error) {
	tagModels := make([]*Tag, len(tags))
	for i, tag := range tags {
		tagModels[i] = NewTag(tag)
	}

	return tagModels, nil
}

func (t *Tag) ID() string {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}
