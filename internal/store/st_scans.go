package store

import (
	"context"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
)

func (s *Store) GetScansForWorkers(ctx context.Context) ([]db.GetScansForWorkersRow, error) {
	return s.query.GetScansForWorkers(ctx)
}

func (s *Store) GetScansForAnalysis(ctx context.Context) ([]db.GetScansForAnalysisRow, error) {
	return s.query.GetScansForAnalysis(ctx)
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
