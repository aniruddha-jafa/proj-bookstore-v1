'use client'
import { useAuthStore } from '@/features/auth/auth-store'

export default function ProfilePage() {
    const { user, isLoading } = useAuthStore()

    if (isLoading) {
        return <div>Loading...</div>
    }
    if (!user) {
        return <div>User not found</div>
    }
    return <div>User: {user.email}</div>
}
