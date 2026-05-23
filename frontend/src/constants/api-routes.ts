/** Browser code must use NEXT_PUBLIC_*; server may use API_URL only. */
const NEXT_PUBLIC_API_URL = `${process.env.NEXT_PUBLIC_API_URL}/api/v1`

const PATHS = {
    // Auth
    AUTH: '/auth',
    LOGIN: '/login',
    SIGNUP: '/signup',
    REFRESH_TOKEN: '/refresh-token',
    LOGOUT: '/logout',
    // User
    USER: '/user',
}

export const API_URL = NEXT_PUBLIC_API_URL

export const API_PATHS = {
    AUTH: {
        LOGIN: `${API_URL}${PATHS.AUTH}${PATHS.LOGIN}`,
        SIGNUP: `${API_URL}${PATHS.AUTH}${PATHS.SIGNUP}`,
        REFRESH_TOKEN: `${API_URL}${PATHS.AUTH}${PATHS.REFRESH_TOKEN}`,
        LOGOUT: `${API_URL}${PATHS.AUTH}${PATHS.LOGOUT}`,
    },
    USER: {
        GET: (id: string) => `${API_URL}${PATHS.USER}/${id}`,
    },
}

export default API_PATHS
