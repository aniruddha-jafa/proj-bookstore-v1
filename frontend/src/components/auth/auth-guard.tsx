'use client'
import { FRONTEND_ROUTES } from '@/constants/frontend-routes'
import { useAuthStore } from '@/features/auth/auth-store'
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
    const { user, accessToken } = useAuthStore()

    useEffect(() => {
        async function bootstrap() {
            // Already authenticated - return
            if (user && accessToken) {
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
    }, [user, accessToken])

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
