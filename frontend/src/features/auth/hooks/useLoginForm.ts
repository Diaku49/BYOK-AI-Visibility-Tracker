import { useCallback, useState } from 'react'
import { useLocation, useNavigate } from 'react-router'
import { ROUTES } from '@/shared/config/routes'
import { authApi } from '@/features/auth/api/authApi'
import { useAuthSession } from '@/features/auth/hooks/useAuthSession'
import { validateEmail, validatePassword, describeAuthError } from '@/features/auth/validation'
import type { AuthErrors } from '@/features/auth/types'

/** Set by useSignupForm when it hands off after a successful registration. */
interface LoginLocationState {
  notice?: string
  email?: string
}

export function useLoginForm() {
  const location = useLocation()
  const handoff = (location.state ?? null) as LoginLocationState | null

  const [values, setValues] = useState({ email: handoff?.email ?? '', password: '' })
  const [errors, setErrors] = useState<AuthErrors>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [notice, setNotice] = useState(handoff?.notice ?? '')

  const { login } = useAuthSession()
  const navigate = useNavigate()

  const setField = useCallback((field: 'email' | 'password', value: string) => {
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
      }

      if (nextErrors.email || nextErrors.password) {
        setErrors(nextErrors)
        return
      }

      setErrors({})
      setNotice('')
      setIsSubmitting(true)

      try {
        const { token } = await authApi.login({
          email: values.email.trim(),
          password: values.password,
        })

        // Persist the JWT before navigating so the next render already sees
        // an active session.
        login(token)
        navigate(ROUTES.home, { replace: true })
      } catch (error) {
        setErrors({ form: describeAuthError(error, 'login') })
      } finally {
        setIsSubmitting(false)
      }
    },
    [values, isSubmitting, login, navigate],
  )

  return { values, errors, notice, isSubmitting, setField, submit }
}
