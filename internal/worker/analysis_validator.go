package worker

import (
	"fmt"

	a "github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
	"github.com/google/uuid"
)

func validateAnalysisRuns(expected []a.RunForAnalysis, returned []a.RunAnalysisResult) error {
	if len(returned) != len(expected) {
		return fmt.Errorf("analyzer returned %d run results, expected %d", len(returned), len(expected))
	}

	expectedIDs := make(map[uuid.UUID]struct{}, len(expected))
	for _, run := range expected {
		expectedIDs[run.ScanRunID] = struct{}{}
	}

	seenIDs := make(map[uuid.UUID]struct{}, len(returned))
	for _, result := range returned {
		if _, ok := expectedIDs[result.ScanRunID]; !ok {
			return fmt.Errorf("analyzer returned unknown scan run ID: %s", result.ScanRunID)
		}

		if _, alreadySeen := seenIDs[result.ScanRunID]; alreadySeen {
			return fmt.Errorf("analyzer returned duplicate scan run ID: %s", result.ScanRunID)
		}

		seenIDs[result.ScanRunID] = struct{}{}
	}

	return nil
}
