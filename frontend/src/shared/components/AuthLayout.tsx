import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { ROUTES } from '@/shared/config/routes'

/** Centered page frame for the auth routes: scanlines, grid, brand mark home. */
export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="relative flex min-h-screen flex-col bg-void">
      <div className="scanlines" aria-hidden="true" />

      <header className="flex h-14 shrink-0 items-center border-b border-line px-4 sm:px-6">
        <Link
          to={ROUTES.home}
          className="font-mono text-sm font-bold tracking-[0.2em] text-accent transition-opacity hover:opacity-70"
        >
          ▣ BYOK
        </Link>
      </header>

      <main className="grid-surface flex flex-1 items-center justify-center px-4 py-12 sm:px-6">
        {children}
      </main>
    </div>
  )
}
