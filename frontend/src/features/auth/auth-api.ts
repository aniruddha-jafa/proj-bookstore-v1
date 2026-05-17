import { apiRequest } from "@/lib/api";
import {
  LoginRequest,
  LoginResponse,
  loginResponseSchema,
  RefreshTokenResponse,
  refreshTokenResponseSchema,
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

export const logout = async (): Promise<ApiResult<void>> => {
  const res = await apiRequest<void>(
    API_ROUTES.AUTH.LOGOUT,
    {
      method: "POST",
      credentials: "include",
    },
  );
  return res;
};

export const refreshToken = async (): Promise<ApiResult<RefreshTokenResponse>> => {
  const res = await apiRequest<RefreshTokenResponse>(
    API_ROUTES.AUTH.REFRESH_TOKEN,
    {
      method: "POST",
      credentials: "include",
    },
    refreshTokenResponseSchema,
  );
  return res;
};