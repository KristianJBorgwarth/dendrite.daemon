package models

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
