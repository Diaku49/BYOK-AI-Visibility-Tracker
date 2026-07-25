package worker

import (
	"log/slog"

	a "github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/gemini"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/openai"
	"github.com/Diaku49/AI-visibility-tracker/internal/store"
	"github.com/google/uuid"
)

const (
	GeiminiP string = "gemini"
	OpenAIP  string = "openai"
)

var (
	ScanRunC  = make(chan *ScanRunJob)
	AnalayzeC = make(chan *a.ScanAnalysisInput)
)

type ScanJob struct {
	BrandName      string
	BrandDomain    string
	ScanRunsResult []provider.RunResponse
	ScanRunLeft    int
}

type ScanRunJob struct {
	ScanID        uuid.UUIDs
	ScanRunID     uuid.UUID
	ProviderKeyID uuid.UUID

	BrandName   string
	BrandDomain string

	EngineID string
	APIKey   string
	Request  provider.RunRequest
}

type WorkerManger struct {
	l                 *slog.Logger
	st                *store.Store
	providerRegistery map[string]provider.Provider
	scans             map[uuid.UUID]*ScanJob
}

func NewWorker(store *store.Store, logger slog.Logger) *WorkerManger {
	pr := make(map[string]provider.Provider)
	pr[GeiminiP] = gemini.NewGeminiProvider()
	pr[OpenAIP] = openai.NewOpenAIProvider()

	return &WorkerManger{
		l:                 &logger,
		st:                store,
		providerRegistery: pr,
		scans:             make(map[uuid.UUID]*ScanJob),
	}
}

func (wm *WorkerManger) Start() {

	// Starting Workers Loop

	// Start Analyze Loop
}

func (wm *WorkerManger) StartSlave() {}

func (wm *WorkerManger) StartAnalyze() {}

func (wm *WorkerManger) MakeScanAnalysisInput(ScanJob) *a.ScanAnalysisInput {

	return &a.ScanAnalysisInput{}
}
