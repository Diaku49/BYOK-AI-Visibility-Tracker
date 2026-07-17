package gemini

import (
	"context"
	"net/http"
	"time"

	p "github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"google.golang.org/genai"
)

type GeminiProvider struct {
	model      string
	httpClient *http.Client
}

func NewGeminiProvider() *GeminiProvider {
	httpcli := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &GeminiProvider{
		httpClient: httpcli,
		model:      "gemini-2.5-flash",
	}
}

func (gp *GeminiProvider) newClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	cli, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     apiKey,
		HTTPClient: gp.httpClient,
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func (gp *GeminiProvider) Run(ctx context.Context, apiKey string, req p.RunRequest) (*p.RunResponse, error) {
	client, err := gp.newClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	model := req.Model
	if model == "" {
		model = gp.model
	}

	result, err := gp.generate(ctx, client, model, req.Prompt)
	if err != nil {
		return nil, err
	}

	return &p.RunResponse{
		EngineID:   p.EngineGemini,
		Model:      model,
		AnswerText: result.Text(),
		Raw:        result,
	}, nil
}

func (gp *GeminiProvider) generate(ctx context.Context, cli *genai.Client, model, text string) (*genai.GenerateContentResponse, error) {
	var temp float32 = 0.6
	config := &genai.GenerateContentConfig{
		Temperature:     &temp,
		MaxOutputTokens: 1024,
	}

	result, err := cli.Models.GenerateContent(
		ctx,
		model,
		genai.Text(text),
		config,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
