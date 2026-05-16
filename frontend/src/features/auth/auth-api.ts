import { apiRequest } from "@/lib/api";
import {
  LoginRequest,
  LoginResponse,
  loginResponseSchema,
} from "./auth-schema";
import API_ROUTES from "@/constants/api-routes";
import type { ApiResult } from "@/lib/api";

export const login = async (
  data: LoginRequest,
): Promise<ApiResult<LoginResponse>> => {
  const res = await apiRequest<LoginResponse>(
    API_ROUTES.AUTH.LOGIN,
    {
      method: "POST",
      body: JSON.stringify(data),
      credentials: "include",
    },
    loginResponseSchema,
  );
  return res;
};
