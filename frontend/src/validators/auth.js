import { z } from "zod";

export const loginSchema = z.object({
    email: z
        .string()
        .min(1, "Email or Username is required")
        .pipe(z.email("Invalid email format").or(z.string().min(3, "Username must be at least 3 characters"))),
    password: z
        .string()
        .min(6, "Password must be at least 6 characters"),
});

export const registerSchema = z.object({
    email: z
        .string()
        .min(1, "Email is required")
        .pipe(z.email("Invalid email format")),
    username: z
        .string()
        .min(3, "Username must be at least 3 characters")
        .max(20, "Username must be at most 20 characters")
        .regex(/^[a-zA-Z0-9_]+$/, "Username can only contain letters, numbers and underscores"),
    password: z
        .string()
        .min(8, "Password must be at least 8 characters")
        .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
        .regex(/[a-z]/, "Password must contain at least one lowercase letter")
        .regex(/[0-9]/, "Password must contain at least one number"),
    confirmPassword: z
        .string()
        .min(1, "Please confirm your password"),
}).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
});
