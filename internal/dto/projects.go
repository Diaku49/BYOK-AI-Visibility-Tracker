package dto

type CreateProject struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Language    string `json:"language" validate:"required"`
	Region      string `json:"region" validate:"required"`
}
