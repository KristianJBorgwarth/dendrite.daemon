package services

import (
	"context"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type index struct {
	notes    []*models.Note
	links    []*models.Link
	cfe      []*models.CustomFronMatter
	tags     []*models.Tag
	noteTags []*models.NoteTag
}

type IIndexRebuilder interface {
	RebuildIndex(ctx context.Context, vaultRoot string) error
}

type indexRebuilder struct {
	uow       *repositories.UnitOfWork
	noteRepo  repositories.INoteRepository
	linkRepo  repositories.ILinkRepository
	tagRepo   repositories.ITagRepository
	indexRepo repositories.IIndexRepository
	cfeRepo   repositories.ICfeRepository
}

func NewIndexRebuilder(
	uow *repositories.UnitOfWork,
	noteRepo repositories.INoteRepository,
	linkRepo repositories.ILinkRepository,
	tagRepo repositories.ITagRepository,
	indexRepo repositories.IIndexRepository,
	cfeRepo repositories.ICfeRepository,
) *indexRebuilder {
	return &indexRebuilder{
		uow:       uow,
		noteRepo:  noteRepo,
		linkRepo:  linkRepo,
		tagRepo:   tagRepo,
		indexRepo: indexRepo,
		cfeRepo:   cfeRepo,
	}
}

func (r *indexRebuilder) RebuildIndex(ctx context.Context, vaultRoot string) error {
	dbctx, err := r.uow.Begin()
	if err != nil {
		return err
	}

	files, err := r.readFiles(vaultRoot)
	if err != nil {
		return err
	}

	index, err := r.buildDBModels(files)
	if err != nil {
		return err
	}

	if err = r.indexRepo.WipeIndex(ctx, dbctx); err != nil {
		return err
	}

	if err = r.buildIndex(ctx, dbctx, index); err != nil {
		return err
	}

	if err = r.uow.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *indexRebuilder) buildIndex(
	ctx context.Context,
	dbctx persistence.IDbContext,
	index *index,
) error {
	if err := r.noteRepo.InsertRange(ctx, dbctx, index.notes); err != nil {
		return err
	}

	if err := r.linkRepo.InsertRange(ctx, dbctx, index.links); err != nil {
		return err
	}

	if err := r.tagRepo.InsertRange(ctx, dbctx, index.tags); err != nil {
		return err
	}

	if err := r.tagRepo.InsertNoteTags(ctx, dbctx, index.noteTags); err != nil {
		return err
	}

	if err := r.cfeRepo.InsertRange(ctx, dbctx, index.cfe); err != nil {
		return err
	}

	return nil
}

func (r *indexRebuilder) readFiles(vault string) ([]*filehandling.File, error) {
	var files []*filehandling.File

	err := filepath.WalkDir(vault, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		slog.Debug("Processing file during index rebuild", "path", path)

		if d.IsDir() {
			if !r.IsValidDirectory(path) {
				slog.Debug("Skipping directory during index rebuild", "path", path)
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".md" {
			slog.Debug("Skipping non-markdown file during index rebuild", "path", path)
			return nil
		}

		pendingFile, err := filehandling.ReadFile(vault, path)
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

func (r *indexRebuilder) IsValidDirectory(path string) bool {
	ignoredDirs := store.GetVaultStore().GetExcludeIndexFiles()
	for part := range strings.SplitSeq(path, string(filepath.Separator)) {
		if slices.Contains(ignoredDirs, part) {
			return false
		}
	}
	return true
}

func (r *indexRebuilder) buildDBModels(files []*filehandling.File) (*index, error) {
	var notes []*models.Note
	var links []*models.Link
	var cfe []*models.CustomFronMatter
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

		mappedCfe, err := models.MapToCfe(note.ID(), file.FrontMatter.Custom)
		if err != nil {
			return nil, err
		}
		cfe = append(cfe, mappedCfe...)
	}

	tags := make([]*models.Tag, 0, len(tagMap))
	for _, tag := range tagMap {
		tags = append(tags, tag)
	}

	return &index{
		notes:    notes,
		links:    links,
		tags:     tags,
		noteTags: noteTags,
		cfe:      cfe,
	}, nil
}
