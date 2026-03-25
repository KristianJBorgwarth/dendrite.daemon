package test

import (
	"strings"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/frontmatter"
)

func TestParseFrontMatter(t *testing.T) {
	input := `title: My Note
template: default
tags: ["test", "note"]
---
This is the content of the note.`

	expected := map[string]string{
		"title":    "My Note",
		"template": "default",
		"tags":     `["test", "note"]`,
	}

	result, err := frontmatter.ParseFrontMatter(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for key, expectedValue := range expected {
		if value, ok := result[key]; !ok || value != expectedValue {
			t.Errorf("expected %s to be %s, got %s", key, expectedValue, value)
		}
	}
}
