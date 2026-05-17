import API_ROUTES from "@/constants/api-routes";
import { UserResponse, userResponseSchema } from "./user-schema";
import { apiRequest, ApiResult } from "@/lib/api";


export const getUser = async (id: string, accessToken: string): Promise<ApiResult<UserResponse>> => {
    const res = await apiRequest<UserResponse>(
        API_ROUTES.USER.GET(id),
        {
            method: "GET",
            headers: {
                Authorization: `Bearer ${accessToken}`,
            },
        },
        userResponseSchema,
    );
    return res;
}