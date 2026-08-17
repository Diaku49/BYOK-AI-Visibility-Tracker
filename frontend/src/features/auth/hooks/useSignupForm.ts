import { useCallback, useState } from 'react'
import { useNavigate } from 'react-router'
import { ROUTES } from '@/shared/config/routes'
import { authApi } from '@/features/auth/api/authApi'
import {
  validateEmail,
  validateName,
  validatePassword,
  describeAuthError,
} from '@/features/auth/validation'
import type { AuthErrors, SignupPayload } from '@/features/auth/types'

const EMPTY = { email: '', password: '', name: '' }

export function useSignupForm() {
  const [values, setValues] = useState(EMPTY)
  const [errors, setErrors] = useState<AuthErrors>({})
  const [isSubmitting, setIsSubmitting] = useState(false)

  const navigate = useNavigate()

  const setField = useCallback((field: keyof typeof EMPTY, value: string) => {
    setValues(current => ({ ...current, [field]: value }))
    // Clear the field error as the user corrects it, and drop any stale
    // submission-level error so it cannot outlive the input it described.
    setErrors(current => ({ ...current, [field]: undefined, form: undefined }))
  }, [])

  const submit = useCallback(
    async (event: React.FormEvent) => {
      event.preventDefault()
      if (isSubmitting) return

      const nextErrors: AuthErrors = {
        email: validateEmail(values.email),
        password: validatePassword(values.password),
        name: validateName(values.name),
      }

      if (nextErrors.email || nextErrors.password || nextErrors.name) {
        setErrors(nextErrors)
        return
      }

      setErrors({})
      setIsSubmitting(true)

      // Omit name entirely when blank rather than sending an empty string.
      const trimmedName = values.name.trim()
      const email = values.email.trim()
      const payload: SignupPayload = {
        email,
        password: values.password,
        ...(trimmedName && { name: trimmedName }),
      }

      try {
        const { message } = await authApi.signup(payload)

        // Signup returns `data: null` — there is no token, so no session can
        // start here. Hand off to login carrying the server's message.
        navigate(ROUTES.login, {
          replace: true,
          state: { notice: message, email },
        })
      } catch (error) {
        setErrors({ form: describeAuthError(error, 'signup') })
      } finally {
        setIsSubmitting(false)
      }
    },
    [values, isSubmitting, navigate],
  )

  return { values, errors, isSubmitting, setField, submit }
}
