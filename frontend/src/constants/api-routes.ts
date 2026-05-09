/** Browser code must use NEXT_PUBLIC_*; server may use API_URL only. */
const BACKEND_API_URL = process.env.BACKEND_API_URL;

const PAGE_PREFIXES = {
  AUTH: "/auth",
};

export const API_ROUTES = {
  AUTH: {
    LOGIN: `${BACKEND_API_URL}${PAGE_PREFIXES.AUTH}/login`,
    SIGNUP: `${BACKEND_API_URL}${PAGE_PREFIXES.AUTH}/signup`,
    REFRESH_TOKEN: `${BACKEND_API_URL}${PAGE_PREFIXES.AUTH}/refresh-token`,
    LOGOUT: `${BACKEND_API_URL}${PAGE_PREFIXES.AUTH}/logout`,
  },
};

export default API_ROUTES;
