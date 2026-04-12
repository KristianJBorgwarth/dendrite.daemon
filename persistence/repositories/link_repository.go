package repositories

import (
	"context"
	"strings"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type ILinkRepository interface {
	Insert(ctx context.Context, dbContext persistence.IDbContext, links []*models.Link) error
	GetByNoteID(ctx context.Context, dbContext persistence.IDbContext, fromNoteID string) ([]*models.Link, error)
	GetBySlug(ctx context.Context, dbContext persistence.IDbContext, targetSlug string) ([]*models.Link, error)
	Search(ctx context.Context, dbContext persistence.IDbContext, query string) ([]*models.Link, error)
	Delete(ctx context.Context, dbContext persistence.IDbContext, fromNoteID string) error
}

type linkRepository struct{}

func NewLinkRepository() ILinkRepository {
	return &linkRepository{}
}

func (r *linkRepository) Insert(ctx context.Context, dbContext persistence.IDbContext, links []*models.Link) error {
	if len(links) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(links))
	args := make([]any, 0, len(links))

	for _, link := range links {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
		args = append(args, link.ID(), link.FromNoteID(), link.TargetSlug(), link.Raw(), link.Display(), link.Line(), link.Col())
	}

	query := "INSERT OR IGNORE INTO link(id, from_note_id, target_slug, raw, display, line, col) VALUES " + strings.Join(placeholders, ",")

	_, err := dbContext.ExecContext(ctx, query, args...)
	return err
}

func (r *linkRepository) GetByNoteID(ctx context.Context, dbContext persistence.IDbContext, fromNoteID string) ([]*models.Link, error) {
	rows, err := dbContext.QueryContext(ctx, "SELECT id, from_note_id, target_slug, raw, display, line, col FROM link WHERE from_note_id = ?", fromNoteID)
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

func (r *linkRepository) GetBySlug(ctx context.Context, dbContext persistence.IDbContext, targetSlug string) ([]*models.Link, error) {
	rows, err := dbContext.QueryContext(ctx, `SELECT id, from_note_id, target_slug, raw, display, line, col FROM link WHERE target_slug = ?`, targetSlug)
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

func (r *linkRepository) Search(ctx context.Context, dbContext persistence.IDbContext, query string) ([]*models.Link, error) {
	rows, err := dbContext.QueryContext(ctx, `SELECT id, from_note_id, target_slug, raw, display, line, col FROM link WHERE raw LIKE ?`, "%"+query+"%")
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

func (r *linkRepository) Delete(ctx context.Context, dbContext persistence.IDbContext, fromNoteID string) error {
	_, err := dbContext.ExecContext(ctx, "DELETE FROM link WHERE from_note_id = ?", fromNoteID)
	return err
}
