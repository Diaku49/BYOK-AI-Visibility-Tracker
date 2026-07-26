package worker

import (
	"context"
	"strings"
	"time"
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
				if err := wc.GetWork(); err != nil {
					wc.l.Error("Failed getting job", "Error", err.Error())
					continue
				}
				continue
			}
		}
	}
}

// Need to implement storing in db
func (wc *WorkerCoordinator) ExecuteScanRun(j *ScanRunTask, retryAttempt int, retryInterval time.Duration) *ScanRunResult {
	var scanResponse ScanRunResult
	// changing scan run state with store method
	wc.st.UpdateStateScanRunByID(j.ScanRunID, "running")

	for ; retryAttempt > 0; retryAttempt-- {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		p := wc.providerRegistry[j.EngineID]
		result, err := p.Run(ctx, j.APIKey, nil, j.Request)
		cancel()

		if err != nil {
			wc.l.Error(err.Error(), "ScanID", j.ScanID, "RunID", j.ScanRunID, "Engine", j.EngineID)
			scanResponse = ScanRunResult{
				scanRunID: j.ScanRunID,
				error:     err.Error(),
			}

			if !isRetryable(err) || retryAttempt == 0 {
				break
			}

			select {
			case <-ctx.Done():
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
		// Update the scan run result.
		wc.st.UpdateScanRunByID(j.ScanRunID, scanResponse)
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
