import { API_PATHS } from '@/constants/api-routes'
import { HTTP_METHOD } from '@/constants/http'
import { apiRequest, ApiResult } from '@/lib/api/api'
import { UserResponse, userResponseSchema } from './user-schema'

export const getUser = async (id: string): Promise<ApiResult<UserResponse>> => {
    const res = await apiRequest<UserResponse>(
        API_PATHS.USER.GET(id),
        {
            method: HTTP_METHOD.GET,
        },
        userResponseSchema
    )
    return res
}

export const deleteUser = async (id: string): Promise<ApiResult<void>> => {
    const res = await apiRequest<void>(API_PATHS.USER.DELETE(id), {
        method: HTTP_METHOD.DELETE,
    })
    return res
}
