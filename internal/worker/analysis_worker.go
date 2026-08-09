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

//********** Needs Complete Refactoring **********

func (wc *WorkerCoordinator) GetAnalysisWork(ctx context.Context) ([]AnalysisTask, error) {
	var tasks []AnalysisTask
	rows, err := wc.st.GetScansForAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("get scans for analysis: %w", err)
	}

	if len(rows) == 0 {
		return nil, nil
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
			errMsg := errMessageListCompetitors
			wc.l.Error(errMsg,
				"scan_id", scanID,
				"project_id", g.projectID,
				"error", err,
			)

			if _, stErr := wc.st.UpdateScanStateByID(ctx, scanID, "failed", &errMsg); stErr != nil {
				wc.l.Error(errMessageUpdateScanStateFailed, "scan_id", scanID, "error", stErr)
			}
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
			errMsg := errMessageNoUsableProviderKey
			wc.l.Error(errMsg,
				"scan_id", scanID,
			)
			if _, stErr := wc.st.UpdateScanStateByID(ctx, scanID, "failed", &errMsg); stErr != nil {
				wc.l.Error(errMessageUpdateScanStateFailed, "scan_id", scanID, "error", stErr)
			}
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

	wc.l.Info("dispatched analysis jobs", "scans", len(scanOrder))

	return tasks, nil
}

func (wc *WorkerCoordinator) ExecuteAnalysis(task *AnalysisTask, apiKey string) {
	ctx := context.Background()
	analysisCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	p, ok := wc.providerRegistry[task.EngineID]
	if !ok {
		wc.l.Error(errMessageUnknownAnalysisEngine, "engine_id", task.EngineID, "scan_id", task.Input.ScanID)
		errMsg := errMessageUnknownAnalysisEngine
		now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		if _, stErr := wc.st.UpdateScanByID(ctx, db.UpdateScanParams{
			ID:            task.Input.ScanID,
			ProjectID:     task.ProjectID,
			Status:        "failed",
			TotalRuns:     task.TotalRuns,
			CompletedRuns: task.CompletedRuns,
			FailedRuns:    task.FailedRuns,
			Error:         &errMsg,
			FinishedAt:    now,
		}); stErr != nil {
			wc.l.Error(errMessageUpdateScanAfterFailure, "scan_id", task.Input.ScanID, "error", stErr)
		}
		return
	}

	result, err := p.AnalyzeScan(analysisCtx, apiKey, task.Input)
	if err != nil {
		wc.l.Error(errMessageAnalysisFailed, "scan_id", task.Input.ScanID, "error", err)

		errMsg := err.Error()
		now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		if _, stErr := wc.st.UpdateScanByID(ctx, db.UpdateScanParams{
			ID:            task.Input.ScanID,
			ProjectID:     task.ProjectID,
			Status:        "failed",
			TotalRuns:     task.TotalRuns,
			CompletedRuns: task.CompletedRuns,
			FailedRuns:    task.FailedRuns,
			Error:         &errMsg,
			FinishedAt:    now,
		}); stErr != nil {
			wc.l.Error(errMessageUpdateScanAfterFailure, "scan_id", task.Input.ScanID, "error", stErr)
		}
		return
	}

	if err := validateAnalysisRuns(task.Input.Runs, result.Runs); err != nil {
		wc.l.Error(errMessageInvalidAnalysisResult, "scan_id", task.Input.ScanID, "error", err)

		errMsg := err.Error()
		now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		if _, stErr := wc.st.UpdateScanByID(ctx, db.UpdateScanParams{
			ID:            task.Input.ScanID,
			ProjectID:     task.ProjectID,
			Status:        "failed",
			TotalRuns:     task.TotalRuns,
			CompletedRuns: task.CompletedRuns,
			FailedRuns:    task.FailedRuns,
			Error:         &errMsg,
			FinishedAt:    now,
		}); stErr != nil {
			wc.l.Error(errMessageUpdateScanAfterFailure, "scan_id", task.Input.ScanID, "error", stErr)
		}
		return
	}

	// Persist per-run analysis results.
	batch := []db.UpdateScanRunAnalysisParams{}
	for _, run := range result.Runs {
		competitorsJSON, _ := json.Marshal(run.CompetitorsMentioned)
		citedDomainsJSON, _ := json.Marshal(run.CitedDomains)

		batch = append(batch, db.UpdateScanRunAnalysisParams{
			ID:                   run.ScanRunID,
			BrandMentioned:       &run.BrandMentioned,
			BrandDomainCited:     &run.BrandDomainCited,
			CompetitorsMentioned: competitorsJSON,
			CitedDomains:         citedDomainsJSON,
		})
	}
	summaryJSON, _ := json.Marshal(result.Summary)
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
		wc.l.Error(errMessagePersistScanAnalysis, "scan_id", task.Input.ScanID, "error", err)
		return
	}

	wc.l.Info("analysis completed", "scan_id", task.Input.ScanID, "runs_analyzed", len(result.Runs))
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
			errMsg := fmt.Sprintf("%s: %v", errMessageDecryptAnalysisKey, err)
			wc.l.Error(errMessageDecryptAnalysisKey, "scan_id", task.Input.ScanID, "error", err)
			if _, stErr := wc.st.UpdateScanStateByID(ctx, task.Input.ScanID, "failed", &errMsg); stErr != nil {
				wc.l.Error(errMessageUpdateScanStateFailed, "scan_id", task.Input.ScanID, "error", stErr)
			}
			continue
		}

		wc.l.Info("analysis job received", "scan_id", task.Input.ScanID, "runs", len(task.Input.Runs))
		wc.ExecuteAnalysis(task, string(apiKey))
	}
}

func (wc *WorkerCoordinator) AnalysisTaskProducer() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		if tasks, err := wc.GetAnalysisWork(ctx); err != nil {
			wc.l.Error(errMessageGetAnalysisJob, "error", err.Error())
		} else {
			for _, task := range tasks {
				AnalysisJobs <- &task
			}
		}

		cancel()
	}
}
