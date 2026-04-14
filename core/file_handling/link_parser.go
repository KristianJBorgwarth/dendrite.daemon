package filehandling

import "strings"

type LinkType int

const (
	Unknown LinkType = iota
	Note
	URL
)

type ParsedLink struct {
	Kind   LinkType
	Target string
}

func ParseLink(s string) ParsedLink {
	if len(s) >= 5 && strings.HasPrefix(s, "[[") && strings.HasSuffix(s, "]]") {
		content := s[2 : len(s)-2]
		if before, _, ok := strings.Cut(content, "|"); ok {
			return ParsedLink{Note, before}
		}
		return ParsedLink{Note, content}
	}
	if strings.HasPrefix(s, "[") {
		if close := strings.Index(s, "]("); close >= 0 && strings.HasSuffix(s, ")") {
			url := s[close+2 : len(s)-1]
			return ParsedLink{URL, url}
		}
	}
	return ParsedLink{Note, s}
}
