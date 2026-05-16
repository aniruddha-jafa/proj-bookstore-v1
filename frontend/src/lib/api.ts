import { z } from "zod";
import MESSAGES from "@/constants/messages";
import { logger } from "@/lib/logger";

export type ApiSuccess<T> = { ok: true; status: number; data: T };
export type ApiFailure = { ok: false; status: number; message: string };
export type ApiResult<T> = ApiSuccess<T> | ApiFailure;

/**
 * Make a request to the API and return the result.
 * @param url - The URL to request.
 * @param init - The init object to pass to the fetch function.
 * @param schema - The schema to parse the response body with.
 * @returns The result of the request.
 */
export async function apiRequest<T>(
  url: string,
  init?: RequestInit,
  schema?: z.ZodType<T>,
): Promise<ApiResult<T>> {
  try {
    const res = await fetch(url, {
      ...init,
      headers: { "Content-Type": "application/json", ...init?.headers },
    });
    // Not all responses have a body
    const body = await res.json().catch(() => null);

    if (!res.ok) {
      return {
        ok: false,
        status: res.status,
        message: body?.message ?? res.statusText ?? "Unknown error",
      };
    }
    // 204 responses don't have a body, so we return null
    if (res.status === 204) {
      return { ok: true, status: res.status, data: null as T };
    }
    // No schema provided, so we return the raw body
    if (!schema) {
      return { ok: true, status: res.status, data: body as unknown as T };
    }
    // Parse the body with the schema
    const parsed = schema.safeParse(body);
    if (!parsed.success) {
      logger.error("API response validation failed", parsed.error.flatten());
      return {
        ok: false,
        status: res.status,
        message: MESSAGES.ERROR_API_PARSING,
      };
    }
    return { ok: true, status: res.status, data: parsed.data };
  } catch {
    throw new Error("network");
  }
}
