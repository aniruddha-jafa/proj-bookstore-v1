'use client'

import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
    NavigationMenu,
    NavigationMenuItem,
    NavigationMenuLink,
    NavigationMenuList,
    navigationMenuTriggerStyle,
} from '@/components/ui/navigation-menu'
import { FRONTEND_ROUTES } from '@/constants/frontend-routes'
import { logout } from '@/features/auth/auth-api'
import { useAuthStore } from '@/features/auth/auth-store'
import { cn } from '@/lib/utils'
import { BookOpen, LogOut, User } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { toast } from 'sonner'

function userInitials(email: string) {
    return email[0]?.toUpperCase() ?? '?'
}

function NavbarBrand() {
    return (
        <NavigationMenu viewport={false} className="w-full">
            <NavigationMenuList>
                <NavigationMenuItem>
                    <NavigationMenuLink asChild>
                        <Link
                            href="/"
                            className={cn(
                                navigationMenuTriggerStyle(),
                                'gap-2 font-semibold'
                            )}
                        >
                            <BookOpen className="size-4" />
                            Bookstore
                        </Link>
                    </NavigationMenuLink>
                </NavigationMenuItem>
            </NavigationMenuList>
        </NavigationMenu>
    )
}

function UserAccountMenu({ email }: { email: string }) {
    const router = useRouter()
    const { logout: logoutStore } = useAuthStore()

    async function handleLogout() {
        const res = await logout()
        if (!res.ok) {
            toast.error('Failed to logout')
            return
        }
        logoutStore()
        router.replace(FRONTEND_ROUTES.LOGIN)
        toast.success('Logged out')
    }

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="gap-2 px-2">
                    <Avatar size="sm">
                        <AvatarFallback>{userInitials(email)}</AvatarFallback>
                    </Avatar>
                    <span className="hidden max-w-40 truncate sm:inline">
                        {email}
                    </span>
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel className="truncate font-normal">
                    {email}
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem asChild>
                    <Link href={FRONTEND_ROUTES.PROFILE}>
                        <User />
                        Profile
                    </Link>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem variant="destructive" onClick={handleLogout}>
                    <LogOut />
                    Log out
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}

export function Navbar() {
    const user = useAuthStore((s) => s.user)

    return (
        <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur">
            <div className="mx-auto flex h-14 max-w-6xl items-center justify-between gap-4 px-4">
                <NavbarBrand />
                <div className="flex shrink-0 items-center">
                    {user ? (
                        <UserAccountMenu email={user.email} />
                    ) : (
                        <Button asChild variant="outline" size="sm">
                            <Link href={FRONTEND_ROUTES.LOGIN}>Sign in</Link>
                        </Button>
                    )}
                </div>
            </div>
        </header>
    )
}
