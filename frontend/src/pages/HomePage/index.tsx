import { MenuButton, SideDrawer } from '@/shared/components'
import { useDisclosure } from '@/shared/hooks'
import { AuthMenu } from '@/features/auth'
import { ProductIntro } from '@/features/visibility'

export default function HomePage() {
  const drawer = useDisclosure()

  return (
    <div className="relative flex min-h-screen flex-col bg-void">
      <div className="scanlines" aria-hidden="true" />

      <header className="sticky top-0 z-30 flex h-14 shrink-0 items-center justify-between border-b border-line bg-void/85 px-4 backdrop-blur-sm sm:px-6">
        <div className="flex items-center gap-4">
          <MenuButton onClick={drawer.toggle} isOpen={drawer.isOpen} />
          <span className="font-mono text-sm font-bold tracking-[0.2em] text-accent">▣ BYOK</span>
        </div>
        <AuthMenu />
      </header>

      <SideDrawer isOpen={drawer.isOpen} onClose={drawer.close} />

      <main className="grid-surface flex flex-1 items-center px-4 py-20 sm:px-6 sm:py-28">
        <div className="mx-auto w-full max-w-6xl">
          <ProductIntro />
        </div>
      </main>

      <footer className="shrink-0 border-t border-line px-4 py-4 sm:px-6">
        <p className="font-mono text-[10px] tracking-[0.15em] text-ink-faint uppercase">
          BYOK AI Visibility Tracker — v0.1.0
        </p>
      </footer>
    </div>
  )
}
