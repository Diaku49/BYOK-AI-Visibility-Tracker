package provider

import "context"

type EngineID string

const (
	EngineGemini EngineID = "gemini"
	EngineGPT    EngineID = "gpt"
)

type Provider interface {
	Run(ctx context.Context, apiKey string, req RunRequest) (*RunResponse, error)
}

type RunRequest struct {
	Prompt       string
	Model        string
	UseWebSearch bool
	Language     string
	Region       string
}

type Citation struct {
	URL   string
	Title string
	Text  string
}

type RunResponse struct {
	EngineID   EngineID
	Model      string
	AnswerText string
	Citations  []Citation
	Raw        any

	InputTokens  int
	OutputTokens int
}
