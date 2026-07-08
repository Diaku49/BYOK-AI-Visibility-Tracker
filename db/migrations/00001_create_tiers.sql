-- +goose Up
CREATE TABLE IF NOT EXISTS tiers (
    name TEXT PRIMARY KEY, -- free, starter, pro, agency

    price_cents INT NOT NULL DEFAULT 0 CHECK (price_cents >= 0),

    max_projects INT NOT NULL DEFAULT 1 CHECK (max_projects > 0),
    max_prompts INT NOT NULL DEFAULT 2 CHECK (max_prompts > 0),
    max_engines INT NOT NULL DEFAULT 1 CHECK (max_engines > 0),
    max_scans_per_month INT NOT NULL DEFAULT 5 CHECK (max_scans_per_month > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO tiers (
    name,
    price_cents,
    max_projects,
    max_prompts,
    max_engines,
    max_scans_per_month
)
VALUES
    ('free', 0, 1, 10, 1, 5),
    ('starter', 900, 2, 50, 2, 30),
    ('pro', 2900, 10, 200, 4, 300),
    ('agency', 7900, 50, 1000, 4, 1000)
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS tiers;