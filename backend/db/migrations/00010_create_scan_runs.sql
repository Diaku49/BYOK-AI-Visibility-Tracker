-- +goose Up
CREATE TABLE IF NOT EXISTS scan_runs (
    id UUID PRIMARY KEY,

    scan_id UUID NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    engine_id TEXT NOT NULL REFERENCES engines(id) ON DELETE RESTRICT,
    prompt_id UUID NOT NULL REFERENCES prompts(id) ON DELETE CASCADE,

    provider_key_id UUID NOT NULL REFERENCES provider_keys(id) ON DELETE RESTRICT,

    try_number INT NOT NULL DEFAULT 1 CHECK (try_number > 0),

    status TEXT NOT NULL DEFAULT 'pending' CHECK (
        status IN ('pending', 'running', 'completed', 'failed', 'canceled')
    ),

    -- provider output
    answer_text TEXT,
    raw_response JSONB,

    -- analyzed result
    brand_mentioned BOOLEAN,
    brand_domain_cited BOOLEAN,

    competitors_mentioned JSONB,
    cited_domains JSONB,

    error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,

    UNIQUE (scan_id, engine_id, prompt_id, try_number)
);

-- +goose Down
DROP TABLE IF EXISTS scan_runs;