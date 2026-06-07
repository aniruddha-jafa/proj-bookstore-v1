import API_PATHS from '@/constants/api-routes'
import { HTTP_METHOD } from '@/constants/http'
import { HTTP_HEADER_NAMES } from '@/constants/http-header-names'
import type { ApiResult } from '@/lib/api/api'
import { apiRequest } from '@/lib/api/api'
import { apiRefreshAcessToken } from '@/lib/api/api-refresh-token'
import {
    LoginRequest,
    LoginResponse,
    loginResponseSchema,
    RefreshTokenResponse,
    SignupRequest,
    SignupResponse,
    signupResponseSchema,
} from './auth-schema'
import { useAuthStore } from './auth-store'

export const login = async (
    data: LoginRequest
): Promise<ApiResult<LoginResponse>> => {
    const res = await apiRequest<LoginResponse>(
        API_PATHS.AUTH.LOGIN,
        {
            method: HTTP_METHOD.POST,
            body: JSON.stringify(data),
        },
        loginResponseSchema
    )
    return res
}

export const signup = async (
    data: SignupRequest
): Promise<ApiResult<SignupResponse>> => {
    const res = await apiRequest<SignupResponse>(
        API_PATHS.AUTH.SIGNUP,
        {
            method: HTTP_METHOD.POST,
            body: JSON.stringify(data),
        },
        signupResponseSchema
    )
    return res
}

export const logout = async (): Promise<ApiResult<void>> => {
    const res = await apiRequest<void>(API_PATHS.AUTH.LOGOUT, {
        method: HTTP_METHOD.POST,
        headers: {
            [HTTP_HEADER_NAMES.CSRF_TOKEN]:
                useAuthStore.getState().csrfToken ?? '',
        },
    })
    return res
}

/**
 * Wraps the refresh token API call – UI should call this instead of directly calling apiRefreshAcessToken.
 */
export const refreshToken = async (): Promise<
    ApiResult<RefreshTokenResponse>
> => {
    return await apiRefreshAcessToken()
}
