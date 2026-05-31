import { UserResponse } from '@/features/user/user-schema'
import { create } from 'zustand'

export const SessionState = {
    UNKNOWN: 'unknown' as const,
    ACTIVE: 'active' as const,
    ENDED: 'ended' as const,
}

interface AuthStore {
    user: UserResponse | null
    accessToken: string | null
    csrfToken: string | null
    isLoading: boolean
    sessionState: (typeof SessionState)[keyof typeof SessionState]
    setUser: (user: UserResponse) => void
    setAccessToken: (accessToken: string) => void
    setCsrfToken: (csrfToken: string) => void
    setSessionState: (
        sessionState: (typeof SessionState)[keyof typeof SessionState]
    ) => void
    setOnLoginSuccess: (
        user: UserResponse,
        accessToken: string,
        csrfToken: string
    ) => void
    logout: () => void
}

export const useAuthStore = create<AuthStore>((set) => ({
    user: null,
    accessToken: null,
    isLoading: true,
    csrfToken: null,
    sessionState: SessionState.UNKNOWN,
    setUser: (user: UserResponse) => set({ user, isLoading: false }),
    setAccessToken: (accessToken: string) => set({ accessToken }),
    setCsrfToken: (csrfToken: string) => set({ csrfToken }),
    setSessionState: (
        sessionState: (typeof SessionState)[keyof typeof SessionState]
    ) => set({ sessionState }),
    setOnLoginSuccess: (
        user: UserResponse,
        accessToken: string,
        csrfToken: string
    ) =>
        set({
            user,
            accessToken,
            csrfToken,
            sessionState: SessionState.ACTIVE,
            isLoading: false,
        }),
    logout: () =>
        set({
            user: null,
            accessToken: null,
            csrfToken: null,
            isLoading: false,
            sessionState: SessionState.ENDED,
        }),
}))
