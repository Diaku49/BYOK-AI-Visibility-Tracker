package worker

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	pollInterval    = 10 * time.Second
	attemptInterval = 5 * time.Second
)

func (wc *WorkerCoordinator) StartScanWorker(c chan *ScanRunTask) {
	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case j := <-c:
			{
				// Do ScanRunTask -- needs to be refactored
				scanResponse := wc.ExecuteScanRun(j, 2, attemptInterval)
				wc.l.Info("Scan ran", "ScanID", scanResponse.scanRunID)
			}
		case <-ticker.C:
			{
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				if err := wc.GetWork(ctx); err != nil {
					wc.l.Error("Failed getting job", "Error", err.Error())
				}
				cancel()
				continue
			}
		}
	}
}

// Need to implement storing in db
func (wc *WorkerCoordinator) ExecuteScanRun(j *ScanRunTask, retryAttempt int, retryInterval time.Duration) *ScanRunResult {
	var scanResponse ScanRunResult
	ctx := context.Background()

	// changing scan run state with store method
	if _, err := wc.st.UpdateStateScanRunByID(ctx, j.ScanRunID, "running"); err != nil {
		wc.l.Error("failed to update scan run state to running", "RunID", j.ScanRunID, "error", err)
	}

	for ; retryAttempt > 0; retryAttempt-- {
		runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		p := wc.providerRegistry[j.EngineID]
		result, err := p.Run(runCtx, j.APIKey, nil, j.Request)
		cancel()

		if err != nil {
			wc.l.Error(err.Error(), "ScanID", j.ScanID, "RunID", j.ScanRunID, "Engine", j.EngineID)
			errMsg := err.Error()
			scanResponse = ScanRunResult{
				scanRunID: j.ScanRunID,
				error:     errMsg,
			}

			if !isRetryable(err) || retryAttempt == 1 {
				// Terminal failure — persist failed state
				if _, stErr := wc.st.UpdateStateScanRunByID(ctx, j.ScanRunID, "failed"); stErr != nil {
					wc.l.Error("failed to update scan run state to failed", "RunID", j.ScanRunID, "error", stErr)
				}
				if stErr := wc.st.IncrementScanFailedRuns(ctx, j.ScanID); stErr != nil {
					wc.l.Error("failed to increment scan failed runs", "ScanID", j.ScanID, "error", stErr)
				}
				break
			}

			select {
			case <-runCtx.Done():
				return &scanResponse
			case <-time.After(retryInterval):
				continue
			}
		}

		scanResponse = ScanRunResult{
			scanRunID: j.ScanRunID,
			result:    *result,
			error:     "",
		}

		// Marshal raw response for storage
		rawJSON, _ := json.Marshal(result.Raw)

		// Update the scan run with results
		answerText := result.AnswerText
		now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		if _, stErr := wc.st.UpdateScanRunByID(ctx, db.UpdateScanRunParams{
			ID:            j.ScanRunID,
			ScanID:        j.ScanID,
			EngineID:      j.EngineID,
			PromptID:      j.PromptID,
			ProviderKeyID: j.ProviderKeyID,
			TryNumber:     j.TryNumber,
			Status:      "completed",
			AnswerText:  &answerText,
			RawResponse: rawJSON,
			FinishedAt:  now,
		}); stErr != nil {
			wc.l.Error("failed to update scan run result", "RunID", j.ScanRunID, "error", stErr)
		}

		wc.l.Info("scan run completed", "scan_id", j.ScanID, "run_id", j.ScanRunID, "engine", j.EngineID)
		if stErr := wc.st.IncrementScanCompletedRuns(ctx, j.ScanID); stErr != nil {
			wc.l.Error("failed to increment scan completed runs", "ScanID", j.ScanID, "error", stErr)
		}
		break
	}

	return &scanResponse
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}

	errText := err.Error()
	return strings.Contains(errText, "429") ||
		strings.Contains(errText, "503") ||
		strings.Contains(errText, "504") ||
		strings.Contains(errText, "RESOURCE_EXHAUSTED") ||
		strings.Contains(errText, "UNAVAILABLE") ||
		strings.Contains(errText, "DEADLINE_EXCEEDED")
}
