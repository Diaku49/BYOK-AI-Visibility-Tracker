package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	a "github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	errMessageListCompetitors        = "failed to list competitors for analysis"
	errMessageUpdateScanStateFailed  = "failed to update scan state to failed"
	errMessageDecryptAnalysisKey     = "failed to decrypt key for analysis"
	errMessageNoUsableProviderKey    = "no usable provider key for analysis, all runs failed"
	errMessageUnknownAnalysisEngine  = "unknown engine for analysis"
	errMessageUpdateScanAfterFailure = "failed to update scan to failed after analysis error"
	errMessageAnalysisFailed         = "analysis failed"
	errMessageInvalidAnalysisResult  = "invalid analysis result"
	errMessagePersistScanAnalysis    = "failed to persist scan analysis"
	errMessageGetAnalysisJob         = "failed getting analysis job"
	errMessageClaimScanForAnalysis   = "failed to claim scan for analysis"
)

func (wc *WorkerCoordinator) GetAnalysisWork(ctx context.Context) ([]AnalysisTask, []error, error) {
	var tasks []AnalysisTask
	var taskErrors []error
	rows, err := wc.st.GetScansForAnalysis(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get scans for analysis: %w", err)
	}

	if len(rows) == 0 {
		return nil, nil, nil
	}

	// Group rows by scan_id. Query is ordered by scan_id so consecutive rows belong to the same scan.
	type scanGroup struct {
		projectID   uuid.UUID
		brandName   string
		brandDomain string
		rows        []db.GetScansForAnalysisRow
	}
	groups := make(map[uuid.UUID]*scanGroup)
	var scanOrder []uuid.UUID

	for _, row := range rows {
		g, ok := groups[row.ScanID]
		if !ok {
			g = &scanGroup{
				projectID:   row.ProjectID,
				brandName:   row.BrandName,
				brandDomain: row.BrandDomain,
			}
			groups[row.ScanID] = g
			scanOrder = append(scanOrder, row.ScanID)
		}
		g.rows = append(g.rows, row)
	}

	for _, scanID := range scanOrder {
		g := groups[scanID]

		competitors, err := wc.st.ListCompetitorsByProject(ctx, g.projectID)
		if err != nil {
			taskErr := fmt.Errorf("%s (scan_id=%s): %w", errMessageListCompetitors, scanID, err)
			errMsg := taskErr.Error()
			if _, stErr := wc.st.UpdateScanStateByID(ctx, scanID, "failed", &errMsg); stErr != nil {
				taskErr = fmt.Errorf("%w; %s: %w", taskErr, errMessageUpdateScanStateFailed, stErr)
			}
			taskErrors = append(taskErrors, taskErr)
			continue
		}

		comps := make([]a.CompetitorForAnalysis, 0, len(competitors))
		for _, c := range competitors {
			domain := ""
			if c.Domain != nil {
				domain = *c.Domain
			}
			comps = append(comps, a.CompetitorForAnalysis{
				Name:   c.Name,
				Domain: domain,
			})
		}

		// Pick the engine and encrypted key from the first completed scan_run.
		var engineID string
		var encryptedKey []byte
		var keyNonce []byte
		var completedRuns int32
		var failedRuns int32

		runs := make([]a.RunForAnalysis, 0, len(g.rows))
		for _, row := range g.rows {
			if row.ScanRunStatus == "completed" {
				completedRuns++

				if engineID == "" {
					engineID = row.EngineID
					encryptedKey = row.EncryptedKey
					keyNonce = row.KeyNonce
				}

				answerText := ""
				if row.AnswerText != nil {
					answerText = *row.AnswerText
				}

				runs = append(runs, a.RunForAnalysis{
					ScanRunID:  row.ScanRunID,
					EngineID:   row.EngineID,
					PromptText: row.PromptText,
					AnswerText: answerText,
				})
			} else {
				failedRuns++
			}
		}

		if engineID == "" {
			taskErr := fmt.Errorf("%s (scan_id=%s)", errMessageNoUsableProviderKey, scanID)
			errMsg := taskErr.Error()
			if _, stErr := wc.st.UpdateScanStateByID(ctx, scanID, "failed", &errMsg); stErr != nil {
				taskErr = fmt.Errorf("%w; %s: %w", taskErr, errMessageUpdateScanStateFailed, stErr)
			}
			taskErrors = append(taskErrors, taskErr)
			continue
		}

		task := AnalysisTask{
			Input: a.ScanAnalysisInput{
				ScanID:      scanID,
				BrandName:   g.brandName,
				BrandDomain: g.brandDomain,
				Competitors: comps,
				Runs:        runs,
			},
			EngineID:      engineID,
			EncryptedKey:  encryptedKey,
			KeyNonce:      keyNonce,
			ProjectID:     g.projectID,
			TotalRuns:     int32(len(g.rows)),
			CompletedRuns: completedRuns,
			FailedRuns:    failedRuns,
		}

		tasks = append(tasks, task)
	}

	return tasks, taskErrors, nil
}

