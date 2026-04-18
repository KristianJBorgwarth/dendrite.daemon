package services

import (
	"context"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"

	filehandling "github.com/KristianJBorgwarth/dendrite.daemon/core/file_handling"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/repositories"
)

type IIndexRebuilder interface {
	RebuildIndex(ctx context.Context, vaultRoot string) error
}

type indexRebuilder struct {
	uow       repositories.UnitOfWork
	noteRepo  repositories.INoteRepository
	linkRepo  repositories.ILinkRepository
	tagRepo   repositories.ITagRepository
	indexRepo repositories.IIndexRepository
}

func NewIndexRebuilder(
	uow repositories.UnitOfWork,
	noteRepo repositories.INoteRepository,
	linkRepo repositories.ILinkRepository,
	tagRepo repositories.ITagRepository,
	indexRepo repositories.IIndexRepository,
) *indexRebuilder {
	return &indexRebuilder{
		uow:      uow,
		noteRepo: noteRepo,
		linkRepo: linkRepo,
		tagRepo:  tagRepo,
	}
}

func (r *indexRebuilder) RebuildIndex(ctx context.Context, vaultRoot string) error {
	files, err := r.readFiles(vaultRoot)
	if err != nil {
		slog.Debug("Failed to read files from vault", "vaultRoot", vaultRoot, "error", err)
		return err
	}

	dbctx, err := r.uow.Begin()
	if err != nil {
		return err
	}

	if err = r.indexRepo.WipeIndex(ctx, dbctx); err != nil {
		r.uow.Rollback()
		slog.Debug("Failed to wipe index, rolling back transaction", "error", err)
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
			if !r.shouldIndexDirectory(path) {
				slog.Debug("Skipping directory during index rebuild", "path", path)
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".md" {
			slog.Debug("Skipping non-markdown file during index rebuild", "path", path)
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

func (r *indexRebuilder) shouldIndexDirectory(path string) bool {
	ignoredDirs := []string{".git", ".templates", "temp", "issues"}
	for part := range strings.SplitSeq(path, string(filepath.Separator)) {
		if slices.Contains(ignoredDirs, part) {
			return false
		}
	}
	return true
}
