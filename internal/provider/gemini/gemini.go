package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
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
		model:      "gemini-flash-latest",
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

func (gp *GeminiProvider) ID() p.EngineID { return p.EngineGemini }

func (gp *GeminiProvider) generate(ctx context.Context, cli *genai.Client, model string, req p.RunRequest) (*genai.GenerateContentResponse, error) {
	var temp float32 = 0.6
	config := &genai.GenerateContentConfig{
		Temperature:     &temp,
		MaxOutputTokens: 3000,
		SystemInstruction: genai.NewContentFromText(
			p.SysUserInstruction,
			genai.RoleUser,
		),
	}

	if req.UseWebSearch {
		config.Tools = []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		}
	}

	result, err := cli.Models.GenerateContent(
		ctx,
		model,
		genai.Text(req.PromptText),
		config,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (gp *GeminiProvider) Run(ctx context.Context, apiKey string, _ *string, req p.RunInput) (*p.RunOutput, error) {
	client, err := gp.newClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	model := req.Model
	if model == "" {
		model = gp.model
	}

	result, err := gp.generate(ctx, client, model, req)
	if err != nil {
		return nil, err
	}

	output := &p.RunOutput{
		EngineID:   p.EngineGemini,
		Model:      model,
		AnswerText: result.Text(),
		Raw:        result,
	}

	if result.UsageMetadata != nil {
		output.Usage = p.TokenUsage{
			InputTokens:  int(result.UsageMetadata.PromptTokenCount),
			OutputTokens: int(result.UsageMetadata.CandidatesTokenCount),
		}
	}

	return output, nil
}

func (gp *GeminiProvider) AnalyzeScan(
	ctx context.Context,
	apiKey string,
	input analyzer.ScanAnalysisInput,
) (*analyzer.ScanAnalysisResult, error) {
	client, err := gp.newClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	prompt, err := analyzer.BuildAnalyzerPrompt(input)
	if err != nil {
		return nil, err
	}

	temp := float32(0.1)

	config := &genai.GenerateContentConfig{
		Temperature:      &temp,
		MaxOutputTokens:  4096,
		ResponseMIMEType: "application/json",
		ResponseSchema:   analyzer.ScanAnalysisSchema(),
		SystemInstruction: genai.NewContentFromText(
			analyzer.AnalyzerSystemInstruction,
			genai.RoleUser,
		),
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		gp.model,
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, err
	}

	text := resp.Text()

	var result analyzer.ScanAnalysisResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("decode analyzer json: %w; raw=%s", err, text)
	}

	return &result, nil
}
