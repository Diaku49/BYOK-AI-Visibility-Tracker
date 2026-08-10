package store

import (
	"context"
	"errors"
	"fmt"

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
		return uuid.Nil, fmt.Errorf("begin create scan transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)

	// Get the project configuration that this scan will execute.
	projectEngines, err := query.ListActiveProjectEnginesForScan(ctx, db.ListActiveProjectEnginesForScanParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("list active project engines for scan: %w", err)
	}
	if len(projectEngines) == 0 {
		return uuid.Nil, ErrNoActiveProjectEngines
	}

	promptIDs, err := query.ListActivePromptIDsForScan(ctx, db.ListActivePromptIDsForScanParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("list active prompts for scan: %w", err)
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
		return uuid.Nil, fmt.Errorf("create scan: %w", err)
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
		return uuid.Nil, fmt.Errorf("create scan runs: %w", err)
	}
	if rowsAffected != int64(totalRuns) {
		return uuid.Nil, errors.New("failed to create all scan runs")
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit create scan transaction: %w", err)
	}

	return scanID, nil
}

func (s *Store) UpdateScanStateByID(ctx context.Context, scanID uuid.UUID, status string, errorMsg *string) (db.Scan, error) {
	scan, err := s.query.UpdateScanStateByID(ctx, db.UpdateScanStateByIDParams{
		ID:     scanID,
		Status: status,
		Error:  errorMsg,
	})
	if err != nil {
		return db.Scan{}, fmt.Errorf("update scan state: %w", err)
	}

	return scan, nil
}

func (s *Store) GetScansForWorkers(ctx context.Context) ([]db.GetEligibleScanRunsForWorkersRow, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin get worker scans transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)
	scanIDs, err := query.ClaimScansForWorkers(ctx)
	if err != nil {
		return nil, fmt.Errorf("claim scans for workers: %w", err)
	}

	if len(scanIDs) == 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("commit empty worker scan transaction: %w", err)
		}
		return nil, nil
	}

	runs, err := query.GetEligibleScanRunsForWorkers(ctx, db.GetEligibleScanRunsForWorkersParams{
		ScanIds: scanIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("get eligible scan runs for workers: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit get worker scans transaction: %w", err)
	}

	return runs, nil
}

func (s *Store) GetScansForAnalysis(ctx context.Context) ([]db.GetScansForAnalysisRow, error) {
	rows, err := s.query.GetScansForAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("get scans for analysis: %w", err)
	}

	return rows, nil
}

func (s *Store) ClaimScanForAnalysis(ctx context.Context, scanID uuid.UUID) (uuid.UUID, error) {
	claimedScanID, err := s.query.ClaimScanForAnalysis(ctx, db.ClaimScanForAnalysisParams{ID: scanID})
	if err != nil {
		return uuid.Nil, fmt.Errorf("claim scan for analysis: %w", err)
	}

	return claimedScanID, nil
}

func (s *Store) UpdateScanByID(ctx context.Context, arg db.UpdateScanParams) (db.Scan, error) {
	scan, err := s.query.UpdateScan(ctx, arg)
	if err != nil {
		return db.Scan{}, fmt.Errorf("update scan: %w", err)
	}

	return scan, nil
}

func (s *Store) IncrementScanCompletedRuns(ctx context.Context, scanID uuid.UUID) error {
	if err := s.query.IncrementScanCompletedRuns(ctx, db.IncrementScanCompletedRunsParams{ID: scanID}); err != nil {
		return fmt.Errorf("increment completed scan runs: %w", err)
	}

	return nil
}

func (s *Store) IncrementScanFailedRuns(ctx context.Context, scanID uuid.UUID) error {
	if err := s.query.IncrementScanFailedRuns(ctx, db.IncrementScanFailedRunsParams{ID: scanID}); err != nil {
		return fmt.Errorf("increment failed scan runs: %w", err)
	}

	return nil
}
