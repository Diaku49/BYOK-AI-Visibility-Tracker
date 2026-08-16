import { AuthLayout } from '@/shared/components'
import { SignupForm } from '@/features/auth'

export default function SignupPage() {
  return (
    <AuthLayout>
      <SignupForm />
    </AuthLayout>
  )
}
