'use client'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
    Field,
    FieldError,
    FieldGroup,
    FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { FRONTEND_ROUTES } from '@/constants/frontend-routes'
import MESSAGES from '@/constants/messages'
import { signup } from '@/features/auth/auth-api'
import {
    signupRequestSchema,
    type SignupRequest,
} from '@/features/auth/auth-schema'
import { zodResolver } from '@hookform/resolvers/zod'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'

export default function SignupPage() {
    const router = useRouter()

    const {
        register,
        handleSubmit,
        formState: { errors, isSubmitting },
    } = useForm<SignupRequest>({
        resolver: zodResolver(signupRequestSchema),
        defaultValues: {
            email: '',
            password: '',
        },
    })

    async function onSubmit(values: SignupRequest) {
        try {
            const res = await signup(values)
            if (!res.ok) {
                toast.error(MESSAGES.SIGNUP_ERROR)
                return
            }
            router.push(FRONTEND_ROUTES.LOGIN)
            toast.success(MESSAGES.SIGNUP_SUCCESS)
        } catch {
            toast.error(MESSAGES.SIGNUP_ERROR)
        }
    }

    return (
        <main className="flex min-h-svh items-center justify-center p-4">
            <form
                className="w-full max-w-md"
                noValidate
                onSubmit={handleSubmit(onSubmit)}
            >
                <Card className="w-full">
                    <CardHeader>
                        <CardTitle>Sign up</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <FieldGroup>
                            <Field>
                                <FieldLabel htmlFor="signup-email">
                                    Email
                                </FieldLabel>
                                <Input
                                    id="signup-email"
                                    type="email"
                                    autoComplete="email"
                                    placeholder="your@email.com"
                                    aria-invalid={
                                        errors.email ? 'true' : 'false'
                                    }
                                    {...register('email')}
                                />
                                <FieldError>{errors.email?.message}</FieldError>
                            </Field>
                            <Field>
                                <FieldLabel htmlFor="signup-password">
                                    Password
                                </FieldLabel>
                                <Input
                                    id="signup-password"
                                    type="password"
                                    autoComplete="new-password"
                                    placeholder="********"
                                    aria-invalid={
                                        errors.password ? 'true' : 'false'
                                    }
                                    {...register('password')}
                                />
                                {errors.password && (
                                    <FieldError>
                                        {errors.password.message}
                                    </FieldError>
                                )}
                            </Field>
                        </FieldGroup>

                        <Button
                            id="signup-submit"
                            type="submit"
                            className="w-full mt-4"
                            disabled={isSubmitting}
                        >
                            Sign up
                        </Button>
                    </CardContent>
                </Card>
            </form>
        </main>
    )
}
