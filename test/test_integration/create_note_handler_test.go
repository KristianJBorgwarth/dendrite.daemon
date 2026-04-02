package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCreateNoteHandler() *handlers.CreateNoteHandler {
	uow := repositories.NewUnitOfWork(Fixture.DB)
	return handlers.NewCreateNoteHandler(uow)
}

func TestCreateNoteHandler_NoTemplate_CreatesNoteFileAndReturnsPath(t *testing.T) {
	handler := newCreateNoteHandler()
	notePath := filepath.Join(t.TempDir(), "my-note.md")

	params, _ := json.Marshal(map[string]any{
		"title": "My Note",
		"path":  notePath,
	})

	result, err := handler.Handle(Fixture.TestContext, params)

	require.NoError(t, err)
	assert.Equal(t, notePath, result)

	var count int
	require.NoError(t, Fixture.DB.QueryRow(`SELECT COUNT(*) FROM notes WHERE slug = ?`, "my-note").Scan(&count))
	assert.Equal(t, 1, count)

	_, statErr := os.Stat(notePath)
	assert.NoError(t, statErr)
}

func TestCreateNoteHandler_WithTemplate_CreatesNoteFileWithTagsAndReturnsPath(t *testing.T) {
	handler := newCreateNoteHandler()
	dir := t.TempDir()

	templatePath := filepath.Join(dir, "template.md")
	require.NoError(t, os.WriteFile(templatePath, []byte("---\ntitle: Template\ntags: [go, testing]\n---\n"), 0o644))

	notePath := filepath.Join(dir, "templated-note.md")
	params, _ := json.Marshal(map[string]any{
		"title":        "Templated Note",
		"path":         notePath,
		"templatePath": templatePath,
	})

	result, err := handler.Handle(Fixture.TestContext, params)

	require.NoError(t, err)
	assert.Equal(t, notePath, result)

	var noteCount int
	require.NoError(t, Fixture.DB.QueryRow(`SELECT COUNT(*) FROM notes WHERE slug = ?`, "templated-note").Scan(&noteCount))
	assert.Equal(t, 1, noteCount)

	rows, err := Fixture.DB.Query(`SELECT name FROM tags WHERE name IN ('go', 'testing')`)
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

func TestCreateNoteHandler_DuplicateSlug_Upserts(t *testing.T) {
	dir := t.TempDir()

	path1 := filepath.Join(dir, "dup-note.md")
	params1, _ := json.Marshal(map[string]any{"title": "Dup Note", "path": path1})
	_, err := newCreateNoteHandler().Handle(Fixture.TestContext, params1)
	require.NoError(t, err)

	path2 := filepath.Join(dir, "dup-note-moved.md")
	params2, _ := json.Marshal(map[string]any{"title": "Dup Note", "path": path2})
	_, err = newCreateNoteHandler().Handle(Fixture.TestContext, params2)
	require.NoError(t, err)

	var count int
	require.NoError(t, Fixture.DB.QueryRow(`SELECT COUNT(*) FROM notes WHERE slug = ?`, "dup-note").Scan(&count))
	assert.Equal(t, 1, count)

	var path string
	require.NoError(t, Fixture.DB.QueryRow(`SELECT path FROM notes WHERE slug = ?`, "dup-note").Scan(&path))
	assert.Equal(t, path2, path)
}

func TestCreateNoteHandler_InvalidJSON_ReturnsError(t *testing.T) {
	_, err := newCreateNoteHandler().Handle(Fixture.TestContext, json.RawMessage(`{invalid json}`))
	assert.Error(t, err)
}

func TestCreateNoteHandler_NonExistentTemplatePath_ReturnsError(t *testing.T) {
	params, _ := json.Marshal(map[string]any{
		"title":        "Ghost Note",
		"path":         filepath.Join(t.TempDir(), "ghost.md"),
		"templatePath": "/non/existent/template.md",
	})

	_, err := newCreateNoteHandler().Handle(Fixture.TestContext, params)
	assert.Error(t, err)
}
