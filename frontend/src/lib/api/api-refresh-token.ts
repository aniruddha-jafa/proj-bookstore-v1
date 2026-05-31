import API_PATHS from '@/constants/api-routes'
import { HTTP_METHOD } from '@/constants/http'
import { HTTP_HEADER_NAMES } from '@/constants/http-header-names'
import MESSAGES from '@/constants/messages'
import {
    CSRFTokenResponse,
    csrfTokenResponseSchema,
    RefreshTokenResponse,
    refreshTokenResponseSchema,
} from '@/features/auth/auth-schema'
import { logout, useAuthStore } from '@/features/auth/auth-store'
import { getUser } from '@/features/user/user-api'
import { apiRequest, ApiResult } from './api'

/**
 * Refreshes the access token and sets it in the auth store.
 * Updates the user in the auth store if it is not set.
 *
 * Kept in /lib/api instead of auth-api.ts to avoid circular dependency.
 * @returns The result of the request.
 */
export const apiRefreshAcessToken = async (): Promise<
    ApiResult<RefreshTokenResponse>
> => {
    // If CSRF token is not set, get it from the API
    if (!useAuthStore.getState().csrfToken) {
        const csrfTokenRes = await apiRequest<CSRFTokenResponse>(
            API_PATHS.AUTH.GET_CSRF_TOKEN,
            {
                method: HTTP_METHOD.GET,
            },
            csrfTokenResponseSchema
        )
        if (!csrfTokenRes.ok || !csrfTokenRes.data?.csrfToken) {
            logout()
            return {
                ok: false,
                status: 401,
                message: MESSAGES.ERROR_UNAUTHORIZED,
            }
        }
        useAuthStore.getState().setCsrfToken(csrfTokenRes.data.csrfToken)
    }

    const csrfToken = useAuthStore.getState().csrfToken
    if (!csrfToken) {
        logout()
        return {
            ok: false,
            status: 401,
            message: MESSAGES.ERROR_UNAUTHORIZED,
        }
    }

    const refreshTokenRes = await apiRequest<RefreshTokenResponse>(
        API_PATHS.AUTH.REFRESH_TOKEN,
        {
            method: HTTP_METHOD.POST,
            headers: {
                [HTTP_HEADER_NAMES.CSRF_TOKEN]: csrfToken,
            },
        },
        refreshTokenResponseSchema
    )
    if (
        !refreshTokenRes.ok ||
        !refreshTokenRes.data?.token ||
        !refreshTokenRes.data?.userId
    ) {
        logout()
        return {
            ok: false,
            status: 401,
            message: MESSAGES.ERROR_UNAUTHORIZED,
        }
    }
    const { token, userId } = refreshTokenRes.data
    useAuthStore.getState().setAccessToken(token)

    // If user is not set, get the user from the API
    if (!useAuthStore.getState().user) {
        const userRes = await getUser(userId)
        if (!userRes.ok) {
            logout()
            return {
                ok: false,
                status: userRes.status,
                message: MESSAGES.ERROR_API_REQUEST_FAILED,
            }
        }
        useAuthStore.getState().setUser(userRes.data)
    }

    // User ID should match with store
    if (userId !== useAuthStore.getState().user?.id) {
        logout()
        return {
            ok: false,
            status: 401,
            message: MESSAGES.ERROR_UNAUTHORIZED,
        }
    }

    return refreshTokenRes
}
