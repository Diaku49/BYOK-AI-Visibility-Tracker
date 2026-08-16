import type { ReactNode } from 'react'

interface AppProvidersProps {
  children: ReactNode
}

/**
 * Single composition point for app-wide context.
 * Add providers here (QueryClientProvider, ThemeProvider, ...) so App.tsx stays
 * a one-liner and provider order is visible in one place.
 */
export default function AppProviders({ children }: AppProvidersProps) {
  return <>{children}</>
}
