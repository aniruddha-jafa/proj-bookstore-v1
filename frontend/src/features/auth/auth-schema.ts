import { z } from 'zod'

// Login

export const loginRequestSchema = z.object({
    email: z.email(),
    password: z.string().min(1, { message: 'Password is required' }),
})

export type LoginRequest = z.infer<typeof loginRequestSchema>

export const loginResponseSchema = z.object({
    id: z.string(),
    email: z.email(),
    /** ISO 8601 format with timezone offset */
    createdAt: z.iso.datetime({ offset: true }),
    updatedAt: z.iso.datetime({ offset: true }),
    token: z.string(),
    csrfToken: z.string(),
})

export type LoginResponse = z.infer<typeof loginResponseSchema>

// Signup

export const signupRequestSchema = z.object({
    email: z.email().min(1, { message: 'Email is required' }),
    password: z.string().min(1, { message: 'Password is required' }),
})

export type SignupRequest = z.infer<typeof signupRequestSchema>

export const signupResponseSchema = z.object({
    id: z.string(),
    email: z.email(),
    createdAt: z.iso.datetime({ offset: true }),
    updatedAt: z.iso.datetime({ offset: true }),
})

export type SignupResponse = z.infer<typeof signupResponseSchema>
// Refresh Token

export const refreshTokenResponseSchema = z.object({
    token: z.string(),
    userId: z.uuid(),
})

export type RefreshTokenResponse = z.infer<typeof refreshTokenResponseSchema>

// CSRF Token

export const csrfTokenResponseSchema = z.object({
    csrfToken: z.string(),
})

export type CSRFTokenResponse = z.infer<typeof csrfTokenResponseSchema>
