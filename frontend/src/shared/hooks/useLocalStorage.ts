import { useCallback, useEffect, useState } from 'react'

function read<T>(key: string, fallback: T): T {
  try {
    const raw = window.localStorage.getItem(key)
    return raw === null ? fallback : (JSON.parse(raw) as T)
  } catch {
    return fallback
  }
}

/**
 * State persisted to localStorage. Reads are lazy, writes are best-effort, and
 * failures fall back to in-memory state so private mode or a full quota can
 * never break rendering.
 */
export function useLocalStorage<T>(key: string, initialValue: T) {
  const [stored, setStored] = useState<T>(() => read(key, initialValue))

  const setValue = useCallback(
    (value: T | ((previous: T) => T)) => {
      setStored(previous => {
        const next = value instanceof Function ? value(previous) : value
        try {
          window.localStorage.setItem(key, JSON.stringify(next))
        } catch {
          // Storage unavailable or full — keep in-memory state authoritative.
        }
        return next
      })
    },
    [key],
  )

  const remove = useCallback(() => {
    try {
      window.localStorage.removeItem(key)
    } catch {
      // Nothing to do — fall through to resetting in-memory state.
    }
    setStored(initialValue)
  }, [key, initialValue])

  // Another tab writing this key updates this one, so logging out in one tab
  // cannot leave a second tab believing it is still authenticated.
  useEffect(() => {
    function handleStorage(event: StorageEvent) {
      if (event.key !== key || event.storageArea !== window.localStorage) return
      setStored(event.newValue === null ? initialValue : read(key, initialValue))
    }

    window.addEventListener('storage', handleStorage)
    return () => window.removeEventListener('storage', handleStorage)
  }, [key, initialValue])

  return [stored, setValue, remove] as const
}
