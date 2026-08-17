# BYOK AI Visibility Tracker

Track whether AI search engines recommend your brand. Queries LLM providers with
your own API keys to measure brand mentions, competitor citations, and overall
visibility.

Current state: project skeleton. The structure and build pipeline are in place;
features are empty placeholders and the landing page renders blank.

## Setup

```bash
npm install
cp .env.example .env.local
npm run dev          # http://localhost:5173
```

## Commands

| Command | What it does |
| --- | --- |
| `npm run dev` | Dev server on port 5173 |
| `npm run build` | `tsc -b`, then `vite build` |
| `npx tsc -b` | Type-check only |
| `npm run preview` | Serve the production build |

`npm run lint` is declared in `package.json` but not functional — `eslint` is not
installed and there is no config. There is no test framework yet.

## Stack

React 19, TypeScript, Vite 6, Tailwind CSS 4, React Router 8.

## Structure

```
src/
├── app/        # entry shell: App, providers, router
├── pages/      # one folder per route, thin composition only
├── features/   # self-contained domains (auth, billing, products, visibility)
├── shared/     # generic UI, hooks, utils, types, base API client
├── styles/
└── main.tsx
```

Dependencies flow one way: `app → pages → features → shared`. Each feature is
reachable only through its `index.ts`. See AGENTS.md for the full conventions.
