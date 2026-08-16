import { Link } from 'react-router-dom'
import { Button, Input } from '@/shared/components'
import { ROUTES } from '@/shared/config/routes'
import AuthShell from '@/features/auth/components/AuthShell'
import FormError from '@/features/auth/components/FormError'
import { useSignupForm } from '@/features/auth/hooks/useSignupForm'
import { PASSWORD_MIN_LENGTH } from '@/features/auth/validation'

export default function SignupForm() {
  const { values, errors, isSubmitting, setField, submit } = useSignupForm()

  return (
    <AuthShell
      command="auth/signup"
      title="Create account"
      subtitle="Bring your own keys. No card required."
      footer={
        <p className="font-mono text-[11px] tracking-wide text-ink-faint">
          Already registered?{' '}
          <Link to={ROUTES.login} className="text-accent transition-opacity hover:opacity-70">
            Login
          </Link>
        </p>
      }
    >
      <form onSubmit={submit} noValidate className="flex flex-col gap-1">
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
          autoComplete="new-password"
          placeholder="••••••••"
          hint={`min ${PASSWORD_MIN_LENGTH}`}
          value={values.password}
          onChange={event => setField('password', event.target.value)}
          error={errors.password}
          disabled={isSubmitting}
        />

        <Input
          label="Name"
          type="text"
          autoComplete="name"
          placeholder="Operator name"
          hint="optional"
          value={values.name}
          onChange={event => setField('name', event.target.value)}
          error={errors.name}
          disabled={isSubmitting}
        />

        <FormError message={errors.form} />

        <Button type="submit" variant="solid" size="lg" className="mt-4 w-full" disabled={isSubmitting}>
          {isSubmitting ? 'Creating…' : 'Create account'}
        </Button>
      </form>
    </AuthShell>
  )
}
