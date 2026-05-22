import { z } from 'zod';
import {
	ZodRoleIdSchema,
	ZodRoleCreateSchema,
	ZodRoleUpdateSchema,
	ZodRoleQuerySchema,
	ZodRoleResponseSchema,
} from '@/mapexos';

/**
 * Role API types
 */
export type RoleId = z.infer<typeof ZodRoleIdSchema>;
export type RoleCreate = z.infer<typeof ZodRoleCreateSchema>;
export type RoleUpdate = z.infer<typeof ZodRoleUpdateSchema>;
export type RoleQuery = z.infer<typeof ZodRoleQuerySchema>;
export type RoleResponse = z.infer<typeof ZodRoleResponseSchema>;
