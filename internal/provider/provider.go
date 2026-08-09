package provider

import (
	"context"

	"github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
)

type EngineID string

const (
	EngineGemini EngineID = "gemini"
	EngineOpenAI EngineID = "openai"
)

type Runner interface {
	ID() EngineID
	Run(ctx context.Context, apiKey string, baseURL *string, input PromptRunRequest) (*PromptRunResult, error)
	AnalyzeScan(ctx context.Context, apiKey string, input analyzer.ScanAnalysisInput) (*analyzer.ScanAnalysisResult, error)
}

type PromptRunRequest struct {
	PromptText   string `json:"prompt_text"`
	Model        string `json:"model"`
	UseWebSearch bool   `json:"use_web_search"`
	Language     string `json:"language"`
	Region       string `json:"region"`
}

type PromptRunResult struct {
	EngineID   EngineID
	Model      string
	AnswerText string
	Citations  []Citation
	Raw        any
	Usage      TokenUsage
}

type TokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type Citation struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

var PlainTextSystemInstruction string = "Answer in plain text only. Do not use markdown formatting, tables, bullet points, headings, or code blocks."
