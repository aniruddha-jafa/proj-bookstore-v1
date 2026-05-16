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
import { loginRequestSchema, type LoginRequest } from "@/features/auth/auth-schema";
import MESSAGES from "@/constants/messages";
import { login } from "@/features/auth/auth-api";


const LOGIN_SUCCESS_URL = "/";

export default function LoginPage() {
  const router = useRouter();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginRequest>({
    resolver: zodResolver(loginRequestSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: LoginRequest) {
    try {
      const res = await login(values);
      if (!res.ok) {
        toast.error(res.message);
        return;
      }
      console.log("DEBUG: Login successful", res.data);
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
