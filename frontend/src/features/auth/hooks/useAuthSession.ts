import { useCallback } from 'react'
import { useLocalStorage } from '@/shared/hooks'
import { STORAGE_KEYS } from '@/shared/config/storage'
import type { AuthSession } from '@/features/auth/types'

interface UseAuthSessionResult extends AuthSession {
  login: (token: string) => void
  logout: () => void
}

/**
 * Owns the persisted session. The only place that writes the auth keys —
 * shared/config/storage.ts holds the key names so the API client can read the
 * token without importing from this feature.
 */
export function useAuthSession(): UseAuthSessionResult {
  const [loggedInFlag, setLoggedInFlag, clearLoggedInFlag] = useLocalStorage<boolean>(
    STORAGE_KEYS.authLoggedIn,
    false,
  )
  const [token, setToken, clearToken] = useLocalStorage<string | null>(
    STORAGE_KEYS.authToken,
    null,
  )

  // Both keys are persisted, but the session only counts as active when they
  // agree — a leftover flag with no token must not unlock the UI.
  const isLoggedIn = loggedInFlag && token !== null

  const login = useCallback(
    (nextToken: string) => {
      setToken(nextToken)
      setLoggedInFlag(true)
    },
    [setToken, setLoggedInFlag],
  )

  const logout = useCallback(() => {
    clearToken()
    clearLoggedInFlag()
  }, [clearToken, clearLoggedInFlag])

  return { isLoggedIn, token, login, logout }
}
