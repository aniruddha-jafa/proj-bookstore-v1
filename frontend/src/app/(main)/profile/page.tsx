'use client'
import { useAuthStore } from '@/features/auth/auth-store'

export default function ProfilePage() {
    const { user, isLoading } = useAuthStore()
    const isAuthenticated = useAuthStore((s) =>
        Boolean(s.accessToken && s.user)
    )

    if (!isAuthenticated) {
        return <div>Please log in to view your profile.</div>
    }
    if (isLoading) {
        return <div>Loading...</div>
    }
    if (!user) {
        return <div>User not found</div>
    }
    return (
        <div className="">
            <h1 className="text-2xl font-bold">Profile</h1>
            <div>User: {user.email}</div>
        </div>
    )
}
