package repositories

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type IIndexRepository interface {
	WipeIndex(ctx context.Context, dbContext persistence.IDbContext) error
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
