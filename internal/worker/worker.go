package worker

import (
	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/gemini"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/openai"
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
	pr := make(map[string]provider.Provider)
	pr[GeiminiP] = gemini.NewGeminiProvider()
	pr[OpenAIP] = openai.NewOpenAIProvider()

	return &Worker{
		providerRegistery: pr,
		scans:             make(map[uuid.UUID]*ScanJob),
	}
}

func (w *Worker) Start() {}
