package repositories

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type ILinkRepository interface {
	GetByNoteID(ctx context.Context, fromNoteID string) ([]*models.Link, error)
	GetBySlug(ctx context.Context, targetSlug string) ([]*models.Link, error)
	Search(ctx context.Context, query string) ([]*models.Link, error)
}

type linkRepository struct {
	dbContext persistence.IDbContext
}

func NewLinkRepository(ctx persistence.IDbContext) ILinkRepository {
	return &linkRepository{dbContext: ctx}
}

func (r *linkRepository) GetByNoteID(ctx context.Context, fromNoteID string) ([]*models.Link, error) {
	rows, err := r.dbContext.QueryContext(ctx, "SELECT id, from_note_id, target_slug, raw, display, line, col FROM links WHERE from_note_id = ?", fromNoteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		var id, fromNoteID, targetSlug, raw, display string
		var line, col int
		if err := rows.Scan(&id, &fromNoteID, &targetSlug, &raw, &display, &line, &col); err != nil {
			return nil, err
		}
		links = append(links, models.NewLink(id, fromNoteID, targetSlug, raw, display, line, col))
	}

	return links, nil
}

func (r *linkRepository) GetBySlug(ctx context.Context, targetSlug string) ([]*models.Link, error) {
	rows, err := r.dbContext.QueryContext(ctx, `SELECT id, from_note_id, target_slug, raw, display, line, col FROM links WHERE target_slug = ?`, targetSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		var id, fromNoteID, targetSlug, raw, display string
		var line, col int
		if err := rows.Scan(&id, &fromNoteID, &targetSlug, &raw, &display, &line, &col); err != nil {
			return nil, err
		}
		links = append(links, models.NewLink(id, fromNoteID, targetSlug, raw, display, line, col))
	}

	return links, nil
}

func (r *linkRepository) Search(ctx context.Context, query string) ([]*models.Link, error) {
	rows, err := r.dbContext.QueryContext(ctx, `SELECT id, from_note_id, target_slug, raw, display, line, col FROM links WHERE raw LIKE ?`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		var id, fromNoteID, targetSlug, raw, display string
		var line, col int
		if err := rows.Scan(&id, &fromNoteID, &targetSlug, &raw, &display, &line, &col); err != nil {
			return nil, err
		}
		links = append(links, models.NewLink(id, fromNoteID, targetSlug, raw, display, line, col))
	}

	return links, nil
}
