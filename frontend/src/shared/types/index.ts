/** Cross-feature primitives only. Domain shapes live in features/<name>/types.ts. */

/**
 * The response envelope every endpoint returns. The HTTP status code carries
 * success/failure; the body carries a human message and an optional payload.
 * `data` is null whenever an endpoint has nothing to return.
 */
export interface ApiEnvelope<T> {
  message: string
  data: T | null
}

export interface Paginated<T> {
  items: T[]
  page: number
  perPage: number
  total: number
}

export type AsyncState = 'idle' | 'loading' | 'success' | 'error'
