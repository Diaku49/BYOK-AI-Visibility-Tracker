-- +goose Up
CREATE TABLE IF NOT EXISTS competitors (
    id UUID PRIMARY KEY,

    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    name TEXT NOT NULL,
    domain TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (project_id, name)
);


-- +goose Down
DROP TABLE IF EXISTS competitors;
