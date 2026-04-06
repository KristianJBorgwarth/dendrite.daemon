package filehandling

import (
	"testing"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTags_ValidFrontMatter_ReturnsTags(t *testing.T) {
	input := `---
title: My Note
tags: ["test", "note"]
---
This is the content of the note.`

	result, err := filehandling.ParseFrontMatter([]byte(input))

	require.NoError(t, err)
	assert.Equal(t, "My Note", result.Title)
	assert.Equal(t, []string{"test", "note"}, result.Tags)
}

func TestParseTags_MissingDelimiter_ReturnsError(t *testing.T) {
	input := `title: My Note
tags: ["test", "note"]
This is the content of the note.`

	_, err := filehandling.ParseFrontMatter([]byte(input))

	assert.ErrorContains(t, err, "missing front matter delimiter")
}
