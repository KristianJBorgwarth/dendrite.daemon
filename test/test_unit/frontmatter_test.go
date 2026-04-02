package frontmatter_test

import (
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/frontmatter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTags_ValidFrontMatter_ReturnsTags(t *testing.T) {
	input := `---
title: My Note
tags: ["test", "note"]
---
This is the content of the note.`

	result, err := frontmatter.ParseTags([]byte(input))

	require.NoError(t, err)
	assert.Equal(t, []string{"test", "note"}, result)
}

func TestParseTags_MissingDelimiter_ReturnsError(t *testing.T) {
	input := `title: My Note
tags: ["test", "note"]
This is the content of the note.`

	_, err := frontmatter.ParseTags([]byte(input))

	assert.ErrorContains(t, err, "missing front matter delimiter")
}
