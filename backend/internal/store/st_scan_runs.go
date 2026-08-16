package store

import (
	"context"
	"fmt"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Store) UpdateStateScanRunByID(ctx context.Context, scanRunID uuid.UUID, state string, errMsg *string) (db.ScanRun, error) {
	scanRun, err := s.query.UpdateScanRunStateByID(ctx, db.UpdateScanRunStateByIDParams{
		ID:     scanRunID,
		Status: state,
		Error:  errMsg,
	})
	if err != nil {
		return db.ScanRun{}, fmt.Errorf("update scan run state: %w", err)
	}

	return scanRun, nil
}

func (s *Store) ClaimScanRun(ctx context.Context, scanRunID, scanID uuid.UUID) (db.ScanRun, error) {
	scanRun, err := s.query.ClaimScanRun(ctx, db.ClaimScanRunParams{
		ID:     scanRunID,
		ScanID: scanID,
	})
	if err != nil {
		return db.ScanRun{}, fmt.Errorf("claim scan run: %w", err)
	}

	return scanRun, nil
}

func (s *Store) UpdateScanRunByID(ctx context.Context, arg db.UpdateScanRunParams) (db.ScanRun, error) {
	scanRun, err := s.query.UpdateScanRun(ctx, arg)
	if err != nil {
		return db.ScanRun{}, fmt.Errorf("update scan run: %w", err)
	}

	return scanRun, nil
}

func (s *Store) MarkScanRunCompleted(ctx context.Context, scanRun db.UpdateScanRunParams) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin mark scan run completed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)
	if _, err := query.UpdateScanRun(ctx, scanRun); err != nil {
		return fmt.Errorf("update completed scan run: %w", err)
	}

	if err := query.IncrementScanCompletedRuns(ctx, db.IncrementScanCompletedRunsParams{ID: scanRun.ScanID}); err != nil {
		return fmt.Errorf("increment completed scan runs: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit mark scan run completed transaction: %w", err)
	}

	return nil
}

func (s *Store) MarkScanRunFailed(ctx context.Context, scanRunID, scanID uuid.UUID, errMsg *string) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin mark scan run failed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)
	if _, err := query.UpdateScanRunStateByID(ctx, db.UpdateScanRunStateByIDParams{
		ID:     scanRunID,
		Status: "failed",
		Error:  errMsg,
	}); err != nil {
		return fmt.Errorf("update failed scan run: %w", err)
	}

	if err := query.IncrementScanFailedRuns(ctx, db.IncrementScanFailedRunsParams{ID: scanID}); err != nil {
		return fmt.Errorf("increment failed scan runs: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit mark scan run failed transaction: %w", err)
	}

	return nil
}

func (s *Store) UpdateScanAnalysis(ctx context.Context, batch []db.UpdateScanRunAnalysisParams, scan db.UpdateScanParams) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin update scan analysis transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)

	for _, scanRun := range batch {
		scanRun.ScanID = scan.ID
		if err := query.UpdateScanRunAnalysis(ctx, scanRun); err != nil {
			return fmt.Errorf("update scan run analysis: %w", err)
		}
	}

	if _, err := query.UpdateScan(ctx, scan); err != nil {
		return fmt.Errorf("update scan analysis summary: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update scan analysis transaction: %w", err)
	}

	return nil
}

func (s *Store) UpdateScanRunAnalysis(ctx context.Context, arg db.UpdateScanRunAnalysisParams) error {
	if err := s.query.UpdateScanRunAnalysis(ctx, arg); err != nil {
		return fmt.Errorf("update scan run analysis: %w", err)
	}

	return nil
}
