import type { ReactNode } from 'react'

interface AuthShellProps {
  /** Terminal path shown in the window chrome, e.g. "auth/login". */
  command: string
  title: string
  subtitle: string
  children: ReactNode
  footer: ReactNode
}

/**
 * Terminal-window chrome shared by the login and signup forms.
 * Layout and framing only — no form state.
 */
export default function AuthShell({
  command,
  title,
  subtitle,
  children,
  footer,
}: AuthShellProps) {
  return (
    <div className="w-full max-w-md animate-fade-in rounded-sharp border border-line bg-surface">
      {/* Window bar */}
      <div className="flex h-9 items-center gap-2 border-b border-line bg-void-deep px-3">
        <span className="h-2 w-2 rounded-full bg-danger" />
        <span className="h-2 w-2 rounded-full bg-warn" />
        <span className="h-2 w-2 rounded-full bg-ok" />
        <span className="ml-2 truncate font-mono text-[10px] tracking-[0.15em] text-ink-faint">
          byok@tracker: ~/{command}
        </span>
      </div>

      <div className="flex flex-col gap-6 p-6 sm:p-8">
        <header className="flex flex-col gap-1.5">
          <h1 className="font-mono text-xl font-extrabold tracking-tight text-ink">
            {title}
            <span className="caret ml-1.5 text-accent" aria-hidden="true" />
          </h1>
          <div className="h-px w-12 bg-accent" />
          <p className="mt-1 text-sm text-ink-dim">{subtitle}</p>
        </header>

        {children}
      </div>

      <div className="border-t border-line px-6 py-4 text-center sm:px-8">{footer}</div>
    </div>
  )
}
