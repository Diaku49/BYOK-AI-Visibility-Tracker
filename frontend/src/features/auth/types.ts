export interface AuthSession {
  isLoggedIn: boolean
  token: string | null
}

export interface LoginPayload {
  email: string
  password: string
}

export interface SignupPayload {
  email: string
  password: string
  /** Optional — the backend falls back to the email local-part. */
  name?: string
}

/** `data` for POST /auth/login. Signup returns `data: null`. */
export interface LoginData {
  token: string
}

/** What the API layer hands back: the server message, plus a token on login. */
export interface AuthResult {
  message: string
  token: string
}

export interface SignupResult {
  message: string
}

export type AuthField = 'email' | 'password' | 'name'

/** Per-field messages plus a `form` slot for whole-submission failures. */
export type AuthErrors = Partial<Record<AuthField | 'form', string>>
