import API_PATHS from '@/constants/api-routes'
import { HTTP_METHOD } from '@/constants/http'
import MESSAGES from '@/constants/messages'
import {
    RefreshTokenResponse,
    refreshTokenResponseSchema,
} from '@/features/auth/auth-schema'
import { logout, useAuthStore } from '@/features/auth/auth-store'
import { getUser } from '@/features/user/user-api'
import { ApiResult, apiRequest } from './api'

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
    const refreshTokenRes = await apiRequest<RefreshTokenResponse>(
        API_PATHS.AUTH.REFRESH_TOKEN,
        {
            method: HTTP_METHOD.POST,
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
