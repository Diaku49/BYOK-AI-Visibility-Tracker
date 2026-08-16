import { useEffect, type ReactNode } from 'react'

interface SideDrawerProps {
  isOpen: boolean
  onClose: () => void
  title?: string
  children?: ReactNode
}

/**
 * Off-canvas panel. Fixed width on desktop so it stays a narrow rail on wide
 * screens; percentage on mobile where a fixed rail would crowd the viewport.
 * Presentation only — the owner holds open state.
 */
export default function SideDrawer({ isOpen, onClose, title = 'Menu', children }: SideDrawerProps) {
  useEffect(() => {
    if (!isOpen) return

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') onClose()
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [isOpen, onClose])

  if (!isOpen) return null

  return (
    <>
      <div
        data-testid="drawer-backdrop"
        onClick={onClose}
        className="fixed inset-0 z-40 bg-void-deep/60"
      />

      <aside
        role="dialog"
        aria-modal="true"
        aria-label={title}
        className="fixed inset-y-0 left-0 z-50 flex w-[72%] animate-drawer-in flex-col border-r border-line bg-void-deep sm:w-64 lg:w-72"
      >
        <header className="flex h-14 shrink-0 items-center justify-between border-b border-line px-4">
          <h2 className="font-mono text-[11px] font-bold tracking-[0.2em] text-ink-dim uppercase">
            {title}
          </h2>
          <button
            type="button"
            onClick={onClose}
            aria-label="Close menu"
            className="flex h-7 w-7 cursor-pointer items-center justify-center rounded-sharp border border-line font-mono text-xs text-ink-dim transition-colors duration-150 hover:border-danger hover:text-danger"
          >
            ✕
          </button>
        </header>

        <nav className="flex-1 overflow-y-auto p-4">{children}</nav>

        <footer className="shrink-0 border-t border-line px-4 py-3">
          <p className="font-mono text-[10px] tracking-[0.15em] text-ink-faint uppercase">
            No routes wired
          </p>
        </footer>
      </aside>
    </>
  )
}
