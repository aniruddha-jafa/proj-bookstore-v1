'use client'
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogContent,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { FRONTEND_ROUTES } from '@/constants/frontend-routes'
import MESSAGES from '@/constants/messages'
import { useAuthStore } from '@/features/auth/auth-store'
import { deleteUser } from '@/features/user/user-api'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'

export default function ProfilePage() {
    const { user, accessToken, logout } = useAuthStore()
    const [isDeleting, setIsDeleting] = useState(false)
    const router = useRouter()
    if (!accessToken || !user) {
        return null
    }

    const handleDeleteUser = async () => {
        setIsDeleting(true)
        const res = await deleteUser(user.id)
        if (!res.ok) {
            toast.error(MESSAGES.ACCOUNT_DELETED_ERROR)
            return
        }
        logout()
        router.replace(FRONTEND_ROUTES.LOGIN)
        toast.success(MESSAGES.ACCOUNT_DELETED_SUCCESS)
        setIsDeleting(false)
    }

    return (
        <div className="">
            <h1 className="text-2xl font-bold">Profile</h1>
            <div>User: {user.email}</div>
            <AlertDialog>
                <AlertDialogTrigger asChild>
                    <Button
                        variant="destructive"
                        className="mt-2"
                        disabled={isDeleting}
                    >
                        Delete
                    </Button>
                </AlertDialogTrigger>
                <AlertDialogContent aria-describedby="delete-user-dialog-description">
                    <AlertDialogHeader>
                        <AlertDialogTitle>
                            Are you sure you want to delete your account?
                        </AlertDialogTitle>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogAction
                            onClick={handleDeleteUser}
                            variant="destructive"
                        >
                            Yes, delete my account
                        </AlertDialogAction>
                        <AlertDialogAction variant="outline">
                            Cancel
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </div>
    )
}
