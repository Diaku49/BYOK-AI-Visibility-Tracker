package analyzer

import "google.golang.org/genai"

func ScanAnalysisSchema() *genai.Schema {
	stringArray := func(description string) *genai.Schema {
		return &genai.Schema{
			Type:        genai.TypeArray,
			Description: description,
			Items: &genai.Schema{
				Type: genai.TypeString,
			},
		}
	}

	runAnalysisSchema := &genai.Schema{
		Type:        genai.TypeObject,
		Description: "Per-scan-run visibility analysis.",
		Properties: map[string]*genai.Schema{
			"scan_run_id": {
				Type:        genai.TypeString,
				Description: "The exact scan_run_id from the input.",
			},
			"brand_mentioned": {
				Type:        genai.TypeBoolean,
				Description: "Whether the target brand is mentioned in the answer text.",
			},
			"brand_domain_cited": {
				Type:        genai.TypeBoolean,
				Description: "Whether the target brand domain is cited in the citations.",
			},
			"competitors_mentioned": stringArray(
				"Competitor names from the provided competitor list that are mentioned in the answer text.",
			),
			"cited_domains": stringArray(
				"Domains extracted from the citations for this run.",
			),
		},
		Required: []string{
			"scan_run_id",
			"brand_mentioned",
			"brand_domain_cited",
			"competitors_mentioned",
			"cited_domains",
		},
		PropertyOrdering: []string{
			"scan_run_id",
			"brand_mentioned",
			"brand_domain_cited",
			"competitors_mentioned",
			"cited_domains",
		},
	}

	competitorSummarySchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"name": {
				Type: genai.TypeString,
			},
			"count": {
				Type: genai.TypeInteger,
			},
		},
		Required: []string{"name", "count"},
	}

	domainSummarySchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"domain": {
				Type: genai.TypeString,
			},
			"count": {
				Type: genai.TypeInteger,
			},
		},
		Required: []string{"domain", "count"},
	}

	summarySchema := &genai.Schema{
		Type:        genai.TypeObject,
		Description: "Aggregate visibility summary for the whole scan.",
		Properties: map[string]*genai.Schema{
			"total_runs": {
				Type: genai.TypeInteger,
			},
			"brand_mentioned_runs": {
				Type: genai.TypeInteger,
			},
			"brand_mention_rate": {
				Type: genai.TypeNumber,
			},
			"brand_domain_cited_runs": {
				Type: genai.TypeInteger,
			},
			"brand_domain_citation_rate": {
				Type: genai.TypeNumber,
			},
			"top_competitors": {
				Type:  genai.TypeArray,
				Items: competitorSummarySchema,
			},
			"top_cited_domains": {
				Type:  genai.TypeArray,
				Items: domainSummarySchema,
			},
			"notes": {
				Type: genai.TypeString,
			},
		},
		Required: []string{
			"total_runs",
			"brand_mentioned_runs",
			"brand_mention_rate",
			"brand_domain_cited_runs",
			"brand_domain_citation_rate",
			"top_competitors",
			"top_cited_domains",
			"notes",
		},
	}

	return &genai.Schema{
		Type:        genai.TypeObject,
		Description: "Complete scan analysis result.",
		Properties: map[string]*genai.Schema{
			"runs": {
				Type:        genai.TypeArray,
				Description: "Analysis for each scan_run.",
				Items:       runAnalysisSchema,
			},
			"summary": summarySchema,
		},
		Required:         []string{"runs", "summary"},
		PropertyOrdering: []string{"runs", "summary"},
	}
}
