-- +goose Up
CREATE TABLE IF NOT EXISTS engines (
    id TEXT PRIMARY KEY, -- gemini, openai, perplexity, xai

    name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO engines (id, name, active)
VALUES
    ('gemini', 'Gemini', true),
    ('openai', 'OpenAI', false),
    ('perplexity', 'Perplexity', false),
    ('xai', 'Grok / xAI', false)
ON CONFLICT (id) DO NOTHING;


-- +goose Down
DROP TABLE IF EXISTS engines;
