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

func (s *Store) CreatePromptForUser(
	ctx context.Context,
	projectID uuid.UUID,
	text string,
	active bool,
	userID uuid.UUID,
) (uuid.UUID, error) {
	promptID := uuid.New()

	prompt, err := s.query.CreatePromptForUser(ctx, db.CreatePromptForUserParams{
		ID:        promptID,
		ProjectID: projectID,
		Text:      text,
		Active:    active,
		UserID:    userID,
	})
	if err != nil {
		if IsNotFound(err) {
			return uuid.Nil, ErrProjectNotFound
		}
		return uuid.Nil, err
	}

	return prompt.ID, nil
}

func (s *Store) CreatePromptsForUser(
	ctx context.Context,
	projectID uuid.UUID,
	inputs []CreatePromptInput,
	userID uuid.UUID,
) ([]uuid.UUID, error) {
	if len(inputs) == 0 {
		return []uuid.UUID{}, nil
	}

	ids := make([]uuid.UUID, len(inputs))
	texts := make([]string, len(inputs))
	actives := make([]bool, len(inputs))
	for i, input := range inputs {
		ids[i] = uuid.New()
		texts[i] = input.Text
		actives[i] = input.Active
	}

	rowsAffected, err := s.query.CreatePromptsForUser(ctx, db.CreatePromptsForUserParams{
		Ids:       ids,
		ProjectID: projectID,
		Texts:     texts,
		Actives:   actives,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrProjectNotFound
	}

	return ids, nil
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
