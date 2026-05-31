package repositories

import (
	"context"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
)

type ICfRepository interface {
	InsertRange(ctx context.Context, dbContext persistence.IDbContext, cfe []*models.CustomFronMatter) error
	Search(ctx context.Context, key string, value string) ([]*models.Note, error)
	Delete(ctx context.Context, dbContext persistence.IDbContext, noteID string) error
}

type cfeRepository struct {
	readDBContext persistence.ReadContext
}

func NewCfeRepository(rdb persistence.ReadContext) ICfRepository {
	return &cfeRepository{readDBContext: rdb}
}

func (r *cfeRepository) InsertRange(
	ctx context.Context,
	dbCtx persistence.IDbContext,
	cfe []*models.CustomFronMatter,
) error {
	statement, err := dbCtx.Prepare(
		`INSERT INTO custom_frontmatter (note_id, key, value)
		VALUES (?, ?, ?) ON CONFLICT DO NOTHING;`)
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

func (r *cfeRepository) Delete(
	ctx context.Context,
	dbCtx persistence.IDbContext,
	noteID string,
) error {
	statement, err := dbCtx.Prepare(`DELETE FROM custom_frontmatter WHERE note_id = ?;`)
	if err != nil {
		return err
	}

	if _, err := statement.ExecContext(ctx, noteID); err != nil {
		return err
	}

	return nil
}

func (r *cfeRepository) Search(
	ctx context.Context,
	key string,
	value string,
) ([]*models.Note, error) {
	rows, err := r.readDBContext.QueryContext(
		ctx,
		`SELECT n.id, n.path, n.title, n.slug, n.created_at, n.updated_at
		FROM note n
		INNER JOIN custom_frontmatter cfe ON n.id = cfe.note_id
		WHERE cfe.key = ? AND cfe.value LIKE ?;`,
		key,
		"%"+value+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var id, title, path, slug, createdAt, updatedAT string
		if err := rows.Scan(&id, &title, &path, &slug, &createdAt, &updatedAT); err != nil {
			return nil, err
		}
		notes = append(notes, models.NewNote(id, title, path, slug, createdAt, updatedAT))
	}

	return notes, nil
}
