/** Browser code must use NEXT_PUBLIC_*; server may use API_URL only. */
const NEXT_PUBLIC_API_URL = `${process.env.NEXT_PUBLIC_API_URL}/api/v1`

const PATHS = {
    AUTH: '/auth',
    USER: '/user',
}

export const API_URL = NEXT_PUBLIC_API_URL

export const API_PATHS = {
    AUTH: {
        LOGIN: `${API_URL}/${PATHS.AUTH}/login`,
        SIGNUP: `${API_URL}/${PATHS.AUTH}/signup`,
        REFRESH_TOKEN: `${API_URL}/${PATHS.AUTH}/refresh-token`,
        LOGOUT: `${API_URL}/${PATHS.AUTH}/logout`,
    },
    USER: {
        GET: (id: string) => `${API_URL}/${PATHS.USER}/${id}`,
    },
}

export default API_PATHS
