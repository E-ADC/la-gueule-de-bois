import type { ApiErrorBody } from './types'

/**
 * Erreur applicative typée, construite à partir du format uniforme
 * `{ "error": "...", "code": "..." }` renvoyé par l'API (voir spec).
 */
export class ApiError extends Error {
  readonly code: string
  readonly status: number

  constructor(message: string, code: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
  }
}

const API_BASE = '/api'

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  /** Corps JSON. Ignoré si `body` (multipart) est fourni. */
  json?: unknown
  /** Corps déjà prêt (ex. FormData pour l'upload de photos). */
  body?: FormData
  signal?: AbortSignal
}

/**
 * Client HTTP minimal basé sur `fetch` natif (pas d'axios, cf. spec).
 * - même origine, cookies de session envoyés via `credentials: "include"`
 * - lève une `ApiError` typée en cas de réponse non-2xx
 * - renvoie `undefined` pour les réponses 204 (No Content)
 */
async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', json, body, signal } = options

  const headers: HeadersInit = {}
  let requestBody: BodyInit | undefined

  if (body) {
    requestBody = body // FormData : ne pas fixer Content-Type, le navigateur gère le boundary
  } else if (json !== undefined) {
    headers['Content-Type'] = 'application/json'
    requestBody = JSON.stringify(json)
  }

  let response: Response
  try {
    response = await fetch(`${API_BASE}${path}`, {
      method,
      headers,
      body: requestBody,
      credentials: 'include',
      signal,
    })
  } catch {
    throw new ApiError('Impossible de contacter le serveur.', 'network_error', 0)
  }

  if (response.status === 204) {
    return undefined as T
  }

  const contentType = response.headers.get('content-type') ?? ''
  const isJson = contentType.includes('application/json')
  const payload = isJson ? await response.json().catch(() => null) : null

  if (!response.ok) {
    const errorBody = payload as ApiErrorBody | null
    throw new ApiError(
      errorBody?.error ?? `Erreur inattendue (${response.status}).`,
      errorBody?.code ?? 'unknown_error',
      response.status,
    )
  }

  return payload as T
}

export const api = {
  get: <T>(path: string, signal?: AbortSignal) => request<T>(path, { method: 'GET', signal }),
  post: <T>(path: string, json?: unknown, signal?: AbortSignal) =>
    request<T>(path, { method: 'POST', json, signal }),
  postForm: <T>(path: string, body: FormData, signal?: AbortSignal) =>
    request<T>(path, { method: 'POST', body, signal }),
  put: <T>(path: string, json?: unknown, signal?: AbortSignal) =>
    request<T>(path, { method: 'PUT', json, signal }),
  putForm: <T>(path: string, body: FormData, signal?: AbortSignal) =>
    request<T>(path, { method: 'PUT', body, signal }),
  del: <T>(path: string, signal?: AbortSignal) => request<T>(path, { method: 'DELETE', signal }),
}
