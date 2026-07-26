package worker

import (
	"fmt"
	"log/slog"

	a "github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
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
	AnalysisJobs = make(chan *a.ScanAnalysisInput)
)

type ScanRunTask struct {
	ScanID        uuid.UUID
	ScanRunID     uuid.UUID
	ProviderKeyID uuid.UUID

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

type WorkerCoordinator struct {
	l                *slog.Logger
	st               *store.Store
	providerRegistry map[string]provider.Runner
}

func NewCoordinator(store *store.Store, logger slog.Logger) *WorkerCoordinator {
	providers := make(map[string]provider.Runner)
	providers[GeminiEngineID] = gemini.NewGeminiProvider()
	providers[OpenAIEngineID] = openai.NewOpenAIProvider()

	return &WorkerCoordinator{
		l:                &logger,
		st:               store,
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

// Should add heartbeat column for scans and store GetWork to handle this
// take available works from db mark them then flow into scan channel
func (wc *WorkerCoordinator) GetWork() error {
	s, err := wc.st.GetScansForWorkers()
	if err != nil {
		return err
	}

	fmt.Printf("data: %v", s)
	return nil
}

func (wc *WorkerCoordinator) PutIntoAnalysis(scanID uuid.UUID, scanBatch *ScanBatch) {
	// runsAnalysis := make([]a.RunForAnalysis, len(scanJob.ScanRunsResult))
	// for _, s := range scanJob.ScanRunsResult {
	// 	run := a.RunForAnalysis{}
	// }

	// input := a.ScanAnalysisInput{
	// 	ScanID:      scanID,
	// 	BrandName:   scanJob.BrandDomain,
	// 	BrandDomain: scanJob.BrandDomain,
	// }
}

func (wc *WorkerCoordinator) StartAnalysisWorker(c chan *a.ScanAnalysisInput) {}
