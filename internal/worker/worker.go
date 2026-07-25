package worker

import (
	"context"
	"log/slog"
	"time"

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

	pollInterval    = 10 * time.Second
	attemptInterval = 5 * time.Second
)

type ScanJob struct {
	BrandName      string
	BrandDomain    string
	ScanRunsResult []ScanRunJobResponse
	ScanRunLeft    int
}

type ScanRunJob struct {
	ScanID        uuid.UUID
	ScanRunID     uuid.UUID
	ProviderKeyID uuid.UUID

	BrandName   string
	BrandDomain string

	EngineID string
	APIKey   string
	Request  provider.RunRequest
}

type ScanRunJobResponse struct {
	scanRunID uuid.UUID
	result    provider.RunResponse
	error     string
}

type WorkerManager struct {
	l                 *slog.Logger
	st                *store.Store
	providerRegistery map[string]provider.Provider
	scans             map[uuid.UUID]*ScanJob
}

func NewWorker(store *store.Store, logger slog.Logger) *WorkerManager {
	pr := make(map[string]provider.Provider)
	pr[GeiminiP] = gemini.NewGeminiProvider()
	pr[OpenAIP] = openai.NewOpenAIProvider()

	return &WorkerManager{
		l:                 &logger,
		st:                store,
		providerRegistery: pr,
		scans:             make(map[uuid.UUID]*ScanJob),
	}
}

func (wm *WorkerManager) Start() {
	// Starting Workers Loop
	for i := 0; i < 0; i++ {
		go wm.StartWorker(ScanRunC)
	}

	// Start Analyze Loop
	wm.StartAnalyze(AnalayzeC)
}

func (wm *WorkerManager) StartWorker(c chan *ScanRunJob) {
	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case j := <-c:
			{
				ticker.Reset(pollInterval)
				// Check existence
				if _, ok := wm.scans[j.ScanID]; !ok {
					wm.l.Error("This Scan does not exist", "ScanID", j.ScanID)
					continue
				}
				wm.scans[j.ScanID].ScanRunLeft--

				// Do ScanRunJob -- needs to be refactored
				scanResponse := wm.Run(j, 2, attemptInterval)
				wm.l.Info("Scan ran successfully", "ScanID", scanResponse.scanRunID)

				// If no run remains flow into Analysis
				if wm.scans[j.ScanID].ScanRunLeft <= 0 {
					scanJob := wm.scans[j.ScanID]
					wm.scans[j.ScanID] = nil
					wm.PutIntoAnalysis(j.ScanID, scanJob)
				}
				continue
			}
		case <-ticker.C:
			{
				if err := wm.GetWork(); err != nil {
					wm.l.Error("Failed getting job", "Error", err.Error())
					continue
				}
				continue
			}
		}
	}
}

// Retry interval needs to be implemented
func (wm *WorkerManager) Run(j *ScanRunJob, retryAttempt int, retryInterval time.Duration) *ScanRunJobResponse {
	var scanResponse ScanRunJobResponse
	for ; retryAttempt > 0; retryAttempt-- {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		p := wm.providerRegistery[j.EngineID]
		result, err := p.Run(ctx, j.APIKey, nil, j.Request)

		if err != nil {
			wm.l.Error(err.Error(), "ScanID", j.ScanID, "RunID", j.ScanRunID, "Engine", j.EngineID)
			scanResponse = ScanRunJobResponse{
				scanRunID: j.ScanRunID,
				result:    *result,
				error:     err.Error(),
			}

			select {
			case <-ctx.Done():
				cancel()
				wm.scans[j.ScanID].ScanRunsResult = append(wm.scans[j.ScanID].ScanRunsResult, scanResponse)
				return &scanResponse
			case <-time.After(retryInterval):
				cancel()
				continue
			}
		}

		// Update the ScanJob
		scanResponse = ScanRunJobResponse{
			scanRunID: j.ScanRunID,
			result:    *result,
			error:     "",
		}
		cancel()
	}
	wm.scans[j.ScanID].ScanRunsResult = append(wm.scans[j.ScanID].ScanRunsResult, scanResponse)

	return &scanResponse
}

// take available works from db mark them then flow into scan channel
func (wm *WorkerManager) GetWork() error {
	return nil
}

func (wm *WorkerManager) PutIntoAnalysis(scanJobID uuid.UUID, scanJob *ScanJob) {}

func (wm *WorkerManager) StartAnalyze(c chan *a.ScanAnalysisInput) {}
