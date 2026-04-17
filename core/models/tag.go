package models

type Tag struct {
	name string
}

func NewTag(name string) *Tag {
	return &Tag{
		name: name,
	}
}

func CreateTags(tags []string) ([]*Tag) {
	tagModels := make([]*Tag, len(tags))
	for i, tag := range tags {
		tagModels[i] = NewTag(tag)
	}

	return tagModels
}

func (t *Tag) Name() string {
	return t.name
}
