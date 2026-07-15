package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateProviderKeyRequest struct {
	EngineID        string `json:"engine_id" validate:"required"`
	Name            string `json:"name" validate:"required,min=1,max=100"`
	Key             string `json:"key" validate:"required,min=1"`
	MonthlyRunLimit *int32 `json:"monthly_run_limit" validate:"omitempty,gt=0"`
}

type UpdateProviderKeyMetadataRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=100"`
	Active          bool   `json:"active"`
	MonthlyRunLimit *int32 `json:"monthly_run_limit" validate:"omitempty,gt=0"`
}

type ProviderKeyResponse struct {
	ID              uuid.UUID `json:"id"`
	EngineID        string    `json:"engine_id"`
	Name            string    `json:"name"`
	Active          bool      `json:"active"`
	MonthlyRunLimit *int32    `json:"monthly_run_limit,omitempty"`
	MonthlyRunsUsed int32     `json:"monthly_runs_used"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ListProviderKeysResponse struct {
	ProviderKeys []ProviderKeyResponse `json:"provider_keys"`
}
