import { z } from 'zod';
import {
	ZodInstanceIdSchema,
	ZodInstanceCreateSchema,
	ZodInstanceUpdateSchema,
	ZodInstanceQuerySchema,
	ZodInstanceResponseSchema,
} from '@/workflows/schemas/instances/instances.schema';

// DTO types
export type InstanceId = z.infer<typeof ZodInstanceIdSchema>;
export type InstanceCreate = z.infer<typeof ZodInstanceCreateSchema>;
export type InstanceUpdate = z.infer<typeof ZodInstanceUpdateSchema>;
export type InstanceQuery = z.infer<typeof ZodInstanceQuerySchema>;
export type InstanceResponse = z.infer<typeof ZodInstanceResponseSchema>;
