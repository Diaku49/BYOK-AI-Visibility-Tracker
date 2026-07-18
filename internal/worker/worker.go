package worker

import (
	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"github.com/google/uuid"
)

type ScanRunJob struct {
	ScanRunID     uuid.UUID
	ProviderKeyID uuid.UUID

	BrandName   string
	BrandDomain string

	EngineID string
	APIKey   string
	Request  provider.RunRequest
}
