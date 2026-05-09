"use client";

import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";

const formSchema = z.object({
    email: z.email(),
    password: z.string().min(1, { message: "Password is required" }),
});

export default function LoginPage() {
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            email: "",
            password: "",
        },
    });

    function onSubmit(values: z.infer<typeof formSchema>) {
        console.log(values);
    }

    return (
        <main className="min-h-svh flex justify-center items-center p-4">
        <form className="w-full max-w-md" onSubmit={form.handleSubmit(onSubmit)}>
            <Card className="w-full">
                <CardHeader>
                    <CardTitle>Login</CardTitle>
                </CardHeader>
                <CardContent>
                    <FieldGroup>
                    <Field >
                        <FieldLabel>Email</FieldLabel>
                        <Input id="login-email" type="email" placeholder="your@email.com" 
                        aria-invalid={form.formState.errors.email ? "true" : "false"}
                        {...form.register("email")}/>
                        <FieldError>{form.formState.errors.email?.message}</FieldError>

                    </Field>
                        <Field>
                            <FieldLabel>Password</FieldLabel>
                            <Input id="login-password" type="password" placeholder="********" 
                            aria-invalid={form.formState.errors.password ? "true" : "false"}
                            {...form.register("password")}/>
                        </Field>
                        {form.formState.errors.password && <FieldError>{form.formState.errors.password.message}</FieldError>}
                    </FieldGroup>
                        <Button id="login-submit" type="submit" className="w-full mt-8">Login</Button>
                </CardContent>
            </Card>
        </form>
    </main>
  );
}