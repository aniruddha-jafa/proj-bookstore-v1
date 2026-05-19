import { UserResponse } from '@/features/user/user-schema'
import { create } from 'zustand'

interface AuthStore {
    user: UserResponse | null
    accessToken: string | null
    isLoading: boolean
    setUser: (user: UserResponse) => void
    setAccessToken: (accessToken: string) => void
    logout: () => void
}

export const useAuthStore = create<AuthStore>((set) => ({
    user: null,
    accessToken: null,
    isLoading: true,
    setUser: (user: UserResponse) => set({ user, isLoading: false }),
    setAccessToken: (accessToken: string) => set({ accessToken }),
    logout: () => set({ user: null, accessToken: null, isLoading: false }),
}))
