package filehandling

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	Title       string
	Slug        string
	FrontMatter FrontMatter
	Content     []string
	Links       *[]ExtractedLink
}

type ExtractedLink struct {
	Raw        string
	Display    string
	TargetSlug string
	Line       int
	Col        int
}

func ReadFile(path string) (*File, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	body = bytes.TrimSpace(body)

	fm, err := ParseFrontMatter(body)
	if err != nil {
		return nil, err
	}

	return &File{
		Title:       fm.Title,
		Slug:        strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		FrontMatter: *fm,
		Content:     strings.Split(string(body), "\n"),
		Links:       ExtractLinks(body),
	}, nil
}

var linkRegex = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)

func ExtractLinks(body []byte) *[]ExtractedLink {
	var links []ExtractedLink

	lineStart := 0
	lineNum := 1

	for i := 0; i <= len(body); i++ {
		if i == len(body) || body[i] == '\n' {
			line := body[lineStart:i]
			matches := linkRegex.FindAllSubmatchIndex(line, -1)
			for _, m := range matches {
				raw := string(line[m[0]:m[1]])
				slug := string(line[m[2]:m[3]])
				var display string
				if m[4] != -1 {
					display = string(line[m[4]:m[5]])
				}
				links = append(links, ExtractedLink{
					TargetSlug: slug,
					Raw:        raw,
					Display:    display,
					Line:       lineNum,
					Col:        m[0] + 1,
				})
			}
			lineStart = i + 1
			lineNum++
		}
	}

	return &links
}
