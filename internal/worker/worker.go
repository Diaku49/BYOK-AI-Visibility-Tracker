package worker

import (
	"context"
	"encoding/json"
	"errors"
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
		wc.l.Info("Scan ran", "ScanID", scanResponse.scanRunID)
	}
}

func (wc *WorkerCoordinator) ScanTaskProducer() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		if err := wc.GetWork(ctx); err != nil {
			wc.l.Error("Failed getting job", "Error", err.Error())
		}
		cancel()
	}
}

func (wc *WorkerCoordinator) ExecuteScanRun(j *ScanRunTask, retryAttempt int, retryInterval time.Duration) *ScanRunResult {
	scanResponse := ScanRunResult{scanRunID: j.ScanRunID}
	ctx := context.Background()

	if _, err := wc.st.ClaimScanRun(ctx, j.ScanRunID, j.ScanID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wc.l.Info("scan run already claimed", "scan_id", j.ScanID, "run_id", j.ScanRunID)
			return &scanResponse
		}

		wc.l.Error(errMessageClaimScanRun, "scan_id", j.ScanID, "run_id", j.ScanRunID, "error", err)
		return &scanResponse
	}

	apiKey, err := wc.keyCipher.Decrypt(j.EncryptedKey, j.KeyNonce)
	if err != nil {
		errMsg := errMessageDecryptProviderKey
		scanResponse.error = errMsg
		wc.l.Error(errMsg, "scan_run_id", j.ScanRunID, "provider_key_id", j.ProviderKeyID, "error", err)
		if stErr := wc.st.MarkScanRunFailed(ctx, j.ScanRunID, j.ScanID, &errMsg); stErr != nil {
			wc.l.Error(errMessageMarkScanRunFailed, "run_id", j.ScanRunID, "error", stErr)
		}
		return &scanResponse
	}

	for ; retryAttempt > 0; retryAttempt-- {
		runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		p, ok := wc.providerRegistry[j.EngineID]
		if !ok {
			errMsg := errMessageUndefinedEngine
			wc.l.Error(errMsg, "ScanID", j.ScanID, "RunID", j.ScanRunID, "Engine", j.EngineID)
			scanResponse = ScanRunResult{
				scanRunID: j.ScanRunID,
				error:     errMsg,
			}

			if stErr := wc.st.MarkScanRunFailed(ctx, j.ScanRunID, j.ScanID, &errMsg); stErr != nil {
				wc.l.Error(errMessageMarkScanRunFailed, "RunID", j.ScanRunID, "error", stErr)
			}

			cancel()
			return &scanResponse
		}

		result, err := p.Run(runCtx, string(apiKey), nil, j.Request)
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
				if stErr := wc.st.MarkScanRunFailed(ctx, j.ScanRunID, j.ScanID, &errMsg); stErr != nil {
					wc.l.Error(errMessageMarkScanRunFailed, "RunID", j.ScanRunID, "error", stErr)
				}
				break
			}

			<-time.After(retryInterval)
			continue
		}

		scanResponse = ScanRunResult{
			scanRunID: j.ScanRunID,
			result:    *result,
			error:     "",
		}

		// Marshal raw response for storage - not gonna be used for anything else, just for logging and debugging
		rawJSON, err := json.Marshal(result.Raw)
		if err != nil {
			wc.l.Error(errMessageMarshalRawResponse, "RunID", j.ScanRunID, "error", err)
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
			wc.l.Error(errMessageMarkScanRunCompleted, "RunID", j.ScanRunID, "error", stErr)
			return &scanResponse
		}

		wc.l.Info("scan run completed", "scan_id", j.ScanID, "run_id", j.ScanRunID, "engine", j.EngineID)
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
