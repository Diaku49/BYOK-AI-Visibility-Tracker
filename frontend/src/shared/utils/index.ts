export function cn(...classes: (string | number | boolean | undefined | null)[]): string {
  return classes.filter(Boolean).join(' ')
}