func (wc *WorkerCoordinator) ExecuteAnalysis(task *AnalysisTask, apiKey string) error {
	ctx := context.Background()
	analysisCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	p, ok := wc.providerRegistry[task.EngineID]
	if !ok {
		return fmt.Errorf("%s: %q", errMessageUnknownAnalysisEngine, task.EngineID)
	}

	result, err := p.AnalyzeScan(analysisCtx, apiKey, task.Input)
	if err != nil {
		return fmt.Errorf("%s: %w", errMessageAnalysisFailed, err)
	}

	if err := validateAnalysisRuns(task.Input.Runs, result.Runs); err != nil {
		return fmt.Errorf("%s: %w", errMessageInvalidAnalysisResult, err)
	}

	// Persist per-run analysis results.
	batch := []db.UpdateScanRunAnalysisParams{}
	for _, run := range result.Runs {
		competitorsJSON, err := json.Marshal(run.CompetitorsMentioned)
		if err != nil {
			return fmt.Errorf("marshal competitors mentioned for scan run %s: %w", run.ScanRunID, err)
		}
		citedDomainsJSON, err := json.Marshal(run.CitedDomains)
		if err != nil {
			return fmt.Errorf("marshal cited domains for scan run %s: %w", run.ScanRunID, err)
		}

		batch = append(batch, db.UpdateScanRunAnalysisParams{
			ID:                   run.ScanRunID,
			BrandMentioned:       &run.BrandMentioned,
			BrandDomainCited:     &run.BrandDomainCited,
			CompetitorsMentioned: competitorsJSON,
			CitedDomains:         citedDomainsJSON,
		})
	}
	summaryJSON, err := json.Marshal(result.Summary)
	if err != nil {
		return fmt.Errorf("marshal scan analysis summary: %w", err)
	}
	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	if err := wc.st.UpdateScanAnalysis(ctx, batch, db.UpdateScanParams{
		ID:            task.Input.ScanID,
		ProjectID:     task.ProjectID,
		Status:        "completed",
		TotalRuns:     task.TotalRuns,
		CompletedRuns: task.CompletedRuns,
		FailedRuns:    task.FailedRuns,
		Summary:       summaryJSON,
		FinishedAt:    now,
	}); err != nil {
		return fmt.Errorf("%s: %w", errMessagePersistScanAnalysis, err)
	}

	return nil
}

func (wc *WorkerCoordinator) StartAnalysisWorker(c chan *AnalysisTask) {
	ctx := context.Background()

	for task := range c {
		if _, err := wc.st.ClaimScanForAnalysis(ctx, task.Input.ScanID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				wc.l.Info("scan already claimed for analysis", "scan_id", task.Input.ScanID)
				continue
			}

			wc.l.Error(errMessageClaimScanForAnalysis, "scan_id", task.Input.ScanID, "error", err)
			continue
		}

		apiKey, err := wc.keyCipher.Decrypt(task.EncryptedKey, task.KeyNonce)
		if err != nil {
			analysisErr := fmt.Errorf("%s: %w", errMessageDecryptAnalysisKey, err)
			wc.l.Error("execute analysis", "scan_id", task.Input.ScanID, "error", analysisErr)
			if failErr := wc.failAnalysis(ctx, task, analysisErr); failErr != nil {
				wc.l.Error(errMessageUpdateScanAfterFailure, "scan_id", task.Input.ScanID, "error", failErr)
			}
			continue
		}

		wc.l.Info("analysis job received", "scan_id", task.Input.ScanID, "runs", len(task.Input.Runs))
		if err := wc.ExecuteAnalysis(task, string(apiKey)); err != nil {
			wc.l.Error("execute analysis", "scan_id", task.Input.ScanID, "error", err)
			if failErr := wc.failAnalysis(ctx, task, err); failErr != nil {
				wc.l.Error(errMessageUpdateScanAfterFailure, "scan_id", task.Input.ScanID, "error", failErr)
			}
			continue
		}

		wc.l.Info("analysis completed", "scan_id", task.Input.ScanID, "runs_analyzed", len(task.Input.Runs))
	}
}

func (wc *WorkerCoordinator) failAnalysis(ctx context.Context, task *AnalysisTask, cause error) error {
	errMsg := cause.Error()
	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	if _, err := wc.st.UpdateScanByID(ctx, db.UpdateScanParams{
		ID:            task.Input.ScanID,
		ProjectID:     task.ProjectID,
		Status:        "failed",
		TotalRuns:     task.TotalRuns,
		CompletedRuns: task.CompletedRuns,
		FailedRuns:    task.FailedRuns,
		Error:         &errMsg,
		FinishedAt:    now,
	}); err != nil {
		return fmt.Errorf("%s: %w", errMessageUpdateScanAfterFailure, err)
	}

	return nil
}

func (wc *WorkerCoordinator) AnalysisTaskProducer() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		if tasks, taskErrors, err := wc.GetAnalysisWork(ctx); err != nil {
			wc.l.Error(errMessageGetAnalysisJob, "error", err)
		} else {
			for _, taskErr := range taskErrors {
				wc.l.Error("prepare analysis job", "error", taskErr)
			}
			for _, task := range tasks {
				AnalysisJobs <- &task
			}
			if len(tasks) > 0 {
				wc.l.Info("dispatched analysis jobs", "scans", len(tasks))
			}
		}

		cancel()
	}
}
