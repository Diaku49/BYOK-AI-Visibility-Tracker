package analyzer

import "github.com/sashabaranov/go-openai/jsonschema"

func OpenAIScanAnalysisSchema() *jsonschema.Definition {
	stringArray := func(description string) jsonschema.Definition {
		return jsonschema.Definition{
			Type:        jsonschema.Array,
			Description: description,
			Items: &jsonschema.Definition{
				Type: jsonschema.String,
			},
		}
	}

	runAnalysisSchema := jsonschema.Definition{
		Type:                 jsonschema.Object,
		Description:          "Per-scan-run visibility analysis.",
		AdditionalProperties: false,
		Properties: map[string]jsonschema.Definition{
			"scan_run_id": {
				Type:        jsonschema.String,
				Description: "The exact scan_run_id from the input.",
			},
			"brand_mentioned": {
				Type:        jsonschema.Boolean,
				Description: "Whether the target brand is mentioned in the answer text.",
			},
			"brand_domain_cited": {
				Type:        jsonschema.Boolean,
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
	}

	competitorSummarySchema := jsonschema.Definition{
		Type:                 jsonschema.Object,
		AdditionalProperties: false,
		Properties: map[string]jsonschema.Definition{
			"name": {
				Type: jsonschema.String,
			},
			"count": {
				Type: jsonschema.Integer,
			},
		},
		Required: []string{"name", "count"},
	}

	domainSummarySchema := jsonschema.Definition{
		Type:                 jsonschema.Object,
		AdditionalProperties: false,
		Properties: map[string]jsonschema.Definition{
			"domain": {
				Type: jsonschema.String,
			},
			"count": {
				Type: jsonschema.Integer,
			},
		},
		Required: []string{"domain", "count"},
	}

	summarySchema := jsonschema.Definition{
		Type:                 jsonschema.Object,
		Description:          "Aggregate visibility summary for the whole scan.",
		AdditionalProperties: false,
		Properties: map[string]jsonschema.Definition{
			"total_runs": {
				Type: jsonschema.Integer,
			},
			"brand_mentioned_runs": {
				Type: jsonschema.Integer,
			},
			"brand_mention_rate": {
				Type: jsonschema.Number,
			},
			"brand_domain_cited_runs": {
				Type: jsonschema.Integer,
			},
			"brand_domain_citation_rate": {
				Type: jsonschema.Number,
			},
			"top_competitors": {
				Type:  jsonschema.Array,
				Items: &competitorSummarySchema,
			},
			"top_cited_domains": {
				Type:  jsonschema.Array,
				Items: &domainSummarySchema,
			},
			"notes": {
				Type: jsonschema.String,
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

	return &jsonschema.Definition{
		Type:                 jsonschema.Object,
		Description:          "Complete scan analysis result.",
		AdditionalProperties: false,
		Properties: map[string]jsonschema.Definition{
			"runs": {
				Type:        jsonschema.Array,
				Description: "Analysis for each scan_run.",
				Items:       &runAnalysisSchema,
			},
			"summary": summarySchema,
		},
		Required: []string{"runs", "summary"},
	}
}
