import { UserResponse } from '@/features/user/user-schema'
import { create } from 'zustand'

interface AuthStore {
    user: UserResponse | null
    accessToken: string | null
    csrfToken: string | null
    isLoading: boolean
    setUser: (user: UserResponse) => void
    setAccessToken: (accessToken: string) => void
    setCsrfToken: (csrfToken: string) => void
    logout: () => void
}

export const useAuthStore = create<AuthStore>((set) => ({
    user: null,
    accessToken: null,
    isLoading: true,
    csrfToken: null,
    setUser: (user: UserResponse) => set({ user, isLoading: false }),
    setAccessToken: (accessToken: string) => set({ accessToken }),
    setCsrfToken: (csrfToken: string) => set({ csrfToken }),
    logout: () =>
        set({
            user: null,
            accessToken: null,
            csrfToken: null,
            isLoading: false,
        }),
}))

// Helper function to logout
export const logout = () => {
    useAuthStore.getState().logout()
}
