-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,

    email TEXT NOT NULL UNIQUE,
    password TEXT,
    name TEXT,

    tier_name TEXT NOT NULL DEFAULT 'free' REFERENCES tiers(name) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;
