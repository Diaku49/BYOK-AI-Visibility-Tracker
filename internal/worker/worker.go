package worker

import (
	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/gemini"
	"github.com/google/uuid"
)

const (
	GeiminiP string = "gemini"
	OpenAIP  string = "openai"
)

type ScanJob struct {
	ScanRuns []ScanRunJob
	ScanLeft int
}

type ScanRunJob struct {
	ScanRunID     uuid.UUID
	ProviderKeyID uuid.UUID

	BrandName   string
	BrandDomain string

	EngineID string
	APIKey   string
	Request  provider.RunRequest
}

type Worker struct {
	providerRegistery map[string]provider.Provider
	scans             map[uuid.UUID]*ScanJob
}

func NewWorker() *Worker {
	return &Worker{
		providerRegistery: make(map[string]provider.Provider),
		scans:             make(map[uuid.UUID]*ScanJob),
	}
}

func (w *Worker) InitRegistery() {
	w.providerRegistery[GeiminiP] = gemini.NewGeminiProvider()
}

func (w *Worker) Start() {}
