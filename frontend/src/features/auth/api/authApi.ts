import { apiClient, ApiError } from '@/shared/api/client'
import type {
  AuthResult,
  LoginData,
  LoginPayload,
  SignupPayload,
  SignupResult,
} from '@/features/auth/types'

const ENDPOINTS = {
  login: '/auth/login',
  signup: '/auth/signup',
} as const

export const authApi = {
  /**
   * POST /auth/login -> { message, data: { token } }
   * Sends no Authorization header — there is no session yet.
   */
  async login(payload: LoginPayload): Promise<AuthResult> {
    const { message, data } = await apiClient.post<LoginData>(ENDPOINTS.login, payload, {
      auth: false,
    })

    const token = data?.token
    if (typeof token !== 'string' || token.length === 0) {
      // 2xx with no usable token means the contract broke upstream — fail
      // loudly rather than persisting an empty session.
      throw new ApiError('Login succeeded but no token was returned', 502)
    }

    return { message: message || 'Authenticated', token }
  },

  /**
   * POST /auth/signup -> { message, data: null }
   * No session is established; the user logs in afterwards.
   */
  async signup(payload: SignupPayload): Promise<SignupResult> {
    const { message } = await apiClient.post<null>(ENDPOINTS.signup, payload, { auth: false })
    return { message: message || 'Account created' }
  },
}
