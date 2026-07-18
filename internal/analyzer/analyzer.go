package analyzer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Analyzer interface {
	AnalyzeScan(ctx context.Context, apiKey string, input ScanAnalysisInput) (*ScanAnalysisResult, error)
}

type ScanAnalysisInput struct {
	ScanID      uuid.UUID `json:"scan_id"`
	BrandName   string    `json:"brand_name"`
	BrandDomain string    `json:"brand_domain"`

	Competitors []CompetitorForAnalysis `json:"competitors"`
	Runs        []RunForAnalysis        `json:"runs"`
}

type CompetitorForAnalysis struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type RunForAnalysis struct {
	ScanRunID  uuid.UUID  `json:"scan_run_id"`
	EngineID   string     `json:"engine_id"`
	PromptText string     `json:"prompt_text"`
	AnswerText string     `json:"answer_text"`
	Citations  []Citation `json:"citations"`
}

type Citation struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ScanAnalysisResult struct {
	Runs    []RunAnalysisResult `json:"runs"`
	Summary ScanSummary         `json:"summary"`
}

type RunAnalysisResult struct {
	ScanRunID uuid.UUID `json:"scan_run_id"`

	BrandMentioned       bool     `json:"brand_mentioned"`
	BrandDomainCited     bool     `json:"brand_domain_cited"`
	CompetitorsMentioned []string `json:"competitors_mentioned"`
	CitedDomains         []string `json:"cited_domains"`
}

type ScanSummary struct {
	TotalRuns int `json:"total_runs"`

	BrandMentionedRuns      int     `json:"brand_mentioned_runs"`
	BrandMentionRate        float64 `json:"brand_mention_rate"`
	BrandDomainCitedRuns    int     `json:"brand_domain_cited_runs"`
	BrandDomainCitationRate float64 `json:"brand_domain_citation_rate"`

	TopCompetitors  []CompetitorSummary `json:"top_competitors"`
	TopCitedDomains []DomainSummary     `json:"top_cited_domains"`

	Notes string `json:"notes"`
}

type CompetitorSummary struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type DomainSummary struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

type ScanAnalysisRequest = ScanAnalysisInput
type ScanAnalysisResponse = ScanAnalysisResult

func BuildAnalyzerPrompt(input ScanAnalysisInput) (string, error) {
	b, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", err
	}

	return "Analyze this scan input and return the structured JSON result:\n\n" + string(b), nil
}
