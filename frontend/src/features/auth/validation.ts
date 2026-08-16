import { ApiError } from '@/shared/api/client'

/** Client-side validation. The backend remains the authority. */

export const PASSWORD_MIN_LENGTH = 8

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function validateEmail(value: string): string | undefined {
  const email = value.trim()
  if (!email) return 'Email is required'
  if (!EMAIL_PATTERN.test(email)) return 'Enter a valid email address'
  return undefined
}

export function validatePassword(value: string): string | undefined {
  if (!value) return 'Password is required'
  if (value.length < PASSWORD_MIN_LENGTH) {
    return `Minimum ${PASSWORD_MIN_LENGTH} characters`
  }
  return undefined
}

/** Name is optional — only a supplied value is constrained. */
export function validateName(value: string): string | undefined {
  const name = value.trim()
  if (!name) return undefined
  if (name.length < 2) return 'Minimum 2 characters'
  return undefined
}

/**
 * Turns a thrown error into a line for the form banner.
 *
 * The API sends a human `message` in every envelope, so prefer it. Two cases
 * override it: transport failures have no server message at all, and 401/403
 * messages are replaced with a neutral line so a wrong-password response cannot
 * confirm which emails are registered.
 */
export function describeAuthError(error: unknown, mode: 'login' | 'signup'): string {
  if (!(error instanceof ApiError)) {
    return mode === 'login' ? 'Could not sign you in' : 'Could not create your account'
  }

  if (error.isNetworkError) {
    return 'Connection failed. Check your network and retry.'
  }

  if (error.status === 401 || error.status === 403) {
    return 'Invalid email or password'
  }

  if (error.message) {
    return error.message
  }

  if (error.status === 409) return 'An account with this email already exists'
  if (error.status === 422) return 'Check your details and try again'
  if (error.status === 429) return 'Too many attempts. Wait a moment and retry.'
  if (error.status >= 500) return 'Server unavailable. Try again shortly.'

  return mode === 'login' ? 'Could not sign you in' : 'Could not create your account'
}
