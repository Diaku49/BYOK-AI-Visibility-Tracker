package store

import (
	"context"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
)

func (s *Store) UpdateStateScanRunByID(ctx context.Context, scanRunID uuid.UUID, state string, errMsg *string) (db.ScanRun, error) {
	return s.query.UpdateScanRunStateByID(ctx, db.UpdateScanRunStateByIDParams{
		ID:     scanRunID,
		Status: state,
		Error:  errMsg,
	})
}

func (s *Store) UpdateScanRunByID(ctx context.Context, arg db.UpdateScanRunParams) (db.ScanRun, error) {
	return s.query.UpdateScanRun(ctx, arg)
}

func (s *Store) UpdateScanRunAnalysis(ctx context.Context, arg db.UpdateScanRunAnalysisParams) error {
	return s.query.UpdateScanRunAnalysis(ctx, arg)
}
