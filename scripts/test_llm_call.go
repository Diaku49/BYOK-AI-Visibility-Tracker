package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Diaku49/AI-visibility-tracker/internal/analyzer"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider"
	"github.com/Diaku49/AI-visibility-tracker/internal/provider/gemini"
	"github.com/google/uuid"
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
	scanID := uuid.New()
	runsForAnalysis := make([]analyzer.RunForAnalysis, 0, len(prompts))

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

		runsForAnalysis = append(runsForAnalysis, analyzer.RunForAnalysis{
			ScanRunID:  uuid.New(),
			EngineID:   string(output.EngineID),
			PromptText: prompt,
			AnswerText: output.AnswerText,
			Citations:  toAnalyzerCitations(output.Citations),
		})

		if i < len(prompts)-1 {
			time.Sleep(betweenRunWait)
		}
	}

	if len(runsForAnalysis) == 0 {
		log.Println("no successful runs to analyze")
		return
	}

	fmt.Printf("\n--- Analysis ---\n")

	analysisInput := analyzer.ScanAnalysisInput{
		ScanID:      scanID,
		BrandName:   envOrDefault("TEST_BRAND_NAME", "Linear"),
		BrandDomain: envOrDefault("TEST_BRAND_DOMAIN", "linear.app"),
		Competitors: []analyzer.CompetitorForAnalysis{
			{Name: envOrDefault("TEST_COMPETITOR_1_NAME", "Jira"), Domain: envOrDefault("TEST_COMPETITOR_1_DOMAIN", "atlassian.com")},
			{Name: envOrDefault("TEST_COMPETITOR_2_NAME", "Asana"), Domain: envOrDefault("TEST_COMPETITOR_2_DOMAIN", "asana.com")},
			{Name: envOrDefault("TEST_COMPETITOR_3_NAME", "ClickUp"), Domain: envOrDefault("TEST_COMPETITOR_3_DOMAIN", "clickup.com")},
		},
		Runs: runsForAnalysis,
	}

	analysisCtx, cancel := context.WithTimeout(ctx, attemptTimeout)
	analysis, err := llm.AnalyzeScan(analysisCtx, apiKey, analysisInput)
	cancel()
	if err != nil {
		log.Fatalf("analysis failed: %v", err)
	}

	analysisJSON, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		log.Fatalf("encode analysis json: %v", err)
	}

	const analysisResultPath = "scripts/analysis_result.json"
	if err := os.WriteFile(analysisResultPath, append(analysisJSON, '\n'), 0644); err != nil {
		log.Fatalf("write analysis result: %v", err)
	}

	fmt.Printf("Analysis result written to %s\n", analysisResultPath)
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
) (*provider.PromptRunResult, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, attemptTimeout)
		output, err := llm.Run(attemptCtx, apiKey, nil, provider.PromptRunRequest{
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

func toAnalyzerCitations(citations []provider.Citation) []analyzer.Citation {
	out := make([]analyzer.Citation, 0, len(citations))
	for _, citation := range citations {
		out = append(out, analyzer.Citation{
			URL:   citation.URL,
			Title: citation.Title,
			Text:  citation.Text,
		})
	}

	return out
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}
