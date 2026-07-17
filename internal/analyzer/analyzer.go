package analyzer

import (
	"context"

	"github.com/google/uuid"
)

type Analyzer interface {
	AnalyzeScan(ctx context.Context, apiKey string, input ScanAnalysisInput) (*ScanAnalysisResult, error)
}

// Input
type ScanAnalysisInput struct {
	BrandName   string
	BrandDomain string

	Competitors []CompetitorForAnalysis
	Runs        []RunForAnalysis
}

type CompetitorForAnalysis struct {
	Name   string
	Domain string
}

type RunForAnalysis struct {
	ScanRunID  uuid.UUID
	EngineID   string
	Prompt     string
	AnswerText string
	Citations  []Citation
}

type Citation struct {
	URL   string
	Title string
	Text  string
}

// Response
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
