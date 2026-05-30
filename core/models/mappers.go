package models

import filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"

func MapToLinkModel(noteID string, extractedLinks []*filehandling.ExtractedLink) []*Link {
	var links []*Link
	for _, link := range extractedLinks {
		links = append(links, CreateLink(noteID, link.TargetSlug, link.Raw, link.Display, link.Line, link.Col))
	}
	return links
}

func MapToCfe(noteID string, extractedCfe map[string]any) []*CustomFronMatter {
	var cfe []*CustomFronMatter
	for key, value := range extractedCfe {
		if valueStr, ok := value.(string); ok {
			cfe = append(cfe, NewCustomFrontMatter(noteID, key, valueStr))
			continue
		} else if valueArr, ok := value.([]string); ok {
			for _, v := range valueArr {
				cfe = append(cfe, NewCustomFrontMatter(noteID, key, v))
			}
			continue
		} else {
			cfe = append(cfe, NewCustomFrontMatter(noteID, key, ""))
		}
	}
	return cfe
}
