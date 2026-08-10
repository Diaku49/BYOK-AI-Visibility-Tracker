//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	openaiProvider "github.com/Diaku49/AI-visibility-tracker/internal/provider/openai"
)

func main() {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is required")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("OPENAI_BASE_URL")), "/")
	if baseURL == "" {
		log.Fatal("OPENAI_BASE_URL is required")
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))

	if len(apiKey) >= 8 {
		fmt.Printf("apiKey prefix=%q\n", apiKey[:8])
	}

	const (
		maxAttempts    = 3
		attemptTimeout = 60 * time.Second
		retryWait      = 15 * time.Second
		betweenRunWait = 15 * time.Second
	)

	prompts := []string{
		"What are the best project management tools for small software teams?",
		"What are good alternatives to Linear for issue tracking?",
		"Which AI visibility tracking tools should a startup consider?",
	}

	if len(os.Args) > 1 {
		prompts = []string{strings.Join(os.Args[1:], " ")}
	}

	llm := openaiProvider.NewOpenAIProvider()
	ctx := context.Background()

	for i, prompt := range prompts {
		fmt.Printf("\n--- Run %d ---\n", i+1)
		fmt.Printf("Prompt: %s\n\n", prompt)

		output, err := runOpenAIWithRetry(
			ctx,
			llm,
			apiKey,
			baseURL,
			model,
			prompt,
			maxAttempts,
			attemptTimeout,
			retryWait,
		)
		if err != nil {
			log.Printf("run %d failed after %d attempts: %v", i+1, maxAttempts, err)
			continue
		}

		fmt.Printf("Engine: %s\n", output.EngineID)
		fmt.Printf("Model: %s\n", output.Model)
		fmt.Printf("Input tokens: %d\n", output.Usage.InputTokens)
		fmt.Printf("Output tokens: %d\n\n", output.Usage.OutputTokens)
		fmt.Println(output.AnswerText)

		if i < len(prompts)-1 {
			time.Sleep(betweenRunWait)
		}
	}
}

func runOpenAIWithRetry(
	ctx context.Context,
	llm *openaiProvider.OpenAIProvider,
	apiKey string,
	baseURL string,
	model string,
	prompt string,
	maxAttempts int,
	attemptTimeout time.Duration,
	retryWait time.Duration,
) (*provider.RunOutput, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, attemptTimeout)
		output, err := llm.Run(
			attemptCtx,
			apiKey,
			&baseURL,
			provider.RunInput{
				PromptText:   prompt,
				Model:        model,
				UseWebSearch: false,
			},
		)
		cancel()

		if err == nil {
			return output, nil
		}

		lastErr = err
		if !isOpenAIRetryable(err) || attempt == maxAttempts {
			return nil, fmt.Errorf("run prompt after attempt %d: %w", attempt, err)
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("wait before retry: %w", ctx.Err())
		case <-time.After(retryWait):
		}
	}

	return nil, fmt.Errorf("run prompt after %d attempts: %w", maxAttempts, lastErr)
}

func isOpenAIRetryable(err error) bool {
	if err == nil {
		return false
	}

	errText := strings.ToLower(err.Error())
	return strings.Contains(errText, "408") ||
		strings.Contains(errText, "429") ||
		strings.Contains(errText, "500") ||
		strings.Contains(errText, "502") ||
		strings.Contains(errText, "503") ||
		strings.Contains(errText, "504") ||
		strings.Contains(errText, "rate limit") ||
		strings.Contains(errText, "timeout") ||
		strings.Contains(errText, "deadline exceeded")
}
