package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
)

var (
	ErrProviderKeyAlreadyExists = errors.New("provider key already exists")
	ErrProviderKeyNotFound      = errors.New("provider key not found")
)

func (s *Store) CreateProviderKey(
	ctx context.Context,
	userID uuid.UUID,
	engineID, name string,
	encryptedKey, keyNonce []byte,
	active bool,
	monthlyRunLimit *int32,
) (uuid.UUID, error) {
	providerKeyID := uuid.New()

	providerKey, err := s.query.CreateProviderKey(ctx, db.CreateProviderKeyParams{
		ID:              providerKeyID,
		UserID:          userID,
		EngineID:        engineID,
		Name:            name,
		EncryptedKey:    encryptedKey,
		KeyNonce:        keyNonce,
		Active:          active,
		MonthlyRunLimit: monthlyRunLimit,
	})
	if err != nil {
		if IsUniqueViolation(err) {
			return uuid.Nil, ErrProviderKeyAlreadyExists
		}
		return uuid.Nil, fmt.Errorf("create provider key: %w", err)
	}

	return providerKey.ID, nil
}

func (s *Store) GetProviderKeyByID(ctx context.Context, providerKeyID uuid.UUID) (db.ProviderKey, error) {
	providerKey, err := s.query.GetProviderKeyByID(ctx, db.GetProviderKeyByIDParams{ID: providerKeyID})
	if err != nil {
		if IsNotFound(err) {
			return db.ProviderKey{}, ErrProviderKeyNotFound
		}
		return db.ProviderKey{}, fmt.Errorf("get provider key: %w", err)
	}

	return providerKey, nil
}

func (s *Store) GetProviderKeyByIDForUser(
	ctx context.Context,
	providerKeyID, userID uuid.UUID,
) (db.ProviderKey, error) {
	providerKey, err := s.query.GetProviderKeyByIDForUser(ctx, db.GetProviderKeyByIDForUserParams{
		ID:     providerKeyID,
		UserID: userID,
	})
	if err != nil {
		if IsNotFound(err) {
			return db.ProviderKey{}, ErrProviderKeyNotFound
		}
		return db.ProviderKey{}, fmt.Errorf("get provider key for user: %w", err)
	}

	return providerKey, nil
}

func (s *Store) ListProviderKeysByUserID(ctx context.Context, userID uuid.UUID) ([]db.ProviderKey, error) {
	providerKeys, err := s.query.ListProviderKeysByUserID(ctx, db.ListProviderKeysByUserIDParams{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("list provider keys for user: %w", err)
	}

	return providerKeys, nil
}

func (s *Store) ListProviderKeysByUserIDAndEngine(
	ctx context.Context,
	userID uuid.UUID,
	engineID string,
) ([]db.ProviderKey, error) {
	providerKeys, err := s.query.ListProviderKeysByUserIDAndEngine(ctx, db.ListProviderKeysByUserIDAndEngineParams{
		UserID:   userID,
		EngineID: engineID,
	})
	if err != nil {
		return nil, fmt.Errorf("list provider keys for user and engine: %w", err)
	}

	return providerKeys, nil
}

func (s *Store) UpdateProviderKeyMetadataForUser(
	ctx context.Context,
	providerKeyID, userID uuid.UUID,
	name string,
	active bool,
	monthlyRunLimit *int32,
) (db.ProviderKey, error) {
	providerKey, err := s.query.UpdateProviderKeyMetadataForUser(ctx, db.UpdateProviderKeyMetadataForUserParams{
		ID:              providerKeyID,
		UserID:          userID,
		Name:            name,
		Active:          active,
		MonthlyRunLimit: monthlyRunLimit,
	})
	if err != nil {
		if IsNotFound(err) {
			return db.ProviderKey{}, ErrProviderKeyNotFound
		}
		return db.ProviderKey{}, fmt.Errorf("update provider key metadata: %w", err)
	}

	return providerKey, nil
}

func (s *Store) DeleteProviderKeyForUser(
	ctx context.Context,
	providerKeyID uuid.UUID,
	userID uuid.UUID,
) error {
	if err := s.query.DeleteProviderKeyForUser(ctx, db.DeleteProviderKeyForUserParams{
		ID:     providerKeyID,
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("delete provider key for user: %w", err)
	}

	return nil
}
