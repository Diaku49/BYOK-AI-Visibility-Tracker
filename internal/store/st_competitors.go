package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
)

var (
	ErrCompetitorAlreadyExists = errors.New("competitor already exists")
	ErrCompetitorNotFound      = errors.New("competitor not found")
)

func (s *Store) CreateCompetitorForUser(
	ctx context.Context,
	projectID uuid.UUID,
	name string,
	domain *string,
	userID uuid.UUID,
) (uuid.UUID, error) {
	competitorID := uuid.New()

	competitor, err := s.query.CreateCompetitorForUser(ctx, db.CreateCompetitorForUserParams{
		ID:        competitorID,
		ProjectID: projectID,
		Name:      name,
		Domain:    domain,
		UserID:    userID,
	})
	if err != nil {
		if IsUniqueViolation(err) {
			return uuid.Nil, ErrCompetitorAlreadyExists
		}
		if IsNotFound(err) {
			return uuid.Nil, ErrProjectNotFound
		}
		return uuid.Nil, err
	}

	return competitor.ID, nil
}

func (s *Store) ListCompetitorsByProjectForUser(
	ctx context.Context,
	projectID, userID uuid.UUID,
) ([]db.Competitor, error) {
	return s.query.ListCompetitorsByProjectForUser(ctx, db.ListCompetitorsByProjectForUserParams{
		ProjectID: projectID,
		UserID:    userID,
	})
}

func (s *Store) ListCompetitorsByProject(ctx context.Context, projectID uuid.UUID) ([]db.Competitor, error) {
	return s.query.ListCompetitorsByProject(ctx, db.ListCompetitorsByProjectParams{
		ProjectID: projectID,
	})
}

func (s *Store) DeleteCompetitorForUser(ctx context.Context, competitorID, userID uuid.UUID) error {
	rowsAffected, err := s.query.DeleteCompetitorForUser(ctx, db.DeleteCompetitorForUserParams{
		ID:     competitorID,
		UserID: userID,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCompetitorNotFound
	}

	return nil
}
