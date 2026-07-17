package gemini

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	httpClient *http.Client
	cli        *genai.Client
}

func NewGeminiProvider() *GeminiProvider {
	httpcli := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &GeminiProvider{
		httpClient: httpcli,
		cli:        nil,
	}
}

func (gp *GeminiProvider) NewGenaiCli(ctx context.Context, apiKey string, httpc *http.Client) (*genai.Client, error) {
	cli, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     apiKey,
		HTTPClient: httpc,
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}
