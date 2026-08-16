# BYOK AI Visibility Tracker

An API-first prototype for measuring how visible a brand is in AI-generated answers.

Companies provide their brand, domain, competitors, prompts, and their own AI-provider keys. The service runs prompts through supported AI providers, analyzes the answers, and stores visibility results that can be compared over time.

This is not a chat wrapper. It is the foundation for answering questions such as:

- Does an AI answer mention the brand?
- Which competitors appear in the answer?
- Is the brand's domain cited?
- How does visibility differ by prompt or AI provider?

## How It Works

```text
Project configuration
  -> scan creation
  -> provider-specific scan runs
  -> structured answer analysis
  -> stored visibility results and scan summary
```

Provider keys are customer-supplied and encrypted before storage. Background workers use durable PostgreSQL state transitions so work can be recovered after a process interruption.

## Repository Layout

```text
backend/   Go API, PostgreSQL access, provider integrations, and workers
frontend/  Planned web dashboard
```

The repository is structured as a monorepo. Each application owns its own dependencies and development commands.

## Current Status

This is an early prototype focused on backend foundations.

- The backend supports Gemini and OpenAI provider integrations.
- PostgreSQL, sqlc queries, scan workers, and analysis workers are implemented.
- API coverage and result retrieval are still evolving.
- The frontend dashboard, reporting, and production hardening are future work.

## Getting Started

See the [backend README](backend/README.md) for local setup, environment variables, migrations, API routes, testing, and Docker usage.
