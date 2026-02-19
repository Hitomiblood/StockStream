import axios, { AxiosError, type AxiosInstance, type InternalAxiosRequestConfig } from 'axios'

export class ApiError extends Error {
  status: number | null

  constructor(message: string, status: number | null = null) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.trim() || 'http://localhost:8080/api/v1'
const TIMEOUT_MS = 10_000
const MAX_RETRIES = 1

function shouldRetry(error: AxiosError): boolean {
  const method = error.config?.method?.toUpperCase()
  const status = error.response?.status
  const networkError = !error.response
  return method === 'GET' && (networkError || (status !== undefined && status >= 500))
}

function normalizeError(error: unknown): ApiError {
  if (!axios.isAxiosError(error)) {
    return new ApiError('Unexpected error while communicating with API')
  }

  const message =
    (error.response?.data as { error?: string } | undefined)?.error ||
    error.message ||
    'Request failed'

  return new ApiError(message, error.response?.status ?? null)
}

const client: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: TIMEOUT_MS,
  headers: {
    'Content-Type': 'application/json'
  }
})

client.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const normalized = { ...config }
  ;(normalized as InternalAxiosRequestConfig & { _retryCount?: number })._retryCount ??= 0
  return normalized
})

client.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const cfg = error.config as (InternalAxiosRequestConfig & { _retryCount?: number }) | undefined

    if (cfg && shouldRetry(error) && (cfg._retryCount ?? 0) < MAX_RETRIES) {
      cfg._retryCount = (cfg._retryCount ?? 0) + 1
      return client.request(cfg)
    }

    return Promise.reject(normalizeError(error))
  }
)

export { client, normalizeError, API_BASE_URL }
