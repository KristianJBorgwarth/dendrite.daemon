package dtos

type Backlink struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Path  string `json:"path"`
	Raw   string `json:"raw"`
	Line  int    `json:"line"`
	Col   int    `json:"col"`
}

func NewBacklink(id, title, slug, path, raw string, col, line int) *Backlink {
	return &Backlink{
		ID:    id,
		Title: title,
		Slug:  slug,
		Path:  path,
		Raw:   raw,
		Line:  line,
		Col:   col,
	}
}
