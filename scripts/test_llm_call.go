package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/gemini"
)

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is required")
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

	ctx := context.Background()

	llm := gemini.NewGeminiProvider()

	for i, prompt := range prompts {
		fmt.Printf("\n--- Run %d ---\n", i+1)
		fmt.Printf("Prompt: %s\n\n", prompt)

		output, err := runWithRetry(ctx, llm, apiKey, prompt, i+1, maxAttempts, attemptTimeout, retryWait)
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

func runWithRetry(
	ctx context.Context,
	llm *gemini.GeminiProvider,
	apiKey string,
	prompt string,
	runNumber int,
	maxAttempts int,
	attemptTimeout time.Duration,
	retryWait time.Duration,
) (*provider.RunOutput, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, attemptTimeout)
		output, err := llm.Run(attemptCtx, apiKey, provider.RunInput{
			PromptText:   prompt,
			UseWebSearch: false,
		})
		cancel()

		if err == nil {
			return output, nil
		}

		lastErr = err
		if !isRetryable(err) || attempt == maxAttempts {
			return nil, err
		}

		log.Printf("run %d attempt %d/%d failed: %v", runNumber, attempt, maxAttempts, err)
		log.Printf("waiting %s before retry...", retryWait)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryWait):
		}
	}

	return nil, lastErr
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}

	errText := err.Error()
	return strings.Contains(errText, "429") ||
		strings.Contains(errText, "503") ||
		strings.Contains(errText, "504") ||
		strings.Contains(errText, "RESOURCE_EXHAUSTED") ||
		strings.Contains(errText, "UNAVAILABLE") ||
		strings.Contains(errText, "DEADLINE_EXCEEDED")
}
