package models

import (
	"errors"
	"log/slog"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
)

func MapToLinkModel(noteID string, extractedLinks []*filehandling.ExtractedLink) []*Link {
	var links []*Link
	for _, link := range extractedLinks {
		links = append(links, CreateLink(noteID, link.TargetSlug, link.Raw, link.Display, link.Line, link.Col))
	}
	return links
}

func MapToCustomFrontmatter(noteID string, extractedCfe map[string]any) ([]*CustomFronMatter, error) {
	var cfe []*CustomFronMatter
	for key, value := range extractedCfe {
		if valueStr, ok := value.(string); ok {
			cfe = append(cfe, NewCustomFrontMatter(noteID, key, valueStr))
			continue
		} else if valueArr, ok := value.([]any); ok {
			for _, v := range valueArr {
				if s, ok := v.(string); ok {
					cfe = append(cfe, NewCustomFrontMatter(noteID, key, s))
				}
			}
			continue
		} else {
			slog.Warn("Unsupported CFE value type, skipping", "key", key, "value", value)
			return nil, errors.New("unsupported CFE value type")
		}
	}
	return cfe, nil
}
