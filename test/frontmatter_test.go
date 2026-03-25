package frontmatter_test

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

func TestExtractTags(t *testing.T) {
	frontMatter := map[string]string{
		"tags": `["test", "note"]`,
	}

	expected := []string{"test", "note"}

	result, err := frontmatter.ExtractTags(frontMatter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d tags, got %d", len(expected), len(result))
	}

	for i, expectedTag := range expected {
		if result[i] != expectedTag {
			t.Errorf("expected tag %d to be %s, got %s", i, expectedTag, result[i])
		}
	}
}
