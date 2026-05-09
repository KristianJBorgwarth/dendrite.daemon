package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/note"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/services"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCreateNoteHandler() *note.CreateNoteHandler {
	return note.NewCreateNoteHandler(
		repositories.NewUnitOfWork(),
		services.NewTagService(repositories.NewTagRepository(persistence.NewReadContext())),
		repositories.NewNoteRepository(persistence.NewReadContext()),
	)
}

func TestCreateNoteHandler_NoTemplate_CreatesNoteFileAndReturnsPath(t *testing.T) {
	// Arrange
	handler := newCreateNoteHandler()

	vaultPath := Fixture.VaultStore.Config.VaultPath()
	notePath := filepath.Join(vaultPath, "my-note.md")
	params, _ := json.Marshal(map[string]any{
		"title":        "My Note",
		"templateName": "",
		"directory":    "",
	})

	// Act
	result, err := handler.Handle(Fixture.TestContext, params)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, notePath, result)

	var count int
	require.NoError(t, Fixture.DB.QueryRow(`SELECT COUNT(*) FROM note WHERE slug = ?`, "my-note").Scan(&count))
	assert.Equal(t, 1, count)

	_, statErr := os.Stat(notePath)
	assert.NoError(t, statErr)
}

func TestCreateNoteHandler_WithTemplate_CreatesNoteFileWithTagsAndReturnsPath(t *testing.T) {
	// Arrange
	handler := newCreateNoteHandler()

	vaultPath := Fixture.VaultStore.Config.VaultPath()
	templateDir := Fixture.VaultStore.Config.TemplateDirectory()
	require.NoError(t, os.MkdirAll(templateDir, 0o755))

	templatePath := filepath.Join(templateDir, "template.md")
	require.NoError(t, os.WriteFile(templatePath, []byte("---\ntitle: Template\ntags: [go, testing]\n---\n"), 0o644))
	defer os.Remove(templatePath)

	subDir := "subdir"
	require.NoError(t, os.MkdirAll(filepath.Join(vaultPath, subDir), 0o755))
	defer os.RemoveAll(filepath.Join(vaultPath, subDir))
	notePath := filepath.Join(vaultPath, subDir, "templated-note.md")
	params, _ := json.Marshal(map[string]any{
		"title":        "Templated Note",
		"templateName": "template.md",
		"directory":    subDir,
	})

	// Act
	result, err := handler.Handle(Fixture.TestContext, params)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, notePath, result)

	var noteCount int
	require.NoError(t, Fixture.DB.QueryRow(`SELECT COUNT(*) FROM note WHERE slug = ?`, "subdir/templated-note").Scan(&noteCount))
	assert.Equal(t, 1, noteCount)

	rows, err := Fixture.DB.Query(`SELECT name FROM tag WHERE name IN ('go', 'testing')`)
	require.NoError(t, err)
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		tags = append(tags, name)
	}
	assert.Len(t, tags, 2)

	_, statErr := os.Stat(notePath)
	assert.NoError(t, statErr)
}

func TestCreateNoteHandler_DuplicateSlug_ReturnsExistingPath(t *testing.T) {
	// Arrange
	handler := newCreateNoteHandler()

	vaultPath := Fixture.VaultStore.Config.VaultPath()
	subDir := "dup-dir"
	require.NoError(t, os.MkdirAll(filepath.Join(vaultPath, subDir), 0o755))
	defer os.RemoveAll(filepath.Join(vaultPath, subDir))

	params, _ := json.Marshal(map[string]any{"title": "Dup Note", "templateName": "", "directory": subDir})
	_, err := handler.Handle(Fixture.TestContext, params)
	require.NoError(t, err)

	// Act — same title and directory produces the same slug, handler returns early
	result, err := handler.Handle(Fixture.TestContext, params)

	// Assert
	require.NoError(t, err)

	var count int
	require.NoError(t, Fixture.DB.QueryRow(`SELECT COUNT(*) FROM note WHERE slug = ?`, "dup-dir/dup-note").Scan(&count))
	assert.Equal(t, 1, count)

	expectedPath := filepath.Join(vaultPath, subDir, "dup-note.md")
	assert.Equal(t, expectedPath, result)
}

func TestCreateNoteHandler_InvalidJSON_ReturnsError(t *testing.T) {
	// Arrange
	handler := newCreateNoteHandler()

	// Act
	_, err := handler.Handle(Fixture.TestContext, json.RawMessage(`{invalid json}`))

	// Assert
	assert.Error(t, err)
}

func TestCreateNoteHandler_NonExistentTemplateName_ReturnsError(t *testing.T) {
	// Arrange
	handler := newCreateNoteHandler()

	params, _ := json.Marshal(map[string]any{
		"title":        "Ghost Note",
		"templateName": "non-existent-template",
		"directory":    "",
	})

	// Act
	_, err := handler.Handle(Fixture.TestContext, params)

	// Assert
	assert.Error(t, err)
}
