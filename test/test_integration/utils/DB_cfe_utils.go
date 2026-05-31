package utils

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

func CreateCfe(
	ctx context.Context,
	dbCtx persistence.IDbContext,
	noteID, key string,
	values []string,
) error {
	for _, value := range values {
		cfe := models.NewCustomFrontMatter(noteID, key, value)
		statement, err := dbCtx.Prepare(
			`INSERT INTO custom_frontmatter (note_id, key, value)
			VALUES (?, ?, ?) ON CONFLICT DO NOTHING;`)
		if err != nil {
			return err
		}

		if _, err := statement.ExecContext(ctx, cfe.NodeID(), cfe.Key(), cfe.Value()); err != nil {
			return err
		}
	}
	return nil
}


