interface MenuButtonProps {
  onClick: () => void
  isOpen: boolean
}

/** Hamburger trigger. Three hairlines that shift to accent on hover. */
export default function MenuButton({ onClick, isOpen }: MenuButtonProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-label="Open menu"
      aria-expanded={isOpen}
      className="group flex h-10 w-10 shrink-0 cursor-pointer flex-col items-center justify-center gap-[5px] rounded-sharp border border-line transition-colors duration-150 hover:border-accent"
    >
      {[0, 1, 2].map(index => (
        <span
          key={index}
          className="block h-px w-4 bg-ink-dim transition-colors duration-150 group-hover:bg-accent"
        />
      ))}
    </button>
  )
}
