# Backend

The backend is an API-first Go service for the BYOK AI Visibility Tracker. It owns authentication, projects, encrypted provider keys, scan creation, background workers, PostgreSQL persistence, and AI-based scan analysis.

The HTTP server and workers start in the same process. A scan is created synchronously, then its provider calls and analysis run asynchronously through database-backed worker state.

## Responsibilities

- Manage users, projects, prompts, competitors, and provider keys.
- Encrypt customer-supplied provider keys before storing them in PostgreSQL.
- Create one `scan_run` for each active prompt and configured provider key.
- Execute pending scan runs with Gemini or OpenAI.
- Analyze completed responses for brand mentions, competitor mentions, and cited domains.
- Persist per-run analysis and a scan-level summary.

## Scan Flow

```text
POST /scan/{projectID}
  -> create a pending scan and its scan runs
  -> scan workers claim and execute pending runs
  -> completed or failed runs are stored
  -> analysis worker claims the finished scan
  -> per-run analysis and scan summary are stored
```

Worker state is stored in PostgreSQL. A worker can recover eligible work after an interrupted process rather than depending on in-memory jobs alone.

## Requirements

- Go `1.26.4` or the version declared in `go.mod`
- PostgreSQL
- [sqlc](https://sqlc.dev/) for generated query code and query validation
- [goose](https://github.com/pressly/goose) for applying migrations

## Configuration

The service reads `DATABASE_URL` directly and uses the `AITRACKER_` prefix for configuration fields.

| Variable | Required | Description |
| --- | --- | --- |
| `DATABASE_URL` | Yes | PostgreSQL connection string. |
| `AITRACKER_MASTER_KEY` | Yes | Base64-encoded 32-byte key used to encrypt provider keys. |
| `AITRACKER_JWT_SECRET` | Recommended | JWT signing secret. The development default is not safe for deployed environments. |
| `AITRACKER_PORT` | No | HTTP port; defaults to `8080`. |

Keep real values in a local `.env` file or your deployment secret manager. Never commit provider keys, database credentials, or the master key.

## Local Development

From this directory:

```sh
go run ./cmd/api
```

Run checks:

```sh
go test ./...
go vet ./...
sqlc vet
```

Regenerate sqlc code after changing a query definition:

```sh
sqlc generate
```

Apply migrations with Goose before starting the API:

```sh
goose -dir db/migrations postgres "$DATABASE_URL" up
```

## HTTP API

Implemented routes currently include:

| Method | Route | Authentication |
| --- | --- | --- |
| `POST` | `/user` | No |
| `POST` | `/user/login` | No |
| `POST` | `/key` | Bearer token |
| `POST` | `/project` | Bearer token |
| `GET` | `/project` | Bearer token |
| `GET` | `/project/{projectID}` | Bearer token |
| `POST` | `/scan/{projectID}` | Bearer token |

Some routes are registered as placeholders while the prototype is still evolving. Result retrieval and the frontend dashboard are not complete yet.

## Project Layout

- `cmd/api` - application entrypoint
- `config` - environment-backed configuration
- `db/migrations` - PostgreSQL migrations
- `db/queries` - sqlc query definitions
- `internal/db` - generated sqlc code
- `internal/store` - database operations and transactions
- `internal/api` - HTTP server, handlers, and authentication
- `internal/provider` - Gemini and OpenAI integrations
- `internal/analyzer` - structured visibility-analysis contracts
- `internal/worker` - scan execution and analysis workers
- `scripts` - manual provider and analysis test scripts

## Docker

The container expects PostgreSQL to be reachable through `DATABASE_URL`. It does not apply migrations automatically.

Build the image from this directory:

```sh
docker build -t ai-visibility-tracker-backend .
```

Run the API:

```sh
docker run --rm -p 8080:8080 \
  -e DATABASE_URL='postgres://user:password@host:5432/database?sslmode=disable' \
  -e AITRACKER_MASTER_KEY='base64-encoded-32-byte-key' \
  -e AITRACKER_JWT_SECRET='replace-with-a-secure-secret' \
  ai-visibility-tracker-backend
```

From the repository root, use `backend` as the build context:

```sh
docker build -f backend/Dockerfile -t ai-visibility-tracker-backend ./backend
```
