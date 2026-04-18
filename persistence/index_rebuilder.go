package persistence

import (
	"context"
	"database/sql"
	"io/fs"
	"log/slog"
	"path/filepath"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
)

type IIndexRebuilder interface {
	RebuildIndex(ctx context.Context, vaultRoot string) error
}

type IndexRebuilder struct {
}

func NewIndexRebuilder() *IndexRebuilder {
	return &IndexRebuilder{}
}

func (r *IndexRebuilder) RebuildIndex(ctx context.Context, vaultRoot string) error {
	files, err := r.ReadFiles(vaultRoot)
	if err != nil {
		slog.Debug("Failed to read files from vault", "vaultRoot", vaultRoot, "error", err)
		return err
	}

	notes, links, tags, noteTags := r.buildDBModels(files)

	tx, err := GetDBContext().DB.Begin()
	if err != nil {
		return err
	}

	if err = r.wipeIndex(ctx, tx); err != nil {
		tx.Rollback()
		slog.Debug("Failed to wipe index, rolling back transaction", "error", err)
		return err
	}

	if err = r.InsertNotes(ctx, tx, notes); err != nil {
		tx.Rollback()
		return err
	}

	if err = r.InsertLinks(ctx, tx, links); err != nil {
		tx.Rollback()
		return err
	}

	if err = r.InsertTags(ctx, tx, tags); err != nil {
		tx.Rollback()
		return err
	}

	if err = r.InsertNoteTags(ctx, tx, noteTags); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (r *IndexRebuilder) InsertNotes(ctx context.Context, tx *sql.Tx, notes []*models.Note) error {
	noteStmt, err := tx.Prepare(`INSERT INTO note (id, title, slug, path) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	for _, note := range notes {
		if _, err = noteStmt.ExecContext(ctx, note.ID(), note.Title(), note.Slug(), note.Path()); err != nil {
			return err
		}
	}

	return nil
}

func (r *IndexRebuilder) InsertLinks(ctx context.Context, tx *sql.Tx, links []*models.Link) error {
	linkStmt, err := tx.Prepare(`INSERT INTO link (id, from_note_id, target_slug, display, raw, line, col) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	for _, link := range links {
		if _, err = linkStmt.ExecContext(ctx, link.ID(), link.FromNoteID(), link.TargetSlug(), link.Display(), link.Raw(), link.Line(), link.Col()); err != nil {
			return err
		}
	}

	return nil
}

func (r *IndexRebuilder) InsertTags(ctx context.Context, tx *sql.Tx, tags []*models.Tag) error {
	tagStmst, err := tx.Prepare(`INSERT INTO tag (name) VALUES (?)`)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if _, err = tagStmst.ExecContext(ctx, tag.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (r *IndexRebuilder) InsertNoteTags(ctx context.Context, tx *sql.Tx, noteTags []*models.NoteTag) error {
	noteTagStmt, err := tx.Prepare(`INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	for _, noteTag := range noteTags {
		if _, err = noteTagStmt.ExecContext(ctx, noteTag.NoteID(), noteTag.TagID()); err != nil {
			return err
		}
	}

	return nil
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

func (r *IndexRebuilder) wipeIndex(ctx context.Context, tx *sql.Tx) error {
	cmd := `DELETE FROM notes; 
	DELETE FROM tags;
	DELETE FROM note_tags;
	DELETE FROM link;`

	_, err := tx.ExecContext(ctx, cmd)
	return err
}
