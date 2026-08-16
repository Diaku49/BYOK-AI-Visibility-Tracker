import { AuthLayout } from '@/shared/components'
import { LoginForm } from '@/features/auth'

export default function LoginPage() {
  return (
    <AuthLayout>
      <LoginForm />
    </AuthLayout>
  )
}
