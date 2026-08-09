package store

import (
	"context"
	"errors"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNoActivePrompts        = errors.New("project has no active prompts")
	ErrNoActiveProjectEngines = errors.New("project has no active provider keys")
)

func (s *Store) CreateScan(ctx context.Context, userID, projectID uuid.UUID) (uuid.UUID, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)

	// Get the project configuration that this scan will execute.
	projectEngines, err := query.ListActiveProjectEnginesForScan(ctx, db.ListActiveProjectEnginesForScanParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		return uuid.Nil, err
	}
	if len(projectEngines) == 0 {
		return uuid.Nil, ErrNoActiveProjectEngines
	}

	promptIDs, err := query.ListActivePromptIDsForScan(ctx, db.ListActivePromptIDsForScanParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		return uuid.Nil, err
	}
	if len(promptIDs) == 0 {
		return uuid.Nil, ErrNoActivePrompts
	}

	totalRuns := len(projectEngines) * len(promptIDs)
	scanID := uuid.New()
	if _, err := query.CreateScan(ctx, db.CreateScanParams{
		ID:            scanID,
		ProjectID:     projectID,
		UserID:        userID,
		Status:        "pending",
		TotalRuns:     int32(totalRuns),
		CompletedRuns: 0,
		FailedRuns:    0,
	}); err != nil {
		if IsNotFound(err) {
			return uuid.Nil, ErrProjectNotFound
		}
		return uuid.Nil, err
	}

	// Data for creating scan runs in batch
	runIDs := make([]uuid.UUID, 0, totalRuns)
	engineIDs := make([]string, 0, totalRuns)
	runPromptIDs := make([]uuid.UUID, 0, totalRuns)
	providerKeyIDs := make([]uuid.UUID, 0, totalRuns)
	for _, projectEngine := range projectEngines {
		for _, promptID := range promptIDs {
			runIDs = append(runIDs, uuid.New())
			engineIDs = append(engineIDs, projectEngine.EngineID)
			runPromptIDs = append(runPromptIDs, promptID)
			providerKeyIDs = append(providerKeyIDs, projectEngine.ProviderKeyID)
		}
	}

	rowsAffected, err := query.CreateScanRuns(ctx, db.CreateScanRunsParams{
		Ids:            runIDs,
		ScanID:         scanID,
		EngineIds:      engineIDs,
		PromptIds:      runPromptIDs,
		ProviderKeyIds: providerKeyIDs,
	})
	if err != nil {
		return uuid.Nil, err
	}
	if rowsAffected != int64(totalRuns) {
		return uuid.Nil, errors.New("failed to create all scan runs")
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return scanID, nil
}

func (s *Store) UpdateScanStateByID(ctx context.Context, scanID uuid.UUID, status string, errorMsg *string) (db.Scan, error) {
	return s.query.UpdateScanStateByID(ctx, db.UpdateScanStateByIDParams{
		ID:     scanID,
		Status: status,
		Error:  errorMsg,
	})
}

func (s *Store) GetScansForWorkers(ctx context.Context) ([]db.GetEligibleScanRunsForWorkersRow, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)
	scanIDs, err := query.ClaimScansForWorkers(ctx)
	if err != nil {
		return nil, err
	}

	if len(scanIDs) == 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}

	runs, err := query.GetEligibleScanRunsForWorkers(ctx, db.GetEligibleScanRunsForWorkersParams{
		ScanIds: scanIDs,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return runs, nil
}

func (s *Store) GetScansForAnalysis(ctx context.Context) ([]db.GetScansForAnalysisRow, error) {
	return s.query.GetScansForAnalysis(ctx)
}

func (s *Store) ClaimScanForAnalysis(ctx context.Context, scanID uuid.UUID) (uuid.UUID, error) {
	return s.query.ClaimScanForAnalysis(ctx, db.ClaimScanForAnalysisParams{ID: scanID})
}

func (s *Store) UpdateScanByID(ctx context.Context, arg db.UpdateScanParams) (db.Scan, error) {
	return s.query.UpdateScan(ctx, arg)
}

func (s *Store) IncrementScanCompletedRuns(ctx context.Context, scanID uuid.UUID) error {
	return s.query.IncrementScanCompletedRuns(ctx, db.IncrementScanCompletedRunsParams{ID: scanID})
}

func (s *Store) IncrementScanFailedRuns(ctx context.Context, scanID uuid.UUID) error {
	return s.query.IncrementScanFailedRuns(ctx, db.IncrementScanFailedRunsParams{ID: scanID})
}
