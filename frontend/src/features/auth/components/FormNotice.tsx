/** Success/info banner. Mirrors FormError but reads as confirmation. */
export default function FormNotice({ message }: { message?: string }) {
  if (!message) return null

  return (
    <p
      role="status"
      className="flex items-start gap-2 rounded-sharp border border-ok/40 bg-ok/5 px-3 py-2 font-mono text-[11px] leading-relaxed text-ok"
    >
      <span aria-hidden="true">✓</span>
      <span>{message}</span>
    </p>
  )
}
