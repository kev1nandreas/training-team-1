import { z } from "zod";

export const loginSchema = z.object({
	email: z
		.string()
		.min(1, "Email is required")
		.email("Invalid email address"),
	password: z.string().min(1, "Password is required"),
});

export const registerSchema = z.object({
	name: z
		.string()
		.min(2, "Name must be at least 2 characters")
		.max(100, "Name too long"),
	email: z
		.string()
		.min(1, "Email is required")
		.email("Invalid email address"),
	password: z
		.string()
		.min(16, "Password must be at least 16 characters")
		.regex(/[A-Z]/, "Password must contain at least one uppercase letter")
		.regex(/[a-z]/, "Password must contain at least one lowercase letter")
		.regex(/[^A-Za-z0-9]/, "Password must contain at least one symbol")
});

export const profileSchema = z.object({
	name: z
		.string()
		.min(2, "Name must be at least 2 characters")
		.max(100, "Name too long"),
	bio: z.string().max(500, "Bio must be 500 characters or fewer").optional(),
});

export const taskSchema = z.object({
	title: z.string().min(1, "Title is required").max(200, "Title too long"),
	description: z.string().max(1000, "Description too long").optional(),
	priority: z.enum(["low", "medium", "high"]),
});

export const searchSchema = z.object({
	query: z
		.string()
		.min(1, "Search term is required")
		.max(100, "Search term too long")
		.regex(/^[a-zA-Z0-9\s\-_.,!?]+$/, "Search contains invalid characters"),
});
