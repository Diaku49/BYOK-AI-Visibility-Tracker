import { forwardRef, type ButtonHTMLAttributes } from 'react'
import { cn } from '@/shared/utils'

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'solid' | 'outline' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
}

const variants: Record<NonNullable<ButtonProps['variant']>, string> = {
  solid:
    'bg-accent text-void border border-accent hover:bg-void hover:text-accent active:translate-y-px',
  outline:
    'bg-transparent text-ink border border-line hover:border-accent hover:text-accent active:translate-y-px',
  ghost: 'bg-transparent text-ink-dim border border-transparent hover:text-accent',
}

const sizes: Record<NonNullable<ButtonProps['size']>, string> = {
  sm: 'h-8 px-3 text-[11px]',
  md: 'h-10 px-5 text-xs',
  lg: 'h-12 px-7 text-sm',
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { variant = 'solid', size = 'md', className, children, ...props },
  ref,
) {
  return (
    <button
      ref={ref}
      className={cn(
        'inline-flex cursor-pointer items-center justify-center gap-2 rounded-sharp',
        'font-mono font-bold tracking-[0.12em] whitespace-nowrap uppercase',
        'transition-colors duration-150',
        'disabled:pointer-events-none disabled:opacity-40',
        variants[variant],
        sizes[size],
        className,
      )}
      {...props}
    >
      {children}
    </button>
  )
})

export default Button
