package store

import (
	"context"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Store) UpdateStateScanRunByID(ctx context.Context, scanRunID uuid.UUID, state string, errMsg *string) (db.ScanRun, error) {
	return s.query.UpdateScanRunStateByID(ctx, db.UpdateScanRunStateByIDParams{
		ID:     scanRunID,
		Status: state,
		Error:  errMsg,
	})
}

func (s *Store) ClaimScanRun(ctx context.Context, scanRunID, scanID uuid.UUID) (db.ScanRun, error) {
	return s.query.ClaimScanRun(ctx, db.ClaimScanRunParams{
		ID:     scanRunID,
		ScanID: scanID,
	})
}

func (s *Store) UpdateScanRunByID(ctx context.Context, arg db.UpdateScanRunParams) (db.ScanRun, error) {
	return s.query.UpdateScanRun(ctx, arg)
}

func (s *Store) MarkScanRunCompleted(ctx context.Context, scanRun db.UpdateScanRunParams) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)
	if _, err := query.UpdateScanRun(ctx, scanRun); err != nil {
		return err
	}

	if err := query.IncrementScanCompletedRuns(ctx, db.IncrementScanCompletedRunsParams{ID: scanRun.ScanID}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) MarkScanRunFailed(ctx context.Context, scanRunID, scanID uuid.UUID, errMsg *string) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)
	if _, err := query.UpdateScanRunStateByID(ctx, db.UpdateScanRunStateByIDParams{
		ID:     scanRunID,
		Status: "failed",
		Error:  errMsg,
	}); err != nil {
		return err
	}

	if err := query.IncrementScanFailedRuns(ctx, db.IncrementScanFailedRunsParams{ID: scanID}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateScanAnalysis(ctx context.Context, batch []db.UpdateScanRunAnalysisParams, scan db.UpdateScanParams) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := s.query.WithTx(tx)

	for _, scanRun := range batch {
		scanRun.ScanID = scan.ID
		if err := query.UpdateScanRunAnalysis(ctx, scanRun); err != nil {
			return err
		}
	}

	if _, err := query.UpdateScan(ctx, scan); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateScanRunAnalysis(ctx context.Context, arg db.UpdateScanRunAnalysisParams) error {
	return s.query.UpdateScanRunAnalysis(ctx, arg)
}
