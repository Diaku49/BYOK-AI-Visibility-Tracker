/** Submission-level error banner. Field errors render inside their Input. */
export default function FormError({ message }: { message?: string }) {
  if (!message) return null

  return (
    <p
      role="alert"
      className="flex items-start gap-2 rounded-sharp border border-danger/40 bg-danger/5 px-3 py-2 font-mono text-[11px] leading-relaxed text-danger"
    >
      <span aria-hidden="true">!</span>
      <span>{message}</span>
    </p>
  )
}
