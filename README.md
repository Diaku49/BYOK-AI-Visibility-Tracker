# BYOK AI Visibility Tracker

An API-first Go prototype for measuring how visible a brand is in AI-generated answers.

Customers configure a project with their brand, domain, competitors, prompts, and their own AI-provider keys. The service runs the prompts through supported providers, analyzes each answer, and stores visibility results for comparison over time.

This is not a chat wrapper. Its purpose is to answer questions such as:

- Was the brand mentioned in an AI answer?
- Which configured competitors were mentioned?
- Was the brand domain cited?
- How does visibility differ across prompts and providers?

## Current Prototype

The current implementation supports Gemini and OpenAI providers, PostgreSQL storage, and a background worker pipeline.

```text
Create scan
  -> create pending scan runs
  -> worker claims and executes each provider call
  -> store completed or failed run results
  -> queue scan analysis
  -> analysis worker claims, analyzes, and saves the final scan summary
```

Provider keys are supplied by customers and encrypted before storage. Workers use database-backed status transitions so pending work can be recovered after an interrupted process.

## Prototype Boundaries

This repository is intentionally early-stage. It currently focuses on backend and worker foundations rather than a complete product experience.

- API endpoints and result retrieval are still being built out.
- A dashboard and reporting UI are not included.
- Provider quotas, retries for analysis, and richer reporting are future work.
- The worker flow is under active development; production hardening and test coverage are still needed.

## Project Layout

- `cmd/api` - application entrypoint
- `config` - environment-backed configuration
- `db/migrations` - PostgreSQL migrations
- `db/queries` - sqlc query definitions
- `internal/db` - generated sqlc code
- `internal/store` - database operations and transactions
- `internal/api` - HTTP server, handlers, and authentication
- `internal/provider` - Gemini and OpenAI provider implementations
- `internal/analyzer` - structured visibility-analysis contracts
- `internal/worker` - scan execution and analysis workers
- `scripts` - manual provider and analysis test scripts

## Development

Generate database code after changing a migration or query:

```sh
sqlc generate
```

Run the available checks:

```sh
go test ./...
go vet ./...
sqlc vet
```

## Docker

The API container expects PostgreSQL to be available through `DATABASE_URL`. It does not run database migrations automatically.

Build the image:

```sh
docker build -t ai-visibility-tracker .
```

Run the API:

```sh
docker run --rm -p 8080:8080 \
  -e DATABASE_URL='postgres://user:password@host:5432/database?sslmode=disable' \
  -e AITRACKER_MASTER_KEY='base64-encoded-32-byte-key' \
  -e AITRACKER_JWT_SECRET='replace-with-a-secure-secret' \
  ai-visibility-tracker
```

`AITRACKER_PORT` is optional and defaults to `8080`. The container exposes port `8080`; map a different host port with Docker if needed.
