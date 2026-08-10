package store

import (
	"context"
	"errors"
	"fmt"

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
		return db.ProjectEngine{}, fmt.Errorf("create project engine: %w", err)
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
		return db.ProjectEngine{}, fmt.Errorf("update project engine: %w", err)
	}

	return projectEngine, nil
}

func (s *Store) DeleteProjectEngineForUser(
	ctx context.Context,
	projectID uuid.UUID,
	engineID string,
	userID uuid.UUID,
) error {
	if err := s.query.DeleteProjectEngineForUser(ctx, db.DeleteProjectEngineForUserParams{
		ProjectID: projectID,
		EngineID:  engineID,
		UserID:    userID,
	}); err != nil {
		return fmt.Errorf("delete project engine: %w", err)
	}

	return nil
}
