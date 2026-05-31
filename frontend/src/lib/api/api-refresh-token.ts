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
import { useAuthStore } from '@/features/auth/auth-store'
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

    let csrfToken = useAuthStore.getState().csrfToken
    if (!csrfToken) {
        const csrfTokenRes = await apiRequest<CSRFTokenResponse>(
            API_PATHS.AUTH.GET_CSRF_TOKEN,
            {
                method: HTTP_METHOD.GET,
            },
            csrfTokenResponseSchema
        )
        if (!csrfTokenRes.ok || !csrfTokenRes.data?.csrfToken) {
            useAuthStore.getState().logout()
            return {
                ok: false,
                status: 401,
                message: MESSAGES.ERROR_UNAUTHORIZED,
            }
        }
        csrfToken = csrfTokenRes.data.csrfToken
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
        useAuthStore.getState().logout()
        return {
            ok: false,
            status: 401,
            message: MESSAGES.ERROR_UNAUTHORIZED,
        }
    }
    const { token, userId } = refreshTokenRes.data

    // If user is not set, get the user from the API
    let user = useAuthStore.getState().user
    if (!user) {
        const userRes = await getUser(userId)

        if (!userRes.ok || !userRes.data) {
            useAuthStore.getState().logout()
            return {
                ok: false,
                status: userRes.status,
                message: MESSAGES.ERROR_UNAUTHORIZED,
            }
        }

        user = userRes.data
    }

    if (user && user.id !== userId) {
        useAuthStore.getState().logout()
        return { ok: false, status: 401, message: MESSAGES.ERROR_UNAUTHORIZED }
    }

    // Success - set logged in state
    useAuthStore.getState().setOnLoginSuccess(user, token, csrfToken)

    return refreshTokenRes
}
