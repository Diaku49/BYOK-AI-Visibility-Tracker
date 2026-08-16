-- +goose Up
CREATE TABLE IF NOT EXISTS provider_keys (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    engine_id TEXT NOT NULL REFERENCES engines(id) ON DELETE RESTRICT,

    name TEXT NOT NULL,

    encrypted_key BYTEA NOT NULL,
    key_nonce BYTEA NOT NULL,

    active BOOLEAN NOT NULL DEFAULT true,

    monthly_run_limit INT CHECK (
        monthly_run_limit IS NULL OR monthly_run_limit > 0
    ),
    monthly_runs_used INT NOT NULL DEFAULT 0 CHECK (
        monthly_runs_used >= 0
    ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (user_id, engine_id, name)
);

-- +goose Down
DROP TABLE IF EXISTS provider_keys;
