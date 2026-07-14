package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
)

var ErrProjectNotFound = errors.New("project not found")

func (s *Store) CreateProject(
	ctx context.Context,
	userID uuid.UUID,
	brandName, domain, language, region string,
) (uuid.UUID, error) {
	projectID := uuid.New()

	project, err := s.query.CreateProject(ctx, db.CreateProjectParams{
		ID:        projectID,
		UserID:    userID,
		BrandName: brandName,
		Domain:    domain,
		Language:  language,
		Region:    region,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return project.ID, nil
}

func (s *Store) GetProjectByID(ctx context.Context, projectID uuid.UUID) (db.Project, error) {
	project, err := s.query.GetProjectByID(ctx, db.GetProjectByIDParams{ID: projectID})
	if err != nil {
		if IsNotFound(err) {
			return db.Project{}, ErrProjectNotFound
		}
		return db.Project{}, err
	}

	return project, nil
}

func (s *Store) GetProjectByIDForUser(ctx context.Context, projectID, userID uuid.UUID) (db.Project, error) {
	project, err := s.query.GetProjectByIDForUser(ctx, db.GetProjectByIDForUserParams{
		ID:     projectID,
		UserID: userID,
	})
	if err != nil {
		if IsNotFound(err) {
			return db.Project{}, ErrProjectNotFound
		}
		return db.Project{}, err
	}

	return project, nil
}

func (s *Store) ListProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]db.Project, error) {
	return s.query.ListProjectsByUserID(ctx, db.ListProjectsByUserIDParams{UserID: userID})
}

func (s *Store) UpdateProject(
	ctx context.Context,
	projectID uuid.UUID,
	brandName, domain, language, region string,
) (db.Project, error) {
	project, err := s.query.UpdateProject(ctx, db.UpdateProjectParams{
		ID:        projectID,
		BrandName: brandName,
		Domain:    domain,
		Language:  language,
		Region:    region,
	})
	if err != nil {
		if IsNotFound(err) {
			return db.Project{}, ErrProjectNotFound
		}
		return db.Project{}, err
	}

	return project, nil
}

func (s *Store) UpdateProjectForUser(
	ctx context.Context,
	projectID, userID uuid.UUID,
	brandName, domain, language, region string,
) (db.Project, error) {
	project, err := s.query.UpdateProjectForUser(ctx, db.UpdateProjectForUserParams{
		ID:        projectID,
		UserID:    userID,
		BrandName: brandName,
		Domain:    domain,
		Language:  language,
		Region:    region,
	})
	if err != nil {
		if IsNotFound(err) {
			return db.Project{}, ErrProjectNotFound
		}
		return db.Project{}, err
	}

	return project, nil
}
