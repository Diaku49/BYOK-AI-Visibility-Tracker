import { Link } from 'react-router-dom'
import { Button, Input } from '@/shared/components'
import { ROUTES } from '@/shared/config/routes'
import AuthShell from '@/features/auth/components/AuthShell'
import FormError from '@/features/auth/components/FormError'
import FormNotice from '@/features/auth/components/FormNotice'
import { useLoginForm } from '@/features/auth/hooks/useLoginForm'

export default function LoginForm() {
  const { values, errors, notice, isSubmitting, setField, submit } = useLoginForm()

  return (
    <AuthShell
      command="auth/login"
      title="Login"
      subtitle="Authenticate to reach your dashboard."
      footer={
        <p className="font-mono text-[11px] tracking-wide text-ink-faint">
          No account?{' '}
          <Link to={ROUTES.signup} className="text-accent transition-opacity hover:opacity-70">
            Create one
          </Link>
        </p>
      }
    >
      <form onSubmit={submit} noValidate className="flex flex-col gap-1">
        {notice && (
          <div className="mb-3">
            <FormNotice message={notice} />
          </div>
        )}

        <Input
          label="Email"
          type="email"
          autoComplete="email"
          placeholder="operator@domain.com"
          value={values.email}
          onChange={event => setField('email', event.target.value)}
          error={errors.email}
          disabled={isSubmitting}
          autoFocus
        />

        <Input
          label="Password"
          type="password"
          autoComplete="current-password"
          placeholder="••••••••"
          value={values.password}
          onChange={event => setField('password', event.target.value)}
          error={errors.password}
          disabled={isSubmitting}
        />

        <FormError message={errors.form} />

        <Button type="submit" variant="solid" size="lg" className="mt-4 w-full" disabled={isSubmitting}>
          {isSubmitting ? 'Authenticating…' : 'Authenticate'}
        </Button>
      </form>
    </AuthShell>
  )
}
