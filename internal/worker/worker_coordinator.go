package worker

import (
	"context"
	"fmt"
	"log/slog"

	a "github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/gemini"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/openai"
	"github.com/Diaku49/AI-visibility-tracker/internal/store"
	"github.com/google/uuid"
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
