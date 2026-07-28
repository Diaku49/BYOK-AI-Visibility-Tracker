package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	a "github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
	"github.com/Diaku49/AI-visibility-tracker/internal/db"
	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/gemini"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/openai"
	"github.com/Diaku49/AI-visibility-tracker/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	GeminiEngineID string = "gemini"
	OpenAIEngineID string = "openai"
)

var (
	ScanRunJobs  = make(chan *ScanRunTask)
	AnalysisJobs = make(chan *AnalysisTask)
)

type ScanRunTask struct {
	ScanID        uuid.UUID
	ScanRunID     uuid.UUID
	PromptID      uuid.UUID
	ProviderKeyID uuid.UUID
	TryNumber     int32

	BrandName   string
	BrandDomain string

	EngineID string
	APIKey   string
	Request  provider.PromptRunRequest
}

type ScanRunResult struct {
	scanRunID uuid.UUID
	result    provider.PromptRunResult
	error     string
}

type ScanBatch struct {
	BrandName      string
	BrandDomain    string
	ScanRunsResult []ScanRunResult
	ScanRunLeft    int
}

type AnalysisTask struct {
	Input    a.ScanAnalysisInput
	EngineID string
	APIKey   string

	ProjectID     uuid.UUID
	TotalRuns     int32
	CompletedRuns int32
	FailedRuns    int32
}

type WorkerCoordinator struct {
	l                *slog.Logger
	st               *store.Store
	keyCipher        *pkg.KeyCipher
	providerRegistry map[string]provider.Runner
}

func NewCoordinator(store *store.Store, logger slog.Logger, keyCipher *pkg.KeyCipher) *WorkerCoordinator {
	providers := make(map[string]provider.Runner)
	providers[GeminiEngineID] = gemini.NewGeminiProvider()
	providers[OpenAIEngineID] = openai.NewOpenAIProvider()

	return &WorkerCoordinator{
		l:                &logger,
		st:               store,
		keyCipher:        keyCipher,
		providerRegistry: providers,
	}
}

func (wc *WorkerCoordinator) Start() {
	// Starting Workers Loop
	for i := 0; i < 2; i++ {
		go wc.StartScanWorker(ScanRunJobs)
	}

	// Start Analyze Loop
	wc.StartAnalysisWorker(AnalysisJobs)
}

func (wc *WorkerCoordinator) GetWork(ctx context.Context) error {
	rows, err := wc.st.GetScansForWorkers(ctx)
	if err != nil {
		return fmt.Errorf("get scans for workers: %w", err)
	}

	for _, row := range rows {
		apiKey, err := wc.keyCipher.Decrypt(row.EncryptedKey, row.KeyNonce)
		if err != nil {
			wc.l.Error("failed to decrypt provider key",
				"scan_run_id", row.ScanRunID,
				"provider_key_id", row.ProviderKeyID,
				"error", err,
			)
			continue
		}

		task := &ScanRunTask{
			ScanID:        row.ScanID,
			ScanRunID:     row.ScanRunID,
			PromptID:      row.PromptID,
			ProviderKeyID: row.ProviderKeyID,
			TryNumber:     row.TryNumber,
			BrandName:     row.BrandName,
			BrandDomain:   row.BrandDomain,
			EngineID:      row.EngineID,
			APIKey:        string(apiKey),
			Request: provider.PromptRunRequest{
				PromptText: row.PromptText,
				Language:   row.Language,
				Region:     row.Region,
			},
		}

		ScanRunJobs <- task
	}

	if len(rows) > 0 {
		wc.l.Info("dispatched scan run tasks", "count", len(rows))
	}

	return nil
}

