package models


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

func (t *Tag) ID() string {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}
