package filehandling

import (
	"os"
)

type File struct{
	Title string
	Slug string
	FrontMatter FrontMatter
	Content[] string
}

func ReadFile(path string) (*File, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fm, err := ParseFrontMatter(content)
	if err != nil {
		return nil, err
	}
}
