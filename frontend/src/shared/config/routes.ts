/**
 * Route paths. Lives in shared/ so both app/ (route definitions) and features/
 * (navigation) can reference it without features importing upward from app/.
 */
export const ROUTES = {
  home: '/',
  login: '/login',
  signup: '/signup',
  product: '/product/:productId',
} as const

export function productPath(productId: string): string {
  return `/product/${productId}`
}
