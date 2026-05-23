import { API_PATHS } from '@/constants/api-routes'
import { apiRequest, ApiResult } from '@/lib/api/api'
import { UserResponse, userResponseSchema } from './user-schema'

export const getUser = async (id: string): Promise<ApiResult<UserResponse>> => {
    const res = await apiRequest<UserResponse>(
        API_PATHS.USER.GET(id),
        {
            method: 'GET',
        },
        userResponseSchema
    )
    return res
}
