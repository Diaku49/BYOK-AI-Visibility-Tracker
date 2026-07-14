package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	UniqueViolationErrorCode = "23505"
)

type Store struct {
	Pool  *pgxpool.Pool
	Query *db.Queries
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
		Pool:  pool,
		Query: db.New(pool),
	}, nil
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == UniqueViolationErrorCode
	}
	return false
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
