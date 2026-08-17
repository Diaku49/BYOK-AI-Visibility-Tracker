import { useNavigate } from 'react-router'
import { Button } from '@/shared/components'
import { ROUTES } from '@/shared/config/routes'
import { useAuthSession } from '@/features/auth/hooks/useAuthSession'

/**
 * Session-dependent header slot: auth entry points when signed out,
 * a logout affordance when signed in.
 */
export default function AuthMenu() {
  const { isLoggedIn, logout } = useAuthSession()
  const navigate = useNavigate()

  if (isLoggedIn) {
    return (
      <Button variant="outline" size="sm" onClick={logout}>
        Log out
      </Button>
    )
  }

  return (
    <div className="flex items-center gap-2">
      <Button variant="ghost" size="sm" onClick={() => navigate(ROUTES.login)}>
        Login
      </Button>
      <Button variant="solid" size="sm" onClick={() => navigate(ROUTES.signup)}>
        Sign up
      </Button>
    </div>
  )
}
