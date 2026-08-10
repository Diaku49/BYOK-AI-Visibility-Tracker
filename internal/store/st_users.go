package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmailNotFound      = errors.New("email not found")
	ErrPasswordNotSet     = errors.New("password not set")
)

func (s *Store) CreateUser(ctx context.Context, email, password, name string) (uuid.UUID, error) {

	userID := uuid.New()
	hashPassword, err := pkg.HashPassword(password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("hash user password: %w", err)
	}

	user, err := s.query.CreateUser(ctx, db.CreateUserParams{
		ID:       userID,
		Email:    email,
		Password: &hashPassword,
		Name:     &name,
	})
	if err != nil {
		if IsUniqueViolation(err) {
			return uuid.Nil, ErrEmailAlreadyExists
		}
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return user.ID, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (uuid.UUID, string, string, error) {
	user, err := s.query.GetUserByEmail(ctx, db.GetUserByEmailParams{Email: email})
	if err != nil {
		if IsNotFound(err) {
			return uuid.Nil, "", "", ErrEmailNotFound
		}
		return uuid.Nil, "", "", fmt.Errorf("get user by email: %w", err)
	}
	if user.Password == nil {
		return uuid.Nil, "", "", ErrPasswordNotSet
	}

	return user.ID, user.TierName, *user.Password, nil
}
