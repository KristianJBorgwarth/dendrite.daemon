package models

import filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"

func MapToLinkModel(noteID string, extractedLinks []*filehandling.ExtractedLink) []*Link {
	var links []*Link
	for _, link := range extractedLinks {
		links = append(links, CreateLink(noteID, link.TargetSlug, link.Raw, link.Display, link.Line, link.Col))
	}
	return links
}
