package persistence

import (
	"context"
	"database/sql"
)

type ReadContext struct{}

func NewReadContext() *ReadContext {
	return &ReadContext{}
}

func (r *ReadContext) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return GetDBContext().DB.ExecContext(ctx, query, args...)
}

func (r *ReadContext) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return GetDBContext().DB.QueryContext(ctx, query, args...)
}

func (r *ReadContext) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return GetDBContext().DB.QueryRowContext(ctx, query, args...)
}
