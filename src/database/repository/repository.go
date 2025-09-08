package repository

import (
	"context"
	"database/sql"
	"strings"
)

const (
	DefaultDateFormat = "2006-01-02_15:04:05"
)

type ConnectionWithContext interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Connection interface {
	Exec(query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func inClause[T any](ids []T, extraArgs ...any) (string, []any) {
	args := make([]any, len(ids)+len(extraArgs))
	ph := make([]string, len(ids))

	if len(extraArgs) > 0 {
		copy(args, extraArgs)
	}

	for i, v := range ids {
		args[i+len(extraArgs)] = v
		ph[i] = "?"
	}

	return strings.Join(ph, ","), args
}
