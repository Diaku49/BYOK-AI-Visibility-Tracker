import { forwardRef, useId, type InputHTMLAttributes } from 'react'
import { cn } from '@/shared/utils'

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string
  error?: string
  hint?: string
}

const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
  { label, error, hint, className, id, ...props },
  ref,
) {
  const generatedId = useId()
  const inputId = id ?? generatedId
  const errorId = `${inputId}-error`
  const hintId = `${inputId}-hint`

  return (
    <div className="flex w-full flex-col gap-2">
      <div className="flex items-baseline justify-between gap-3">
        <label
          htmlFor={inputId}
          className="font-mono text-[10px] font-bold tracking-[0.2em] text-ink-dim uppercase"
        >
          {label}
        </label>
        {hint && (
          <span id={hintId} className="font-mono text-[10px] tracking-wider text-ink-faint uppercase">
            {hint}
          </span>
        )}
      </div>

      <input
        ref={ref}
        id={inputId}
        aria-invalid={error ? true : undefined}
        aria-describedby={cn(error && errorId, hint && hintId) || undefined}
        className={cn(
          'h-11 w-full rounded-sharp border bg-void px-3',
          'font-mono text-sm text-ink placeholder:text-ink-faint',
          'transition-colors duration-150 outline-none',
          'focus:border-accent focus:bg-void-deep',
          error ? 'border-danger' : 'border-line hover:border-line-bright',
          className,
        )}
        {...props}
      />

      {/* Reserved height stops the form shifting when errors appear */}
      <span
        id={errorId}
        role={error ? 'alert' : undefined}
        className="min-h-[14px] font-mono text-[10px] tracking-wide text-danger"
      >
        {error}
      </span>
    </div>
  )
})

export default Input
