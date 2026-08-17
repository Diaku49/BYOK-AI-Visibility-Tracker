import { apiClient, ApiError } from '@/shared/api/client'
import type {
  AuthResult,
  LoginPayload,
  SignupPayload,
  SignupResult,
} from '@/features/auth/types'

const ENDPOINTS = {
  signup: '/user',
  login: '/user/login',
} as const

export const authApi = {
  /**
   * POST /user/login -> { message, data: "<jwt>" }
   * Sends no Authorization header — there is no session yet.
   */
  async login(payload: LoginPayload): Promise<AuthResult> {
    const { message, data } = await apiClient.post<string>(ENDPOINTS.login, payload, {
      auth: false,
    })

    const token = data
    if (typeof token !== 'string' || token.length === 0) {
      // 2xx with no usable token means the contract broke upstream — fail
      // loudly rather than persisting an empty session.
      throw new ApiError('Login succeeded but no token was returned', 502)
    }

    return { message: message || 'Authenticated', token }
  },

  /**
   * POST /user -> { message, data?: null }
   * No session is established; the user logs in afterwards.
   */
  async signup(payload: SignupPayload): Promise<SignupResult> {
    const { message } = await apiClient.post<null>(ENDPOINTS.signup, payload, { auth: false })
    return { message: message || 'Account created' }
  },
}
