package models

import "github.com/google/uuid"

type Tag struct {
	id   string
	name string
}

func NewTag(id string, name string) *Tag {
	return &Tag{
		id:   id,
		name: name,
	}
}

func NewTags(tags []string) ([]*Tag, error) {

	tagModels := make([]*Tag, len(tags))
	for i, tag := range tags {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		tagModels[i] = NewTag(id.String(), tag)
	}

	return tagModels, nil
}

func (t *Tag) ID() string {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}
