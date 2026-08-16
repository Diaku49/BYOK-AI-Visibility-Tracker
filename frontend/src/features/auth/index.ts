/** PUBLIC API — other code imports only from '@/features/auth'. */
export { default as AuthMenu } from '@/features/auth/components/AuthMenu'
export { default as LoginForm } from '@/features/auth/components/LoginForm'
export { default as SignupForm } from '@/features/auth/components/SignupForm'
export { useAuthSession } from '@/features/auth/hooks/useAuthSession'
export type {
  AuthSession,
  LoginPayload,
  SignupPayload,
  AuthResult,
  SignupResult,
} from '@/features/auth/types'
