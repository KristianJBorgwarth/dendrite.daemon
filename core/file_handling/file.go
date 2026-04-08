package filehandling

import (
	"bytes"
	"os"
)

type File struct{
	Title string
	Slug string
	FrontMatter FrontMatter
	Links []string
	Content[] string
}

func ReadFile(path string) (*File, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	content = bytes.TrimSpace(content)

	fm, err := ParseFrontMatter(content)
	if err != nil {
		return nil, err
	}

	links, err := ResolveLinks(content)
	if err != nil {
		return nil, err
	}

	return &File{
		Title: fm.Title,
		Slug: Slugify(fm.Title),
		FrontMatter: *fm,
		Links: links,
	}, nil
}

func ResolveLinks(file []byte) ([]string, error) {
	// Placeholder for link resolution logic
	return []string{}, nil
}
