package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
)

var (
	ErrProjectEngineAlreadyExists = errors.New("project engine already exists")
	ErrProjectEngineNotFound      = errors.New("project engine not found")
)

func (s *Store) CreateProjectEngineForUser(
	ctx context.Context,
	projectID uuid.UUID,
	engineID string,
	providerKeyID, userID uuid.UUID,
) (db.ProjectEngine, error) {
	projectEngine, err := s.query.CreateProjectEngineForUser(ctx, db.CreateProjectEngineForUserParams{
		ProjectID:     projectID,
		EngineID:      engineID,
		ProviderKeyID: providerKeyID,
		UserID:        userID,
	})
	if err != nil {
		if IsUniqueViolation(err) {
			return db.ProjectEngine{}, ErrProjectEngineAlreadyExists
		}
		if IsNotFound(err) {
			return db.ProjectEngine{}, ErrProjectEngineNotFound
		}
		return db.ProjectEngine{}, err
	}

	return projectEngine, nil
}

func (s *Store) UpdateProjectEngineForUser(
	ctx context.Context,
	projectID uuid.UUID,
	engineID string,
	providerKeyID, userID uuid.UUID,
) (db.ProjectEngine, error) {
	projectEngine, err := s.query.UpdateProjectEngineForUser(ctx, db.UpdateProjectEngineForUserParams{
		ProviderKeyID: providerKeyID,
		ProjectID:     projectID,
		EngineID:      engineID,
		UserID:        userID,
	})
	if err != nil {
		if IsNotFound(err) {
			return db.ProjectEngine{}, ErrProjectEngineNotFound
		}
		return db.ProjectEngine{}, err
	}

	return projectEngine, nil
}

func (s *Store) DeleteProjectEngineForUser(
	ctx context.Context,
	projectID uuid.UUID,
	engineID string,
	userID uuid.UUID,
) error {
	return s.query.DeleteProjectEngineForUser(ctx, db.DeleteProjectEngineForUserParams{
		ProjectID: projectID,
		EngineID:  engineID,
		UserID:    userID,
	})
}
