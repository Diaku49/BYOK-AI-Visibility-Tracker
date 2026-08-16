/**
 * localStorage keys. Lives in shared/ so the base API client can read the JWT
 * without importing from features/ (which would invert the layer rules).
 */
export const STORAGE_KEYS = {
  authLoggedIn: 'byok.auth.loggedIn',
  authToken: 'byok.auth.token',
} as const

/**
 * Reads the persisted JWT outside React. useLocalStorage writes JSON, so the
 * raw value is a quoted string — parse it back before use.
 */
export function readAuthToken(): string | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEYS.authToken)
    if (!raw) return null
    const parsed: unknown = JSON.parse(raw)
    return typeof parsed === 'string' && parsed.length > 0 ? parsed : null
  } catch {
    return null
  }
}
