import { readAuthToken } from '@/shared/config/storage'
import type { ApiEnvelope } from '@/shared/types'

/** Local Go API default. Deployments override this in their environment. */
const DEFAULT_BASE_URL = 'http://localhost:8080'

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? DEFAULT_BASE_URL

export class ApiError extends Error {
  constructor(
    /** Server-supplied `message` when present, otherwise a generic fallback. */
    message: string,
    readonly status: number,
    readonly isNetworkError = false,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

interface RequestOptions extends Omit<RequestInit, 'body'> {
  body?: unknown
  /** Set false for endpoints that must not carry the JWT. */
  auth?: boolean
}

function isEnvelope(value: unknown): value is ApiEnvelope<unknown> {
  return typeof value === 'object' && value !== null && 'message' in value
}

/**
 * Base fetch wrapper. Every endpoint answers with `{ message, data }` and a
 * meaningful status code. This returns the whole envelope — `message` matters
 * on success too (signup returns nothing else) — and raises ApiError carrying
 * the server's `message` whenever the status is not ok.
 *
 * Feature `api/` folders build on this rather than calling fetch directly.
 */
async function request<T>(path: string, options: RequestOptions = {}): Promise<ApiEnvelope<T>> {
  const { body, headers, auth = true, ...rest } = options

  const token = auth ? readAuthToken() : null

  let response: Response
  try {
    response = await fetch(`${BASE_URL}${path}`, {
      ...rest,
      credentials: 'include',
      headers: {
        Accept: 'application/json',
        ...(body !== undefined && { 'Content-Type': 'application/json' }),
        ...(token && { Authorization: `Bearer ${token}` }),
        ...headers,
      },
      ...(body !== undefined && { body: JSON.stringify(body) }),
    })
  } catch {
    // fetch only rejects on transport failure — no status exists here.
    throw new ApiError('Connection failed. Check your network and retry.', 0, true)
  }

  // 204 carries no body at all, so there is no envelope to parse.
  if (response.status === 204) {
    if (!response.ok) throw new ApiError('Request failed', response.status)
    return { message: '', data: null }
  }

  let payload: unknown = null
  try {
    payload = await response.json()
  } catch {
    // Non-JSON body (proxy error page, gateway timeout). Fall through to the
    // status check so the caller still gets the right status.
    payload = null
  }

  const envelope = isEnvelope(payload) ? payload : null

  if (!response.ok) {
    throw new ApiError(envelope?.message ?? `Request failed: ${path}`, response.status)
  }

  return {
    message: envelope?.message ?? '',
    data: (envelope?.data ?? null) as T | null,
  }
}

export const apiClient = {
  get: <T>(path: string, options?: RequestOptions) =>
    request<T>(path, { ...options, method: 'GET' }),
  post: <T>(path: string, body?: unknown, options?: RequestOptions) =>
    request<T>(path, { ...options, method: 'POST', body }),
  patch: <T>(path: string, body?: unknown, options?: RequestOptions) =>
    request<T>(path, { ...options, method: 'PATCH', body }),
  delete: <T>(path: string, options?: RequestOptions) =>
    request<T>(path, { ...options, method: 'DELETE' }),
}
