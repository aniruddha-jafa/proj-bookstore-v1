import API_PATHS from '@/constants/api-routes'
import type { ApiResult } from '@/lib/api'
import { apiRequest } from '@/lib/api'
import { LoginRequest, LoginResponse, loginResponseSchema } from './auth-schema'

export const login = async (
    data: LoginRequest
): Promise<ApiResult<LoginResponse>> => {
    const res = await apiRequest<LoginResponse>(
        API_PATHS.AUTH.LOGIN,
        {
            method: 'POST',
            body: JSON.stringify(data),
        },
        loginResponseSchema
    )
    return res
}

export const logout = async (): Promise<ApiResult<void>> => {
    const res = await apiRequest<void>(API_PATHS.AUTH.LOGOUT, {
        method: 'POST',
    })
    return res
}
