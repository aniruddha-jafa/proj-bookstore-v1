"use client";

import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { loginFormSchema } from "@/schemas/auth";
import type { z } from "zod";
import API_ROUTES from "@/constants/api-routes";
import MESSAGES from "@/constants/messages";

type LoginFormValues = z.infer<typeof loginFormSchema>;

const LOGIN_SUCCESS_URL = "/";

export default function LoginPage() {
  const router = useRouter();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: LoginFormValues) {
    try {
      const res = await fetch(API_ROUTES.AUTH.LOGIN, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: values.email,
          password: values.password,
        }),
      });

      if (!res.ok) {
        const body = await res.json();
        let loginFailedMessage = `Login failed (${res.status})`;
        if (body.message) {
          loginFailedMessage = body.message;
        }
        console.error("Login failed:", loginFailedMessage);
        toast.error(loginFailedMessage);
        return;
      }

      router.push(LOGIN_SUCCESS_URL);
      router.refresh();
    } catch {
      toast.error(MESSAGES.LOGIN_ERROR);
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
            <CardTitle>Login</CardTitle>
          </CardHeader>
          <CardContent>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="login-email">Email</FieldLabel>
                <Input
                  id="login-email"
                  type="email"
                  autoComplete="email"
                  placeholder="your@email.com"
                  aria-invalid={errors.email ? "true" : "false"}
                  {...register("email")}
                />
                <FieldError>{errors.email?.message}</FieldError>
              </Field>
              <Field>
                <FieldLabel htmlFor="login-password">Password</FieldLabel>
                <Input
                  id="login-password"
                  type="password"
                  autoComplete="current-password"
                  placeholder="********"
                  aria-invalid={errors.password ? "true" : "false"}
                  {...register("password")}
                />
                {errors.password && (
                  <FieldError>{errors.password.message}</FieldError>
                )}
              </Field>
            </FieldGroup>

            <Button
              id="login-submit"
              type="submit"
              className="w-full mt-4"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Signing in…" : "Login"}
            </Button>
          </CardContent>
        </Card>
      </form>
    </main>
  );
}
