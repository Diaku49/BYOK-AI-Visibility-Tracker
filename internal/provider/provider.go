package provider

import (
	"context"
)

type EngineID string

const (
	EngineGemini EngineID = "gemini"
	EngineOpenAI EngineID = "openai"
)

type Provider interface {
	ID() EngineID
	Run(ctx context.Context, apiKey string, baseURL *string, input RunInput) (*RunOutput, error)
}

type RunInput struct {
	PromptText   string `json:"prompt_text"`
	Model        string `json:"model"`
	UseWebSearch bool   `json:"use_web_search"`
	Language     string `json:"language"`
	Region       string `json:"region"`
}

type RunOutput struct {
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

type RunRequest = RunInput
type RunResponse = RunOutput

var SysUserInstruction string = "Answer in plain text only. Do not use markdown formatting, tables, bullet points, headings, or code blocks."
