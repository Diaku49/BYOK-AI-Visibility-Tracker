package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	pollInterval    = 10 * time.Second
	attemptInterval = 5 * time.Second
)
var (
	errMessageClaimScanRun         = "failed to claim scan run"
	errMessageDecryptProviderKey   = "failed to decrypt provider key"
	errMessageUndefinedEngine      = "Undefined Engine"
	errMessageMarshalRawResponse   = "failed to marshal raw response"
	errMessageMarkScanRunCompleted = "failed to mark scan run as completed"
	errMessageMarkScanRunFailed    = "failed to mark scan run as failed"
)

func (wc *WorkerCoordinator) StartScanWorker(c chan *ScanRunTask) {
	for task := range c {
		scanResponse := wc.ExecuteScanRun(task, 2, attemptInterval)
		if scanResponse.err != nil {
			if errors.Is(scanResponse.err, pgx.ErrNoRows) {
				wc.l.Info("scan run already claimed", "scan_id", task.ScanID, "run_id", task.ScanRunID)
				continue
			}

			wc.l.Error("execute scan run", "scan_id", task.ScanID, "run_id", task.ScanRunID, "error", scanResponse.err)
			continue
		}

		for _, warning := range scanResponse.warnings {
			wc.l.Warn("scan run completed with warning", "scan_id", task.ScanID, "run_id", task.ScanRunID, "error", warning)
		}

		wc.l.Info("scan run completed", "scan_id", task.ScanID, "run_id", task.ScanRunID, "engine", task.EngineID)
	}
}

func (wc *WorkerCoordinator) ScanTaskProducer() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		if count, err := wc.GetWork(ctx); err != nil {
			wc.l.Error("get scan run jobs", "error", err)
		} else if count > 0 {
			wc.l.Info("dispatched scan run tasks", "count", count)
		}
		cancel()
	}
}

func (wc *WorkerCoordinator) ExecuteScanRun(j *ScanRunTask, retryAttempt int, retryInterval time.Duration) *ScanRunResult {
	scanResponse := ScanRunResult{}
	ctx := context.Background()

	if _, err := wc.st.ClaimScanRun(ctx, j.ScanRunID, j.ScanID); err != nil {
		scanResponse.err = fmt.Errorf("%s: %w", errMessageClaimScanRun, err)
		return &scanResponse
	}

	apiKey, err := wc.keyCipher.Decrypt(j.EncryptedKey, j.KeyNonce)
	if err != nil {
		runErr := fmt.Errorf("%s: %w", errMessageDecryptProviderKey, err)
		scanResponse.err = wc.markScanRunFailed(ctx, j, runErr)
		return &scanResponse
	}

	for ; retryAttempt > 0; retryAttempt-- {
		runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		p, ok := wc.providerRegistry[j.EngineID]
		if !ok {
			cancel()
			scanResponse.err = wc.markScanRunFailed(ctx, j, errors.New(errMessageUndefinedEngine))
			return &scanResponse
		}

		result, err := p.Run(runCtx, string(apiKey), nil, j.Request)
		cancel()

		if err != nil {
			if !isRetryable(err) || retryAttempt == 1 {
				scanResponse.err = wc.markScanRunFailed(ctx, j, fmt.Errorf("run provider request: %w", err))
				break
			}

			<-time.After(retryInterval)
			continue
		}

		// Raw responses are diagnostic only; keep the successful result when serialization fails.
		rawJSON, err := json.Marshal(result.Raw)
		if err != nil {
			scanResponse.warnings = append(scanResponse.warnings, fmt.Errorf("%s: %w", errMessageMarshalRawResponse, err))
		}

		// Update the scan run with results
		answerText := result.AnswerText
		now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		if stErr := wc.st.MarkScanRunCompleted(ctx, db.UpdateScanRunParams{
			ID:            j.ScanRunID,
			ScanID:        j.ScanID,
			EngineID:      j.EngineID,
			PromptID:      j.PromptID,
			ProviderKeyID: j.ProviderKeyID,
			TryNumber:     j.TryNumber,
			Status:        "completed",
			AnswerText:    &answerText,
			RawResponse:   rawJSON,
			FinishedAt:    now,
		}); stErr != nil {
			scanResponse.err = fmt.Errorf("%s: %w", errMessageMarkScanRunCompleted, stErr)
			return &scanResponse
		}

		break
	}

	return &scanResponse
}

func (wc *WorkerCoordinator) markScanRunFailed(ctx context.Context, task *ScanRunTask, cause error) error {
	errMsg := cause.Error()
	if err := wc.st.MarkScanRunFailed(ctx, task.ScanRunID, task.ScanID, &errMsg); err != nil {
		return fmt.Errorf("%w; %s: %w", cause, errMessageMarkScanRunFailed, err)
	}

	return cause
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
