package store

import (
	"context"
	"errors"
	"fmt"

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
		return db.Prompt{}, fmt.Errorf("get prompt for user: %w", err)
	}

	return prompt, nil
}

func (s *Store) ListPromptsByProjectForUser(
	ctx context.Context,
	projectID, userID uuid.UUID,
) ([]db.Prompt, error) {
	prompts, err := s.query.ListPromptsByProjectForUser(ctx, db.ListPromptsByProjectForUserParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("list prompts for project and user: %w", err)
	}

	return prompts, nil
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
		return db.Prompt{}, fmt.Errorf("update prompt for user: %w", err)
	}

	return prompt, nil
}
