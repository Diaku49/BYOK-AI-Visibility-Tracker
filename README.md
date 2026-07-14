# BYOK AI Visibility Tracker

API-first Go prototype for tracking how visible a brand is in AI-generated answers.

Customers provide their brand name, competitors, prompts, and their own provider API keys. The system runs those prompts through AI/search providers such as OpenAI, Gemini, Grok, or Perplexity, then stores and compares the results over time.

The goal is not to wrap chat APIs. The value is in measuring brand visibility across answer engines:

- whether the customer's brand appears in responses
- whether competitors appear in responses
- which domains are cited
- how visibility changes across prompts, providers, and time

## Current Scope

This repo is currently focused on the backend foundations:

- Go API structure
- PostgreSQL schema and migrations
- sqlc-generated database access
- store-layer methods around generated queries
- user, project, prompt, competitor, provider key, and project engine data models

A dashboard may be added later, but the first version is intended to be API-first.

## Project Layout

- `cmd/api` - API entrypoint
- `config` - application configuration
- `db/migrations` - PostgreSQL migrations
- `db/queries` - sqlc query definitions
- `internal/db` - sqlc-generated Go code
- `internal/store` - hand-written store layer around generated queries
- `internal/api` - HTTP API handlers, middleware, and server setup
- `internal/provider` - provider integration boundary
- `internal/parser` - answer parsing boundary
- `internal/worker` - background processing boundary

## Development

Generated database code is managed by sqlc:

```sh
sqlc generate
```

Common checks:

```sh
go test ./...
go vet ./...
sqlc vet
```
