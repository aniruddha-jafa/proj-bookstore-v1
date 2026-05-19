import type { Metadata } from 'next'
import './globals.css'
import { Geist } from 'next/font/google'
import { cn } from '@/lib/utils'
import { ToasterDismissable } from '@/components/ui/toast-dismissable'

const geist = Geist({ subsets: ['latin'], variable: '--font-sans' })

export const metadata: Metadata = {
    title: 'Bookstore',
    description: 'Bookstore',
}

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode
}>) {
    return (
        <html lang="en" className={cn('font-sans', geist.variable)}>
            <body>
                {children}
                <ToasterDismissable />
            </body>
        </html>
    )
}
