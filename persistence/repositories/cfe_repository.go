package repositories

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type ICfeRepository interface {
	InsertRange(ctx context.Context, dbContext persistence.IDbContext, cfe []*models.CustomFronMatter) error
}

type cfeRepository struct {
	readDBContext persistence.ReadContext
}

func NewCfeRepository(rdb persistence.ReadContext) ICfeRepository {
	return &cfeRepository{readDBContext: rdb}
}

func (r *cfeRepository) InsertRange(
	ctx context.Context,
	dbCtx persistence.IDbContext,
	cfe []*models.CustomFronMatter,
) error {
	statement, err := dbCtx.Prepare(
		`INSERT INTO custom_frontmatter (note_id, key, value)
		VALUES (?, ?, ?) ON CONFLICT DO NOTHING;`,)
	if err != nil {
		return err
	}

	for _, c := range cfe {
		if _, err := statement.ExecContext(ctx, c.NodeID(), c.Key(), c.Value()); err != nil {
			return err
		}
	}

	return nil
}
