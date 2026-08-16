import { useNavigate } from 'react-router-dom'
import { Button } from '@/shared/components'
import { ROUTES } from '@/shared/config/routes'
import { PRODUCT } from '@/features/visibility/constants'

export default function ProductIntro() {
  const navigate = useNavigate()

  return (
    <section className="flex w-full max-w-3xl animate-fade-in flex-col items-start gap-8">
      {/* Status strip — signals "instrument", not "brochure" */}
      <div className="flex items-center gap-3 rounded-sharp border border-line bg-surface px-3 py-1.5">
        <span className="h-1.5 w-1.5 shrink-0 bg-ok" />
        <span className="font-mono text-[10px] tracking-[0.22em] text-ink-dim uppercase">
          {PRODUCT.status}
        </span>
      </div>

      <div className="flex flex-col gap-5">
        <h1 className="font-mono text-4xl leading-[1.05] font-extrabold tracking-tight text-ink sm:text-5xl lg:text-6xl">
          {PRODUCT.name}
          <span className="caret ml-2 text-accent" aria-hidden="true" />
        </h1>

        {/* Accent rule: one deliberate horizontal break */}
        <div className="h-px w-24 bg-accent" />

        <p className="max-w-2xl text-base leading-relaxed text-ink-dim sm:text-lg">
          {PRODUCT.description}
        </p>
      </div>

      <div className="flex flex-col items-start gap-3 sm:flex-row sm:items-center">
        <Button variant="solid" size="lg" onClick={() => navigate(ROUTES.signup)}>
          Get Started
        </Button>
        <span className="font-mono text-[10px] tracking-[0.15em] text-ink-faint uppercase">
          Bring your own keys — no card required
        </span>
      </div>
    </section>
  )
}
