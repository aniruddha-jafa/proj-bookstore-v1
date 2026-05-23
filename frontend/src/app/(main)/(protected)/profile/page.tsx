'use client'
import { useAuthStore } from '@/features/auth/auth-store'

export default function ProfilePage() {
    const { user, accessToken } = useAuthStore()

    if (!accessToken || !user) {
        return null
    }

    return (
        <div className="">
            <h1 className="text-2xl font-bold">Profile</h1>
            <div>User: {user.email}</div>
        </div>
    )
}
