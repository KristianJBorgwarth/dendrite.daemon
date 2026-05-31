package integration_test

import (
	"encoding/json"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/dtos"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/note"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/KristianJBorgwarth/dendrite.daemon/test/test_integration/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newGetNotesByCfeHandler() *note.GetNotesByCfeHandler {
	return note.NewGetNotesByCfeHandler(
		repositories.NewCfeRepository(persistence.NewReadContext()),
	)
}

func Test_Handle_ReturnsNotesWithMacthingCfe(t *testing.T) {
	// Arrange
	handler := newGetNotesByCfeHandler()
	noteID := "test-note-id"
	cfeKey := "author"
	cfeValues := []string{"John Doe", "Jane Smith"}
	utils.CreateNote(Fixture.TestContext, Fixture.DB, noteID, "Test Note", "test-note.md", "test");

	utils.CreateCfe(Fixture.TestContext, Fixture.DB, noteID, 
		cfeKey, cfeValues)

	params, _ := json.Marshal(map[string]any{
		"key":   "author",
		"value": "J",
	})

	// Act
	result, err := handler.Handle(Fixture.TestContext, params)

	require.NoError(t, err)
	require.NotNil(t, result)

	noteDtos, ok := result.([]*dtos.NoteDto)
	require.True(t, ok)
	assert.Len(t, noteDtos, 2)
}

func Test_Handle_ReturnsEmptyListWhenNoMatchingCfe(t *testing.T) {
	// Arrange
	handler := newGetNotesByCfeHandler()
	noteID := "test-note-id-2"
	cfeKey := "category"
	cfeValues := []string{"Tech", "Lifestyle"}
	utils.CreateNote(Fixture.TestContext, Fixture.DB, noteID, "Another Test Note", "another-test-note.md", "test");

	utils.CreateCfe(Fixture.TestContext, Fixture.DB, noteID, 
		cfeKey, cfeValues)

	params, _ := json.Marshal(map[string]any{
		"key":   "category",
		"value": "NonExistingCategory",
	})

	// Act
	result, err := handler.Handle(Fixture.TestContext, params)

	require.NoError(t, err)
	require.NotNil(t, result)

	noteDtos, ok := result.([]*dtos.NoteDto)
	require.True(t, ok)
	assert.Len(t, noteDtos, 0)
}
