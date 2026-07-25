package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
	p "github.com/Diaku49/AI-visibility-tracker/internal/provider"
	openai "github.com/sashabaranov/go-openai"
)

type debugTransport struct {
	base http.RoundTripper
}

func (t debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf("HTTP %s %s\n", req.Method, req.URL.String())

	return t.base.RoundTrip(req)
}

type OpenAIProvider struct {
	model          string
	webSearchModel string
	httpClient     *http.Client
}

func NewOpenAIProvider() *OpenAIProvider {
	httpcli := &http.Client{
		Timeout: 30 * time.Second,
		Transport: debugTransport{
			base: http.DefaultTransport,
		},
	}

	return &OpenAIProvider{
		model:          openai.GPT4Dot1Mini,
		webSearchModel: "gpt-5-search-api",
		httpClient:     httpcli,
	}
}

func NewOpenAiProvider() *OpenAIProvider {
	return NewOpenAIProvider()
}

func (oap *OpenAIProvider) NewClient(baseURL *string, apiKey string) (*openai.Client, error) {
	config := openai.DefaultConfig(apiKey)
	config.HTTPClient = oap.httpClient

	if baseURL != nil {
		config.BaseURL = *baseURL
	}

	return openai.NewClientWithConfig(config), nil
}

func (oap *OpenAIProvider) ID() p.EngineID { return p.EngineOpenAI }

func (oap *OpenAIProvider) generate(
	ctx context.Context,
	cli *openai.Client,
	model string,
	req p.RunRequest,
) (*openai.ChatCompletionResponse, error) {
	request := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: p.SysUserInstruction,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.PromptText,
			},
		},
	}
	request.MaxCompletionTokens = 3000

	result, err := cli.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (oap *OpenAIProvider) Run(
	ctx context.Context,
	apiKey string,
	baseURL *string,
	req p.RunInput,
) (*p.RunOutput, error) {
	client, err := oap.NewClient(baseURL, apiKey)
	if err != nil {
		return nil, err
	}

	model := req.Model
	if model == "" {
		model = oap.model
		if req.UseWebSearch {
			model = oap.webSearchModel
		}
	}

	result, err := oap.generate(ctx, client, model, req)
	if err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	message := result.Choices[0].Message
	if message.Refusal != "" {
		return nil, fmt.Errorf("openai refused the request: %s", message.Refusal)
	}

	return &p.RunOutput{
		EngineID:   p.EngineOpenAI,
		Model:      model,
		AnswerText: message.Content,
		Raw:        result,
		Usage: p.TokenUsage{
			InputTokens:  result.Usage.PromptTokens,
			OutputTokens: result.Usage.CompletionTokens,
		},
	}, nil
}

func (oap *OpenAIProvider) AnalyzeScan(
	ctx context.Context,
	apiKey string,
	input analyzer.ScanAnalysisInput,
) (*analyzer.ScanAnalysisResult, error) {
	client, err := oap.NewClient(nil, apiKey)
	if err != nil {
		return nil, err
	}

	prompt, err := analyzer.BuildAnalyzerPrompt(input)
	if err != nil {
		return nil, err
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: oap.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: analyzer.AnalyzerSystemInstruction,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature:         0.1,
			MaxCompletionTokens: 4096,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:        "scan_analysis",
					Description: "Complete scan visibility analysis result.",
					Schema:      analyzer.OpenAIScanAnalysisSchema(),
					Strict:      true,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no analyzer choices")
	}

	message := resp.Choices[0].Message
	if message.Refusal != "" {
		return nil, fmt.Errorf("openai refused the analyzer request: %s", message.Refusal)
	}

	var result analyzer.ScanAnalysisResult
	if err := json.Unmarshal([]byte(message.Content), &result); err != nil {
		return nil, fmt.Errorf("decode analyzer json: %w; raw=%s", err, message.Content)
	}

	return &result, nil
}
