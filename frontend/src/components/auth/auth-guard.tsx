'use client'
import { FRONTEND_ROUTES } from '@/constants/frontend-routes'
import { SessionState, useAuthStore } from '@/features/auth/auth-store'
import { apiRefreshAcessToken } from '@/lib/api/api-refresh-token'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { Spinner } from '../ui/spinner'

const AUTH_STATE = {
    AUTHENTICATED: 'authenticated',
    UNAUTHENTICATED: 'unauthenticated',
    LOADING: 'loading',
}

export function AuthGuard({ children }: { children: React.ReactNode }) {
    const router = useRouter()
    const [authState, setAuthState] = useState(AUTH_STATE.LOADING)
    const { user, accessToken, sessionState } = useAuthStore()

    useEffect(() => {
        async function bootstrap() {
            // If the user logged out, don't try to refresh
            if (sessionState === SessionState.ENDED) {
                setAuthState(AUTH_STATE.UNAUTHENTICATED)
                return
            }

            // Already authenticated - return
            if (user && accessToken && sessionState === SessionState.ACTIVE) {
                setAuthState(AUTH_STATE.AUTHENTICATED)
                return
            }

            // Try refreshing the token e.g. on page refresh
            const res = await apiRefreshAcessToken()
            if (!res.ok) {
                setAuthState(AUTH_STATE.UNAUTHENTICATED)
                return
            }
            setAuthState(AUTH_STATE.AUTHENTICATED)
        }
        bootstrap()
    }, [user, accessToken, sessionState])

    useEffect(() => {
        if (authState === AUTH_STATE.UNAUTHENTICATED) {
            router.replace(FRONTEND_ROUTES.LOGIN)
        }
    }, [authState, router])

    if (authState === AUTH_STATE.LOADING) {
        return (
            <div className="flex h-100 w-full items-center justify-center">
                <Spinner className="size-6" />
            </div>
        )
    }

    if (authState === AUTH_STATE.UNAUTHENTICATED) {
        return null
    }

    return children
}
