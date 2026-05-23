import { Navbar } from '@/components/layout/navbar'

export default function MainLayout({
    children,
}: Readonly<{
    children: React.ReactNode
}>) {
    return (
        <>
            <Navbar />
            <main className="mx-auto max-w-6xl px-4 py-6">{children}</main>
        </>
    )
}
