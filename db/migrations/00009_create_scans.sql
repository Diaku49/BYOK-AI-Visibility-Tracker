-- +goose Up
CREATE TABLE IF NOT EXISTS scans (
    id UUID PRIMARY KEY,

    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    status TEXT NOT NULL DEFAULT 'pending' CHECK (
        status IN ('pending', 'running', 'completed', 'failed', 'canceled')
    ),

    tries_per_prompt INT NOT NULL DEFAULT 1 CHECK (tries_per_prompt > 0),

    total_runs INT NOT NULL DEFAULT 0 CHECK (total_runs >= 0),
    completed_runs INT NOT NULL DEFAULT 0 CHECK (completed_runs >= 0),
    failed_runs INT NOT NULL DEFAULT 0 CHECK (failed_runs >= 0),

    summary JSONB,
    error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS scans;
