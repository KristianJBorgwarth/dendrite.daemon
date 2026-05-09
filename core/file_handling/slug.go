package filehandling

import (
	"bytes"
	"path/filepath"
	"strings"
)

func Slugify(s string) string {
	var slug bytes.Buffer
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			slug.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			slug.WriteRune(r + 32) // Convert to lowercase
		} else {
			slug.WriteRune('-')
		}
	}
	return slug.String()
}

func SlugRelative(base, path string) string {
	rel := strings.TrimPrefix(path, base+string(filepath.Separator))
	rel = strings.TrimSuffix(rel, filepath.Ext(rel))
	return rel
}
