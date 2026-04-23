package repositories

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type IIndexRepository interface {
	GetBuildFlag(ctx context.Context) (bool, error)
	WipeIndex(ctx context.Context, dbContext persistence.IDbContext) error
	SetBuildFlag(ctx context.Context, dbContext persistence.IDbContext, rebuilding bool) error
}

type indexRepository struct {
	readDBContext persistence.ReadContext
}

func NewIndexRepository(rdb persistence.ReadContext) IIndexRepository {
	return &indexRepository{readDBContext: rdb}
}

func (r *indexRepository) WipeIndex(ctx context.Context, dbContext persistence.IDbContext) error {
	cmd := `DELETE FROM note; 
	DELETE FROM tag;
	DELETE FROM note_tag;
	DELETE FROM link;`

	_, err := dbContext.ExecContext(ctx, cmd)
	return err
}

func (r *indexRepository) GetBuildFlag(ctx context.Context) (bool, error) {
	cmd := `SELECT rebuilding, rebuilding_at FROM index_build_flag WHERE id = 'X'`

	row := r.readDBContext.QueryRowContext(ctx, cmd)

	var rebuilding bool
	err := row.Scan(&rebuilding)
	if err != nil {
		return true, err
	}

	return rebuilding, nil
}

func (r *indexRepository) SetBuildFlag(ctx context.Context, dbContext persistence.IDbContext, rebuilding bool) error {
	cmd := `INSERT INTO index_build_flag (id, rebuilding)` +
		`VALUES ('X', ?)` +
		`ON CONFLICT (id) DO UPDATE SET rebuilding = EXCLUDED.rebuilding, rebuilding_at = EXCLUDED.rebuilding_at`

	_, err := dbContext.ExecContext(ctx, cmd, rebuilding)
	return err
}
