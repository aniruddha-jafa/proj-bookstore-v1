import API_PATHS from '@/constants/api-routes'
import { FRONTEND_ROUTES } from '@/constants/frontend-routes'
import MESSAGES from '@/constants/messages'
import { logout, useAuthStore } from '@/features/auth/auth-store'
import { logger } from '@/lib/logger'
import { z } from 'zod'
import { apiRefreshAcessToken } from './api-refresh-token'

export type ApiSuccess<T> = { ok: true; status: number; data: T }
export type ApiFailure = { ok: false; status: number; message: string }
export type ApiResult<T> = ApiSuccess<T> | ApiFailure

type ApiOptions = RequestInit & {
    _skipRetryOnAuthFailure?: boolean
}

const AUTH_PATHS_NO_REFRESH_RETRY = [
    API_PATHS.AUTH.REFRESH_TOKEN,
    API_PATHS.AUTH.LOGIN,
    API_PATHS.AUTH.SIGNUP,
    API_PATHS.AUTH.LOGOUT,
] as const

/**
 * Make a request to the API and return the result.
 * @param url - The URL to request.
 * @param options - The init object to pass to the fetch function.
 * @param schema - The schema to parse the response body with.
 * @returns The result of the request.
 */
export async function apiRequest<T>(
    url: string,
    options: ApiOptions = {},
    schema?: z.ZodType<T>
): Promise<ApiResult<T>> {
    options.headers = {
        'Content-Type': 'application/json',
        Accept: 'application/json',
        ...(options.headers as Record<string, string>),
    }

    const { accessToken } = useAuthStore.getState()
    if (accessToken) {
        options.headers['Authorization'] = `Bearer ${accessToken}`
    }

    try {
        const res = await fetch(url, {
            // Send cookies with request
            ...options,
            credentials: 'include',
        })

        const isUnauthorized = res.status === 401
        const shouldSkipRetry =
            options._skipRetryOnAuthFailure ||
            AUTH_PATHS_NO_REFRESH_RETRY.some((path) => url.startsWith(path))

        // Unauthorized and already retried, return error
        if (isUnauthorized && shouldSkipRetry) {
            const body = await res.json().catch(() => null)
            const message = body?.message ?? MESSAGES.ERROR_UNAUTHORIZED
            return {
                ok: false,
                status: 401,
                message,
            }
        }

        // Unauthorized and not retried, try to refresh token
        if (isUnauthorized && !shouldSkipRetry) {
            options._skipRetryOnAuthFailure = true

            const refreshTokenRes = await apiRefreshAcessToken()

            if (!refreshTokenRes.ok || !refreshTokenRes.data?.token) {
                logout()
                if (typeof window !== 'undefined') {
                    window.location.href = FRONTEND_ROUTES.LOGIN
                }
                return {
                    ok: false,
                    status: 401,
                    message: MESSAGES.ERROR_UNAUTHORIZED,
                }
            }
            const { token: newAccessToken } = refreshTokenRes.data
            useAuthStore.getState().setAccessToken(newAccessToken)
            options.headers['Authorization'] = `Bearer ${newAccessToken}`

            return await apiRequest<T>(url, {
                ...options,
            })
        }

        // Not all responses have a body
        const body = await res.json().catch(() => null)

        if (!res.ok) {
            return {
                ok: false,
                status: res.status,
                message: body?.message ?? res.statusText ?? 'Unknown error',
            }
        }
        // 204 responses don't have a body
        if (res.status === 204) {
            return { ok: true, status: res.status, data: null as T }
        }

        // No schema provided, so return the raw body
        if (!schema) {
            return { ok: true, status: res.status, data: body as unknown as T }
        }

        // Parse the body with the schema
        const parsed = schema.safeParse(body)
        if (!parsed.success) {
            logger.error(
                'API response validation failed',
                parsed.error.flatten()
            )
            return {
                ok: false,
                status: res.status,
                message: MESSAGES.ERROR_INVALID_API_RESPONSE,
            }
        }
        return { ok: true, status: res.status, data: parsed.data }
    } catch (error) {
        logger.error('API request error: ', error)
        return {
            ok: false,
            status: 500,
            message: MESSAGES.ERROR_API_REQUEST_FAILED,
        }
    }
}
