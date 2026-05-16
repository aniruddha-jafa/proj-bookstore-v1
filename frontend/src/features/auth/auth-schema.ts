import { z } from "zod";

export const loginRequestSchema = z.object({
  email: z.email(),
  password: z.string().min(1, { message: "Password is required" }),
});

export type LoginRequest = z.infer<typeof loginRequestSchema>;

export const loginResponseSchema = z.object({
  id: z.string(),
  email: z.string(),
  /** ISO 8601 format with timezone offset */
  createdAt: z.iso.datetime({ offset: true }),
  updatedAt: z.iso.datetime({ offset: true }),
  token: z.string(),
  refreshToken: z.string(),
});

export type LoginResponse = z.infer<typeof loginResponseSchema>;