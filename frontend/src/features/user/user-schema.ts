import z from 'zod'

export const userResponseSchema = z.object({
    id: z.uuid(),
    email: z.email(),
    createdAt: z.iso.datetime({ offset: true }),
    updatedAt: z.iso.datetime({ offset: true }),
})

export type UserResponse = z.infer<typeof userResponseSchema>
