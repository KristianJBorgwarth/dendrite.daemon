package utils

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

func CreateNote(
	ctx context.Context,
	dbCtx persistence.IDbContext,
	noteID,
	path,
	title,
	slug string,
) error {
	statement, err := dbCtx.Prepare(`
	INSERT INTO note (id, title, path, slug, created_at, updated_at)
	VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	ON CONFLICT DO NOTHING;`)
	if err != nil {
		return err
	}

	if _, err := statement.ExecContext(ctx, noteID, title, path, slug); err != nil {
		return err
	}

	return nil
}
