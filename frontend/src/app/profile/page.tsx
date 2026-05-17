"use client";
import { useUserStore } from "@/features/user/user-auth-store";

export default function ProfilePage() {
    const { user, isLoading } = useUserStore();

    if (isLoading) {
        return <div>Loading...</div>;
    }
    if (!user) {
        return <div>User not found</div>;
    }
    return <div>User: {user.email}</div>;
}