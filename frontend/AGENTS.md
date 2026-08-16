# AGENTS.md

Conventions for this repository. Read before adding files.

## Commands

```bash
npm run dev      # dev server on port 5173
npm run build    # tsc -b, then vite build
npx tsc -b       # type-check only
npm run preview  # serve the production build
```

`npm run lint` is declared in `package.json` but does not work: `eslint` is not
installed and there is no config. Use `npx tsc -b` to verify. There is no test
framework — do not assume one exists.

## Layer rules

Four layers, one direction of dependency:

```
app → pages → features → shared
```

- `app/` — entry, providers, router. Composition only, no business logic.
- `pages/` — one folder per route, `index.tsx` inside. Thin: compose features,
  read route params, set page layout. No data fetching, no business rules.
- `features/` — self-contained business domains. All logic lives in `hooks/`.
- `shared/` — feature-agnostic. Generic UI, generic hooks, utils, base API client.

Never import upward: `shared/` must not import from `features/`, and `features/`
must not import from `pages/` or `app/`.

## Feature boundaries

Each feature has this exact shape:

```
features/<name>/
├── index.ts        # PUBLIC API — the only file other code may import
├── components/     # presentation, feature-specific
├── hooks/          # ALL logic — no separate "containers"
├── api/            # feature-scoped API calls, built on shared/api/client
└── types.ts        # domain shapes
```

`index.ts` is the boundary. Other code imports `@/features/auth`, never
`@/features/auth/hooks/useLoginForm`. Inside a feature, deep imports are fine.

Features do not import each other. If two need the same thing, it belongs in
`shared/`, or one of them owns it and exposes it via `index.ts`.

## Logic placement

Logic goes in a hook, not a component. Components take props and render.
A component with `useState` + a `fetch` + validation is doing three jobs — pull
them into a hook in the same feature and leave the JSX behind.

## Imports

Always `@/...` (mapped to `src/*` in both `tsconfig.json` and `vite.config.ts`).
No relative parent traversal (`../../`).

## Styling

Tailwind v4, configured in `src/styles/index.css`. Design tokens belong in the
`@theme` block there — consume them as utilities (`bg-surface`, `text-accent`),
never as hardcoded hex in components.

## API layer

`shared/api/client.ts` owns fetch, credentials, and error shape. Feature `api/`
folders call `apiClient`, never `fetch` directly.

BYOK note: user LLM API keys never touch the browser — no localStorage, no
client-side headers. Proxy provider calls through the backend.
