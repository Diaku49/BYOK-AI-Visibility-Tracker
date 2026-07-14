package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const UniqueViolationErrorCode = "23505"

type Store struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func Connect(ctx context.Context, dbURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Store{
		pool:  pool,
		query: db.New(pool),
	}, nil
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == UniqueViolationErrorCode
	}
	return false
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
