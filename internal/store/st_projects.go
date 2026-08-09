package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/Diaku49/AI-visibility-tracker/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrProjectNotFound = errors.New("project not found")

func (s *Store) CreateProject(
	ctx context.Context,
	userID uuid.UUID,
	project dto.CreateProject,
) (uuid.UUID, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)

	// making project
	ID := uuid.New()
	projectID, err := query.CreateProject(ctx, db.CreateProjectParams{
		ID:        ID,
		UserID:    userID,
		BrandName: project.BrandName,
		Domain:    project.Domain,
		Language:  project.Language,
		Region:    project.Region,
	})
	if err != nil {
		return uuid.Nil, err
	}

	// check Engine IDs and ProviderKeys if this key belongs to this engineID and userID
	for _, provider := range project.Providers {
		providerKeyID, err := uuid.Parse(provider.ProviderKeyID)
		if err != nil {
			return uuid.Nil, err
		}

		if _, err := query.GetProviderKeyByIDForUserAndEngine(ctx, db.GetProviderKeyByIDForUserAndEngineParams{
			ID:       providerKeyID,
			UserID:   userID,
			EngineID: provider.EngineID,
		}); err != nil {
			if IsNotFound(err) {
				return uuid.Nil, ErrProviderKeyNotFound
			}
			return uuid.Nil, err
		}

		if _, err := query.CreateProjectEngineForUser(ctx, db.CreateProjectEngineForUserParams{
			ProjectID:     projectID,
			EngineID:      provider.EngineID,
			ProviderKeyID: providerKeyID,
			UserID:        userID,
		}); err != nil {
			return uuid.Nil, err
		}
	}

	// making prompts
	for _, prompt := range project.Prompts {
		if _, err := query.CreatePromptForUser(ctx, db.CreatePromptForUserParams{
			ID:        uuid.New(),
			ProjectID: projectID,
			Text:      prompt,
			Active:    true,
			UserID:    userID,
		}); err != nil {
			return uuid.Nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return projectID, nil
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
