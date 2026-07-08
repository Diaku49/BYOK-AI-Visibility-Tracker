-- +goose Up
CREATE TABLE IF NOT EXISTS project_engines (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    engine_id TEXT NOT NULL REFERENCES engines(id) ON DELETE RESTRICT,

    provider_key_id UUID NOT NULL REFERENCES provider_keys(id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (project_id, engine_id)
);


-- +goose Down
DROP TABLE IF EXISTS project_engines;