func (wc *WorkerCoordinator) GetAnalysisWork(ctx context.Context) error {
	rows, err := wc.st.GetScansForAnalysis(ctx)
	if err != nil {
		return fmt.Errorf("get scans for analysis: %w", err)
	}

	if len(rows) == 0 {
		return nil
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
			wc.l.Error("failed to list competitors for analysis",
				"scan_id", scanID,
				"project_id", g.projectID,
				"error", err,
			)
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

		// Pick the engine and API key from the first completed scan_run.
		var engineID string
		var apiKey string
		var completedRuns int32
		var failedRuns int32

		runs := make([]a.RunForAnalysis, 0, len(g.rows))
		for _, row := range g.rows {
			if row.ScanRunStatus == "completed" {
				completedRuns++

				if engineID == "" {
					decrypted, err := wc.keyCipher.Decrypt(row.EncryptedKey, row.KeyNonce)
					if err != nil {
						wc.l.Error("failed to decrypt key for analysis",
							"scan_id", scanID,
							"scan_run_id", row.ScanRunID,
							"error", err,
						)
						continue
					}
					engineID = row.EngineID
					apiKey = string(decrypted)
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
			wc.l.Error("no usable provider key for analysis, all runs failed",
				"scan_id", scanID,
			)
			continue
		}

		task := &AnalysisTask{
			Input: a.ScanAnalysisInput{
				ScanID:      scanID,
				BrandName:   g.brandName,
				BrandDomain: g.brandDomain,
				Competitors: comps,
				Runs:        runs,
			},
			EngineID:      engineID,
			APIKey:        apiKey,
			ProjectID:     g.projectID,
			TotalRuns:     int32(len(g.rows)),
			CompletedRuns: completedRuns,
			FailedRuns:    failedRuns,
		}

		AnalysisJobs <- task
	}

	wc.l.Info("dispatched analysis jobs", "scans", len(scanOrder))

	return nil
}

func (wc *WorkerCoordinator) ExecuteAnalysis(task *AnalysisTask) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	p, ok := wc.providerRegistry[task.EngineID]
	if !ok {
		wc.l.Error("unknown engine for analysis", "engine_id", task.EngineID, "scan_id", task.Input.ScanID)
		return
	}

	analyzer, ok := p.(a.Analyzer)
	if !ok {
		wc.l.Error("provider does not support analysis", "engine_id", task.EngineID, "scan_id", task.Input.ScanID)
		return
	}

	result, err := analyzer.AnalyzeScan(ctx, task.APIKey, task.Input)
	if err != nil {
		wc.l.Error("analysis failed", "scan_id", task.Input.ScanID, "error", err)

		errMsg := err.Error()
		now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		if _, stErr := wc.st.UpdateScanByID(ctx, db.UpdateScanParams{
			ID:             task.Input.ScanID,
			ProjectID:      task.ProjectID,
			Status:         "failed",
			TotalRuns:      task.TotalRuns,
			CompletedRuns:  task.CompletedRuns,
			FailedRuns:     task.FailedRuns,
			Error:          &errMsg,
			FinishedAt:     now,
		}); stErr != nil {
			wc.l.Error("failed to update scan to failed after analysis error", "scan_id", task.Input.ScanID, "error", stErr)
		}
		return
	}

	// Persist per-run analysis results.
	for _, run := range result.Runs {
		competitorsJSON, _ := json.Marshal(run.CompetitorsMentioned)
		citedDomainsJSON, _ := json.Marshal(run.CitedDomains)

		if err := wc.st.UpdateScanRunAnalysis(ctx, db.UpdateScanRunAnalysisParams{
			ID:                   run.ScanRunID,
			BrandMentioned:       &run.BrandMentioned,
			BrandDomainCited:     &run.BrandDomainCited,
			CompetitorsMentioned: competitorsJSON,
			CitedDomains:         citedDomainsJSON,
		}); err != nil {
			wc.l.Error("failed to persist run analysis", "scan_run_id", run.ScanRunID, "error", err)
		}
	}

	// Persist scan summary and mark completed.
	summaryJSON, _ := json.Marshal(result.Summary)
	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	if _, err := wc.st.UpdateScanByID(ctx, db.UpdateScanParams{
		ID:            task.Input.ScanID,
		ProjectID:     task.ProjectID,
		Status:        "completed",
		TotalRuns:     task.TotalRuns,
		CompletedRuns: task.CompletedRuns,
		FailedRuns:    task.FailedRuns,
		Summary:       summaryJSON,
		FinishedAt:    now,
	}); err != nil {
		wc.l.Error("failed to update scan with summary", "scan_id", task.Input.ScanID, "error", err)
		return
	}

	wc.l.Info("analysis completed", "scan_id", task.Input.ScanID, "runs_analyzed", len(result.Runs))
}

func (wc *WorkerCoordinator) StartAnalysisWorker(c chan *AnalysisTask) {
	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case task := <-c:
			{
				wc.l.Info("analysis job received", "scan_id", task.Input.ScanID, "runs", len(task.Input.Runs))
				wc.ExecuteAnalysis(task)
			}
		case <-ticker.C:
			{
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				if err := wc.GetAnalysisWork(ctx); err != nil {
					wc.l.Error("failed getting analysis job", "error", err.Error())
				}
				cancel()
				continue
			}
		}
	}
}
