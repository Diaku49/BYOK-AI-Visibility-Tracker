package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
)

var ErrPromptNotFound = errors.New("prompt not found")

type CreatePromptInput struct {
	Text   string
	Active bool
}

func (s *Store) GetPromptByIDForUser(ctx context.Context, promptID, userID uuid.UUID) (db.Prompt, error) {
	prompt, err := s.query.GetPromptByIDForUser(ctx, db.GetPromptByIDForUserParams{
		ID:     promptID,
		UserID: userID,
	})
	if err != nil {
		if IsNotFound(err) {
			return db.Prompt{}, ErrPromptNotFound
		}
		return db.Prompt{}, err
	}

	return prompt, nil
}

func (s *Store) ListPromptsByProjectForUser(
	ctx context.Context,
	projectID, userID uuid.UUID,
) ([]db.Prompt, error) {
	return s.query.ListPromptsByProjectForUser(ctx, db.ListPromptsByProjectForUserParams{
		ProjectID: projectID,
		UserID:    userID,
	})
}

func (s *Store) UpdatePromptForUser(
	ctx context.Context,
	promptID uuid.UUID,
	text string,
	active bool,
	userID uuid.UUID,
) (db.Prompt, error) {
	prompt, err := s.query.UpdatePromptForUser(ctx, db.UpdatePromptForUserParams{
		Text:   text,
		Active: active,
		ID:     promptID,
		UserID: userID,
	})
	if err != nil {
		if IsNotFound(err) {
			return db.Prompt{}, ErrPromptNotFound
		}
		return db.Prompt{}, err
	}

	return prompt, nil
}
