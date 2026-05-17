/** Browser code must use NEXT_PUBLIC_*; server may use API_URL only. */
const NEXT_PUBLIC_API_URL = `${process.env.NEXT_PUBLIC_API_URL}/api/v1`;

export const PAGE_PREFIXES = {
  AUTH: "/auth",
  USER: "/user",
};

export const API_ROUTES = {
  AUTH: {
    LOGIN: `${NEXT_PUBLIC_API_URL}${PAGE_PREFIXES.AUTH}/login`,
    SIGNUP: `${NEXT_PUBLIC_API_URL}${PAGE_PREFIXES.AUTH}/signup`,
    REFRESH_TOKEN: `${NEXT_PUBLIC_API_URL}${PAGE_PREFIXES.AUTH}/refresh-token`,
    LOGOUT: `${NEXT_PUBLIC_API_URL}${PAGE_PREFIXES.AUTH}/logout`,
  },
  USER: {
    GET: (id: string) => `${NEXT_PUBLIC_API_URL}${PAGE_PREFIXES.USER}/${id}`,
  },
};

export default API_ROUTES;
