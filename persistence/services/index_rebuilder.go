package services

import (
	"database/sql"
	"io/fs"
	"path/filepath"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
)

type IIndexRebuilder interface {
	RebuildIndex(vaultRoot string) error
}

type IndexRebuilder struct {
	db *sql.DB
}

func (r *IndexRebuilder) RebuildIndex(vaultRoot string) error {
	files, err := r.ReadFiles(vaultRoot)
	if err != nil {
		return err
	}

	notes, links, tags, noteTags := r.buildDBModels(files)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if err = r.wipeIndex(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = r.InsertData(tx, notes, links, tags, noteTags); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (r *IndexRebuilder) InsertData(
	tx *sql.Tx, 
	notes []*models.Note, 
	links []*models.Link, 
	tags []*models.Tag, 
	noteTags []*models.NoteTag) error {
	panic("not implemented")
}

func (r *IndexRebuilder) ReadFiles(vault string) ([]*filehandling.File, error) {
	var files []*filehandling.File

	err := filepath.WalkDir(vault, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		pendingFile, err := filehandling.ReadFile(path)
		if err != nil {
			return err
		}

		files = append(files, pendingFile)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *IndexRebuilder) buildDBModels(files []*filehandling.File) ([]*models.Note, []*models.Link, []*models.Tag, []*models.NoteTag) {
	var notes []*models.Note
	var links []*models.Link
	tagMap := make(map[string]*models.Tag)
	var noteTags []*models.NoteTag

	for _, file := range files {
		note := models.CreateNote(file.Path, file.Title, file.Slug)
		notes = append(notes, note)
		for _, t := range models.CreateTags(file.FrontMatter.Tags) {
			tagMap[t.Name()] = t
		}
		noteTags = append(noteTags, models.CreateNoteTags(note.ID(), file.FrontMatter.Tags)...)
		links = append(links, models.MapToLinkModel(note.ID(), file.ExtractedLinks)...)
	}

	tags := make([]*models.Tag, 0, len(tagMap))
	for _, tag := range tagMap {
		tags = append(tags, tag)
	}

	return notes, links, tags, noteTags
}

func (r *IndexRebuilder) wipeIndex(tx *sql.Tx) error {
	cmd := `DELETE FROM notes; 
	DELETE FROM tags;
	DELETE FROM note_tags;
	DELETE FROM link;`

	_, err := tx.Exec(cmd)
	return err
}
