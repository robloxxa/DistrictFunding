package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func QueryOneRowToAddrStruct[T any](ctx context.Context, db *pgxpool.Pool, query string, arguments ...any) (*T, error) {
	rows, err := db.Query(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}

	t, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[T])
	if err != nil {
		return nil, err
	}

	return t, nil
}

func Exec(ctx context.Context, db *pgxpool.Pool, query string, arguments ...any) error {
	_, err := db.Exec(ctx, query, arguments)
	return err
}
