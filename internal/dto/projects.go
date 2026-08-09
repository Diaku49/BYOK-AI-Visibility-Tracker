package dto

type CreateProject struct {
	BrandName string     `json:"brand_name" validate:"required"`
	Domain    string     `json:"domain" validate:"required"`
	Language  string     `json:"language" validate:"required"`
	Region    string     `json:"region" validate:"required"`
	Providers []Provider `json:"provider_keys"`
	Prompts   []string   `json:"prompts"`
}

type Provider struct {
	EngineID      string `json:"engine_id"`
	ProviderKeyID string `json:"provider_key_id"`
}
