package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmailNotFound      = errors.New("email not found")
)

func (s *Store) SignUpUser(ctx context.Context, email, password, name string) (uuid.UUID, error) {

	userID := uuid.New()
	hashPassword, err := HashPassword(password)
	if err != nil {
		return uuid.Nil, err
	}

	user, err := s.Query.CreateUser(ctx, db.CreateUserParams{
		ID:       userID,
		Email:    email,
		Password: &hashPassword,
		Name:     &name,
	})
	if err != nil {
		if IsUniqueViolation(err) {
			return uuid.Nil, ErrEmailAlreadyExists
		}
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (s *Store) LoginUser(ctx context.Context, email, password string) (uuid.UUID, string, error) {

	return uuid.Nil, "", nil
}
